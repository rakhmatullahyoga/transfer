package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"transfer/transfer"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/joho/godotenv/autoload"
)

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type transferHttpHandler struct {
	uc transfer.TransferUsecase
}

func NewTransferHttpHandler(uc transfer.TransferUsecase) *transferHttpHandler {
	return &transferHttpHandler{
		uc: uc,
	}
}

func (h *transferHttpHandler) Route() (r *chi.Mux) {
	r = chi.NewMux()
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Get("/accounts", h.inquiry)
	r.Post("/transactions", h.transfer)
	r.Post("/callback", h.callback)
	return
}

func writeResponse(w http.ResponseWriter, res Response, status int) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}

func mapError(err error) (res Response, httpStatus int) {
	res = Response{
		Message: err.Error(),
	}
	switch err {
	case transfer.ErrParseInput, transfer.ErrInvalidTransactionStatus:
		httpStatus = http.StatusBadRequest
	case transfer.ErrAccountNotFound:
		httpStatus = http.StatusNotFound
	case transfer.ErrInvalidStatusTransition:
		httpStatus = http.StatusUnprocessableEntity
	default:
		res.Message = "internal server error"
		httpStatus = http.StatusInternalServerError
	}
	return
}

func writeError(w http.ResponseWriter, err error) {
	res, status := mapError(err)
	writeResponse(w, res, status)
}

func (h *transferHttpHandler) inquiry(w http.ResponseWriter, r *http.Request) {
	val, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		err := transfer.ErrParseInput
		writeError(w, err)
		return
	}

	acc, err := h.uc.Inquiry(r.Context(), val.Get("accountNumber"))
	if err != nil {
		writeError(w, err)
		return
	}

	writeResponse(w, Response{Data: acc}, http.StatusOK)
}

func (h *transferHttpHandler) transfer(w http.ResponseWriter, r *http.Request) {
	var params transfer.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		writeError(w, transfer.ErrParseInput)
		return
	}

	trx, err := h.uc.Transfer(r.Context(), params)
	if err != nil {
		writeError(w, err)
		return
	}

	writeResponse(w, Response{Data: trx}, http.StatusOK)
}

func (h *transferHttpHandler) callback(w http.ResponseWriter, r *http.Request) {
	var params transfer.CallbackRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		writeError(w, transfer.ErrParseInput)
		return
	}

	err := h.uc.ProcessCallback(r.Context(), params)
	if err != nil {
		writeError(w, err)
		return
	}

	writeResponse(w, Response{Message: "success update transaction"}, http.StatusOK)
}

UNAME := $(shell uname)
ARCH  := $(shell arch)

tidy:
	go mod tidy

bin:
	@mkdir -p bin

clean: bin
	rm -f bin/transfer

compile: clean
	go build -o bin/transfer main.go

run: compile
	./bin/transfer

env:
	cp env.sample .env

bin/migrate: bin
ifeq ($(UNAME), Linux)
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar zxf - --directory /tmp && cp /tmp/migrate bin/
else ifeq ($(UNAME), Darwin)
ifeq ($(ARCH), arm64) # for Apple processor macs
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.darwin-arm64.tar.gz | tar zxf - --directory /tmp && cp /tmp/migrate bin/
else
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.darwin-amd64.tar.gz | tar zxf - --directory /tmp && cp /tmp/migrate bin/
endif
else
	@echo "Your OS is not supported."
endif

GOOS ?= linux
ARCH ?= amd64

test:
	go list ./... |grep -v vendor | xargs go test -v

build:
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(ARCH) go build -a -installsuffix cgo -o assignment-messagebird .

run-docker-from-local:
	@$(MAKE) build
	docker build -t assignment-messagebird -f Dockerfile.scratch . --build-arg apikey=$(APIKEY)
	docker run -it -p 8080:80 --rm --name assignment-messagebird assignment-messagebird

run-locally:
	go run main.go serve -H 127.0.0.1 -p 8080 --apikey $(APIKEY)

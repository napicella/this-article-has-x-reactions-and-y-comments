.PHONY: fmt build

fmt:
	cd ./src; go fmt

build: fmt
	cd ./src; GOOS=linux GOARCH=amd64 go build -o ../bin/update-article com.napicella


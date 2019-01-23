.PHONY: init install-deps test build pack ship run clean



export CUR_DIR=$(shell pwd)
OS?=$(shell uname -s | tr '[:upper:][:lower:]' '[:lower:][:lower:]')

init:
	go get -u github.com/golang/dep/cmd/dep

install-deps:
	rm -rf ./vendor
	dep ensure -v

build: 
	rm -rf ./dist
	mkdir dist
	GOOS=$(OS) GOARCH=amd64 go build  -o ./dist/strmproc ./...

clean:
	rm -rf ./dist/ 

pack:
	docker build -f ./build/strmproc.Dockerfile -t gattal/strmproc:latest .

upload:
	docker push gattal/strmproc:latest

run:
	docker run --name strmproc --net net_basic gattal/strmproc:latest
	
ship: init test pack upload clean


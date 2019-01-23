FROM gattal/goalpine:latest AS build-env


WORKDIR /go/src/github.com/govinda-attal/fabric-couchdb-to-kafka/strmproc
COPY . .

RUN mkdir dist 

RUN dep ensure -v

RUN	GOOS=linux GOARCH=amd64 go build -o ./dist/strmproc ./...

FROM alpine:3.7
RUN apk -U add ca-certificates

WORKDIR /app
COPY --from=build-env /go/src/github.com/govinda-attal/fabric-couchdb-to-kafka/strmproc/dist/ /app/
CMD /app/strmproc
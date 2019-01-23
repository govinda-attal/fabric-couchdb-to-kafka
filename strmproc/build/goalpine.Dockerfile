FROM golang:1.10-alpine
RUN apk add --no-cache curl bash git openssh
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN apk add --no-cache gcc musl-dev
FROM golang:1.10.3-alpine3.7 as builder

RUN apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep

WORKDIR /go/src/app
COPY ./ ./
RUN dep ensure -vendor-only

RUN go build -o ./output/app

ENTRYPOINT ["/go/src/app/output/app"]
#docker build -t="your-name/image-name" .
#docker push your-name/image-name

FROM golang:1.10.3-alpine3.7 as build

RUN apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep
# RUN apk update
# RUN apk add dep


WORKDIR /go/src/app
COPY ./ ./

RUN dep ensure -vendor-only

RUN go build -o ./output/app

FROM alpine
WORKDIR /app
COPY --from=build /go/src/app/output /app
EXPOSE 80
VOLUME ["/app/dist"]

ENTRYPOINT ["/app/app"]
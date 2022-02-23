FROM golang:alpine AS builder
ENV BUILD_PATH /client

WORKDIR $BUILD_PATH

COPY ./ $BUILD_PATH

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $BUILD_PATH/app $BUILD_PATH/cmd/client/main.go

FROM alpine:3.15.0 as runner

COPY --from=builder /client/app /go/bin/app

CMD ["/go/bin/app"]

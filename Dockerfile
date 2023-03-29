FROM golang:1.17.1-alpine as builder
ADD . /src
RUN cd /src && go build -o app

FROM alpine:latest
WORKDIR /src
COPY --from=builder /src/app /src/
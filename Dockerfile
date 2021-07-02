FROM golang:1.16.5-alpine3.13 as builder
RUN apk add git
WORKDIR /build/hpong/
COPY . /build/hpong/
ENV CGO_ENABLED=0
RUN go mod download && \
    go build -o ./hpong ./app/main.go && \
    chmod +x ./hpong

FROM scratch
#FROM alpine:latest
COPY --from=builder /build/hpong/hpong /go/bin/
EXPOSE 8080
ENTRYPOINT ["/go/bin/hpong"]
#CMD ["/go/bin/hpong"]
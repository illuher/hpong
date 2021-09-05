# Build stage
FROM golang:1.17 as builder
ENV appname=hpong
WORKDIR /build/app/
COPY . /build/app/
ENV CGO_ENABLED=0
RUN go build -o ./$appname ./app/main.go && \
    chmod +x ./$appname

# Alpine, test & debug
FROM alpine:3.14.2 as alpine
COPY --from=builder /build/app/$appname /bin/$appname
ENTRYPOINT ["hpong"]

# Minimal scratch
FROM scratch as scratch
COPY --from=builder /build/app/$appname /bin/$appname
ENTRYPOINT ["hpong"]
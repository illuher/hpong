# Build stage
FROM golang:1.17 as builder
ENV appname=hpong
WORKDIR /build/
COPY . /build/
RUN mkdir -p /tmp/binary/
ENV CGO_ENABLED=0
RUN go build -o /tmp/binary/$appname ./app/main.go && \
    chmod +x /tmp/binary/$appname

# Alpine, test & debug
FROM alpine:3.14.2 as alpine
COPY --from=builder /tmp/binary/$appname /bin/$appname
ENTRYPOINT ["hpong"]

# Minimal scratch
FROM scratch as scratch
COPY --from=builder /tmp/binary/$appname /bin/$appname
ENTRYPOINT ["hpong"]
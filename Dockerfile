# Build stage
FROM golang:1.19 as builder
ENV appname=hpong
WORKDIR /build/
COPY /src/ /build/
RUN mkdir -p /tmp/binary/
RUN CGO_ENABLED=0 go build -o /tmp/binary/$appname ./app/main.go && \
    chmod +x /tmp/binary/$appname

# Alpine, test & debug
FROM alpine:3.16.2 as alpine
COPY --from=builder /tmp/binary/$appname /bin/$appname
ENTRYPOINT ["hpong"]

# Minimal scratch
FROM scratch as scratch
COPY --from=builder /tmp/binary/$appname /bin/$appname
ENTRYPOINT ["hpong"]
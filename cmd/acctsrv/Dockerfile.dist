FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git ca-certificates
ADD . /src
RUN cd /src && go get -d -v && CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-w -s" -o server .

FROM scratch
WORKDIR /app
COPY --from=builder /src/server /app/server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/app/server"]

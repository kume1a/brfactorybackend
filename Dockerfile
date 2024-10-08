FROM golang:1.23-alpine AS builder

WORKDIR /build
RUN apk update && apk upgrade && \
  apk add --no-cache ca-certificates && \
  update-ca-certificates

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

FROM alpine AS app

WORKDIR /app

COPY --from=builder /build/app .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt \
                    /etc/ssl/certs/

ENV BRFACTORY_ENV=production

COPY dockerentrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
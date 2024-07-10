FROM golang-1.20-alpine AS builder

WORKDIR /build
RUN apk update && apk upgrade && \
  apk add --no-cache ca-certificates && \
  update-ca-certificates

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

FROM scratch

WORKDIR /app

COPY --from=builder /build/app .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt \
                    /etc/ssl/certs/

ENTRYPOINT [ "./app" ]
CMD [ "serve", "--http=0.0.0.0:8080" ]
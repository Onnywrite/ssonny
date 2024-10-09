FROM golang:1.23.1-alpine3.20 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o ./bin/sso ./cmd/main.go

FROM alpine:3.20 AS runner
WORKDIR /lib/sso

COPY --from=builder /app/bin ./
RUN adduser -DH ssousr && chown -R ssousr: /lib/sso && chmod -R 700 /lib/sso

USER ssousr
 
CMD [ "./sso" ]
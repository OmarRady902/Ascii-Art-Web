# Stage 1: build
FROM golang:1.25.6-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ascii-art-web .

# Stage 2: run
FROM alpine:latest

LABEL author="Omar rady & Hazem elmahdy" version="1.0" description="ASCII Art web"

WORKDIR /app

COPY --from=builder /app/ascii-art-web .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/banners ./banners

EXPOSE 8080

CMD ["./ascii-art-web"]
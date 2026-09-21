FROM golang:1.22-alpine AS builder
WORKDIR /src

COPY go.mod ./
COPY . .
# No go.sum is checked in (built without local network/Go access); resolve
# and pin dependencies at image build time instead.
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/app

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /out/app ./app

EXPOSE 8080
ENTRYPOINT ["./app"]

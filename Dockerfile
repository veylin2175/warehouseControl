FROM golang:1.23.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /warehouse-control ./cmd/warehouse-control

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /warehouse-control /app/warehouse-control

COPY ./config ./config
COPY ./static ./static

RUN chmod +x /app/warehouse-control

EXPOSE 8080

CMD ["/app/warehouse-control"]
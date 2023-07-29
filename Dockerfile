FROM golang:1.20-alpine AS builder

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

RUN go build -o raychat .

FROM golang:1.20-alpine AS runner

WORKDIR /app

COPY --from=builder /app/raychat .

EXPOSE 8080

CMD ["/app/raychat"]

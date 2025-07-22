FROM golang:1.23.4 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o rate-limiter cmd/main.go

FROM scratch
COPY --from=build /app/rate-limiter .
COPY --from=build /app/.env.docker .env
CMD ["./rate-limiter"]
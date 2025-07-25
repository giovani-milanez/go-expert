FROM golang:1.23.4 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o stress-test cmd/main.go

FROM scratch
COPY --from=build /app/stress-test .
ENTRYPOINT ["./stress-test"]
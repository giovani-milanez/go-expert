FROM golang:1.23.4 as build
WORKDIR /app
COPY ./service-b .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o serviceb
RUN apt-get update && apt-get install -y ca-certificates openssl

FROM scratch
COPY --from=build /app/serviceb .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV OTEL_COLLECTOR_URL="otel-collector:4317"
CMD ["./serviceb"]
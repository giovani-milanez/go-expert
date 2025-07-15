FROM golang:1.23.4 as build
WORKDIR /app
COPY ./service-a .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o servicea
RUN apt-get update && apt-get install -y ca-certificates openssl

FROM scratch
COPY --from=build /app/servicea .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV OTEL_COLLECTOR_URL="otel-collector:4317"
ENV SERVICEB_URL="service-b:8080"
CMD ["./servicea"]
FROM golang:1.23.4 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o cloudrun
RUN apt-get update && apt-get install -y ca-certificates openssl

FROM scratch
COPY --from=build /app/cloudrun .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["./cloudrun"]
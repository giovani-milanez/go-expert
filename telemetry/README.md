## Como rodar:

docker-compose up -d

depois enviar uma requisao POST, exemplo:

POST http://localhost:8081/clima HTTP/1.1
Host: localhost:8081
Content-Type: application/json

{
    "cep":"01001000"
}


Serviço A (responsável pelo input):        Rodando na porta 8081
Serviço B (responsável pela orquestração): Rodando na porta 8080
Zipkin UI                                : http://localhost:9411/
OTEL gRPC                                : localhost:4317

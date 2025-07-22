## Como rodar:

docker-compose up -d --build

Os parametros encontram-se no arquivo .env.docker

Web Server com Rate Limiter : localhost:8080
Redis                       : localhost:6379


Ao realizar GET localhost:8080 devera retornar string 'Hello, World!' e HTTP 200 em caso de sucesso.


Caso a quantidade de requisicoes por segundo for maior que IP_MAX_REQ_PER_SECOND ou TOKEN_MAX_REQ_PER_SECOND (se fornecido o Header API_KEY) ira retornar HTTP 429 (Too many requests).
Apos esse periodo o usuario nao podera realizar requisicoes pela duracao de IP_BLOCK_DURATION/TOKEN_BLOCK_DURATION. Ao expirar esse periodo, os dados de contagem sao resetados.

O codigo esta utilizando Redis para persistencia dos dados do limiter via RedisRateLimiterStorage.
Caso deseje utilizar armazenamento em memoria basta instanciar NewMemoryRateLimiterStorage() ao inves de NewRedisRateLimiterStorage().

Para rodar os testes automatizados:

go test ./...

Atualmente esta utilizando persistencia em memoria.
Para testes com Redis, verificar metodo getStorage no arquivo RateLimiterMiddleware_test.go

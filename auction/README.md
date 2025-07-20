## Como rodar:

docker-compose up -d


Nota: Foi adicionado arquivo cmd/auction/auction.http para testes

## Teste para validar se o fechamento está acontecendo de forma automatizada

go test -timeout 30s -run ^TestAutomaticallyCloseAuction$ giovani-milanez/go-expert/auction/cmd/auction
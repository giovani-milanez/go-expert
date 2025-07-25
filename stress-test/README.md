## Como rodar:

# Cria imagem
docker build -t stress-test .

# Rodar a aplicacao
docker run -it stress-test --url=http://google.com --requests=1000 --concurrency=10


Nota: Para cancelar basta apertar CTRL+C (desde que seja rodado com docker run -it)
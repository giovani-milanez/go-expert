## Como rodar

docker-compose up -d

cd cmd/ordersystem

go run .

  

## Como Testar

Web - Porta 8000:

veja os arquivos api/create_order.http e api/list_orders.http

e.g.

GET http://localhost:8000/order HTTP/1.1

POST http://localhost:8000/order HTTP/1.1

{...}

  

gRPC - Porta 50051:

  

go install github.com/ktr0731/evans@latest

evans -r repl

package pb

service OrderService

call ListOrder

call CreateOrder

  
  

GraphQL - Porta 8080:

Acessar http://localhost:8080/

Utilize esse exemplo:

  

    mutation createOrder {
    
    createOrder(input: { id: "aaa2", Price: 12.2, Tax: 2 }) {
    
    id
    
    Price
    
    Tax
    
    FinalPrice
    
    }
    
    }
    
      
    
    query listOrder {
    
    orders {
    
    id
    
    Price
    
    Tax
    
    FinalPrice
    
    }
    
    }

## RabitMQ

Rodando na porta 15672 credenciais guest/guest

Cria automaticamente fila 'orders' e associa a amq.direct

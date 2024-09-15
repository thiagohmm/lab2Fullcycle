RUN DOCKER COMPOSE
docker compose up --build
após o comando terminar envie o post abaixo


ENDEREÇO DE ACESSSO
<http://127.0.0.1:8080/weather>

Exemplo de Uso:

POST http://127.0.0.1:8080/weather HTTP/1.1
Host: 127.0.0.1:8080
Content-Type: application/json

{
  "cep": "03142010"

}
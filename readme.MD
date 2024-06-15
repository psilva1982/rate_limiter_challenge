# Rate Limiter
Implementação básica de um middleware para rate limiter em Go que utiliza Redis para armazenamento das informações de acesso e Mysql para armazenar dos dados do usuário e está configurado para limitar tanto por IP quanto por token. A aplicação pode ser testada com diferentes cenários de carga e ajustada conforme necessário.

## 1 - Configuração do ambiente
Para configurar o ambiente altere as variáveis no arquivo .env conforme sua necessidade. As variáveis necessárias para executar o projeto são:

### Variáveis do Rate Limiter ###
 - **IP_RATE_LIMIT**: Define o número de requisições por segundo por IP
 - **TOKEN_RATE_LIMIT**: Define o número de requisições por segundo por API_TOKEN
 - **BLOCK_DURATION**: Define em **segundos** o tempo de bloqueio das requisições 

### Variáveis de Acesso a Recursos ###

Altere somente se tiver necessidade, estão setadas para utilizar o Redis e Mysql.

- **REDIS_ADDR:** - Endereço do servidor Redis
- **REDIS_PASSWORD:** - Senha de acesso do Redis
- **REDIS_DB:** - Redis Database
- **MYSQL_HOST:** - Endereço do servidor Mysql
- **MYSQL_ROOT_PASSWORD:** - Senha de root do mysql
- **MYSQL_DATABASE:** - Database Mysql
- **MYSQL_USER:** - Usuário do banco de dados
- **MYSQL_PASSWORD:** - Senha do usuário do banco de dados
- **MYSQL_PORT:** - Porta de acesso do Mysql

**Exemplo de arquivo .env**

```
# Limites de requisições por IP e Token
IP_RATE_LIMIT=5
TOKEN_RATE_LIMIT=10
BLOCK_DURATION=5 # em segundos

# Redis
REDIS_ADDR=redis_ratelimiter:6379
REDIS_PASSWORD=
REDIS_DB=0

# Mysql
MYSQL_HOST=mysql_ratelimiter
MYSQL_ROOT_PASSWORD=root
MYSQL_DATABASE=rate_limiter
MYSQL_USER=user
MYSQL_PASSWORD=password
MYSQL_PORT=3306
```

## 2 - Executando o projeto 

Faça o clone do repositório: 

```
git clone XXXXXXX
```

Inicie o projeto utilizando docker compose

```
docker compose -f "deployments/docker-compose.yml" up
```

Pronto, agora é só testar o projeto acessando o endereço http://localhost:8080 

Para utilizar uma API_KEY é necessário criar um usuário, utilizando a rota http://localhost:8080/signup/. Após criar o usuário é possível obter o API_TOKEN no endereço http://localhost:8080/get-api-key/

## 3 - Documentação (Swagger)

Depois de executar o projeto, o endereço da documentação estará disponível em [http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html)

## 4 - Teste automatizado da solução. 

Para facilitar os testes, existem dois testes automatizados. 

1) **TestIPRateLimiter** - É o teste de acesso **por IP**. Para este teste, foi criado um laço para execução das requisições até o limite permitido. Durante está etapa, o teste verifica se todas as requisições retornaram **status code 200**. Execido o limite, o teste fará **mais duas requisições**, que deverão retornar **status code 429** demonstrando que as requisições foram bloqueadas pelo middelware. Ainda na execução do teste, foi acrescentado um **time.Sleep()** com a duração definida na variável **BLOCK_DURATION** do arquivo **.env**. Passado o tempo de bloqueio definido na variável o teste realizará mais uma chamada para o endpoint que voltará a retornar **status code 200**. 

2) **TestTokenRateLimiter** - É o teste de acesso **por Token**. Neste caso o acesso ao endpoint é feito utilizando um **http header o Token de acesso**. Assim como no teste acima, criado um laço para execução das requisições até o limite permitido. Durante está etapa, o teste verifica se todas as requisições retornaram **status code 200**. Execido o limite, o teste fará **mais duas requisições**, que deverão retornar **status code 429** demonstrando que as requisições foram bloqueadas pelo middelware. Ainda na execução do teste, foi acrescentado um **time.Sleep()** com a duração definida na variável **BLOCK_DURATION** do arquivo **.env**. Passado o tempo de bloqueio definido na variável o teste realizará mais uma chamada para o endpoint que voltará a retornar **status code 200**. 

**DICA** 
- Após os **dois retornos 429**, o teste espera o tempo definido na variável **BLOCK_DURATION** para executar uma nova requisição. Então para efeitos de teste não coloque um valor muito alto nesta variável. 

Para executar o teste basta digitar:
```
docker compose -f "deployments/docker-compose.yml" run --rm app sh -c "go test -v main_test.go"
```
# Nagini Api

Este projeto está sendo criado com a intenção de estudar e ao mesmo tempo criar uma API para um ecommerce com cenários reais

- [Recursos](#recursos)
    - [Documentação APIs (Em Breve)](#documentação-apis)
- [Desenvolvimento](#desenvolvimento)
    - [Requisitos](#requisitos)
    - [Instalação](#instalação)
        - [Docker](#docker-compose)
    - [Configuração](#configuração)
    - [Testes](#Testes)

### Documentação APIs

Para a documentação do projeto será utilizado o [Swagger](https://swagger.io/). Ferramenta que provê interface para testes.

![swagger](./docs/images/swagger.png)

Por padrão a documentação está disponível no endpoint `/swagger-ui.html#/`.

### Catálogo de erros

| Erro | Descrição           | Ocorre quando                                                  |
| ---- | ------------------- | -------------------------------------------------------------- |
|  400 | Bad Request         | Os dados enviados no request estão inválidos                   |
|  404 | Not Found           | O recurso não foi encontrado                                   |
|  500 | Internal Error      | Acontece um erro interno no módulo                             |

## Desenvolvimento

### Requisitos

```

* Golang
* Docker
* Docker Compose
* PostgreSQL 14

```

### Instalação

#### Docker compose:

Acessar a pasta raiz do projeto e executar:

```

https://docs.docker.com/compose/install/
docker-compose up -d

```

### Configuração

Lista de variáveis de ambiente necessárias para a execução da aplicação

| Variável               | Descrição                             |   Tipo   | Obrigatório |  Valor Padrão   |
| ---------------------- | ------------------------------------- | :------: | :---------: | :-------------: |
| DB_NAME          | Nome do banco de dados                |  Texto   |     Não     |    nagini-api    |
| DB_USERNAME      | Usuário para conexão de dados         |  Texto   |     Não     |    nagini-api    |
| DB_PASSWORD      | Senha do usuário para acesso ao banco |  Texto   |     Não     |    nagini-api    |
| DB_HOST          | Host para acesso ao Banco             |  Texto   |     Não     |    localhost    |
| DB_PORT          | Porta para acesso ao Banco            | Numérico |     Não     |      5432       |
| APPLICATION_PORT          | Porta para acesso a Aplicação            | Numérico |     Não     |      8000       |


### Testes

```bash
# unit tests
--

```

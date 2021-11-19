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

Lista de variáveis de ambiente necessárias para a execução da aplicação (Pode ser encontrado um arquivo de exemplo um arquivo chamado env.example na pasta root do seu projeto)

| Variável               | Descrição                             |   Tipo   | Obrigatório |  Valor Padrão   |
| ---------------------- | ------------------------------------- | :------: | :---------: | :-------------: |
| DB_NAME          | Nome do banco de dados                |  Texto   |     Não     |    nagini-api    |
| DB_USERNAME      | Usuário para conexão de dados         |  Texto   |     Não     |    nagini-api    |
| DB_PASSWORD      | Senha do usuário para acesso ao banco |  Texto   |     Não     |    nagini-api    |
| DB_HOST          | Host para acesso ao Banco             |  Texto   |     Não     |    localhost    |
| DB_PORT          | Porta para acesso ao Banco            | Numérico |     Não     |      5432       |
| APPLICATION_PORT          | Porta para acesso a Aplicação            | Numérico |     Não     |      8000       |

### Migrates
É usado uma biblioteca de migrate para criação das roles defaults (USER, ADMIN e SUPPORT)
Os migrations se encontram no diretório ```root/db/migrate``` em aquivos ```.sql``` com o padrão de nome ```[<Data em Mili Segundos>_<Nome>]``` ou criar o arquivo por linha de comando:
``` 'pgmgr migration [\<Nome>]'```
<br>
Para executar o migrate basta executar o comando ``` 'pgmgr db migrate'```
<br>
É importante lembrar que as credenciais do seu banco de dados esteja configurado no arquivo ```.pgmgr.json``` contendo um arquivo de exemplo chamado ```pgmgr.json.example```

### Testes

```bash
# unit tests
--

```

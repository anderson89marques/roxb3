# Roxb3
Esse projeto realiza ingestão, o processamento e a exposição de dados agregados de negociações da B3.

## Principais dependências do projeto

- Golang 1.25
- Gin-gonic
- Goroutines 
- Postgres
- Docker

## Arquitetura do projeto
Para organização lógica e de pastas me inspirei nos projetos que já trabalho e que usam arquitetura hexagonal, mas não segui totalmente pois esse desafio não necessitava.

- `internal/`: Código da aplicação
  - `services/`: Diretório onde fica a camada de serviço.
  - `infra/`: Diretório onde ficam os pacotes de infraestrutura, por exemplo configuração
  - `core/`: Diretório onde fica o core da aplicação como dominio, interfaces(ports), funções de utilidades e etc.
  - `adapters/`: Comunicação com o mundo externo, seja expondo interface(drivers) para ser consumido, ou consumindo recursos externos(drivens)  
- `docs/`: Diretório onde ficam os arquivos do swagger.
- `files/`: Diretório onde ficam os arquivos a serem importatos no formato csv.
- `migrations/`: Arquivos com queries para o banco de dados
- `cmd/`: Os dois entrypoints da aplicação, ingestão e api rest.
- `Dockerfile`: Arquivos de migração do banco de dados
- `Makefile`: Comandos do projeto
- `docker-compose.yml`: Arquivo de configuração usado pelo ```docker compose```

Como o desafio tem um grande volume de dados a serem ingeridos, optei por usar goroutines com o padrão worker pool para fazer a ingestão dos dados para o banco postgres de forma paralela.

Para consulta das informações eu optei por criar um materialized view, permitindo fazer consultas performaticas aos dados 
exigidos no desafio.

## Rodando a aplicação

Para rodar a aplicação você precisar ter o ```Docker``` e o ```Docker Compose``` instalados.
O projeto usa o ```make``` também para consolidar os principais comandos.

### Rodando a Ingestao de Stocks

Existe um diretório chamado ```files``` onde devem ficar os arquivos a serem importados.
Nesse repo já existem uns arquivos que podem ser usados.
Abaixo como executar a ingestão dentro.

```
$ make run-ingest
```

Isso vai rodar o ```cli``` e executar o processo de ingestão dos arquivos que estão no diretório ```files```.
Lembrando que essa importação acontece de forma paralela com goroutines + worker pool.

### Rodando a interface de consulta HTTP REST.

```
$ make run-api
```

Pronto, agora basta acessar a rota ```http://localhost:8080/docs/index.html``` para visualizar a documentação interativa e poder
fazer uso da API desenvolvida.
Existe um endpoint que recebe o ativo(ticker) a ser consultado, que é um parâmetro obrigatório e um parâmetro opcional chamado data inicio(data_inicio).

EX: 
```
$ curl http://localhost:8080/roxb3/stocks?ticker=A1AP34&data_inicio=2025-08-21
```

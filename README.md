# DevGroup



### [Docs](https://www.kucharsky.ru/api)&nbsp;&nbsp;&nbsp;&nbsp;[Schema](https://www.postgresqltutorial.com/postgresql-getting-started/postgresql-sample-database/)

# Getting started
## Build the app
```
docker compose up
```

# Additional
## Make Targets

- **make dev**: Starts app in dev mode, need to install fresh
- **make dev-db**: Starts databases for dev
- **make prod**: docker compose up
- **make migrate**: Migrations, uses goose
- **make migrate-down**: Migration down
- **make api**: Generates html from openapi, uses redocly 

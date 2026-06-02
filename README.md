# GameApp

# Migrations

```bash

sql-migrate up -env="production" -config=infrastructure/dbconfig.yml -limit=1
sql-migrate down -env="production" -config=infrastructure/dbconfig.yml -limit=1
sql-migrate status -env="production" -config=infrastructure/dbconfig.yml

```
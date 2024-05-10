## Command to access DB in Docker Container

```
docker exec -it mobile-wallet-provider-golang-vue-postgresql-db-1 psql -U postgres -d mydb
```

## Starting Mailhog locally (MacOS)

```
brew install mailhog
brew services start mailhog
```
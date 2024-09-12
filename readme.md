# Message API
## Info
This API handles message business for the automatic message sending system.

### Preparation Steps
```shell script
# generate a postgresql container with dummy data.

sh ./test/scripts/test_db.sh

# run this docker command to generate a redis.

docker run --name redis-test -d -p 6379:6379 -e REDIS_PASSWORD=12345 redis:latest redis-server --requirepass 12345

```

### Run
```shell script
# generate docs for every single swagger update.
cd src # go to the src directory.
swag init -g main.go

# run application at local.
go run main.go --ENV=qa # possible ENV opts: qa, prod
```

### Required Packages
```shell script
# swaggo library.
go get -u github.com/swaggo/swag/cmd/swag
go install github.com/swaggo/swag/cmd/swag@latest
```

### Additional Info
```shell script
# swagger url: http://localhost:4001/message-api/swagger/index.html
```


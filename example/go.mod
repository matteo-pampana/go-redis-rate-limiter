module github.com/matteo-pampana/go-redis-rate-limiter/example

go 1.16

replace github.com/matteo-pampana/go-redis-rate-limiter => ./..

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/go-redis/redis/v8 v8.11.5
	github.com/matteo-pampana/go-redis-rate-limiter v0.0.0-20220401212637-de52b26d94e0
)

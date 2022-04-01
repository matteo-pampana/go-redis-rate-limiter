# Go-Redis Rate Limiter

This repository contains a Golang's implementation of a rate limiter based on a Redis DB.

This implementation is based on the Redis Best Practices for rate limiting [explained here](https://redis.com/redis-best-practices/basic-rate-limiting/#:~:text=This%20service%20states%20that%20it,after%20one%20minute%20as%20well.).

## Run the example

To run the example you must have docker installed on your machine. 

Then, go to the `example` folder and run:
```bash
docker-compose up
```

Docker compose will spin up a redis instance and a gin-based webserver.

Go to your favorite browser, or do a `cURL` (or open postman, whatever you like) at the following URL:
```
http://localhost:8081/rate-limiter?name=<name>
```

The rate limiter is based on the `name` query param. 
You can call it with the same name for 10 times, then you have to wait 1 minute to call the endpoint again.

---
Made with ❤️ &nbsp; by <a href="https://github.com/matteo-pampana"/>Matteo Pampana</a>
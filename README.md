# gin-mongo-api-hazelcast

Build a REST API with Golang and MongoDB - Gin-gonic Version

This repository shows the source code for building a user management application with Golang using the Gin-gonic framework and MongoDB And Hazelcast .

## First:
$- docker-compose up -d



## Second: 
## Run This Command Into mongodb-1 shell
$- docker exec -it mongodb-1 sh 
$- mongod --replSet rs0 --bind_ip_all
$- mongosh --eval "rs.initiate({_id: \"rs0\",members: [{_id: 0, host: \"mongodb\"},{_id: 1, host: \"mongodb2\"}]})"

## Three:
$- docker exec -it mongodb-2 sh
$- mongod --replSet rs0 --bind_ip_all

## Four:
$- docker-compose up -d mongoimport

## Five:
$- go run main.go

###  It is integrated with Hazelcast And when the program is running, it automatically reads the data from the database using a goroutine and maps it into Hazelcast. Also, if new data is added, it is automatically added to the Hazelcast map.
سثث


# go-distributed-sys

Check out the article [Building Microservices with Event Sourcing/CQRS in Go using gRPC, NATS Streaming and CockroachDB](https://medium.com/@shijuvar/building-microservices-with-event-sourcing-cqrs-in-go-using-grpc-nats-streaming-and-cockroachdb-983f650452aa)
## Technologies Used: 
* Go
* NATS Streaming
* gRPC
* CockroachDB


## Compile Proto files
Run the command below from the nats-streaming directory:

protoc -I pb/ pb/*.proto --go_out=plugins=grpc:pb

## Set up CockroachDB

#### Create user
cockroach user set shijuvar --insecure

#### Create Database
cockroach sql --insecure -e 'CREATE DATABASE ordersdb'

#### Grant privileges to the shijuvar user
cockroach sql --insecure -e 'GRANT ALL ON DATABASE ordersdb TO shijuvar'

### Start CockroachDB Cluster 

#### Start node 1:
cockroach start --insecure \
--store=ordersdb-1 \
--host=localhost \
--background

#### Start node 2:
cockroach start --insecure \
--store=ordersdb-2 \
--host=localhost \
--port=26258 \
--http-port=8081 \
--join=localhost:26257 \
--background

#### Start node 3:
cockroach start --insecure \
--store=ordersdb-3 \
--host=localhost \
--port=26259 \
--http-port=8082 \
--join=localhost:26257 \
--background

#### Start a SQL Shell for CockroachDB:
cockroach sql \
--url="postgresql://shijuvar@localhost:26257/ordersdb?sslmode=disable";

## Run NATS Streaming Server
nats-streaming-server \
--store file \
--dir ./data \
--max_msgs 0 \
--max_bytes 0

## Basic Workflow in the example:
1. A client app post an Order to an HTTP API.
2. An HTTP API (**orderservice**) receives the order, then executes a command onto Event Store, which is an immutable log of events, to create an event via its gRPC API (**eventstore**). 
3. The Event Store API executes the command and then publishes an event "order-created" to NATS Streaming server to let other services know that an event is created.
4. The Payment service (**paymentservice**) subscribes the event “order-created”, then make the payment, and then create an another event “order-payment-debited” via Event Store API. 
5. The Query syncing workers (**orderquery-store1 and orderquery-store2 as queue subscribers**) are also subscribing the event “order-created” that synchronise the data models to provide state of the aggregates for query views.
6. The Event Store API executes a command onto Event Store to create an event “order-payment-debited” and publishes an event to NATS Streaming server to let other services know that the payment has been debited.
7. The restaurant service (**restaurantservice**) finally approves the order.
8. A Saga coordinator manages the distributed transactions and makes void transactions on failures (to be implemented). 


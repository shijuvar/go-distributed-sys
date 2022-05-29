
# Guidance for building distributed systems and Microservices in Go

## Articles
* [Building Event-Driven Distributed Systems in Go with gRPC, NATS JetStream and CockroachDB](https://shijuvar.medium.com/building-event-driven-distributed-systems-in-go-with-grpc-nats-jetstream-and-cockroachdb-c4b899c8636d)

## Technologies Used: 
* Go
* NATS JetStream
* gRPC
* CockroachDB
* Go kit
* Zipkin (for distributed tracing)


## Compile Proto files
Run the command below from the eventstream directory:

```
protoc eventstore/*.proto \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.
```		


## Set up CockroachDB 

#### Set up CockroachDB  Cluster with three nodes

```
cockroach start \
--insecure \
--store=orders-1 \
--listen-addr=localhost:26257 \
--http-addr=localhost:8080 \
--join=localhost:26257,localhost:26258,localhost:26259 \
--background
```

```
cockroach start \
--insecure \
--store=orders-2 \
--listen-addr=localhost:26258 \
--http-addr=localhost:8081 \
--join=localhost:26257,localhost:26258,localhost:26259 \
--background
```
```
cockroach start \
--insecure \
--store=orders-3 \
--listen-addr=localhost:26259 \
--http-addr=localhost:8082 \
--join=localhost:26257,localhost:26258,localhost:26259 \
--background
```

#### cockroach init command to perform a one-time initialization of the cluster
```
cockroach init --insecure --host=localhost:26257
```

#### Start a SQL Shell for CockroachDB:
```
cockroach sql --insecure --host=localhost:26257
```

#### Create user
```
cockroach user set shijuvar --insecure
```

#### Create Databases
```
cockroach sql --insecure -e 'CREATE DATABASE eventstoredb'
```

```
cockroach sql --insecure -e 'CREATE DATABASE ordersdb'
```

#### Grant privileges to the shijuvar user
```
cockroach sql --insecure -e 'GRANT ALL ON DATABASE ordersdb TO shijuvar'
```
```
cockroach sql --insecure -e 'GRANT ALL ON DATABASE eventstoredb TO shijuvar'
```

## Run NATS JetStream Server 

```
nats-server -js
```

## Run JetStream with config file

```
nats-server -c js.conf
```

// js.conf file
```
jetstream {
    // jetstream data will be in /data/jetstream-server/jetstream
    store_dir: "/data/jetstream-server"

    // 1GB
    max_memory_store: 1073741824

    // 10GB
    max_file_store: 10737418240
}

http_port: 8222
```

## Prerequisites for running the eventstream demo

* Run the NATS JetStream Server
* Run CockroachDB
* Run bootstrapper app (eventstream/bootstrapper) as a one time activity for creating stream in JetStream and create two databases in CockroachDB
* In order to use distributed tracing with gokit-ordersvc (order service with Go kit), run Zipkin distributed tracing system  

## Basic Workflow in the example (eventstream directory):

1. A client app post an Order to an HTTP API (ordersvc)
2. An HTTP API (ordersvc) receives the order, then executes a command onto Event Store, which is an immutable log of events of domain events, to create an event via its gRPC API (eventstoresvc).
3. The Event Store API executes the command and then publishes an event "ORDERS.created" to NATS JetStream server to let other services know that a domain event is created.
4. The Payment worker (paymentworker) subscribes the event "ORDERS.created", then make the payment, and then create another event "ORDERS.paymentdebited" via Event Store API.
5. The Event Store API executes a command onto Event Store to create an event "ORDERS.paymentdebited" and publishes an event to NATS JetStream server to let other services know that the payment has been debited.
6. The Query synchronising worker (querymodelworker) subscribes the event "ORDERS.created" that synchronise the query data model to provide state of the aggregates for query views.
7. The review worker (reviewworker) subscribes the event "ORDERS.paymentdebited" and finally approves the order, and then create another event "ORDERS.approved" via Event Store API.
8. A Saga coordinator manages the distributed transactions and makes void transactions on failures (to be implemented)


## Training and Consulting in India
As a Consulting Solutions Architect, I do provide [training and consulting on Go programming language and distributed systems architectures](https://github.com/shijuvar/shijuvar/blob/master/masterclass.md), in India.


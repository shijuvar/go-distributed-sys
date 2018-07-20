# go-distributed-sys

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

	


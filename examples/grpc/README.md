## Compile proto files
Run the command below from the gRPC directory:

protoc --proto_path=pb pb/*.proto --go_out=pb \
--go_opt=paths=source_relative --go-grpc_out=pb \
--go-grpc_opt=paths=source_relative 

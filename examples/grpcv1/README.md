## Compile proto files
Run the command below from the grpcv1 directory:

protoc --proto_path=customer customer/*.proto --go_out=customer \
--go_opt=paths=source_relative --go-grpc_out=customer \
--go-grpc_opt=paths=source_relative

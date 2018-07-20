build:
	protoc --proto_path pb/ pb/*.proto --go_out=plugins=grpc:pb
buildcs:
	protoc --proto_path pb/ pb/*.proto --csharp_out=plugins=grpc:pb
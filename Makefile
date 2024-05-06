generate_grpc_code:
	protoc --go_out=cyclops --go_opt=paths=source_relative --go-grpc_out=cyclops --go-grpc_opt=paths=source_relative cyclops.proto 

swag:
	swag init --parseDependency

up:
	docker-compose up --build

gen_proto:
	protoc --proto_path=proto proto/user.proto --go_out=internal --go-grpc_out=internal

run:
	go run ./cmd
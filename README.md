## gRPC ATUH

### compile protobuf file 

1. Install the protocol compiler plugins for Go using the following commands:
```
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```
2. Update your PATH so that the protoc compiler can find the plugins:
`$ export PATH="$PATH:$(go env GOPATH)/bin"`

### Regenerate gRPC code
Before you can use the new service method, you need to recompile the updated .proto file.

`$ protoc --proto_path=proto proto/*.proto  --go_out=:pb --go-grpc_out=:pb`
OR
`$ make gen`



### requirements

- install Evans :
`$ go install github.com/ktr0731/evans@latest`

- Install the `uuid-ossp` extension postgresql:
`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`


##### set environment
```
export POSTGRES_DEVELOP_DB_USERNAME=myuser
export POSTGRES_DEVELOP_DB_PASSWORD=mypassword
export POSTGRES_DEVELOP_DB_NAME=test_db
```


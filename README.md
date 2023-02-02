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

```
$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/usermgmt.proto
```
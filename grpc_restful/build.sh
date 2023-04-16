protoc -I./proto -I./proto/googleapis  --go_out=plugins=grpc:. proto/helloworld.proto
protoc -I./proto -I./proto/googleapis  --swagger_out=logtostderr=true:.  ./proto/helloworld.proto
protoc -I./proto -I./proto/googleapis --grpc-gateway_out=logtostderr=true:. ./proto/helloworld.proto
go mod tidy
go build

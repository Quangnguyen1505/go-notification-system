**GRPC**
```bash
protoc --proto_path=. --proto_path="D:\Notification System\Lib Grpc Swagger\googleapis" --proto_path="D:\Notification System\Lib Grpc Swagger\grpc-gateway" --go_out=./gen --go_opt=paths=source_relative --go-grpc_out=./gen --go-grpc_opt=paths=source_relative --grpc-gateway_out=./gen --grpc-gateway_opt=paths=source_relative notification.proto
```
//go:generate protoc --go-grpc_out=../ --go_out=../ --proto_path=../../service-auth ../../service-auth/grpc/auth_service.proto
package grpc

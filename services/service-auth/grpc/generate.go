//go:generate protoc --go-grpc_out=../ --go_out=../ auth_service.proto
//go:generate protoc --go-grpc_out=../ --go_out=../ --proto_path=../../service-orgn ../../service-orgn/grpc/orgn_service.proto
package grpc

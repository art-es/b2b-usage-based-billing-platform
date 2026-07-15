# Auth Service

## Code generation

Run code generation
```sh
go generate ./...
```

To support mock generation use that command on the top of testing file
```go
//go:generate mockgen -source=file.go -destination=file_mock_test.go -package=$GOPACKAGE
```

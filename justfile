# generate client definitions
generate:
    protoc --go_out=pkg/svc/.generated --go_opt=paths=source_relative --go-grpc_out=pkg/svc/.generated --go-grpc_opt=paths=source_relative -I ./pkg/proto ./pkg/proto/*.proto

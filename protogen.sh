#!/bin/bash

protoc -Isrc/proto --go_opt=module=github.com/rexlx/portping --go_out=./src/proto/ \
--go-grpc_opt=module=github.com/rexlx/portping --go-grpc_out=./src/proto/ src/proto/*.proto

go build -o server ./src/server/*.go
go build -o client ./src/client/*.go

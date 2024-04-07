#!/bin/bash

mkdir bin
go build -o bin/sophie-gateway cmd/gateway/main.go
go build -o bin/sophie-file cmd/file/main.go
go build -o bin/sophie-gen cmd/gen/main.go
go build -o bin/sophie-logstash cmd/logstash/main.go
go build -o bin/sophie-schedule cmd/schedule/manager.go
go build -o bin/sophie-schedule-worker cmd/schedule/worker.go
go build -o bin/sophie-system cmd/system/main.go
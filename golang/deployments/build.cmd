@echo off
go build -ldflags="-s -w" -o basl.exe cmd/main.go
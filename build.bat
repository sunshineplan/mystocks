@echo off
go build -ldflags "-s -w"
npm run build

export GOOS=windows GOARCH=amd64
go build windows-config.go
export GOOS=linux GOARCH=amd64
go build linux-config.go

all:
	go fmt .
	go vet .
	golint .

build:
	go build -o handlauf-macos .
	env GOOS=linux GOARCH=amd64 go build -o handlauf-linux

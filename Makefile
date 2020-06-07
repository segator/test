dependency:
	go get -d -v

build-server: dependency
    go generate
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o transcoderd

build-worker: dependency
	go generate
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o transcoderw

build: build-server build-worker


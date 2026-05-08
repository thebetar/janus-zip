run:
	go run cmd/janus-zip/main.go $(FILE)

test:
	go test -v ./...

build:
	go build -o janus-zip cmd/janus-zip/main.go
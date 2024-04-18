dev: 
	go build main.go

build: dev
	@podman build -t "save:local" .

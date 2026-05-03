VERSION ?= dev
LDFLAGS  = -X github.com/mrmackaniel/launchpad/cmd.Version=$(VERSION)

build:
	go build -ldflags "$(LDFLAGS)" -o launchpad .

run:
	go run . $(ARGS)

test:
	go test ./...

docker-build:
	docker build -t launchpad .

docker-run:
	docker run --rm launchpad $(ARGS)

clean:
	rm -f launchpad

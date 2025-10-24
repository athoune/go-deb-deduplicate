build: bin
	go build -o ./bin/rear ./cli/rear
	go build -o ./bin/unar ./cli/unar
	go build -o ./bin/rexz ./cli/rexz
	go build -o ./bin/unxz ./cli/unxz
	go build -o ./bin/deduplicate ./cli/deduplicate

bin:
	mkdir -p bin

test:
	go test -timeout 10s -cover ./archive
	go test -timeout 10s -cover ./chunker
	go test -timeout 10s -cover ./deb
	go test -timeout 10s -cover ./warehouse

docker:
	docker build -t deb-deduplicate .

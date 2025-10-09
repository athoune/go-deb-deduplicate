ar: bin
	go build -o ./bin/rear ./cli/rear
	go build -o ./bin/unar ./cli/unar

bin:
	mkdir -p bin

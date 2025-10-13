ar: bin
	go build -o ./bin/rear ./cli/rear
	go build -o ./bin/unar ./cli/unar
	go build -o ./bin/rexz ./cli/rexz
	go build -o ./bin/unxz ./cli/unxz

bin:
	mkdir -p bin

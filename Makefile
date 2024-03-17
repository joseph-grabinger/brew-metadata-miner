all: clean run

clean:
	rm -rf ./tmp

run:
	go run . 

test:
	go test -v ./...

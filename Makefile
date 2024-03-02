all: clean run

clean:
	rm -rf ./out
	rm -rf ./homebrew-core

run:
	go run . 

build:
	go build -o bin/gored

run: build
	./bin/gored

build:
	go build -o ./bin/rt

render: build
	./bin/rt > image.ppm
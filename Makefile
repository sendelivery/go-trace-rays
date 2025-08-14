build:
	go build -o ./bin/rt

render: build
	./bin/rt -complex > image.ppm

render-parallel: build
	./bin/rt -parallel -complex > image.ppm

debug: build
	./bin/rt

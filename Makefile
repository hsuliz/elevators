build-image-local:
	docker build -t elevators:latest .

run-image-local: build-image-local
	docker run --rm -p 8080:8080 elevators:latest

build-image:
	docker build --provenance=false --platform linux/amd64 -t elevators:latest .

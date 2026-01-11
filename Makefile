build-image-local:
	docker build -t elevators-app .

run-image-local:
	docker run -it -p 8080:8080 elevators-app --rm
SHELL=/bin/bash

test:
	docker build --target test -t test .
	docker run --rm test

develop:
	docker build --target develop -t develop .
	docker run -it --mount type=bind,source="$$(pwd)",target=/app --mount type=bind,source="/tmp",target=/tmp --rm develop

run:
	docker build --target run -t run .
	docker run --rm run
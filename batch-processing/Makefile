SHELL=/bin/bash

test:
	docker build --target test -t test .
	docker run --rm test

run:
	docker build --target run -t run .
	docker run --rm run
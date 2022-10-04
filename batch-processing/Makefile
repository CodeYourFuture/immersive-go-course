SHELL=/bin/bash
.PHONY: run test develop

outputs:
	mkdir -p outputs

run: outputs
	docker build --target run -t run .
	docker run \
		--env-file docker_env \
		--mount type=bind,source="$$(echo $$HOME)/.aws",target=/root/.aws \
		--mount type=bind,source="$$(pwd)/outputs",target=/outputs \
		--rm run

test: outputs
	docker build --target test -t test .
	docker run \
		--env-file docker_env \
		--rm test

develop: outputs
	docker build --target develop -t develop .
	docker run -it \
		--env-file docker_env \
	    --mount type=bind,source="$$(go env GOCACHE)",target=/root/.cache/go-build \
		--mount type=bind,source="$$(echo $$HOME)/.aws",target=/root/.aws \
		--mount type=bind,source="$$(pwd)",target=/app \
		--mount type=bind,source="$$(pwd)/outputs",target=/outputs \
		--mount type=bind,source="/tmp",target=/tmp \
		--rm develop
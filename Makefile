webpack :
	cd frontend && yarn run bundle

js-lint :
	cd frontend && yarn run lint

fmt :
	go fmt ./...
	cd frontend && yarn run prettier

run :
	# Using `concurrently` from npm to run both front-end builds and back-end at
	# the same time.
	cd frontend && \
	npx concurrently "cd .. && go run main.go" "yarn run bundle" \
		--names go,webpack \
		--prefix-colors green.bold,blue.bold

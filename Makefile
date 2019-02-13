webpack :
	cd frontend && yarn run bundle

js-lint :
	cd frontend && yarn run lint

fmt :
	go fmt ./...
	cd frontend && yarn run prettier

run :
	npx concurrently "go run main.go" "cd frontend && yarn run bundle" \
		--names go,webpack \
		--prefix-colors green.bold,blue.bold

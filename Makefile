build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/handleCRCToken handleCRCToken/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/registerWebhook registerWebhook/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/subscribeWebhook subscribeWebhook/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/eventHandler eventHandler/main.go

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: deploy
deploy: clean build
	sls deploy --verbose

.PHONY: build test run

build:
	go vet ./...
	GOARCH=arm64 GOOS=linux go build -ldflags "-s -w" -o bootstrap ./cmd/lambda
	zip web-observer.zip bootstrap
	rm -rf bootstrap
	mkdir -p build
	mv web-observer.zip ./build/
	GOARCH=arm64 GOOS=linux go build -ldflags "-s -w" -o bootstrap ./cmd/bot/lambda
	zip bot.zip bootstrap
	rm -rf bootstrap
	mv bot.zip ./build/

deploy:
	make build
	sam deploy --no-fail-on-empty-changeset --no-confirm-changeset  --capabilities CAPABILITY_NAMED_IAM --stack-name web-observer-test --template-file ./deployment/observer.yml

test:
	go test ./... -v

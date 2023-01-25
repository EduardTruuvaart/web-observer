.PHONY: build test run

build:
	go vet ./...
	GOARCH=arm64 GOOS=linux go build -ldflags "-s -w" -o bootstrap ./cmd/tracker
	zip tracker.zip bootstrap
	rm -rf bootstrap
	mkdir -p build
	mv tracker.zip ./build/

deploy:
	make build
	sam deploy --guided --no-fail-on-empty-changeset --no-confirm-changeset --region eu-west-1 --profile personal --stack-name tracker-test --template-file ./deployment/tracker.yml

test:
	go test ./... -v

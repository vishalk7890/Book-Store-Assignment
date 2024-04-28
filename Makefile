BUILD_PATH = build
BIN_PATH = build/bin
HEALTH_PATH = health
UNIT_TEST_DIR=$(shell go list ./... | grep -v /cmd)
LOCAL_ENV = \
	PORT=8080 \
	GIN_MODE=debug \
	DBHost="localhost" \
	DBPort=5432 \
	DBUser="postgres" \
	DB_PASSWORD="admin" \
	DBName="bookStore"

.PHONY: run
run:
	env $(LOCAL_ENV) go run ./cmd/main.go

.PHONY: build
build:
	go build -o $(BIN_PATH)/app ./cmd/main.go

.PHONY: build-local
build-local:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -v -ldflags '-s -w' -a -tags netgo -installsuffix netgo -o ${BIN_PATH}/bootstrap ./cmd/main.go


.PHONY: all test clean
test:
	go install gotest.tools/gotestsum@latest
	# gotestsum --jsonfile report.json -- -coverprofile=coverage.out -race -v $(UNIT_TEST_DIR)
	go run gotest.tools/gotestsum@latest
	go mod tidy

.PHONY: clean
clean:
	rm -rf $(BUILD_PATH)

.PHONY: lint
lint:
	go fmt ./...
	go vet ./...

.PHONY: staticcheck
staticcheck:
	go fmt ./...
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...
	go mod tidy

.PHONY: gosec
gosec:
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec -fmt=text ./...
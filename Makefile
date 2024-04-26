# BUILD_PATH = build
# BIN_PATH = build/bin
# HEALTH_PATH = health
# #UNIT_TEST_DIR=$(shell go list ./... | grep -v /cmd)
# LOCAL_ENV = \
# 	PORT=8080 \
# 	GIN_MODE=debug \
# 	DBHost="localhost" \
# 	DBPort=5432 \
# 	DBUser="postgres" \ 
# 	DBPassword="admin" \ 
# 	DBName="bookStore" \ 


# .PHONY: run
# run:
# 	$(LOCAL_ENV) \
# 	go run ./cmd/

# # .PHONY: build
# # build:
# # 	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -ldflags '-s -w' -a -tags netgo -installsuffix netgo -o ${BIN_PATH}/bootstrap ./cmd/api/main.go

# # .PHONY: build-local
# # build-local:
# # 	GOOS=linux GOARCH=amd64 go build -o ${BIN_PATH}/app ./cmd/api/main.go

# # .PHONY: all test clean
# # test:
# # 	make unit-test
# # 	make staticcheck
# # 	make govulncheck
# # 	go mod tidy

# # .PHONY: unit-test
# # unit-test:
# # 	go install gotest.tools/gotestsum@latest
# # 	gotestsum --jsonfile report.json -- -coverprofile=coverage.out -race -v $(UNIT_TEST_DIR)
# # 	go mod tidy

# # .PHONY: staticcheck
# # staticcheck:
# # 	go fmt ./...
# # 	go install honnef.co/go/tools/cmd/staticcheck@latest
# # 	staticcheck ./...
# # 	go mod tidy

# # .PHONY: govulncheck
# # govulncheck:
# # 	go install golang.org/x/vuln/cmd/govulncheck@latest
# # 	govulncheck ./...
# # 	go mod tidy

# # .PHONY: gosec
# # gosec:
# # 	go install github.com/securego/gosec/v2/cmd/gosec@latest
# # 	gosec -fmt=text ./...

# # .PHONY: gosec-pipeline
# # gosec-pipeline:
# # 	go install github.com/securego/gosec/v2/cmd/gosec@latest
# # 	gosec -fmt=sonarqube -out=gosec-report.json -no-fail ./...

# # .PHONY: sonar_strict
# # sonar_strict:
# # 	cd /ns && \
# # 	make init DEFAULT_DIR=${BITBUCKET_CLONE_DIR}/app/$(appName) && \
# # 	make scan ARGS="--projectVersion=${VERSION} --login=${SONARCLOUD_TOKEN} --url=${SONARCLOUD_URL} --debug --strict"

# # .PHONY: sonar
# # sonar:
# # 	cd /ns && \
# # 	make init DEFAULT_DIR=${BITBUCKET_CLONE_DIR}/app/$(appName) && \
# # 	make scan ARGS="--projectVersion=${VERSION} --login=${SONARCLOUD_TOKEN} --url=${SONARCLOUD_URL} --debug"

# # .PHONY: deploy
# # deploy:
# # 	make build -B || exit 1
# # 	mkdir -p deploy/bin
# # 	zip -9 -j deploy/bin/lambda.zip build/bin/bootstrap
# # 	echo "start to upload management/$$currentVersion.zip"
# # 	aws s3 cp deploy/bin/lambda.zip s3://${ARTIFACTS_PATH}/management/$$currentVersion.zip
# # 	echo "successfully uploaded management/$$currentVersion.zip"


# # local-deploy:
# # 	make build -B
# # 	mkdir -p deploy/bin
# # 	zip -9 -j deploy/bin/lambda.zip build/bin/app
# # 	aws s3 cp deploy/bin/lambda.zip s3://${ARTIFACTS_PATH_LOCAL}/management/$(shell date +"%Y%m%d%H%M%S").zip



BUILD_PATH = build
BIN_PATH = build/bin
HEALTH_PATH = health
LOCAL_ENV = \
	PORT=8080 \
	GIN_MODE=debug \
	DBHost="localhost" \
	DBPort=5432 \
	DBUser="postgres" \ 
	DBPassword="admin" \ 
	DBName="bookStore" \ 

.PHONY: run
run:
	$(LOCAL_ENV) \
	go run ./cmd/main.go

.PHONY: build
build:
	go build -o $(BIN_PATH)/app ./cmd/main.go

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -rf $(BUILD_PATH)

.PHONY: lint
lint:
	go fmt ./...
	golint ./...
	go vet ./...

.PHONY: staticcheck
staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...
	go mod tidy

.PHONY: mockgen
mockgen:
	mockery --name=Service
	mockery --name=Repository

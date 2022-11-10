-include $(shell curl -sSL -o .build-harness "https://cloudposse.tools/build-harness"; echo .build-harness)

all: init deps lint build test

deps: init
	@go install gotest.tools/gotestsum@latest

fmt:
	@gofumpt -w -l .

lint:
	@golangci-lint -v run

build:
	echo build

test:
	@mkdir reports
	@gotestsum --format testname --junitfile ./reports/unit-tests.xml --jsonfile ./reports/unit-tests.json -- -coverprofile=cover.out ./...
	@go tool cover -func=./reports/cover.out

sonar:
	@sonar-scanner -Dsonar.organization=roberthorrox -Dsonar.projectKey=RobertHorrox_gofstools -Dsonar.sources=. -Dsonar.host.url=https://sonarcloud.io
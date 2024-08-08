ARTIFACT_NAME := tgbot

build:
	@go build -o bin/${ARTIFACT_NAME}/${ARTIFACT_NAME} cmd/${ARTIFACT_NAME}/main.go 

run-build:
	bin/${ARTIFACT_NAME}/${ARTIFACT_NAME}

# run:
# 	@go run cmd/${ARTIFACT_NAME}/main.go

run:
	air --build.cmd "go build -o ./bin/${ARTIFACT_NAME}/${ARTIFACT_NAME} cmd/${ARTIFACT_NAME}/main.go" --build.bin "./bin/${ARTIFACT_NAME}/${ARTIFACT_NAME}"

test:
	@go test ./...

clean:
	@rm -rf bin/${ARTIFACT_NAME}

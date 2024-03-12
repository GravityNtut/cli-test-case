VERSION 0.7
FROM golang:1.20.4-alpine
build-cli:
    GIT CLONE git@github.com:BrobridgeOrg/gravity-cli.git /gravity-cli
    WORKDIR /gravity-cli
    RUN go build -cover
    SAVE ARTIFACT gravity-cli AS LOCAL ./
test-case-docker:
    WORKDIR /cli-test-case
    COPY . .
    ARG GOCOVERDIR=./coverage_data
    RUN go mod download
    RUN go test -p 1 ./...
    RUN go tool covdata textfmt -i=$GOCOVERDIR -o coverage_result.txt
    SAVE ARTIFACT coverage_result.txt AS LOCAL ./
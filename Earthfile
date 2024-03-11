VERSION 0.7
FROM golang:1.20.4-alpine
build-cli:
    GIT CLONE git@github.com:BrobridgeOrg/gravity-cli.git /gravity-cli
    WORKDIR /gravity-cli
    RUN go build
    SAVE ARTIFACT gravity-cli AS LOCAL ./tmp/output
test-case-docker:
    WORKDIR /cli-test-case
    COPY . .
    RUN go mod download
    RUN go test -p 1 ./...
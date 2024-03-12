VERSION 0.7
FROM golang:1.20.4-alpine
build-cli:
    GIT CLONE git@github.com:BrobridgeOrg/gravity-cli.git /gravity-cli
    WORKDIR /gravity-cli
    RUN go build -cover
    SAVE ARTIFACT gravity-cli AS LOCAL ./

dispatcher-image:
    GIT CLONE git@github.com:BrobridgeOrg/gravity-dispatcher.git /gravity-dispatcher
    SAVE IMAGE gravity-dispatcher:latest

integration-test: 
    BUILD +build-cli
    FROM earthly/dind:alpine-3.18-docker-23.0.6-r4
    # install go 1.21.3 for go test in earthly/dind container
    RUN apk update && apk upgrade --available \
        && apk add --no-cache ca-certificates tzdata curl bash net-tools \
        && wget https://golang.org/dl/go1.21.3.linux-amd64.tar.gz \
        && tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz \
        && rm go1.21.3.linux-amd64.tar.gz
    ENV PATH=/usr/local/go/bin:$PATH

    WORKDIR /cli-test-case
    COPY . .
    ARG GOCOVERDIR=./coverage_data
    RUN go mod download

    WITH DOCKER \
        --compose docker-compose.yaml 
        # --load gravity-dispatcher=+dispatcher-image
        RUN go test -p 1 ./...
    END

    # RUN go test -p 1 ./...
    RUN go tool covdata textfmt -i=$GOCOVERDIR -o coverage_result.txt
    SAVE ARTIFACT coverage_result.txt AS LOCAL ./

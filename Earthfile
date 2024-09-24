VERSION 0.7
FROM golang:1.21.3-alpine
build-cli:
    # 47ca8b652eb4b3032934fb348d07b0f7c889d687 固定commit 版本進行測試
    # c16350f4cb4658ff1c204c6782caa2d8e12f9e29
    GIT CLONE --branch 47ca8b652eb4b3032934fb348d07b0f7c889d687 git@github.com:BrobridgeOrg/gravity-cli.git /gravity-cli
    WORKDIR /gravity-cli
    RUN go build -cover
    SAVE ARTIFACT gravity-cli AS LOCAL ./

dep:
    FROM earthly/dind:alpine-3.18-docker-23.0.6-r4
    # install go 1.21.3 for go test in earthly/dind container
    RUN apk update && apk upgrade --available \
        && apk add --no-cache ca-certificates tzdata curl bash net-tools \
        && wget https://golang.org/dl/go1.21.3.linux-amd64.tar.gz \
        && tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz \
        && rm go1.21.3.linux-amd64.tar.gz
    ENV PATH=/usr/local/go/bin:$PATH
    
    WORKDIR /cli-test-case  
    COPY go.mod go.sum ./
    RUN go mod download
    COPY +build-cli/gravity-cli .
    COPY . .
lint:
    FROM +dep
    RUN curl -sfL https://github.com/golangci/golangci-lint/releases/download/v1.55.1/golangci-lint-1.55.1-linux-amd64.tar.gz | \
            tar zx -C /usr/local/bin/ --strip-components=1 golangci-lint-1.55.1-linux-amd64/golangci-lint && \
        curl -sfL https://github.com/goreleaser/goreleaser/releases/download/v1.22.1/goreleaser_Linux_x86_64.tar.gz | tar -xz -C /usr/local/bin/ && \
        chmod +x /usr/local/bin/goreleaser
    COPY go.mod go.sum ./
    RUN go mod download
    COPY .golangci.yml ./
    RUN golangci-lint --version && \
        golangci-lint run --timeout 5m0s ./...

integration-test: 
    FROM +dep
    RUN mkdir -p coverage_data
    ARG GOCOVERDIR=/cli-test-case/coverage_data
    WITH DOCKER \
        --compose docker-compose.yaml
        RUN go test -p 1 ./...  || true
    END
    # 下載輸出測試結果所需相依檔案，並輸出測試Coverage報告
    RUN go get github.com/BrobridgeOrg/gravity-cli/cmd
    RUN go tool covdata textfmt -i=$GOCOVERDIR -o coverage_result.txt
    RUN go tool cover -func=coverage_result.txt
    RUN go tool cover -html=coverage_result.txt -o coverage_result.html
    SAVE ARTIFACT coverage_result.html AS LOCAL .
ci:
    BUILD +integration-test
    BUILD +lint
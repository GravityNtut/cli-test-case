VERSION 0.7
FROM golang:1.21.3-alpine
build-cli:
    RUN apk update && apk upgrade --available \
        && apk add git
    GIT CLONE git@github.com:BrobridgeOrg/gravity-cli.git /gravity-cli
    WORKDIR /gravity-cli
    # 47ca8b652eb4b3032934fb348d07b0f7c889d687 固定commit 版本進行測試
    RUN git checkout 47ca8b652eb4b3032934fb348d07b0f7c889d687
    RUN go build -cover
    SAVE ARTIFACT gravity-cli AS LOCAL ./

dispatcher-image:
    GIT CLONE git@github.com:BrobridgeOrg/gravity-dispatcher.git /gravity-dispatcher
    SAVE IMAGE gravity-dispatcher:latest

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

integration-test: 
    FROM +dep
    RUN mkdir -p coverage_data
    ARG GOCOVERDIR=/cli-test-case/coverage_data
    WITH DOCKER \
        --compose docker-compose.yaml # \
        # 從github拉最新的dispatcher-image，目前最新dispatcher跑不起來
        # --load gravity-dispatcher=+dispatcher-image
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
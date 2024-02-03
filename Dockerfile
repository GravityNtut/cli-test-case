
FROM golang:1.20.4
FROM docker:dind


WORKDIR /test_case
COPY . . 

RUN go mod tidy

CMD [ "go", "test", "./..." ]

FROM golang:1.11.2-alpine as builder

ENV GO111MODULE=on

ARG SOURCE_COMMIT
ARG VERSION=latest
ARG SOURCE_BRANCH=master
ARG USER=xxxcoltxxx

WORKDIR /go/src/github.com/xxxcoltxxx/smsc-balance-exporter
COPY . .

# Install external dependcies
RUN apk add --no-cache ca-certificates curl git

# Compile binary
RUN CGO_ENABLED=0 GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o smsc_balance_exporter -ldflags " \
        -X github.com/prometheus/common/version.Revision=${SOURCE_COMMIT} \
        -X github.com/prometheus/common/version.Version=${VERSION} \
        -X github.com/prometheus/common/version.Branch=${SOURCE_BRANCH} \
        -X github.com/prometheus/common/version.BuildDate=$(date +'%Y-%m-%d_%H:%M:%S') \
        -X github.com/prometheus/common/version.BuildUser=${USER} \
    "

# Copy compiled binary to clear Alpine Linux image
FROM alpine:latest

ARG VERSION=latest

LABEL maintainer="Aleksandr Paramonov<xxxcoltxxx@gmail.com>"
LABEL version="${VERSION}"
LABEL description="Balance exporter for https://smsc.ru service"

WORKDIR /
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/xxxcoltxxx/smsc-balance-exporter .
COPY static ./static
RUN chmod +x smsc_balance_exporter

EXPOSE 9601

CMD ["./smsc_balance_exporter"]

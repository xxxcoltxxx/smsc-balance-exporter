# build stage
FROM golang:1.11 AS build-env
ADD . /src
RUN cd /src && go get ./...
RUN cd /src && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /src/exporter .

# final stage
FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build-env /src/exporter /
ENTRYPOINT [ "/exporter" ]
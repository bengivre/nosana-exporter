FROM --platform=linux/amd64 golang:alpine AS build-stage

WORKDIR /usr/local/go/src/build

COPY src/* go.mod go.sum ./
RUN env GOOS=linux GOARCH=amd64 go build -o nosana_exporter

FROM  --platform=linux/amd64 alpine

WORKDIR /usr/local/bin

COPY --from=build-stage /usr/local/go/src/build/nosana_exporter ./

ENTRYPOINT ["/usr/local/bin/nosana_exporter"]
EXPOSE 8995
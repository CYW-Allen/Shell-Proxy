FROM golang:1.21.3-alpine3.18

COPY ./src /app/src
COPY ./test/sh /app/build

RUN mkdir dist

WORKDIR /app/src

RUN go mod download

RUN go build -o /app/build/server

WORKDIR /app/build

ENTRYPOINT [ "/bin/sh", "-c", "cp /app/build/server /app/dist/server && /app/build/server" ]
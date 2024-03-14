FROM golang:1.21.3-alpine3.18

COPY ./src /app/src
COPY ./demo /app/build

# For exporting the artifact
# RUN mkdir dist

WORKDIR /app/src

RUN go mod download

RUN go build -o /app/build/server

WORKDIR /app/build

# For exporting the artifact
# ENTRYPOINT [ "/bin/sh", "-c", "cp /app/build/server /app/dist/server && /app/build/server" ]
ENTRYPOINT [ "/bin/sh", "-c", "/app/build/server" ]
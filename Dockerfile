# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

ADD . /app
WORKDIR /app
COPY * ./
RUN go mod download
RUN cd cmd/server && go build -o /grpc-demo
EXPOSE 8080
CMD [ "/grpc-demo" ]

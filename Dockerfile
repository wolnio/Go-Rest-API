# syntax=docker/dockerfile:1

FROM golang:1.17.5

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o hello .

EXPOSE 8000

CMD [ "/app/hello" ]
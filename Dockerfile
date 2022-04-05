FROM golang:1.16-alpine3.15 as development
WORKDIR /app
<<<<<<< HEAD
=======

>>>>>>> 41b06dcb6f21789e611960798f39ab40d3916cb7
COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o main ./cmd/main.go

EXPOSE 9090

CMD ["./main"]


FROM golang:1.16-alpine3.15 as pre-production
WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal

RUN go build -ldflags "-s -w" -o main  ./cmd/main.go
EXPOSE 9090

CMD ["./main"]


FROM alpine as production

COPY --from=pre-production  /app/main /

EXPOSE 9090
CMD ["/main"]

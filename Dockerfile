FROM golang:1.16-alpine3.15 as development
WORKDIR /app

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
COPY ./docs ./docs
RUN go build -ldflags "-s -w" -o main  ./cmd/main.go


FROM alpine as production

COPY --from=pre-production  /app/main /

EXPOSE 9090
CMD ["/main"]

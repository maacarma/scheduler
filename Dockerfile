FROM golang:1.22.3-alpine

WORKDIR /scheduler

ADD . .
RUN go mod download

RUN CGO_ENABLED=0 go build -o /usr/local/bin/scheduler ./cmd/scheduler

EXPOSE 7187

ENTRYPOINT ["scheduler"]
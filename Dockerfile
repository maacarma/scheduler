FROM golang:1.22.3-alpine

WORKDIR /scheduler
COPY go.mod go.sum /scheduler/
RUN go mod download

ADD . .
RUN CGO_ENABLED=0 go build -o /usr/local/bin/scheduler ./cmd/scheduler

EXPOSE 7187

ENTRYPOINT ["scheduler"]
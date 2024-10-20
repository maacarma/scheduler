# Build Stage
FROM golang:1.22.3-alpine AS build-stage

WORKDIR /scheduler

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o scheduler-bin ./cmd/scheduler

# Final Stage
FROM alpine:3

WORKDIR /scheduler

COPY --from=build-stage /scheduler/config /scheduler/config

COPY  --from=build-stage /scheduler/pkg/services/tasks/store/postgres/sql /scheduler/pkg/services/tasks/store/postgres/sql

COPY --from=build-stage /scheduler/scheduler-bin /usr/local/bin/scheduler

EXPOSE 7187

ENTRYPOINT [ "scheduler" ]

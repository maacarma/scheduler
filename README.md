# Scheduler

[![Go Report Card](https://goreportcard.com/badge/github.com/maacarma/scheduler?v=1)](https://goreportcard.com/report/github.com/maacarma/scheduler) 
[![Github Workflow](https://github.com/maacarma/scheduler/actions/workflows/go.yaml/badge.svg)](https://github.com/maacarma/scheduler/actions/workflows/go.yaml/badge.svg)

This easy-to-use application lets you schedule recurring and non-recurring HTTP tasks. It supports multiple databases, making it adaptable to any architecture. Tasks can be scheduled using cron expressions or human-readable strings and can be stopped with ease.


## 🚀 Installation

✨ Version `v1.0.0` is available ✨

`docker pull gogree/scheduler`

## 💡 Supported Features

### Generical

* **Multiple databases:** Supports for famous databases like MongoDB, PostgreSQL 

### Highly Configurable

* **Tailor API calls:** Customize your task API requests with headers, authentication, JSON payloads, and more.
* **Flexible scheduling:** Schedule tasks using cron expressions or simple human-readable intervals (e.g., 1 minute, 1 day 3 hours).
* **Robust stop conditions:** Control task execution based on end dates, recurrence count or instant stopping.
* **Multi Zonal UTC** Accepts time configurations based on UTC.

### Monitoring and Alerting (⏰ will be there soon)

* **Real-time notifications:** Receive Slack or email alerts when tasks fail.
* **Detailed logging:** Track historical records of API calls for analysis.
* **Customized alerts:** Set up alerts based on specific conditions to meet your monitoring needs.


## Running 

* to specify the database, you can pass the evn variable `DATABASE` with the value `mongo` or `postgres` to the docker container
* to specify the database connection string, you can pass the evn variable `MONGO_URL` or `POSTGRES_URL` to the docker container

### Deploying to Kubernetes
* sample yaml attached [sample-k8s.yaml](https://github.com/maacarma/scheduler/blob/main/examples/sample-k8s-deployment.yaml)

### Usage
* sample curl attached [sample-curls.md](https://github.com/maacarma/scheduler/examples/sample-curls.md)


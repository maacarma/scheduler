
# Sample curl commands for usage and testing

### Get all the tasks

```bash
$ curl --location "http://localhost:7187/tasks"
```

### Get all the tasks in a namespace
```bash
$ export namespace=mynamespace
$ curl --location "http://localhost:7187/tasks/n/$namespace"
```

### Create a new sample task
```bash
$ curl --location 'http://localhost:7187/tasks' \
--header 'Content-Type: application/json' \
--data '{
    "url": "htts://google.com",
    "method": "GET",
    "namespace": "default",
    "headers": {
        "header1": ["headerv1"]
    },
    "interval": "10s",
    "start_unix": 1725216780,
    "end_unix": 1725216840
}'
```

### Delete an existing task 
```bash
$ export task_id=1
$ curl --location --request DELETE "http://localhost:7187/tasks/$task_id"
```

### Toggle the status
```bash
$ export task_id=1
$ curl --location --request PUT "http://localhost:7187/tasks/$task/status"
```

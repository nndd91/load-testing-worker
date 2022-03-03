# Readme
Small workflow to help load test temporal cluster. You can start multiple StarterWorkflow to have parallel processing.

## Setup
1. install go in your system
2. in `config/config.go`, update your hostPort

## Commands
Register namespace
```
tctl --ns load-test n re
```

Running Worker:
```
go run cmd/main.go
```

Start workflow (Change the first int to the number of workflows you want to start):
```
tctl --ns load-test wf start -wt StarterWorkflow -tq load-test -i 1 -i false -i false
```

To load temporal cluster with timers, we can set the boolean flags to on:
```
tctl --ns load-test wf start -wt StarterWorkflow -tq load-test -i 5 -i true -i true
```


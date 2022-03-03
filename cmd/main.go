package main

import (
	"fmt"
	"github.com/nndd91/load-testing-worker/config"
	"github.com/nndd91/load-testing-worker/loadtesting"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/zap"
)

// Starts a worker process. Any err should result in a panic or fatal and crash the pods.
func NewTemporalWorker(logger *zap.Logger) *worker.Worker {
	// The client and worker are heavyweight objects that should be created once per process.
	serviceClient, err := client.NewClient(client.Options{
		Namespace: config.Namespace,
		HostPort:  config.HostPort,
	})
	if err != nil {
		panic(fmt.Sprintf("unable to create temporal service client: %v", err))
	}

	temporalWorker := worker.New(serviceClient, config.TaskQueue, worker.Options{})

	// Register workflows and activities
	temporalWorker.RegisterWorkflow(loadtesting.StarterWorkflow)
	temporalWorker.RegisterWorkflow(loadtesting.HelloWorldWorkflow)

	temporalWorker.RegisterActivity(loadtesting.HelloWorldActivity)
	temporalWorker.RegisterActivity(loadtesting.StartManyWorkflows)

	logger.Info("Loaded workflows and activities. Creating worker..")
	return &temporalWorker
}

func main() {
	logger, _ := zap.NewDevelopment()
	w := NewTemporalWorker(logger)

	err := (*w).Run(worker.InterruptCh())

	if err != nil {
		logger.Error("Error from worker", zap.Error(err))
	}
}

package loadtesting

import (
	"context"
	"fmt"
	"github.com/nndd91/load-testing-worker/config"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	temporalSdk "go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
	"time"
)

type StartWorkflowParams struct {
	WorkflowName string
	Namespace    string
	Options      client.StartWorkflowOptions
	WorkflowArgs []interface{}
}

func startWorkflow(ctx context.Context, params StartWorkflowParams) (client.WorkflowRun, error) {
	c, err := client.NewClient(client.Options{
		HostPort:  config.HostPort,
		Namespace: config.Namespace,
	})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	wf, err := c.ExecuteWorkflow(ctx, params.Options, params.WorkflowName, params.WorkflowArgs...)
	return wf, err
}

func StartManyWorkflows(ctx context.Context, num int, useTimer bool, useRecursiveTimer bool) error {
	for i := 0; i < num; i++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		_, err := startWorkflow(ctx, StartWorkflowParams{
			WorkflowName: "HelloWorldWorkflow",
			Namespace:    config.Namespace,
			Options: client.StartWorkflowOptions{
				TaskQueue: config.TaskQueue,
			},
			WorkflowArgs: []interface{}{
				fmt.Sprintf("%v", i),
				useTimer,
				useRecursiveTimer,
			},
		})

		if (err != nil) && !temporalSdk.IsWorkflowExecutionAlreadyStartedError(err) {
			return err
		}
		activity.RecordHeartbeat(ctx, i)

	}
	return nil
}

// StarterWorkflow allow us to run many hello world loadtesting
func StarterWorkflow(ctx workflow.Context, num int, useTimer bool, useRecursiveTimer bool) (string, error) {
	logger := workflow.GetLogger(ctx)
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Hour * 6,
		HeartbeatTimeout:       time.Minute,
		RetryPolicy: &temporalSdk.RetryPolicy{
			InitialInterval:    time.Second * 30,
			BackoffCoefficient: 2,
			MaximumInterval:    time.Minute * 5,
		},
	})

	err := workflow.ExecuteActivity(ctx, StartManyWorkflows, num, useTimer, useRecursiveTimer).Get(ctx, nil)
	if err != nil {
		logger.Error("Activity failed.", zap.Error(err))
		return "", err
	}

	return "Completed!", nil
}

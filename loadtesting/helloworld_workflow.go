package loadtesting

import (
	"go.temporal.io/sdk/temporal"
	"time"

	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
)

func runRecursiveTimer(
	ctx workflow.Context,
	selector *workflow.Selector,
	d time.Duration,
	onTimerDone func(),
) {
	timer := workflow.NewTimer(ctx, d)

	onTimerDone()

	(*selector).AddFuture(timer, func(f workflow.Future) {
		runRecursiveTimer(ctx, selector, d, onTimerDone)
	})
}

// HelloWorldWorkflow is the "HelloWorld" loadtesting
func HelloWorldWorkflow(ctx workflow.Context, customerID string, useTimer bool, useRecursiveTimer bool) (string, error) {

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		// Most activity are api calls and they should finish within second. 1 min should be reasonable.
		StartToCloseTimeout: time.Minute * 1,
		// When activity are stuck or retried, we want them to keep retrying up to 30 mins.
		ScheduleToCloseTimeout: time.Minute * 30,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 10,
			BackoffCoefficient: 2,
		},
	})

	logger := workflow.GetLogger(ctx)

	var result string

	err := workflow.ExecuteActivity(ctx, HelloWorldActivity, customerID).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", zap.Error(err))
		return "", err
	}

	endWorkflow := false
	var signalVal string
	signalChan := workflow.GetSignalChannel(ctx, "test-signal")
	selector := workflow.NewSelector(ctx)

	if useTimer {
		timer1 := workflow.NewTimer(ctx, time.Minute*1)
		selector.AddFuture(timer1, func(fut workflow.Future) {})

		timer2 := workflow.NewTimer(ctx, time.Minute*3)
		selector.AddFuture(timer2, func(fut workflow.Future) {})

		timer3 := workflow.NewTimer(ctx, time.Minute*10)
		selector.AddFuture(timer3, func(fut workflow.Future) {
			endWorkflow = true
		})
	}

	if useRecursiveTimer {
		runRecursiveTimer(ctx, &selector, time.Minute*60, func() {
			// Simulate executing activity inside callback func
			_ = workflow.ExecuteActivity(ctx, HelloWorldActivity, "in call back func").Get(ctx, nil)
		})
	}

	selector.AddReceive(signalChan, func(channel workflow.ReceiveChannel, more bool) {
		channel.Receive(ctx, &signalVal)
	})

	for !endWorkflow {
		selector.Select(ctx)
		err = workflow.ExecuteActivity(ctx, HelloWorldActivity, customerID).Get(ctx, &result)
		if err != nil {
			return "", err
		}
	}

	logger.Info("HelloWorld loadtesting completed.", zap.Any("result", result))

	return result, nil
}

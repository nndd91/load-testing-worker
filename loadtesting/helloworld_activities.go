package loadtesting

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
)

// HelloWorldActivity is the activity triggered from "HelloWorld" loadtesting
func HelloWorldActivity(ctx context.Context, id string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "CustomerID:", id)

	// Simulate api calls
	time.Sleep(time.Second * 30)

	// Simulate api response load
	sampleText := "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."

	return "Some arbitrary value to represent a payload. This will be of some length so that we can simulate load on the history database." + sampleText, nil
}

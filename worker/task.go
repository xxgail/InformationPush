package worker

import (
	"context"
	"fmt"
	"github.com/RichardKnop/machinery/v1/tasks"
)

var asyncTaskMap map[string]interface{}

const (
	HelloWorldTaskName          = "HelloWorldTask"
	DeleteAppShareImageTaskName = "DeleteAppShareImageTask"
)

func HelloWorld() error {
	fmt.Println("Hello GaiGai!")
	return nil
}

func SendHelloWorldTask(ctx context.Context, taskName string) {
	args := make([]tasks.Arg, 0)
	task, _ := tasks.NewSignature(taskName, args)
	task.RetryCount = 3
	_, _ = AsyncTaskCenter.SendTaskWithContext(ctx, task)
}

func initAsyncTaskMap() {
	asyncTaskMap = make(map[string]interface{})
	asyncTaskMap[HelloWorldTaskName] = HelloWorld
}

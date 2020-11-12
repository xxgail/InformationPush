package main

import (
	"InformationPush/worker"
	"context"
)

func main() {
	worker.SendHelloWorldTask(context.Background(), "HelloWorldTask")
}

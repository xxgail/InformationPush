package worker

import (
	"InformationPush/lib/redislib"
	"context"
	"fmt"
	"github.com/RichardKnop/machinery/example/tracers"
	"github.com/RichardKnop/machinery/v1"
	mchConf "github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	opentracing_log "github.com/opentracing/opentracing-go/log"
	"reflect"
	"time"
)

var (
	AsyncTaskCenter *machinery.Server
)

func startServer() (*machinery.Server, error) {
	cnf := &mchConf.Config{
		Broker:        "redis://127.0.0.1:6379",
		DefaultQueue:  "ServerTasksQueue",
		ResultBackend: "redis://127.0.0.1:6379",
	}

	server, _ := machinery.NewServer(cnf)

	tasksMap := map[string]interface{}{
		"sayHello":    SayHello,
		"longTask":    LongRunningTask,
		"sendMessage": SendMessage,
	}

	return server, server.RegisterTasks(tasksMap)
}

func NewAsyncTaskWorker() error {
	consumerTag := "TestWorker"

	cleanup, err := tracers.SetupTracer(consumerTag)
	if err != nil {
		log.FATAL.Fatalln("Unable to instantiate a tracer:", err)
	}
	defer cleanup()

	AsyncTaskCenter, err = startServer()
	if err != nil {
	}

	worker := AsyncTaskCenter.NewWorker(consumerTag, 0)
	errorhandler := func(err error) {
		log.ERROR.Println("I am an error handler:", err)
	}
	preTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("I am a start of tasks handler for:", signature.Name)
	}

	postTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("I am an end of tasks handler for:", signature.Name)
	}
	worker.SetPostTaskHandler(postTaskHandler)
	worker.SetErrorHandler(errorhandler)
	worker.SetPreTaskHandler(preTaskHandler)

	return worker.Launch()
}

func SendTest() error {
	cleanup, err := tracers.SetupTracer("sender")
	if err != nil {
		log.FATAL.Fatalln("Unable to instantiate a tracer:", err)
	}
	defer cleanup()

	server, err := startServer()
	if err != nil {
		return err
	}

	var (
		SayOne, SayTwo, SayThree tasks.Signature
		longRunningTask          tasks.Signature
	)

	var initTasks = func() {
		SayOne = tasks.Signature{
			Name: "sayHello",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: "One111111",
				},
			},
		}

		SayTwo = tasks.Signature{
			Name: "sayHello",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: "Two222222",
				},
			},
		}

		SayThree = tasks.Signature{
			Name: "sayHello",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: "Three333333",
				},
			},
		}

		longRunningTask = tasks.Signature{
			Name: "longTask",
		}
	}

	span, ctx := opentracing.StartSpanFromContext(context.Background(), "send")
	defer span.Finish()

	batchID := uuid.New().String()
	span.SetBaggageItem("batch.id", batchID)
	span.LogFields(opentracing_log.String("batch.id", batchID))

	log.INFO.Println("Starting batch:", batchID)
	/*
	 * First, let's try sending a single task
	 */
	initTasks()

	log.INFO.Println("Single task:")

	asyncResult, err := server.SendTaskWithContext(ctx, &SayOne)
	if err != nil {
		return fmt.Errorf("Could not send task: %s", err.Error())
	}

	results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		return fmt.Errorf("Getting task result failed with error: %s", err.Error())
	}
	log.INFO.Printf("1 + 1 = %v\n", tasks.HumanReadableResults(results))

	//initTasks()
	//asyncResult, err = server.SendTaskWithContext(ctx, &longRunningTask)
	//if err != nil {
	//	return fmt.Errorf("Could not send task: %s", err.Error())
	//}
	//
	//results, err = asyncResult.Get(time.Duration(time.Millisecond * 5))
	//if err != nil {
	//	return fmt.Errorf("Getting long running task result failed with error: %s", err.Error())
	//}
	//log.INFO.Printf("Long running task returned = %v\n", tasks.HumanReadableResults(results))

	return nil
}

func Execute() string {
	cleanup, err := tracers.SetupTracer("sender")
	if err != nil {
		log.FATAL.Fatalln("Unable to instantiate a tracer:", err)
	}
	defer cleanup()

	server, err := startServer()
	if err != nil {

	}

	var (
		SayOne tasks.Signature
	)
	str, _ := redislib.GetClient().SPop(context.Background(), "SendMessage").Result()
	fmt.Println(str, reflect.TypeOf(str))
	var initTasks = func() {
		SayOne = tasks.Signature{
			Name: "sendMessage",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: str,
				},
			},
		}
	}

	span, ctx := opentracing.StartSpanFromContext(context.Background(), "send")
	defer span.Finish()

	batchID := uuid.New().String()
	span.SetBaggageItem("batch.id", batchID)
	span.LogFields(opentracing_log.String("batch.id", batchID))

	log.INFO.Println("Starting batch:", batchID)
	/*
	 * First, let's try sending a single task
	 */
	initTasks()

	log.INFO.Println("Single task:")

	asyncResult, err := server.SendTaskWithContext(ctx, &SayOne)
	if err != nil {
		fmt.Errorf("Could not send task: %s", err.Error())
	}

	results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		fmt.Errorf("Getting task result failed with error: %s", err.Error())
	}
	log.INFO.Printf(tasks.HumanReadableResults(results))

	// 发送消息
	result := tasks.HumanReadableResults(results)
	return result
}

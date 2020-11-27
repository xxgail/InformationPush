package worker

import (
	"InformationPush/common"
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
	"github.com/spf13/viper"
	"time"
)

var (
	AsyncTaskCenter *machinery.Server
	AsyncTaskMap    map[string]interface{}
)

var (
	sendMessage = "sendMessage"
)

func initAsyncTaskMap() {
	AsyncTaskMap = make(map[string]interface{})
	AsyncTaskMap[sendMessage] = SendMessageTask
}

func startServer() (*machinery.Server, error) {
	cnf := &mchConf.Config{
		Broker:        "redis://" + viper.GetString("redis.password") + "@" + viper.GetString("redis.addr"),
		DefaultQueue:  "ServerTasksQueue",
		ResultBackend: "redis://" + viper.GetString("redis.password") + "@" + viper.GetString("redis.addr"),
	}

	server, _ := machinery.NewServer(cnf)

	initAsyncTaskMap()

	return server, server.RegisterTasks(AsyncTaskMap)
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

	worker := AsyncTaskCenter.NewWorker(consumerTag, 10)
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

// task init
func workTaskInit(name string, args []tasks.Arg) *tasks.Signature {
	return &tasks.Signature{
		Name: name,
		Args: args,
	}
}

// single task
func singleTask(signature *tasks.Signature) {
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "send")
	defer span.Finish()

	batchID := uuid.New().String()
	span.SetBaggageItem("batch.id", batchID)
	span.LogFields(opentracing_log.String("batch.id", batchID))

	log.INFO.Println("Starting batch:", batchID)
	/*
	 * First, let's try sending a single task
	 */
	log.INFO.Println("Single task:")

	asyncResult, err := AsyncTaskCenter.SendTaskWithContext(ctx, signature)
	if err != nil {
		fmt.Errorf("Could not send task: %s", err.Error())
	}
	asyncResult.GetState().IsSuccess()

	results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		fmt.Errorf("Getting task result failed with error: %s", err.Error())
	}
	log.INFO.Printf(tasks.HumanReadableResults(results))
}

// group task
func groupTask(signatures []*tasks.Signature) {
	fmt.Println(common.GetFileLineNum())
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "PushJobGroup")
	defer span.Finish()

	batchID := uuid.New().String()
	span.SetBaggageItem("batch.id", batchID)
	span.LogFields(opentracing_log.String("batch.id", batchID))

	log.INFO.Println("Group of tasks (parallel execution):")

	group, err := tasks.NewGroup(signatures...)
	if err != nil {
		fmt.Errorf("Error creating group: %s", err.Error())
	}

	asyncResults, err := AsyncTaskCenter.SendGroupWithContext(ctx, group, 10)
	if err != nil {
		fmt.Errorf("Could not send group: %s", err.Error())
	}

	for _, asyncResult := range asyncResults {
		results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
		if err != nil {
			fmt.Errorf("Getting task result failed with error: %s", err.Error())
		}
		common.GetFileLineNum()
		log.INFO.Printf("%v\n", tasks.HumanReadableResults(results))
	}
}

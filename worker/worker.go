package worker

import (
	"github.com/RichardKnop/machinery/v1"
	mchConf "github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
)

var (
	AsyncTaskCenter *machinery.Server
)

func init() {
	tc, err := NewTaskCenter()
	if err != nil {
		panic(err)
	}
	AsyncTaskCenter = tc
}

func NewTaskCenter() (*machinery.Server, error) {
	cnf := &mchConf.Config{
		Broker:        "redis://127.0.0.1:6379",
		DefaultQueue:  "ServerTasksQueue",
		ResultBackend: "redis://127.0.0.1:6379",
	}
	server, err := machinery.NewServer(cnf)
	if err != nil {
		return nil, err
	}
	return server, server.RegisterTasks(asyncTaskMap)
}

func NewAsyncTaskWorker(concurrency int) *machinery.Worker {
	consumerTag := "TestWorker"

	worker := AsyncTaskCenter.NewWorker(consumerTag, concurrency)
	errorhandler := func(err error) {
		log.ERROR.Println("I am an error handler:", err)
	}
	preTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("I am a start of task handler for:", signature.Name)
	}

	postTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("I am an end of task handler for:", signature.Name)
	}
	worker.SetPostTaskHandler(postTaskHandler)
	worker.SetErrorHandler(errorhandler)
	worker.SetPreTaskHandler(preTaskHandler)
	return worker
}

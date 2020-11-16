package main

import (
	"fmt"
	"reflect"
)

//var (
//	app *cli.App
//)
//
//func init() {
//	// Initialise a CLI app
//	app = cli.NewApp()
//	app.Name = "machinery"
//	app.Usage = "machinery worker and send example tasks with machinery send"
//	app.Author = "Richard Knop"
//	app.Email = "risoknop@gmail.com"
//	app.Version = "0.0.0"
//}
//
//func main() {
//	app.Commands = []cli.Command{
//		{
//			Name:  "worker",
//			Usage: "launch machinery worker",
//			Action: func(c *cli.Context) error {
//				if err := send(); err != nil {
//					return cli.NewExitError(err.Error(), 1)
//				}
//				return nil
//			},
//		},
//	}
//}

//{
//	Name:  "worker",
//	Usage: "launch machinery worker",
//	Action: func(c *cli.Context) error {
//		if err := worker(); err != nil {
//			return cli.NewExitError(err.Error(), 1)
//		}
//		return nil
//	},
//},
//{
//	Name:  "send",
//	Usage: "send example tasks ",
//	Action: func(c *cli.Context) error {
//		if err := send(); err != nil {
//			return cli.NewExitError(err.Error(), 1)
//		}
//		return nil
//	},
//},
//}
//
//	// Run the CLI app
//	app.Run(os.Args)
//}
//
//var (
//	AsyncTaskCenter *machinery.Server
//)
//
//func StartServer() (*machinery.Server, error) {
//	cnf := &mchConf.Config{
//		Broker:        "redis://127.0.0.1:6379",
//		DefaultQueue:  "ServerTasksQueue",
//		ResultBackend: "redis://127.0.0.1:6379",
//	}
//
//	server,_ := machinery.NewServer(cnf)
//
//	tasksMap := map[string]interface{}{
//		"sayHello" :    worker2.SayHello,
//		"longTask" :    worker2.LongRunningTask,
//		"sendMessage" : worker2.SendMessage,
//	}
//
//	return server, server.RegisterTasks(tasksMap)
//}
//
//func worker() error {
//	fmt.Println("-------------")
//	consumerTag := "TestWorker"
//
//	cleanup, err := tracers.SetupTracer(consumerTag)
//	if err != nil {
//		log.FATAL.Fatalln("Unable to instantiate a tracer:", err)
//	}
//	defer cleanup()
//
//	AsyncTaskCenter,err = StartServer()
//	if err != nil {
//	}
//
//	worker := AsyncTaskCenter.NewWorker(consumerTag, 0)
//	errorhandler := func(err error) {
//		log.ERROR.Println("I am an error handler:", err)
//	}
//	preTaskHandler := func(signature *tasks.Signature) {
//		log.INFO.Println("I am a start of tasks handler for:", signature.Name)
//	}
//
//	postTaskHandler := func(signature *tasks.Signature) {
//		log.INFO.Println("I am an end of tasks handler for:", signature.Name)
//	}
//	worker.SetPostTaskHandler(postTaskHandler)
//	worker.SetErrorHandler(errorhandler)
//	worker.SetPreTaskHandler(preTaskHandler)
//
//	return worker.Launch()
//}
//
//
//func send() error {
//	cleanup, err := tracers.SetupTracer("sender")
//	if err != nil {
//		log.FATAL.Fatalln("Unable to instantiate a tracer:", err)
//	}
//	defer cleanup()
//
//	server, err := StartServer()
//	if err != nil {
//
//	}
//
//	var (
//		SayOne, SayTwo, SayThree 			tasks.Signature
//		longRunningTask 					tasks.Signature
//	)
//
//	var initTasks = func() {
//		SayOne = tasks.Signature{
//			Name: "sayHello",
//			Args: []tasks.Arg{
//				{
//					Type: "string",
//					Value: "One111111",
//				},
//			},
//		}
//
//		SayTwo = tasks.Signature{
//			Name: "sayHello",
//			Args: []tasks.Arg{
//				{
//					Type: "string",
//					Value: "Two222222",
//				},
//			},
//		}
//
//		SayThree = tasks.Signature{
//			Name: "sayHello",
//			Args: []tasks.Arg{
//				{
//					Type: "string",
//					Value: "Three333333",
//				},
//			},
//		}
//
//		longRunningTask = tasks.Signature{
//			Name: "longTask",
//		}
//	}
//
//	span, ctx := opentracing.StartSpanFromContext(context.Background(), "send")
//	defer span.Finish()
//
//	batchID := uuid.New().String()
//	span.SetBaggageItem("batch.id", batchID)
//	span.LogFields(opentracing_log.String("batch.id", batchID))
//
//	log.INFO.Println("Starting batch:", batchID)
//	/*
//	 * First, let's try sending a single task
//	 */
//	initTasks()
//
//	log.INFO.Println("Single task:")
//
//	asyncResult, err := server.SendTaskWithContext(ctx, &SayOne)
//	if err != nil {
//		return fmt.Errorf("Could not send task: %s", err.Error())
//	}
//
//	results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
//	if err != nil {
//		return fmt.Errorf("Getting task result failed with error: %s", err.Error())
//	}
//	log.INFO.Printf(tasks.HumanReadableResults(results))
//
//	//initTasks()
//	//asyncResult, err = server.SendTaskWithContext(ctx, &longRunningTask)
//	//if err != nil {
//	//	return fmt.Errorf("Could not send task: %s", err.Error())
//	//}
//	//
//	//results, err = asyncResult.Get(time.Duration(time.Millisecond * 5))
//	//if err != nil {
//	//	return fmt.Errorf("Getting long running task result failed with error: %s", err.Error())
//	//}
//	//log.INFO.Printf("Long running task returned = %v\n", tasks.HumanReadableResults(results))
//
//	return nil
//}

func main() {
	m := make(map[string]interface{})
	m["title"] = "标题"
	m["content"] = "content"
	m["pushId"] = []string{"aaa"}
	for _, v := range m {
		fmt.Println(v)
		fmt.Println(reflect.TypeOf(v))
	}
}

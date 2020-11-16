package worker

import (
	"InformationPush/common"
	"context"
	"fmt"
	"github.com/RichardKnop/machinery/example/tracers"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	opentracing_log "github.com/opentracing/opentracing-go/log"
	"reflect"
	"time"
)

func PushJob(param map[string]interface{}, retryCount int, delayTime int64) string {
	cleanup, err := tracers.SetupTracer("TestWorker")
	if err != nil {
		log.FATAL.Fatalln("Unable to instantiate a tracer:", err)
	}
	defer cleanup()

	server, err := startServer()
	if err != nil {

	}

	var args []tasks.Arg
	var arg tasks.Arg
	for k, v := range param {
		arg = tasks.Arg{
			Name:  k,
			Type:  reflect.TypeOf(v).Name(),
			Value: v,
		}
		args = append(args, arg)
	}
	signature := workTaskInit("sendMessage", args)

	fmt.Println(common.GetFileLineNum())

	// 重复次数
	signature.RetryCount = retryCount
	eta := time.Now().UTC().Add(time.Second * time.Duration(delayTime))
	signature.ETA = &eta

	fmt.Println(common.GetFileLineNum())
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

	asyncResult, err := server.SendTaskWithContext(ctx, &signature)
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

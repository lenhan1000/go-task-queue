package main

import (
	"context"
	"fmt"
	"os"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	opentracing_log "github.com/opentracing/opentracing-go/log"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/google/uuid"
	"github.com/urfave/cli"

	testtasks "github.com/lenhan1000/go-task-queue/tasks"
)

var (
	app        *cli.App
	configPath string
)

func handleError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func init() {
	app = cli.NewApp()
	app.Name = "pc-part-crawler"
	app.Usage = "general purpose workers"
	app.Author = "Nhan Phan"
	app.Email = "nhantp1@gmail.com"
	app.Version = "0.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "c",
			Value:       "",
			Destination: &configPath,
			Usage:       "Path to a configuration file",
		},
	}
}

func main() {
	app.Commands = []cli.Command{
		{
			Name:  "worker",
			Usage: "lauch machinery worker",
			Action: func(c *cli.Context) error {
				if err := worker(); err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
		{
			Name:  "send",
			Usage: "send test tasks ",
			Action: func(c *cli.Context) error {
				if err := send(); err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
	}
	app.Run(os.Args)
}

func loadConfig() (*config.Config, error) {
	if configPath != "" {
		return config.NewFromYaml(configPath, true)
	}
	return config.NewFromEnvironment(true)
}

func startServer() (*machinery.Server, error) {
	cnf, err := loadConfig()
	handleError(err, "error")

	server, err := machinery.NewServer(cnf)
	handleError(err, "error")

	tasks := map[string]interface{}{
		"test": testtasks.TestTask,
	}

	return server, server.RegisterTasks(tasks)
}

func worker() error {
	consumerTag := "general_worker"

	server, err := startServer()
	handleError(err, "error")

	worker := server.NewWorker(consumerTag, 10)

	errorHandler := func(err error) {
		log.ERROR.Println("Error Handler: ", err)
	}

	preTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("Pre Task Handler:", signature.Name)
	}

	postTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("Post Task Handler: ", signature.Name)
	}
	worker.SetErrorHandler(errorHandler)
	worker.SetPostTaskHandler(postTaskHandler)
	worker.SetPreTaskHandler(preTaskHandler)

	return worker.Launch()
}

func workers() error {
	consumerTag := "general_worker"

	server, err := startServer()
	handleError(err, "error")
	for i := 1; i < 5; i++ {

	}
	worker := server.NewWorker(consumerTag, 0)

	errorHandler := func(err error) {
		log.ERROR.Println("Error Handler: ", err)
	}

	preTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("Pre Task Handler:", signature.Name)
	}

	postTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("Post Task Handler: ", signature.Name)
	}
	worker.SetErrorHandler(errorHandler)
	worker.SetPostTaskHandler(postTaskHandler)
	worker.SetPreTaskHandler(preTaskHandler)

	return worker.Launch()
}

func send() error {
	server, err := startServer()
	handleError(err, "err")

	var (
		testTask tasks.Signature
	)

	var initTasks = func() {
		testTask = tasks.Signature{
			Name: "test",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: "lmao",
				},
				{
					Type:  "string",
					Value: "lmao2",
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

	initTasks()

	log.INFO.Println("Single tasks:")

	asyncResult, err := server.SendTaskWithContext(ctx, &testTask)
	handleError(err, "LUL")

	results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
	handleError(err, "LUL2")
	log.INFO.Printf("%v", tasks.HumanReadableResults(results))

	return nil
}

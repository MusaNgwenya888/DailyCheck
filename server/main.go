package main

import (
	"DailyCheckSfn/utility"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// This endpoint list step functions executions
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3001",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/listexecutions", func(c *fiber.Ctx) error {
		status := c.Query("status")
		arn := c.Query("arn")
		results := listSfnExecutions(status, arn)
		strResults := strings.Join(results, " ")

		return c.SendString(strResults)
	})

	log.Fatal(app.Listen(":4001"))
}

func listSfnExecutions(statusInput string, sfnArn string) []string {
	a := utility.NewAwssession("")
	var inputs []string

	// Specify the ARN of the state machine to re-run failed executions for
	stateMachineArn := sfnArn
	status := statusInput

	// List the failed executions for the state machine
	listInput := &sfn.ListExecutionsInput{
		StateMachineArn: aws.String(stateMachineArn),
		StatusFilter:    aws.String(status),
	}

	resp, err := a.Sfc.ListExecutions(listInput)
	if err != nil {
		fmt.Println("Error listing executions:", err)
	}

	// Loop through the executions and retrieve the input for each one
	for _, execution := range resp.Executions {
		// Specify the ARN of the execution to retrieve the input for
		executionArn := *execution.ExecutionArn

		// Call the DescribeExecution method to retrieve information about the execution
		descInput := &sfn.DescribeExecutionInput{
			ExecutionArn: aws.String(executionArn),
		}

		descResp, err := a.Sfc.DescribeExecution(descInput)
		if err != nil {
			fmt.Println("Error describing execution:", err)
			continue
		}

		input := *descResp.Input
		inputs = append(inputs, input)
	}
	return inputs
}

// Step Functions alternative - Go example.
//
// Submit an order fulfillment workflow, wait for completion.
// No Step Functions state machine, no JSON definitions, no AWS.
//
// Usage:
//
//	export AXME_API_KEY="your-key"
//	go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AxmeAI/axme-sdk-go/axme"
)

func main() {
	client, err := axme.NewClient(axme.ClientConfig{
		APIKey: os.Getenv("AXME_API_KEY"),
	})
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	ctx := context.Background()

	intentID, err := client.SendIntent(ctx, map[string]any{
		"intent_type":   "intent.workflow.process.v1",
		"to_agent":      "agent://myorg/production/workflow-processor",
		"workflow_id":   "WF-2026-0073",
		"workflow_type": "order_fulfillment",
		"steps": []string{
			"validate_payment",
			"reserve_inventory",
			"schedule_shipping",
			"send_confirmation",
		},
		"order_id": "ORD-2026-99821",
	}, axme.RequestOptions{})
	if err != nil {
		log.Fatalf("send intent: %v", err)
	}
	fmt.Printf("Intent submitted: %s\n", intentID)

	result, err := client.WaitFor(ctx, intentID, axme.ObserveOptions{})
	if err != nil {
		log.Fatalf("wait: %v", err)
	}
	fmt.Printf("Final status: %v\n", result["status"])
}

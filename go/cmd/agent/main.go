// Workflow processor agent - Go example.
//
// Listens for intents via SSE, processes order fulfillment workflows,
// resumes with completion results.
//
// Usage:
//
//	export AXME_API_KEY="<agent-key>"
//	go run ./cmd/agent/
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AxmeAI/axme-sdk-go/axme"
)

const agentAddress = "stepfn-alt-processor-demo"

func handleIntent(ctx context.Context, client *axme.Client, intentID string) error {
	intentData, err := client.GetIntent(ctx, intentID, axme.RequestOptions{})
	if err != nil {
		return fmt.Errorf("get intent: %w", err)
	}

	intent, _ := intentData["intent"].(map[string]any)
	if intent == nil {
		intent = intentData
	}
	payload, _ := intent["payload"].(map[string]any)
	if payload == nil {
		payload = map[string]any{}
	}
	if pp, ok := payload["parent_payload"].(map[string]any); ok {
		payload = pp
	}

	workflowID, _ := payload["workflow_id"].(string)
	if workflowID == "" {
		workflowID = "unknown"
	}
	workflowType, _ := payload["workflow_type"].(string)
	if workflowType == "" {
		workflowType = "unknown"
	}
	orderID, _ := payload["order_id"].(string)
	if orderID == "" {
		orderID = "unknown"
	}
	stepsRaw, _ := payload["steps"].([]any)

	fmt.Printf("  Workflow: %s\n", workflowID)
	fmt.Printf("  Type: %s\n", workflowType)
	fmt.Printf("  Order: %s\n", orderID)
	fmt.Printf("  Steps: %d\n", len(stepsRaw))

	// Execute each step
	for i, stepRaw := range stepsRaw {
		step, _ := stepRaw.(string)
		fmt.Printf("  [%d/%d] %s...\n", i+1, len(stepsRaw), step)
		time.Sleep(1 * time.Second)
	}

	result := map[string]any{
		"action":          "complete",
		"workflow_id":     workflowID,
		"order_id":        orderID,
		"steps_completed": len(stepsRaw),
		"order_status":    "fulfilled",
		"completed_at":    time.Now().UTC().Format(time.RFC3339),
	}

	_, err = client.ResumeIntent(ctx, intentID, result, axme.RequestOptions{})
	if err != nil {
		return fmt.Errorf("resume intent: %w", err)
	}
	fmt.Printf("  Workflow %s completed. Order status: fulfilled\n", workflowID)
	return nil
}

func main() {
	apiKey := os.Getenv("AXME_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: AXME_API_KEY not set.")
	}

	client, err := axme.NewClient(axme.ClientConfig{APIKey: apiKey})
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	ctx := context.Background()

	fmt.Printf("Agent listening on %s...\n", agentAddress)
	fmt.Println("Waiting for intents (Ctrl+C to stop)")

	intents, errCh := client.Listen(ctx, agentAddress, axme.ListenOptions{})

	go func() {
		for err := range errCh {
			log.Printf("Listen error: %v", err)
		}
	}()

	for delivery := range intents {
		intentID, _ := delivery["intent_id"].(string)
		status, _ := delivery["status"].(string)

		if intentID == "" {
			continue
		}

		if status == "DELIVERED" || status == "CREATED" || status == "IN_PROGRESS" {
			fmt.Printf("[%s] Intent received: %s\n", status, intentID)
			if err := handleIntent(ctx, client, intentID); err != nil {
				fmt.Printf("  Error processing intent: %v\n", err)
			}
		}
	}
}

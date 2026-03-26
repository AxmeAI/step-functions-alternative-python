# Step Functions Alternative - Python

AXME is a coordination infrastructure for durable execution of long-running intents across distributed systems.

Step Functions gives you workflow orchestration on AWS. But you pay with JSON state machine definitions, AWS lock-in, and Lambda coupling. Need to add a human approval step? Build a custom callback pattern. Need agents from different services? Build cross-account IAM. AXME replaces the JSON state machine with a simple Python API and works anywhere.

> **Alpha** - Built with [AXME](https://github.com/AxmeAI/axme) (AXP Intent Protocol).
> [cloud.axme.ai](https://cloud.axme.ai) - [hello@axme.ai](mailto:hello@axme.ai)

---

## The Problem

You need to orchestrate a multi-step workflow. Your options:

### AWS Step Functions (JSON state machine + Lambda functions)

```json
{
  "Comment": "Order fulfillment workflow",
  "StartAt": "ValidatePayment",
  "States": {
    "ValidatePayment": {
      "Type": "Task",
      "Resource": "arn:aws:lambda:us-east-1:123456:function:validate-payment",
      "Next": "ReserveInventory",
      "Retry": [{"ErrorEquals": ["States.ALL"], "MaxAttempts": 3}]
    },
    "ReserveInventory": {
      "Type": "Task",
      "Resource": "arn:aws:lambda:us-east-1:123456:function:reserve-inventory",
      "Next": "ScheduleShipping"
    },
    "ScheduleShipping": {
      "Type": "Task",
      "Resource": "arn:aws:lambda:us-east-1:123456:function:schedule-shipping",
      "Next": "SendConfirmation"
    },
    "SendConfirmation": {
      "Type": "Task",
      "Resource": "arn:aws:lambda:us-east-1:123456:function:send-confirmation",
      "End": true
    }
  }
}
```

**Plus:** 4 Lambda functions, IAM roles, CloudFormation/SAM/CDK templates, AWS account, CloudWatch for monitoring. Want human approval? Add a callback pattern with SQS + API Gateway + DynamoDB.

### AXME (4 lines, managed service)

```python
intent_id = client.send_intent({
    "intent_type": "intent.workflow.process.v1",
    "to_agent": "agent://myorg/production/workflow-processor",
    "payload": {"workflow_type": "order_fulfillment", "order_id": "ORD-2026-99821"},
})
result = client.wait_for(intent_id)
```

---

## Quick Start

### Python

```bash
pip install axme
export AXME_API_KEY="your-key"   # Get one: axme login
```

```python
from axme import AxmeClient, AxmeClientConfig
import os

client = AxmeClient(AxmeClientConfig(api_key=os.environ["AXME_API_KEY"]))

intent_id = client.send_intent({
    "intent_type": "intent.workflow.process.v1",
    "to_agent": "agent://myorg/production/workflow-processor",
    "payload": {
        "workflow_id": "WF-2026-0073",
        "workflow_type": "order_fulfillment",
        "steps": ["validate_payment", "reserve_inventory", "schedule_shipping", "send_confirmation"],
        "order_id": "ORD-2026-99821",
    },
})

print(f"Submitted: {intent_id}")
result = client.wait_for(intent_id)
print(f"Done: {result['status']}")
```

### TypeScript

```bash
npm install @axme/axme
```

```typescript
import { AxmeClient } from "@axme/axme";

const client = new AxmeClient({ apiKey: process.env.AXME_API_KEY! });

const intentId = await client.sendIntent({
  intentType: "intent.workflow.process.v1",
  toAgent: "agent://myorg/production/workflow-processor",
  payload: {
    workflowId: "WF-2026-0073",
    workflowType: "order_fulfillment",
    steps: ["validate_payment", "reserve_inventory", "schedule_shipping", "send_confirmation"],
    orderId: "ORD-2026-99821",
  },
});

console.log(`Submitted: ${intentId}`);
const result = await client.waitFor(intentId);
console.log(`Done: ${result.status}`);
```

---

## More Languages

| Language | Directory | Install |
|----------|-----------|---------|
| [Python](python/) | `python/` | `pip install axme` |
| [TypeScript](typescript/) | `typescript/` | `npm install @axme/axme` |
| [Go](go/) | `go/` | `go get github.com/AxmeAI/axme-sdk-go` |

---

## Step Functions vs AXME

| | AWS Step Functions | AXME |
|---|---|---|
| **Workflow definition** | JSON state machine (Amazon States Language) | Simple Python/TS/Go API |
| **Vendor lock-in** | AWS only (Lambda, IAM, CloudFormation) | Cloud-agnostic, managed service |
| **Human approval** | Build it yourself (callback + SQS + API Gateway) | Built-in (8 task types) |
| **Multi-agent** | Cross-account IAM complexity | First-class (agents, services, humans) |
| **Infrastructure** | Lambda + IAM + CloudWatch + CloudFormation | SDK + API key |
| **Setup time** | Hours (IAM, Lambda deploy, state machine) | Minutes (SDK + API key) |
| **Delivery modes** | Push only (Lambda invoke) | 5 modes (SSE, poll, push, inbox, internal) |
| **Observability** | CloudWatch + X-Ray | Real-time SSE lifecycle stream |
| **Best for** | AWS-native microservice orchestration | Agent-era workflows: services + agents + humans |

Step Functions is tightly integrated with AWS. AXME is better when you want a simple API without cloud lock-in, or when your workflows involve human approvals and multi-agent coordination.

---

## How It Works

```
+-----------+  send_intent()   +----------------+   deliver    +-----------+
|           | ---------------> |                | -----------> |           |
|  Client   |                  |   AXME Cloud   |              | Workflow  |
|           | <- wait_for() -- |   (platform)   | <- resume()  | Processor |
|           |                  |                |  with result |  (agent)  |
+-----------+                  |   retries,     |              |           |
                               |   timeouts,    |              | 4 steps:  |
                               |   delivery     |              | pay/inv/  |
                               +----------------+              | ship/conf |
                                                               +-----------+
```

1. Client submits a **workflow intent** with steps and order details
2. Platform **delivers** it to the workflow processor agent
3. Agent **executes** each step (validate, reserve, ship, confirm) and **resumes** with result
4. Client **observes** lifecycle events in real time (SSE stream)
5. Platform handles retries, timeouts, and delivery guarantees

---

## Run the Full Example

### Prerequisites

```bash
# Install CLI (one-time)
curl -fsSL https://raw.githubusercontent.com/AxmeAI/axme-cli/main/install.sh | sh
# Open a new terminal, or run the "source" command shown by the installer

# Log in
axme login

# Install Python SDK
pip install axme
```

### Terminal 1 - submit the scenario

```bash
axme scenarios apply scenario.json
# Note the intent_id in the output
```

### Terminal 2 - start the agent

Get the agent key after scenario apply:

```bash
# macOS
cat ~/Library/Application\ Support/axme/scenario-agents.json | grep -A2 stepfn-alt-processor-demo

# Linux
cat ~/.config/axme/scenario-agents.json | grep -A2 stepfn-alt-processor-demo
```

Run in your language of choice:

```bash
# Python
AXME_API_KEY=<agent-key> python agent.py

# TypeScript (requires Node 20+)
cd typescript && npm install
AXME_API_KEY=<agent-key> npx tsx agent.ts

# Go
cd go && go run ./cmd/agent/
```

### Verify

```bash
axme intents get <intent_id>
# lifecycle_status: COMPLETED
```

---

## Related

- [AXME](https://github.com/AxmeAI/axme) - project overview
- [AXP Spec](https://github.com/AxmeAI/axp-spec) - open Intent Protocol specification
- [AXME Examples](https://github.com/AxmeAI/axme-examples) - 20+ runnable examples across 5 languages
- [AXME CLI](https://github.com/AxmeAI/axme-cli) - manage intents, agents, scenarios from the terminal

---

Built with [AXME](https://github.com/AxmeAI/axme) (AXP Intent Protocol).

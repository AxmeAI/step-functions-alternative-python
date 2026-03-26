"""
Workflow processor agent - processes order fulfillment workflows.

Listens for intents via SSE. Simulates a 4-step order fulfillment
workflow and resumes with completion results.

Usage:
    export AXME_API_KEY="<agent-key>"
    python agent.py
"""

import os
import sys
import time

sys.stdout.reconfigure(line_buffering=True)

from axme import AxmeClient, AxmeClientConfig


AGENT_ADDRESS = "stepfn-alt-processor-demo"


def handle_intent(client, intent_id):
    """Process order fulfillment workflow - 4 steps."""
    intent_data = client.get_intent(intent_id)
    intent = intent_data.get("intent", intent_data)
    payload = intent.get("payload", {})
    if "parent_payload" in payload:
        payload = payload["parent_payload"]

    workflow_id = payload.get("workflow_id", "unknown")
    workflow_type = payload.get("workflow_type", "unknown")
    order_id = payload.get("order_id", "unknown")
    steps = payload.get("steps", [])

    print(f"  Workflow: {workflow_id}")
    print(f"  Type: {workflow_type}")
    print(f"  Order: {order_id}")
    print(f"  Steps: {len(steps)}")

    # Execute each step
    for i, step in enumerate(steps, 1):
        print(f"  [{i}/{len(steps)}] {step}...")
        time.sleep(1)

    result = {
        "action": "complete",
        "workflow_id": workflow_id,
        "order_id": order_id,
        "steps_completed": len(steps),
        "order_status": "fulfilled",
        "completed_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
    }

    client.resume_intent(intent_id, result)
    print(f"  Workflow {workflow_id} completed. Order status: {result['order_status']}")


def main():
    api_key = os.environ.get("AXME_API_KEY", "")
    if not api_key:
        print("Error: AXME_API_KEY not set.")
        print("Run the scenario first: axme scenarios apply scenario.json")
        print("Then get the agent key from ~/.config/axme/scenario-agents.json")
        sys.exit(1)

    client = AxmeClient(AxmeClientConfig(api_key=api_key))

    print(f"Agent listening on {AGENT_ADDRESS}...")
    print("Waiting for intents (Ctrl+C to stop)\n")

    for delivery in client.listen(AGENT_ADDRESS):
        intent_id = delivery.get("intent_id", "")
        status = delivery.get("status", "")

        if not intent_id:
            continue

        if status in ("DELIVERED", "CREATED", "IN_PROGRESS"):
            print(f"[{status}] Intent received: {intent_id}")
            try:
                handle_intent(client, intent_id)
            except Exception as e:
                print(f"  Error processing intent: {e}")


if __name__ == "__main__":
    main()

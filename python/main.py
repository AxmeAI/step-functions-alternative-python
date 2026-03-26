"""
Step Functions alternative - submit an order fulfillment workflow.

Step Functions needs JSON state machine definitions and AWS infrastructure.
AXME needs one intent.

Usage:
    pip install axme
    export AXME_API_KEY="your-key"
    python main.py
"""

import os
from axme import AxmeClient, AxmeClientConfig


def main():
    client = AxmeClient(
        AxmeClientConfig(api_key=os.environ["AXME_API_KEY"])
    )

    # Submit a workflow - replaces Step Functions state machine + Lambda functions
    intent_id = client.send_intent(
        {
            "intent_type": "intent.workflow.process.v1",
            "to_agent": "agent://myorg/production/workflow-processor",
            "payload": {
                "workflow_id": "WF-2026-0073",
                "workflow_type": "order_fulfillment",
                "steps": [
                    "validate_payment",
                    "reserve_inventory",
                    "schedule_shipping",
                    "send_confirmation",
                ],
                "order_id": "ORD-2026-99821",
            },
        }
    )
    print(f"Intent submitted: {intent_id}")

    # Observe lifecycle events in real time (SSE stream, no polling)
    print("Watching lifecycle...")
    for event in client.observe(intent_id):
        status = event.get("status", "")
        print(f"  [{status}] {event.get('event_type', '')}")
        if status in ("COMPLETED", "FAILED", "TIMED_OUT", "CANCELLED"):
            break

    # Fetch final state
    intent = client.get_intent(intent_id)
    print(f"\nFinal status: {intent['intent']['lifecycle_status']}")


if __name__ == "__main__":
    main()

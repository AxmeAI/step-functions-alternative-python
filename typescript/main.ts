/**
 * Step Functions alternative - submit an order fulfillment workflow.
 *
 * Step Functions needs JSON state machine definitions and AWS infrastructure.
 * AXME needs one intent.
 *
 * Usage:
 *   npm install @axme/axme
 *   export AXME_API_KEY="your-key"
 *   npx tsx main.ts
 */

import { AxmeClient } from "@axme/axme";

async function main() {
  const client = new AxmeClient({ apiKey: process.env.AXME_API_KEY! });

  // Submit a workflow - replaces Step Functions state machine + Lambda functions
  const intentId = await client.sendIntent({
    intentType: "intent.workflow.process.v1",
    toAgent: "agent://myorg/production/workflow-processor",
    payload: {
      workflowId: "WF-2026-0073",
      workflowType: "order_fulfillment",
      steps: [
        "validate_payment",
        "reserve_inventory",
        "schedule_shipping",
        "send_confirmation",
      ],
      orderId: "ORD-2026-99821",
    },
  });
  console.log(`Intent submitted: ${intentId}`);

  // Wait for completion - no polling, no state machine, no AWS
  const result = await client.waitFor(intentId);
  console.log(`Final status: ${result.status}`);
}

main().catch(console.error);

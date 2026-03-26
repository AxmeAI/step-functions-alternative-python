/**
 * Workflow processor agent - TypeScript example.
 *
 * Listens for intents via SSE, processes order fulfillment workflows,
 * resumes with completion results.
 *
 * Usage:
 *   export AXME_API_KEY="<agent-key>"
 *   npx tsx agent.ts
 */

import { AxmeClient } from "@axme/axme";

const AGENT_ADDRESS = "stepfn-alt-processor-demo";

async function handleIntent(client: AxmeClient, intentId: string) {
  const intentData = await client.getIntent(intentId);
  const intent = intentData.intent ?? intentData;
  let payload = intent.payload ?? {};
  if (payload.parent_payload) {
    payload = payload.parent_payload;
  }

  const workflowId = payload.workflow_id ?? "unknown";
  const workflowType = payload.workflow_type ?? "unknown";
  const orderId = payload.order_id ?? "unknown";
  const steps: string[] = payload.steps ?? [];

  console.log(`  Workflow: ${workflowId}`);
  console.log(`  Type: ${workflowType}`);
  console.log(`  Order: ${orderId}`);
  console.log(`  Steps: ${steps.length}`);

  // Execute each step
  for (let i = 0; i < steps.length; i++) {
    console.log(`  [${i + 1}/${steps.length}] ${steps[i]}...`);
    await new Promise((r) => setTimeout(r, 1000));
  }

  const result = {
    action: "complete",
    workflow_id: workflowId,
    order_id: orderId,
    steps_completed: steps.length,
    order_status: "fulfilled",
    completed_at: new Date().toISOString(),
  };

  await client.resumeIntent(intentId, result, { ownerAgent: AGENT_ADDRESS });
  console.log(`  Workflow ${workflowId} completed. Order status: ${result.order_status}`);
}

async function main() {
  const apiKey = process.env.AXME_API_KEY;
  if (!apiKey) {
    console.error("Error: AXME_API_KEY not set.");
    process.exit(1);
  }

  const client = new AxmeClient({ apiKey });

  console.log(`Agent listening on ${AGENT_ADDRESS}...`);
  console.log("Waiting for intents (Ctrl+C to stop)\n");

  for await (const delivery of client.listen(AGENT_ADDRESS)) {
    const intentId = delivery.intent_id;
    const status = delivery.status;

    if (!intentId) continue;

    if (["DELIVERED", "CREATED", "IN_PROGRESS"].includes(status)) {
      console.log(`[${status}] Intent received: ${intentId}`);
      try {
        await handleIntent(client, intentId);
      } catch (e) {
        console.error(`  Error processing intent: ${e}`);
      }
    }
  }
}

main().catch(console.error);

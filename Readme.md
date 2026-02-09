# Purpose of Building this project
The purpose of Building this project is that an backend services which accepts the events and deliver the notification asynchoronously, reliable and at scale , even in the presence of failures.


Ek aisa backend banana hai jo:

Event accept kare

Turant response de

Notification baad mein bheje

Fail ho to retry kare

Crash ho to data na lose kare

👉 Ye real-world problem hai.

🎯 2️⃣ Why NOT a simple CRUD app?

✍️ Likho:

CRUD synchronous hota hai

Failure handling nahi hoti

Retry logic nahi hota

Real production problems nahi cover hote

Is project ka goal:

Failure ko normal banana

📥 3️⃣ Inputs (System kya lega?)

✍️ Event input:

{
  "event_type": "USER_SIGNUP",
  "recipient": "user@email.com",
  "channel": "email",
  "payload": {}
}


Likho:

Event type

Recipient

Delivery channel

Payload (flexible)

📤 4️⃣ Outputs (System kya karega?)

✍️ Expected behavior:

HTTP request → 202 Accepted

Notification eventually delivered

Ya failure ke baad retry

Final failure me mark as FAILED

Important:

Client ko turant response mile, delivery guarantee later

🔁 5️⃣ Sync vs Async (VERY IMPORTANT)

✍️ Decision:

HTTP layer = synchronous

Delivery = asynchronous

Reason likho:

Email slow hota hai

Webhook unreliable hota hai

User wait nahi karega

🧱 6️⃣ Core Components (boxes draw karo)

✍️ Draw this:

Client
  ↓
HTTP API
  ↓
Database (Event Store)
  ↓
Worker Engine
  ↓
Notifier (Email/Webhook)


Likho:

DB is source of truth

Workers are stateless

💥 7️⃣ Failure Scenarios (IMPORTANT SECTION)

✍️ Likho (bullet points):

Email provider down

Network timeout

Duplicate request

Server crash during processing

Worker restart

For each failure:
👉 System should retry, not lose data

🔂 8️⃣ Retry Strategy (Concept only)

✍️ Likho:

Retry with exponential backoff

Max retries = N

After that → DEAD / FAILED

No code yet. Just idea.

🔐 9️⃣ Idempotency (One-liner)

✍️ Likho:

Same event should not create multiple notifications.

Reason:

Client retries

Network flakiness

🧠 10️⃣ Non-Goals (This is senior thinking)

✍️ Likho what we are NOT doing:

No UI

No auth (for now)

No real email provider initially

No Kafka initially

This keeps scope sane.

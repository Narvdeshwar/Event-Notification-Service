# Interview Questions: Event-Driven Notification Service

This document contains a comprehensive set of interview questions and a deep dive into the project's architecture and flow.

---

## 🏗️ Deep Architecture & Project Flow

This system is built as a **Distributed Task Queue** using PostgreSQL as the persistent engine. Here is the step-by-step lifecycle of a notification:

### 1. The Ingestion Phase (Sync)
*   **Action:** A client sends a POST request to `/events`.
*   **Logic:** The `API Handler` validates the JSON. If valid, it calls the `NotificationService.Enqueue()`.
*   **Persistence:** A record is inserted into the `notifications` table with status `PENDING`.
*   **Response:** The API returns `202 Accepted` immediately. The client doesn't wait for the email to be sent.

### 2. The Discovery Phase (Polling)
*   **Component:** The `Poller` runs in a background goroutine.
*   **Strategy:** It queries the DB every few seconds using `FOR UPDATE SKIP LOCKED`.
    *   `FOR UPDATE`: Locks the rows so no other poller instance can pick them up.
    *   `SKIP LOCKED`: If another instance is already processing a row, this instance simply skips it instead of waiting.
*   **Handoff:** The poller fetches a batch (e.g., 10-50 records) and pushes them into a **Go Channel**.

### 3. The Execution Phase (Workers)
*   **Component:** Multiple `Worker` goroutines (configurable) are listening to the Go Channel.
*   **Processing:**
    *   A worker pulls a job from the channel.
    *   It calls the `Notifier` (Email/Webhook).
    *   **Success:** If the notifier succeeds, the worker updates the DB status to `SENT`.
    *   **Failure:** If it fails, the worker calculates the **Exponential Backoff** (e.g., attempt 1 = 2s, attempt 2 = 4s) and updates `next_retry_at`.

### 4. The Resilience Phase (Recovery)
*   **Stuck Jobs:** If a worker crashes *after* picking up a job but *before* finishing it, the job stays in `PROCESSING` status.
*   **Recovery:** A separate routine periodically checks for jobs that have been in `PROCESSING` for more than 5 minutes and resets them to `PENDING` so they can be retried.
*   **DLQ:** If a job fails more than `MaxRetries` (e.g., 5 times), it is moved to the **Dead Letter Queue (DLQ)** for manual investigation.

---

## 🟢 Basic Level (Junior / Entry)

## 🟢 Basic Level (Junior / Entry)
*Focus: Project overview, fundamental Go concepts, and basic asynchronous patterns.*

### 1. What is the core problem this project solves?
- **Answer:** It solves the issue of slow or unreliable notification delivery. Instead of making a user wait for an email to be sent (synchronous), it accepts the request immediately, stores it in a database, and processes it in the background (asynchronous).

### 2. Why use HTTP status code `202 Accepted`?
- **Answer:** `202` signals that the request has been received and validated but is not yet processed. This is the standard for asynchronous APIs.

### 3. What is the role of the Database in this project?
- **Answer:** The database acts as a **Persistent Queue** and the **Source of Truth**. It ensures that if the service crashes, no events are lost because they are permanently stored before being processed.

### 4. Explain the difference between "Sync" and "Async" in the context of this project.
- **Answer:**
    - **Sync:** The API layer where the client waits for a response from our server.
    - **Async:** The Worker layer where notifications are sent to external providers (Email/Webhooks) without the client waiting.

---

## 🟡 Intermediate Level (Mid-Level)
*Focus: Implementation details, Concurrency, and Resilience.*

### 5. Explain the "Poller and Worker" pattern used here.
- **Answer:**
    - **Poller:** A background routine that periodically queries the DB for `PENDING` jobs and pushes them into a Go channel.
    - **Worker:** Multiple goroutines that listen to that channel, process the jobs, and update the DB status. This decouples data fetching from execution and allows for easy horizontal scaling.

### 6. Why did you use `SELECT ... FOR UPDATE SKIP LOCKED`?
- **Answer:** This is critical for high availability. It allows multiple instances of the service to run simultaneously. `FOR UPDATE` locks the rows so other workers don't grab them, and `SKIP LOCKED` ensures workers don't hang waiting for a lock, instead moving to the next available job.

### 7. How does the Retry Strategy work?
- **Answer:** If a notification fails, the system increments the `Attempts` counter and calculates a `next_retry_at` time (exponential or linear backoff). The job remains in the DB, and the Poller will pick it up again only when the current time exceeds the `next_retry_at` timestamp.

### 8. What is the purpose of the `RecoverStuckJob` logic?
- **Answer:** If a worker picks up a job (sets it to `PROCESSING`) but then the server crashes, that job remains stuck. `RecoverStuckJob` scans for jobs that have been in `PROCESSING` status for an unusually long time (e.g., > 5 mins) and resets them to `PENDING`.

---

## 🔴 Advanced Level (Senior / Architect)
*Focus: Scalability, Trade-offs, and Distributed Systems.*

### 9. What are the pros and cons of using PostgreSQL as a Queue vs. Kafka?
- **Answer:**
    - **Postgres Pros:** Simple architecture, ACID compliance (no lost data), easy to query status.
    - **Postgres Cons:** Performance degrades as the table grows; locking overhead at very high scale (10k+ events/sec).
    - **Kafka:** Better for massive scale and event streaming, but adds significant operational complexity and makes "exactly-once" delivery harder to manage.

### 10. How do you ensure Idempotency?
- **Answer:** By using a unique ID (like a UUID) for every event. If a client retries a request with the same ID, the database's unique constraint prevents us from creating a duplicate notification.

### 11. How would you handle a "Thundering Herd" problem if 1 million notifications are scheduled for the same second?
- **Answer:**
    - Implement a **Rate Limiter** on the Poller to fetch in small batches.
    - Use a **Leaky Bucket** or **Token Bucket** algorithm to ensure external providers (SendGrid, Twilio) aren't overwhelmed.

### 12. Discuss the trade-offs of the Go Channel size in the Poller.
- **Answer:**
    - A small channel (or unbuffered) provides backpressure; the poller won't fetch more if workers are slow.
    - A large channel can act as a buffer for bursts but risks losing queued jobs if the process crashes (since jobs are moved from DB to memory).

---

## 🚀 Scenario-Based Questions

### 13. "An external email API is returning 429 Too Many Requests. How does your system react?"
- **Answer:** The system should treat `429` as a retryable error. The worker will update the `next_retry_at` time and put the job back in the queue, allowing for a cooldown period.

### 14. "We need to prioritize Password Reset emails over Marketing Newsletters. How would you change your SQL query?"
- **Answer:** Add a `priority` column (integer) to the `notifications` table and change the Poller query to `ORDER BY priority DESC, created_at ASC`.

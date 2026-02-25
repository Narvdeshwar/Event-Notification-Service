# Interview Questions: Event-Driven Notification Service

This document contains a comprehensive set of interview questions and a deep dive into the project's architecture and flow.

---

## 🏗️ Deep Architecture & Project Flow

This system is an **Event-Driven Asynchronous Worker System**. It is designed to handle high-write loads and guarantee delivery even if individual components fail.

### 1. Ingestion Phase (The Gatekeeper)
*   **Layer:** Synchronous HTTP API (Go Gin Framework).
*   **Process:**
    1.  The client hits the `POST /events` endpoint.
    2.  The `API Handler` performs **Schema Validation** (checking types, mandatory fields).
    3.  **Atomic Persistence:** The `NotificationService` generates a unique `UUID` and inserts the event into the `notifications` table with `status = 'PENDING'`.
    4.  **Feedback Loop:** Once the DB confirms the write, the API returns a `202 Accepted`.
*   **Why?** Writing to a local DB is orders of magnitude faster than sending an email via a remote provider. This keeps the client's "Time to First Byte" (TTFB) very low.

### 2. The Bridge (Poller & Go Channels)
*   **Layer:** Background Control Loop.
*   **Core Logic:** The `Poller` is a single goroutine that runs a `time.Ticker`.
*   **Transactional Fetching:**
    ```sql
    SELECT * FROM notifications
    WHERE status = 'PENDING' AND (next_retry_at IS NULL OR next_retry_at <= NOW())
    ORDER BY created_at ASC
    LIMIT 50
    FOR UPDATE SKIP LOCKED
    ```
    *   **FOR UPDATE:** Tells Postgres "I am working on these rows, don't let anyone else touch them."
    *   **SKIP LOCKED:** This is the "Secret Sauce." It prevents workers from waiting on each other. If Instance A has locked the first 50 rows, Instance B will automatically skip them and grab rows 51-100.
*   **The Channel:** Once fetched, the Poller pushes these records into a **Buffered Go Channel**. This channel acts as an in-memory pressure valve.

### 3. Execution Phase (The Worker Pool)
*   **Layer:** Concurrent Worker Goroutines.
*   **Mechanism:** On startup, the system spawns **N worker goroutines** (e.g., 20 Workers).
*   **Competition:** All workers share the same Go Channel. Go's scheduler ensures that jobs are distributed fairly among them.
*   **Notifier Interface:** We use an `Interface` for the Notifier. This allows us to swap a "Mock Email Service" for "SendGrid" or "Twilio" without changing the core worker logic.

### 4. Failure & Retry Strategy (Exponential Backoff)
*   **The Problem:** Most failures (Network timeout, API Rate limit) are temporary.
*   **The Solution:**
    1.  If a send fails, we increment the `attempts` column.
    2.  We calculate the **Backoff Delay**: `delay = 2^attempts * seconds`.
    3.  We set the `status` back to `PENDING` and update `next_retry_at`.
    4.  The Poller will ignore this job until the backoff time has passed.

### 5. The "Safety Net" (Recovery & DLQ)
*   **Cleanup Routine:** A separate background task scans for jobs stuck in `PROCESSING` for > 5 minutes. This handles cases where a worker crashed mid-transaction.
*   **Dead Letter Queue (DLQ):** Once a job hits `attempts = 5`, it is marked as `FAILED` or moved to a `failed_notifications` table. This prevents "Poison Pills" (invalid data that causes crashes) from retrying forever and wasting resources.

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

## 🛠️ Engineering & Operational Excellence
*Focus: Logging, Scaling, Testing, and Metrics.*

### 15. How do you monitor the health of this system?
- **Answer:** We use the `internal/metrics` package to track the number of successful vs failed notifications. In a production environment, we would expose these via a `/metrics` endpoint for **Prometheus** to scrape, allowing us to build **Grafana** dashboards to visualize the throughput and error rates.

### 16. How did you (or would you) test this system?
- **Answer:**
    - **Unit Tests:** Mock the `Store` and `Notifier` interfaces to test the `NotificationService` and `Worker` logic in isolation.
    - **Integration Tests:** Use **Docker Compose** or `Testcontainers` to run a real Postgres instance and verify the `SKIP LOCKED` logic and retry state transitions.
    - **Load Testing:** Use tools like `k6` or `Locust` to flood the `/events` endpoint and monitor how the DB handles the write load.

### 17. How is the configuration managed?
- **Answer:** It uses environment variables (handled in `internal/config`) to follow **12-Factor App** principles. This allows the same Docker image to run in Dev, Staging, and Production by simply changing the env vars (like `DB_URL` and `WORKER_COUNT`).

### 18. How would you secure this API?
- **Answer:** Since it's currently open, I would add a **Middleware** layer for:
    - **API Key/JWT Auth:** To ensure only authorized services can trigger notifications.
    - **Rate Limiting:** To prevent a single client from flooding the system and affecting others.

---

## 🐹 Go Runtime & Concurrency Deep Dive
*Focus: Language specifics that interviewers love.*

### 19. Why use `context.Context` throughout the service?
- **Answer:** To handle **Timeouts** and **Cancellations**. If the API request is cancelled by the client, the context propagates that cancellation to the DB query, saving system resources. In workers, it helps in graceful shutdowns.

### 20. What happens if a Worker goroutine panics?
- **Answer:** If not recovered, a panic in a single goroutine will crash the **entire application**. In a production settings, I would add a `defer recover()` block inside each worker's `Start()` loop to log the error and keep the other workers alive.

### 21. Is your `NotificationService` thread-safe?
- **Answer:** Yes, because it is **stateless**. It doesn't store data in local variables; it pushes everything to the database. In Go, stateless services are naturally thread-safe for concurrent use.

---

## ⚖️ Trade-offs & Where This System Fails
*Every architecture has a breaking point. A senior candidate knows exactly where theirs is.*

### 22. What are the key Trade-offs of using a DB as a Queue?
- **Trade-off 1: Consistency vs. Throughput.** By using PostgreSQL, we gain **ACID transactions** (we never lose a notification). However, we trade off **latency**, as a DB query is slower than an in-memory broadcast.
- **Trade-off 2: Simplicity vs. Scalability.** Using a DB is easy to deploy and monitor. But as you scale to millions of events per hour, the `SKIP LOCKED` overhead and vacuuming of the `notifications` table become a performance bottleneck.

### 23. Where will this project fail (Scaling Limits)?
- **Database CPU Bottleneck:** When the `notifications` table grows to millions of rows, the Poller's `SELECT` queries (even with indexes) will consume significant CPU, slowing down both ingestion and processing.
- **Poller Latency:** If the Poller runs every 5 seconds, even if a worker is free, a notification might sit in the DB for 4.9 seconds before being picked up.
- **Connection Exhaustion:** If you scale the number of Workers too high, you might run out of PostgreSQL database connections.

### 24. How would you fix these failures for 100x traffic?
- **Switch to Kafka/RabbitMQ:** Move the "Queue" part out of the DB into a dedicated message broker.
- **Database Partitioning:** Partition the `notifications` table by `created_at` or `status` to keep the active index size small.
- **Change Data Capture (CDC):** Instead of polling the DB, use a tool like **Debezium** to listen to the Postgres WAL (Write Ahead Log) and push events to workers in real-time.

---

---

## 🧠 Technical Deep Dive: Why and How

This section provides a granular explanation of the logic blocks that make up your service.

### 1. The Poller: The "Heartbeat" of Orchestration
The Poller is the bridge between **Persistence** (Disk) and **Execution** (Memory).
*   **The Ticker Loop:** It uses a `time.NewTicker`. This ensures it pulses at a steady rhythm (e.g., every 2 seconds). Unlike a simple `time.Sleep`, a Ticker maintains a more consistent cadence even if a single pass takes longer.
*   **Batching for Efficiency:** It doesn't pick up 1 task at a time. It uses a `LIMIT 50`. Each database query has "overhead" (TCP handshake inside a pool, query parsing). Batching amortizes that overhead, allowing you to process 1,000 tasks with only 20 DB trips.
*   **Handoff & Backpressure:** Once jobs are fetched, the Poller pushes them into a **Buffered Go Channel**.
    *   **The "Pressure Valve":** If your Workers are slow (e.g., the Email API is down), the channel will fill up. Go's runtime will then **block** the Poller from fetching more. This prevents your service from running out of RAM (OOM) by "stuffing" millions of records into memory when they can't be processed.

### 2. The Worker Pool: Managed Concurrency
Instead of creating a new goroutine for every request (which is dangerous at scale), we use a "Fixed Worker Pool."
*   **The Problem:** Unbounded goroutines lead to context-switching overhead and potential memory exhaustion.
*   **The Solution:** On startup, we spawn **N workers** (long-lived goroutines).
*   **Channel Competition:** All workers do `job := <- channel`. The Go Scheduler efficiently wakes up an idle worker when a job arrives.
*   **Pool Lifecycle:** Workers are kept alive for the lifetime of the program. This saves the cost of creating and destroying goroutines constantly.

### 3. Postgres `SKIP LOCKED`: High-Performance Concurrency
In traditional SQL, if Instance A and Instance B both query `SELECT * FROM notifications`, they might get the same row. If they both try to lock it, one waits.
*   **The Magic:** `FOR UPDATE SKIP LOCKED` marks a row as "Under maintenance" and tells the *other* process: "If you see this row is locked, don't wait—just skip it and find a free one."
*   **The Outcome:** This allows you to scale your backend horizontally (run 10 copies of the app) and be **100% certain** that no email will ever be sent twice by two different servers.

### 4. Exponential Backoff: Strategic Resilience
We don't retry "immediately" because the cause of failure is often a temporary system overload.
*   **The Math:** We use `2^attempts`.
    *   Attempt 1: 2s delay.
    *   Attempt 2: 4s delay.
    *   Attempt 3: 8s delay.
*   **The Reasoning:** This delay gives the external system (like SendGrid) breathing room to recover. It avoids a "Retry Storm" where your system accidentally DDoS's its own providers.

### 5. Graceful Shutdown with `sync.WaitGroup`
Shutdown is the most dangerous part of an async system. You don't want to kill a process while a worker is halfway through an email.
*   **The Counter:** `sync.WaitGroup` tracks active jobs.
*   **The Logic:**
    1.  Main receives a `SIGTERM`.
    2.  It sends a "signal" to the Poller to stop fetching (via `context.Cancel`).
    3.  It calls `wg.Wait()`.
    4.  The Workers finish their *current* job and then exit.
*   **The Result:** No "Ghost Tasks" – every notification is either fully sent or safely stored in the DB for the next restart.

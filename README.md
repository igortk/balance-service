# 📘 Balance Service

**Balance Service** is a core microservice in the exchange system that manages user balances.  
It listens to order events via **RabbitMQ**, processes them atomically, and updates balances accordingly in the database.

---

## 📌 Features

- 🧾 **Get User Balance** – Return full user balance (available + locked).
- 💰 **Emit User Balance** – Manually credit a user’s balance in a given currency (for testing or admin purposes).
- 🔄 **Handle Order Events** – Listens to `OrderUpdateEvent` via RabbitMQ and processes balance updates (BUY/SELL logic).
- 🔐 **ACID Transactions** – All updates are safely wrapped in SQL transactions for consistency.

---

## 🐇 RabbitMQ Event Handling

### 📥 Subscribed Queue: `order.events`

**Event:** `OrderUpdateEvent` (sent from Order Service)  
**Format (proto):**
```proto
message OrderUpdateEvent {
  string id = 1;
  Order order = 2;
  MatchedUser matched_user = 4;
  error.Error error = 3;
}

🛡️ Transaction Safety (ACID)

Balance updates are wrapped in transactional logic using:

func (cl *Client) UpdateBalancesTx(ctx context.Context, db *sql.DB, users ...*User) (err error)

If any error occurs — the entire operation is rolled back to ensure consistency.


🧰 Tech Stack

🐹 Go (Golang)

🐇 RabbitMQ

📦 PostgreSQL

⚡ gRPC + Protocol Buffers

☁️ Docker (soon)

✅ SQL Transactions

👨‍💻 Maintainers

Made with ❤️ by the Exchange Platform TeamPRs, issues, and forks welcome!
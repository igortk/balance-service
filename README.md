# 📘 Balance Service

**Balance Service** is a microservice in the exchange system that manages user balances.  
It listens to order events/requests via **RabbitMQ**, processes them atomically, and updates balances accordingly in the database.

---

## 📌 Features

- 🧾 **Get User Balance** – Return full user balance (available + locked).
- 💰 **Emit User Balance** – Manually credit a user’s balance in a given currency (for testing or admin purposes).
- 🔄 **Handle Order Events** – Listens to `OrderUpdateEvent` via RabbitMQ and processes balance updates (BUY/SELL logic).
- 🔐 **ACID Transactions** – All updates are safely wrapped in SQL transactions for consistency.

---
### 🧰 Tech Stack

- 🐹 **Go (Golang)** — core business logic and services
- 🐇 **RabbitMQ** — asynchronous communication between services
- 📦 **PostgreSQL** — relational storage for balances and orders
- ⚡ **gRPC + Protobuf** — efficient binary communication format
- ✅ **SQL Transactions** — ensures ACID-compliant balance updates
---
## 🐇 RabbitMQ Request/Event Handling
### 🔄 **Handle Order Events**

**Event:** `OrderUpdateEvent` (sent from Order Processing Service)  
**Format (proto):**
```proto
message OrderUpdateEvent {
  string id = 1;
  Order order = 2;
  MatchedUser matched_user = 4;
  error.Error error = 3;
}
```

#### 📦 Order Structure:
```proto
message Order {
  string order_id = 1;
  string user_id = 2;
  string pair = 3; // e.g., "USD/EUR"
  double init_volume = 4;
  double fill_volume = 5;
  double init_price = 6;
  OrderStatus status = 7;
  Direction direction = 8;
  int64 updatedDate = 9;
  int64 createdDate = 10;
}
```

#### 👤 Matched User:
```proto
message MatchedUser {
  string user_id = 1;
  double volume = 2;
  double price = 3;
}
```
#### 🔢 Enums
```proto
enum OrderStatus {
  ORDER_STATUS_UNDEFINED = 0;
  ORDER_STATUS_NEW = 1;
  ORDER_STATUS_MATCHED = 2;
  ORDER_STATUS_DONE = 3;
  ORDER_STATUS_REMOVED = 4;
}

enum Direction {
  ORDER_DIRECTION_UNDEFINED = 0;
  ORDER_DIRECTION_BUY = 1;
  ORDER_DIRECTION_SELL = 2;
}
```
---

### 🔄 Business Logic

The balance logic depends on the order status and direction:

#### 🟡 On `ORDER_STATUS_NEW`
The user's funds are **locked** for trade execution:

- **BUY**  
  Lock `quote` currency:  
  `locked_balance = init_price × init_volume`

- **SELL**  
  Lock `base` currency:  
  `locked_balance = init_volume`

#### 🟢 On `ORDER_STATUS_MATCHED`
Funds are transferred between users:

- **BUYER**
  - 🔓 Unlock `quote` currency (reduce locked)
  - 💰 Increase available balance in `base` currency (e.g., BTC)

- **SELLER**
  - 🔓 Unlock `base` currency (reduce locked)
  - 💰 Increase available balance in `quote` currency (e.g., USD)

### 💰 **Emit User Balance**
```proto
message EmitBalanceByUserIdRequest{
  string id = 1;
  string user_id = 2;
  string currency_name = 3;
  double amount = 4;
}

message EmitBalanceByUserIdResponse{
  string id = 1;
  string user_id = 2;
  Balance balance = 3;
  error.Error error = 4;
}
```
### **Get User Balance**
```proto
message GetBalanceByUserIdRequest{
  string id = 1;
  string user_id = 2;
}

message GetBalanceByUserIdResponse{
  string id = 1;
  string user_id = 2;
  repeated Balance user_balance = 3;
  error.Error error = 4;
}

message Balance{
  string currency = 1;
  double balance = 2;
  double locked_balance = 3;
  int64 updated_date = 4;
}
```
---

### 🛡️ Transaction Safety (ACID)

Balance updates are wrapped in transactional logic using:

```go
func (cl *Client) UpdateBalancesTx(ctx context.Context, db *sql.DB, users ...*User) (err error)
```
If any error occurs — the entire operation is rolled back to ensure consistency.

---

### 🐳 Running via docker
_Coming soon_

---

Made with ❤️ by the Ihor Tkachenko, issues, and forks welcome!
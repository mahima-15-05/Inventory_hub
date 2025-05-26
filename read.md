#  Inventory Control Hub

The **Inventory Control Hub** is a centralized system designed to efficiently manage and track inventory operations, including restocking and unstocking of products. It ensures real-time synchronization between stock levels and order transactions, preventing discrepancies and maintaining data integrity.

---

## 🔧 Key Modules

### 1.  Product Module
- Manages the catalog of available products.
- Ensures **uniqueness** of each product; duplicate entries are strictly prevented.
- Stores essential product details such as name, SKU, price, and quantity in stock.

### 2.  Purchase Order Module
- Handles incoming stock by processing purchase orders.
- When a purchase order is created, the system **increments** the inventory count for the respective products.

### 3.  Sales Order Module
- Manages customer sales transactions.
- When a sales order is placed, the system **decrements** the available inventory accordingly.

---

## 🚀 Getting Started

### 🔧 Prerequisites

- Go installed (e.g., Go 1.20+)
- MySQL running (recommended: XAMPP server)
- Postman (optional, for API testing)

---

Ensure the following are installed on your system:
- [Go](https://golang.org/doc/install)
- [MySQL](https://www.mysql.com/) (via XAMPP or other)
- [Postman](https://www.postman.com/) *(optional for API testing)*

---

### 📂 Project Setup

```bash
# 1. Clone the repository
git clone https://github.com/mahima-15-05/Inventory_hub.git
cd Inventory_hub

# 2. Environment file (.env) is already included

# 3. Install Go dependencies
go mod tidy

# 4. Install 'air' if not already installed 
go install github.com/cosmtrek/air@latest

# 6. Start the application (Choose one)
# a) With live reload
air

# b) Without live reload
go run main.go
# 🛒 Go Ecommerce Backend

A scalable and secure backend service for an ecommerce platform, built with **Go**, **Gin**, and **MongoDB**. This project demonstrates real-world API design, secure authentication, modular architecture, and middleware usage for a production-ready ecommerce backend.

## 🚀 Features

### 🧑 User Authentication
- `POST /users/signup`: Register a new user with validation.
- `POST /users/login`: Authenticate and receive a JWT token.

### 🛍️ Product Management
- `GET /users/productview`: View all available products.
- `GET /users/search`: Search products by query.
- `POST /admin/addproduct`: Add a product (Admin-only).

### 🛒 Cart Management
- `GET /addtocart`: Add item to user's cart.
- `GET /removeitem`: Remove item from cart.
- `GET /listcart`: List items currently in the cart.
- `GET /cartcheckout`: Buy all items in the cart.
- `GET /instantbuy`: Instantly purchase a single item.

### 📦 Address Management
- `POST /addaddress`: Add a new address to user profile.
- `PUT /edithomeaddress`: Edit the saved home address.
- `PUT /editworkaddress`: Edit the saved work address.
- `DELETE /deleteaddresses`: Delete saved addresses.

> 🔒 Protected routes are secured using JWT authentication middleware.

---

## 🛠️ Tech Stack & Skills Demonstrated

| Technology | Purpose |
|-----------|---------|
| **Go (Golang)** | Core backend programming language |
| **Gin** | Fast HTTP router for routing and middleware |
| **MongoDB** | NoSQL database for storing users, products, carts |
| **JWT (github.com/dgrijalva/jwt-go)** | Secure token-based user authentication |
| **Validator (go-playground/validator)** | Struct validation during signup, login, and product creation |
| **dotenv (joho/godotenv)** | Environment-based configuration for secure key and port management |
| **Golang Crypto** | Password hashing and secure token signing |
| **Clean Architecture** | Layered and modular project structure |
| **Middleware** | Authentication and logging with Gin's middleware chain |

---

## 📁 Project Structure (Simplified)

```

go-ecommerce-backend/
├── controllers/        # Handler logic for each route
├── models/             # MongoDB schema models
├── middleware/         # JWT Authentication
├── routes/             # Public and protected routes setup
├── utils/              # Helpers like token creation and hashing
├── .env                # Environment variables (PORT, MONGODB\_URI, etc.)
├── main.go             # Application entrypoint
└── go.mod              # Dependency management

````

---

## 🧪 How to Run Locally

### Prerequisites:
- Go >= 1.20
- MongoDB running locally or on cloud
- `.env` file with:

```env
PORT=8000
MONGODB_URI=mongodb://localhost:27017
JWT_SECRET=your_jwt_secret
````

### Installation:

```bash
git clone https://github.com/sharukh010/go-ecommerce-backend.git
cd go-ecommerce-backend
go mod tidy
go run main.go
```

---

## 📬 API Testing

You can use tools like **Postman** or **cURL**. For authenticated routes, add:

```
Authorization: Bearer <JWT_TOKEN>
```

---

## 📌 Why This Project?

This backend was built to deepen my expertise in:

* Writing RESTful APIs in Go using Gin
* Using MongoDB with the Go driver
* Structuring backend applications with clean layering
* Implementing middleware for logging and authentication
* Building secure and modular ecommerce applications

---

## 📚 Future Improvements

* Role-based access control (RBAC) for Admin vs User
* Pagination and filters for product search
* Payment gateway integration
* Dockerfile for containerized deployment
* Unit and integration tests

---

## 📧 Contact

Made with ❤️ by [Sharukh010](https://github.com/sharukh010)


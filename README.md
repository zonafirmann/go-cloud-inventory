# ☁️ Go Cloud Inventory: Enterprise API Boilerplate

![Go Version](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-18-4169E1?style=flat&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)
![License](https://img.shields.io/badge/License-MIT-green.svg)

A high-performance, production-ready REST API built for scalable inventory management. Designed with **Clean Architecture** principles and containerized using **Docker Multi-Stage Builds** for seamless cloud deployment.

## 🚀 Key Features

* **Enterprise Architecture:** Strict separation of concerns (Models, Repository, Handlers, Config).
* **ACID Transactions:** Financial-grade checkout logic utilizing PostgreSQL row-level locking (`FOR UPDATE`) to prevent race conditions during high-concurrency sales.
* **Smart Analytics Engine:** Built-in rule-based algorithm generating automated stock alerts (Critical/Overstock) without relying on paid third-party APIs.
* **Cloud-Native Deployment:** Fully orchestrated with `docker-compose`, decoupling the Go API and PostgreSQL into isolated, secure containers.

## 🛠️ Tech Stack

* **Language:** Golang (Go 1.22)
* **Database:** PostgreSQL 18
* **Driver:** `pgx/v5` (High-performance native driver)
* **DevOps:** Docker & Docker Compose (Multi-stage lightweight image)

## 📂 Project Structure

```text
.
├── config/                  # Environment & Database connections
├── handlers/                # HTTP routing and API endpoints
├── models/                  # Domain entities and JSON definitions
├── repository/              # Direct database interactions and business rules
├── Dockerfile               # Multi-stage image build instructions
├── docker-compose.yml       # Container orchestration (API + DB)
└── main.go                  # Application entry point
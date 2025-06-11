# GASHA WAF - Built to Evolve, Designed to Protect

![Platform](https://img.shields.io/badge/platform-Docker%20%7C%20Go%20%7C%20Python-blue.svg)
![License](https://img.shields.io/badge/license-Apache%202.0-green.svg)

![LOGO](./image/40.png)

## Overview

**GASHA WAF** is an AI-powered Web Application Firewall engineered to protect modern web applications from evolving cyber threats, especially **injection-based attacks** like SQLi, XSS, and more. It fuses rule-based security with machine learning models to detect and mitigate threats in real time.

This project was developed as a final capstone submission for the Software Engineering And Computing Technology Program at AAiT.

## 🚀 Features

* **AI-Driven Threat Detection**: Inline detection using RF.
* **Online Learning**: Continual update and adjustment to new traffic patterns.
* **Custom Rule Engine**: Add custom WAF rules per application.
* **Rate Limiting**: Stop brute-force and DoS-style attacks.
* **Real-Time Logging & Dashboard**: Centralized logs and statistics via frontend.
* **Notification System**: Email and dashboard-based alerts for critical events.

## 🧠 Architecture

GASHA WAF consists of multiple microservices:

* **Interceptor**: Captures and forwards HTTP requests.
* **ML Server**: Classifies traffic and detects anomalies.
* **Backend API**: Core WAF engine and configuration service.
* **Frontend Dashboard**: Visualization, rule management, and alert configuration.

![Architecture Diagram](./image/Arc.jpg)

---

## 📦 Tech Stack

* **Backend**: Go (Gin Framework)
* **Machine Learning**: Python
* **Database**: PostgreSQL
* **Frontend**: React + Vite
* **Infrastructure**: Docker, GitHub Actions

---

## 🐳 Installation (Docker-Based)

### Prerequisites

* Docker
* linux OS

### Step-by-Step

#### 1. Pull the images:

```bash
docker pull natnaelcrypto/ml
docker pull natnaelcrypto/waf-db
docker pull natnaelcrypto/waf-backend
docker pull natnaelcrypto/interceptor
docker pull natnaelcrypto/waf-frontend
```

#### 2. Environment Variables

#### Backend

```
DB_HOST=pg-xxxxxxxxxx
DB_PORT=5432
DB_USER=username
DB_PASSWORD=********
DB_NAME=waf
DB_SSLMODE=require/disable
WSKEY=***
JWT_SECRET_KEY=***
```

#### Frontend

```
VITE_BACKEND_URL=https://ipaddress-of-backend:port
```

#### Interceptor

```
BACKENDURL=https://ipaddress-of-backend:8443
MLHOSTURL=https://ipaddress-of-ml:port:8090
WSKEY=***********
```

#### ML Service

```
DEBUG=False
BACKEND_API_URL =http://ipaddress-of-backend:8090/
```
#### Run the containers in order

```bash
docker network create waf_network
```

Default postgres database
```bash
docker run -d --name postgres --network waf_network --network-alias db -v waf_pgdata:/var/lib/postgresql/data -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=****** -e POSTGRES_DB=waf --restart unless-stopped postgres:16-alpine

docker restart postgres
```
To Run Backend
```bash
docker run -d --name backend --network waf_network --network-alias backend -e DB_USER=****** -e DB_PASSWORD=****** -e DB_HOST=****** -e DB_PORT=****** -e DB_NAME=waf -e DB_SSLMODE=require -e WSKEY=****** -e JWT_SECRET_KEY=****** --restart unless-stopped -p 8443:8443 natnaelcrypto/waf-backend:latest
```
To Run Frontend
```bash
docker run -d --name frontend --network host -e VITE_BACKEND_URL=http://ipaddress-of-backend:8443 -p 5173:5173 --restart unless-stopped natnaelcrypto/waf-frontend:latest
```
To Run Interceptor
```bash 
docker run -d --name interceptor --network host -e BACKENDURL=https://backend:8443 -e MLHOSTURL=https://ml:8090/ -e WSKEY=****** --restart unless-stopped -p 80:80 -p 443:443 natnaelcrypto/interceptor:latest
```
To Run Ml
```bash
docker run -d --name ml -p 8090:8090 -e DEBUG=true -e BACKEND_API_URL=https://backend:8443/ natnaelcrypto/ml:latest
```
---

## 🛡️ Detection Capabilities

* SQL Injection
* Cross-Site Scripting (XSS)
* Command Injection
* LDAP Injection
* XML Injection
* NoSQL Injection
* File Inclusion
* Path Traversal
* Server-Side Template Injection (SSTI)
* Rate-based attacks (via built-in rate limiter)

All attacks are **logged**, **classified**, and **blocked** in real-time.

---

## 📊 Dashboard & Analytics

* **Comprehensive stats dashboard** for:

  * Traffic volume
  * Attack types
  * Request origins
  * Model vs Rule-based detections
  * Notifications & alert settings

![Dashboard Screenshots](./image/Screenshot%202025-05-25%20211331.png)

---

## 🧪 CI/CD & Testing

* CI/CD managed through GitHub Actions.
* Automated testing pipeline ensures deployment safety.

---

## 📄 License

Licensed under the [Apache License 2.0](LICENSE).

---

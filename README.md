# AI-Powered Web Application Firewall (WAF)

## Overview
The **AI-Powered Web Application Firewall (WAF)** is a next-generation security solution designed to protect web applications from evolving cyber threats. This project combines traditional WAF capabilities with artificial intelligence (AI) to provide enhanced threat detection, traffic analysis, and real-time protection. 

This project is the capstone submission for [Your Degree/Program Name] and was developed over the course of four months.

## Features
- **Real-Time Traffic Monitoring**: Captures and logs HTTP requests for analysis.
- **AI-Driven Threat Detection**: Utilizes machine learning models (XGBoost, CNN, Random Forest) for identifying and mitigating attacks like SQL Injection, XSS, and DDoS.
- **Seamless Integration**: Requires minimal setup—only the target web application's IP address and port.
- **Custom WAF Engine**: Built from scratch using Golang, optimized for performance with multi-threading.
- **Embedded Database**: SQLite for efficient logging and analysis of traffic.
- **User-Friendly Deployment**: No modification to existing web servers required.

## Key Technologies
- **Backend**: Golang (for the WAF engine)
- **Machine Learning**: Python (for AI model training and integration)
- **Database**: Badgerdb (for real-time logging)
- **Deployment**: Virtual Machine for portability and scalability.

## System Architecture
1. **Traffic Interception**: Golang-based engine intercepts HTTP/S traffic.
2. **Request Analysis**: Incoming traffic is analyzed using both rule-based logic and AI models.
3. **Decision Making**: The AI model determines whether a request is benign or malicious.
4. **Response Handling**: Malicious requests are blocked, logged, and reported; benign requests are forwarded to the web server.


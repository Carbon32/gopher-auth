<p align="center">
    <img
    src="https://readme-typing-svg.demolab.com/?font=Fira+Code&weight=700&size=48&duration=3000&pause=1000&color=58A6FF&center=true&vCenter=true&width=500&lines=Gopher%20Auth"
    alt="Go Authentication System"
    />
</p>

<p align="center">
    A Go-based authentication system built with PostgreSQL, designed to handle user accounts, secure authentication, session management, and account verification.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Language-Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/Database-PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL" />
  <img src="https://img.shields.io/badge/Container-Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker" />
</p>

---

## Introduction

Originally, this project wasn't supposed to be an authentication system. The initial goal was to build a more complete backend containing various systems for managing products, orders, and other related functionality.

The project was being developed for a local application, but due to some questionable decisions along the way, the Go component has evolved into a dedicated authentication system. And here we are. Despite the change in direction, the project includes several features designed to provide a solid foundation for user authentication and account management.

## Features

### User Account Management

- **Account Registration:** Create new user accounts.
- **User Login:** Authenticate using either a username or an e-mail address.
- **Profile Modification:** Update personal account information.

### JWT-Based Authentication

- **Access Tokens:** Authenticate requests using JSON Web Tokens (JWT).
- **Protected Routes:** Restrict access to endpoints that require authentication.
- **Token-Based Authorization:** Establish an authentication mechanism for protected resources.

### Refresh Token Management

- **Session Tracking:** Track user sessions through refresh tokens.
- **Session Control:** Provide better control over active authentication sessions.

### E-mail Verification

- **Verification Links:** Send e-mail verification links to users.
- **Token Expiration:** Ensure verification links expire after a configured period.
- **Sensitive Information Verification:** Apply verification steps when modifying sensitive account information, such as e-mail addresses and passwords.

### Rate Limiting

- **Abuse Prevention:** Limit requests to reduce abusive activity.
- **Account-Based Tracking:** Use account information to help identify excessive requests.
- **IP-Based Tracking:** Track request activity by IP address to help prevent abuse.

Alongside these core features, the project includes various smaller features and implementation details that help tie everything together.

---

## Requirements

Before running the project, make sure you have the following installed:

- [Go](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/download/)
- Bash, to execute the provided shell scripts

If you plan to run the project using Docker, make sure [Docker](https://www.docker.com/) is installed and running.

---

## Getting Started

Follow these steps to configure and run the project locally.

### 1. Configure Environment Variables

Create a `.env` file in the project root to configure the application's environment.

The available configuration options can be found in `docker.env.example`. Use it as a reference when setting up your environment.

Make sure all required variables are properly configured before starting the application, including the settings needed to connect to your PostgreSQL database.

> **Note:** Review the example configuration carefully and adjust the values to match your environment.

### 2. Set Up Database Migrations

If this is your first time running the project against your database, execute the migration script:

```bash
bash migrate.sh
```

This initializes the database schema by applying the project's migrations.

### 3. Run the Application

Once your environment is configured and the database is ready, start the application using:

```bash
bash run.sh
```

The application should now start using the configuration provided in your `.env` file.

---

## Running with Docker

If you prefer Docker, the project includes a dedicated script for running the application in a containerized environment.

Execute:

```bash
bash docker.sh
```

Before doing so, make sure Docker is installed and running and that the required environment configuration is in place.

For Docker-specific configuration options, refer to `docker.env.example`.

---

## Frontend Integration

Connecting this authentication system to a frontend is a separate topic.

Frontend integration, including handling authentication requests, managing tokens, and interacting with protected endpoints, will likely be covered in another repository.

For now, this project focuses on the backend authentication layer and the functionality required to manage users and their sessions.

---

*I believe the project is suitable for real-world use, provided that it is properly integrated with the necessary supporting services and configured for your specific environment.*
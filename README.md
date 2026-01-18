# ClipLink

## Overview

ClipLink is a full-featured URL shortener that allows users to:
- Register and log in.
- Shorten long URLs to compact, shareable links.
- Delete shortened URLs.
- Redirect users from shortened links to their original destinations.
- Track usage within user limits.

Built using Go (Golang) for the backend and React + Vite for the frontend,
it leverages MongoDB for storage and uses JWT (JSON Web Tokens) for authentication.
ClipLink can serve the frontend directly from the Go backend for easy deployment.

---

## Table of Contents

- [Features](#features)
- [Technologies Used](#technologies-used)
- [Installation](#installation)
- [Configuration](#configuration)
- [Frontend Build & Deployment](#frontend-build--deployment)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
  - [User Registration](#user-registration)
  - [User Login](#user-login)
  - [Shortening a URL](#shortening-a-url)
  - [Deleting a URL](#deleting-a-url)
  - [Redirecting](#redirecting)
- [Frontend Fetch Examples](#frontend-fetch-examples)
- [Production Deployment Notes](#production-deployment-notes)
- [End Note](#end-note)

---

## Features

- User registration and login with hashed passwords (SHA-256).
- URL shortening with a 48-hour TTL (configurable).
- Maximum of 5 active URLs per user (configurable).
- Ability to delete shortened URLs.
- Redirection from shortened URLs to original URLs.
- JWT-based authentication for secure API access.
- Frontend served via Go backend for production convenience.
- CORS configuration for frontend-backend integration.
- Safe URL normalization and validation.

---

## Technologies Used

- **Backend:** Go (Golang)
- **Frontend:** React + Vite
- **Database:** MongoDB
- **Authentication:** JSON Web Tokens (JWT)
- **Other Libraries:** crypto/sha256, encoding/json, net/http, go.mongodb.org/mongo-driver

---

## Installation

### Prerequisites

1. Install Go (1.22 or later).
2. Install Node.js (16+ recommended) for building the frontend.
3. Install MongoDB and run an instance locally or remotely.

### Steps

1. Clone the repository:

```bash
git clone https://github.com/DemonSlayer256/ClipLink.git
cd ClipLink
```

2. Install Go dependencies:

```bash
go mod tidy
```

3. Install frontend dependencies:

```bash
cd frontend
npm install
```

4. Configure environment variables in .env:

```
- MONGO_URI="mongodb://localhost:27017"
- SECURE_KEY="your_jwt_secret_key"
- CORS_URL="http://localhost:5173"
```

---

## Frontend Build & Deployment

1. Build the React frontend:

```bash
cd frontend
npm run build
```

This generates a dist/ folder with static files.

2. Serve the frontend via Go backend:

- The Go backend is set to serve static files and API endpoints from `dist/` folder.
- Start the server:

```bash
go run main.go
```

3. Access the app:

- Go to http://localhost:8080 to see your frontend served by Go.

---

## Configuration

Adjust these in your .env or Go code:

- **MongoDB URI:** `MONGO_URI`
- **JWT Secret Key:** `SECURE_KEY`
- **CORS URL:** `CORS_URL` (React dev or deployed domain)
- **Max URLs per user:** `max_url_limit` in Go
- **URL TTL:** `max_TTL` in Go (in hours)

---

## Usage

### Register a User

```bash
POST /register
Content-Type: application/json

{
  "user": "username",
  "pass": "yourpassword"
}
```

### Login

```bash
POST /login
Content-Type: application/json

{
  "user": "username",
  "pass": "yourpassword"
}
```

- Response includes a JWT token: `{"token": "your_jwt_token"}`

### Shorten a URL

```bash
POST /shorten
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "url": "https://example.com"
}
```

### Delete a URL

```bash
POST /delete
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "code": "shortenedCode"
}
```

### Redirect

- Access shortened URL: http://localhost:8080/<shortenedCode>
- Automatically redirects to original URL.

---

## API Endpoints

### User Registration

- `POST /register`
- Body: `{"user": "username", "pass": "password"}`
- Responses:
  - 201 Created: User registered
  - 409 Conflict: User already exists
  - 400 Bad Request: Invalid JSON

### User Login

- `POST /login`
- Body: `{"user": "username", "pass": "password"}`
- Responses:
  - 200 OK: Returns JWT token
  - 401 Unauthorized: Invalid credentials

### Shorten URL

- `POST /shorten`
- Headers: `Authorization: Bearer <token>`
- Body: `{"url": "https://example.com"}`

### Delete URL

- `POST /delete`
- Headers: `Authorization: Bearer <token>`
- Body: `{"code": "shortenedCode"}`

### Redirect

- `GET /<shortenedCode>`
- Responses:
  - 302 Found: Redirect to original URL
  - 404 Not Found: URL expired or invalid

---

## Frontend Fetch Examples

Here are examples using `import.meta.env.VITE_API_URL` in React:

```javascript
const API_URL = import.meta.env.VITE_API_URL;

// Login
const loginUser = async (username, password) => {
  const response = await fetch(`${API_URL}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ user: username, pass: password }),
  });
  return await response.json();
};

// Shorten URL
const shortenUrl = async (url, token) => {
  const response = await fetch(`${API_URL}/shorten`, {
    method: "POST",
    headers: { 
      "Content-Type": "application/json",
      "Authorization": `Bearer ${token}`
    },
    body: JSON.stringify({ url }),
  });
  return await response.json();
};
```

---

## Production Deployment Notes

1. **Single Backend Deployment:**
   - Go backend serves both API and static frontend.
   - Ensures no CORS issues and easier hosting.

2. **Temporary Public Hosting:**
   - Use services like `ngrok` for 48-hour public access without a domain.

3. **Environment Variables:**
   - Keep `SECURE_KEY` secret.
   - MongoDB URI should be secured if using remote database.

4. **Frontend Build:**
   - Always run `npm run build` before deploying to production.

5. **Recommended Hosting:**
   - VPS, cloud instance (AWS EC2, DigitalOcean), or containerized Docker deployment.

6. **Logging & Debugging:**
   - Go logs errors in the terminal.
   - Monitor MongoDB TTL index for automatic URL expiration.

---

## End Note

ClipLink provides a **secure, scalable, and production-ready URL shortening service**.  
It demonstrates:

- Serving React frontend via Go.
- JWT authentication for secure APIs.
- MongoDB TTL indexes for automatic link expiration.
- Easy deployment and temporary public sharing.

> Extend ClipLink with analytics, custom domains, or unlimited links per user.

# Thank You for Using ClipLink!

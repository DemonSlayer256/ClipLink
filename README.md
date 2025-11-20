# ClipLink

## Overview

ClipLink allows users to register, log in, shorten URLs, 
delete shortened URLs, and redirect users from shortened links to their 
original destinations. Built using Go, it utilizes MongoDB for data storage 
and supports user authentication through JSON Web Tokens (JWT).

## Table of Contents

- [Features](#features)
- [Technologies Used](#technologies-used)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
  - [User Registration](#user-registration)
  - [User Login](#user-login)
  - [Shortening a URL](#shortening-a-url)
  - [Deleting a URL](#deleting-a-url)
  - [Redirecting](#redirecting)
  - [End Note](#endnote)

## Features

- User registration and login with hashed passwords.
- URL shortening capabilities with expiration times.
- Ability to delete shortened URLs.
- Redirection from shortened URLs to original destinations.
- JWT for user authentication and session management.

## Technologies Used

- Go (Golang) for the backend.
- MongoDB for data storage.
- JWT (JSON Web Tokens) for user authentication.
- bufio, encoding/json, net/http, and other standard libraries for functionality.

## Installation

### Prerequisites

1. Install Go (1.22 or later).
2. Install MongoDB.
3. Set up a MongoDB instance (either locally or remotely).

### Steps

1. Clone the repository:
  ```bash
  git clone https://github.com/Dome
  cd ClipLink
  ```

2. Install the necessary Go modules:
  ```bash
  go mod tidy
  ```

3. Create a .env file in the root of the project with the MongoDB URI:
  ```bash
  MONGODB_URI="mongodb://localhost:27017"
  ```

4. Run the service:
  ```bash
   go run main.go
   ```

The server will start by default on port **8080**.

## Configuration

You can configure various parameters in the code:

- **MongoDB URI**: Adjust this in the .env file.
- **TTL for URLs**: Modify the max_TTL variable (in hours) for how long a shortened URL is valid.

## Usage

### Registering a User

To register a new user, send a POST request to /register with the following JSON body:
```bash
{
  "user": "username",
  "password": "yourpassword"
}
```

### Logging In

To log in, send a POST request to /login with the user's credentials:
```bash
{
  "username": "username",
  "password": "yourpassword"
}
```

You will receive a token upon successful authentication.

### Shortening a URL

Send a POST request to /shorten with the following JSON body and include the JWT token in the Authorization header:
```bash
{
  "url": "https://example.com"
}
```

### Deleting a URL

To delete a shortened URL, send a DELETE request to /shorten with the JSON body containing the shortened code, and include the JWT token in the Authorization header:
```bash
{
  "code": "shortenedCode"
}
```

### Redirecting

Request the shortened URL (e.g., http://localhost:8080/yourShortCode) to be redirected to the original URL.

## API Endpoints

### User Registration

- **Endpoint**: /register
- **Method**: POST
- **Body**:
```bash
{
  "user": "username",
  "password": "yourpassword"
}
```
- **Responses**:
    - 201 Created: Successful registration.
    - 409 Conflict: User already exists.
    - 400 Bad Request: Invalid JSON data.

### User Login

- **Endpoint**: /login
- **Method**: POST
- **Body**:
```bash
{
  "username": "username",
  "password": "yourpassword"
}
```
- **Responses**:
    - 200 OK: Successful login with a token.
    - 401 Unauthorized: Invalid credentials.

### Shortening a URL

- **Endpoint**: /shorten
- **Method**: POST
- **Headers**: Authorization: Bearer <token>
- **Body**:
```bash
{
  "url": "https://example.com"
}
```

### Deleting a URL

- **Endpoint**: /shorten
- **Method**: DELETE
- **Headers**: Authorization: Bearer <token>
- **Body**:
```bash
{
  "code": "shortenedCode"
}
```

### Redirecting

- **Endpoint**: /{shortened}
- **Method**: GET
- **Responses**:
    - 302 Found: Redirects to the link associated with the shortened code
    - 404 NotFound: The code has expired or the page is not found

## End Note

 ClipLink provides a simple, secure, and extensible URL‑shortening service. By leveraging Go's performance, MongoDB's flexible storage, and JWT‑based authentication, the project demonstrates a production‑ready backend that can be easily deployed, customized, and integrated into larger systems.
Feel free to fork, extend the feature set (e.g., analytics, custom domains) or contribute improvements via pull requests.

#!/bin/bash

# Configuration
API_URL="http://localhost:8080"  # Update with your API URL
USERNAME="newuser"                 # Username for registration
PASSWORD="securepassword"          # Password for registration
SHORT_URL="http://example.com"     # URL to shorten

# Step 1: Register a new user
echo "Registering a new user..."
register_response=$(curl -s -X POST "$API_URL/register" \
  -H "Content-Type: application/json" \
  -d "{\"user\":\"$USERNAME\", \"pass\":\"$PASSWORD\"}")

echo "Response from registration: $register_response"

# Step 2: Log in to receive a JWT token
echo "Logging in..."
login_response=$(curl -s -X POST "$API_URL/login" \
  -H "Content-Type: application/json" \
  -d "{\"user\":\"$USERNAME\", \"pass\":\"$PASSWORD\"}")
  
# Extract JWT token from login response (assume it's in the form `{"token": "your_token_here"}`)
TOKEN=$(echo $login_response | jq -r .token)

if [ "$TOKEN" == "null" ]; then
  echo "Failed to log in: $login_response"
  exit 1
fi

echo "Successfully logged in. Token: $TOKEN"

# Step 3: Shorten a URL using the JWT token
echo "Shortening URL..."
shorten_response=$(curl -s -X POST "$API_URL/shorten" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"url\":\"$SHORT_URL\"}")

echo "Response from shortening: $shorten_response"

# Optional: Extract and print the shortened URL if your response contains it
# Here we assume the response contains a field 'shortened_url'
SHORTENED_URL=$(echo $shorten_response | jq -r .short_url)

if [ "$SHORTENED_URL" != "null" ]; then
  echo "Shortened URL: $SHORTENED_URL"
else
  echo "Failed to shorten URL: $shorten_response"
fi

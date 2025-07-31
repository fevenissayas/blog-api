# Blog API Documentation

## Base URL

```
http://localhost:<port>
```
Replace `<port>` with the port your server is running on (usually 8080).

---

## Endpoints

### Register User

- **URL:** `/register`
- **Method:** `POST`
- **Description:** Register a new user account.

#### Request Body

Send as JSON:

```json
{
  "username": "string, min 3, max 20, required",
  "email": "string, valid email, required",
  "password": "string, min 6, required",
  "role": "user | admin (optional, default: user)",
  "bio": "string (optional)",
  "profile_picture": "string (optional, URL or path)",
  "contact_info": "string (optional)"
}
```

#### Example Request

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "johndoe@example.com",
    "password": "12345678
  }'
```

#### Responses

- **201 Created**
    ```json
    {
      "message": "Succesfully Registered User"
    }
    ```

- **400 Bad Request**
    ```json
    {
      "message": "invalid input"
    }
    ```

- **422 Unprocessable Entity**
    ```json
    {
      "message": "<validation error details>"
    }
    ```

- **500 Internal Server Error**
    ```json
    {
      "message": "<error details>"
    }
    ```

#### Validation Rules

- `username`: Required, 3-20 characters.
- `email`: Required, must be valid email format.
- `password`: Required, minimum 6 characters.
- `role`: Optional, `"user"` or `"admin"`.
- Optional fields: `bio`, `profile_picture`, `contact_info`.

#### Notes

- Passwords are never returned in responses.
- The user is stored in a MongoDB collection (implied by use of `primitive.ObjectID`).

---

## Data Model

### User

| Field           | Type                | Description                        |
|-----------------|---------------------|------------------------------------|
| id              | string (ObjectID)   | Unique identifier (generated)      |
| username        | string              | Username (3-20 chars, required)    |
| email           | string              | Email address (required)           |
| password        | string              | Password (6+ chars, required)      |
| role            | string (`user`/`admin`) | User role (default: user)      |
| created_at      | string (ISO 8601)   | When account was created           |
| updated_at      | string (ISO 8601)   | Last updated time                  |
| bio             | string              | User bio (optional)                |
| profile_picture | string              | URL/path to profile image (optional)|
| contact_info    | string              | Additional contact info (optional) |

---

## Error Codes

- **400**: Malformed JSON or missing body fields.
- **422**: Fails validation (e.g., invalid email, too short username/password).
- **500**: Internal server/database error.

---

## Example Success Response

```json
{
  "message": "Succesfully Registered User"
}
```

---

## Example Validation Error

```json
{
  "message": "Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag"
}
```

---

## Getting Started

1. Run the server with `go run main.go`.
2. Use the `/register` endpoint to create users.

---

## API Testing with Postman

You can test this API using the following [Postman collection](https://lively-crescent-132029.postman.co/workspace/73182311-5a50-4ab2-a96a-9bade0810ab5/collection/44609334-c9dcbde7-f805-4589-adaa-29ac37fcbee2?action=share&source=copy-link&creator=44609334).

1. Click the link above to open the collection in Postman.
2. Import it to your workspace and start testing the available endpoints!
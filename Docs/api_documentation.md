# Blog API Documentation

## Base URL

```
http://localhost:<port>
```
Replace `<port>` with your server’s port (e.g. `8080`).

---

## Endpoints

### 1. Register User

- **URL:** `/auth/register`
- **Method:** `POST`
- **Description:** Register a new user.

#### Request Body

```json
{
  "username": "string, min 3, max 20, required",
  "email": "string, valid email, required",
  "password": "string, min 6, required",
  "bio": "string (optional)",
  "profile_picture": "string (optional, URL or path)",
  "contact_info": "string (optional)"
}
```

#### Example Success Response (201)

```json
{
  "message": "User created successfully"
}
```

#### Error Responses

- **400 Bad Request**
```json
{
  "error": "Email is already registered"
}
```

- **422 Unprocessable Entity**
```json
{
  "error": "Invalid email format"
}
```

- **500 Internal Server Error**
```json
{
  "error": "failed to create user: <details>"
}
```

---

### 2. Login User

- **URL:** `/auth/login`
- **Method:** `POST`
- **Description:** Log in a user and receive JWT tokens.

#### Request Body

```json
{
  "email": "string, valid email, required",
  "password": "string, required"
}
```

#### Success Response (200 OK)

```json
{
  "access_token": "<JWT access token>",
  "refresh_token": "<JWT refresh token>"
}
```

#### Error Responses

- **401 Unauthorized**
```json
{
  "error": "incorrect email or password"
}
```

- **422 Unprocessable Entity**
```json
{
  "error": "email and password are required"
}
```

---

### 3. Refresh Token

- **URL:** `/auth/refresh`
- **Method:** `POST`
- **Description:** Issues a new access and refresh token using a valid refresh token.

#### Request Headers

```http
Authorization: Bearer <refresh_token>
```

#### Example Request

```bash
curl -X POST http://localhost:8080/auth/refresh   -H "Authorization: Bearer <your_refresh_token>"
```
#### Success Response (200 OK)

```json
{
  "access_token": "<new_access_token>",
  "refresh_token": "<new_refresh_token>"
}
```

#### Error Responses

- **401 Unauthorized**
    - Missing or malformed Authorization header
    - Expired, revoked, or invalid refresh token

```json
{
  "error": "invalid refresh token: token is expired"
}
```

```json
{
  "error": "refresh token has been revoked"
}
```

```json
{
  "error": "user not found"
}
```

---

#### Security Notes

- Old refresh tokens are revoked after issuing new ones.
- Token validation checks:
  - Token exists in DB
  - Token matches the one stored
  - Token is not expired or revoked
- User info is fetched by ID in the token payload.


---

## User Model

| Field           | Type                | Description                          |
|----------------|---------------------|--------------------------------------|
| id             | string (ObjectID)   | MongoDB-generated ID                 |
| username       | string              | Required, 3-20 characters            |
| email          | string              | Required, valid email                |
| password       | string              | Required, stored as bcrypt hash      |
| role           | string              | Default: `user`                      |
| is_verified    | bool                | Defaults to false                    |
| bio            | string              | Optional                             |
| profile_picture| string              | Optional URL/path                    |
| contact_info   | string              | Optional                             |
| created_at     | string (ISO 8601)   | Timestamp of account creation        |
| updated_at     | string (ISO 8601)   | Timestamp of last update             |

---

## Common Status Codes

- `200 OK`: Success
- `201 Created`: User registered
- `400 Bad Request`: Invalid request or duplicate
- `401 Unauthorized`: Invalid credentials/token
- `422 Unprocessable Entity`: Invalid/missing input
- `500 Internal Server Error`: Server/database error

---








## Getting Started

1. Run the server with:
    ```bash
    go run main.go
    ```
2. Use `/auth/register` to register a user.
3. Use `/auth/login` to get tokens.
4. Use `/auth/refresh` to refresh your access token.

## API Testing with Postman

You can test this API using the following [Postman collection](https://lively-crescent-132029.postman.co/workspace/73182311-5a50-4ab2-a96a-9bade0810ab5/collection/44609334-c9dcbde7-f805-4589-adaa-29ac37fcbee2?action=share&source=copy-link&creator=44609334).

1. Click the link above to open the collection in Postman.
2. Import it to your workspace and start testing the available endpoints!
# Blog API Documentation

The **Blogs API** provides a comprehensive set of endpoints for managing blog posts and related features. It supports creating, retrieving, updating, and deleting blogs, as well as advanced functionality such as filtering, searching, and generating AI-powered content suggestions. Additionally, it handles blog likes to facilitate user interaction and engagement.

## Main Features

- **Blog Management:** Create, read, update, and delete blog posts.
    
- **Filtering & Searching:** Search blogs by tags, dates, titles, user IDs, and sort them by popularity or creation date.
    
- **AI Suggestions:** Generate AI-powered suggestions to improve blog content, flow, grammar, and structure.
    
## API Testing with Postman

You can test this API using the following [Postman collection](https://lively-crescent-132029.postman.co/workspace/73182311-5a50-4ab2-a96a-9bade0810ab5/collection/44609334-c9dcbde7-f805-4589-adaa-29ac37fcbee2?action=share&source=copy-link&creator=44609334).

## Route Structure

- `GET /blogs/` — Retrieve paginated list of blogs.
    
- `GET /blogs/:id` — Retrieve a single blog by ID.
    
- `POST /blogs/` — Create a new blog (requires authentication).
    
- `PUT /blogs/:id` — Update an existing blog by ID (requires authentication).
    
- `DELETE /blogs/:id` — Delete a blog by ID (requires authentication).
    
- `GET /blogs/filter` — Filter blogs by criteria (requires authentication).
    
- `POST /blogs/aisuggestion` — Get AI-generated suggestions for blog improvement (requires authentication).
    
- `GET /blogs/search` — Search blogs with query parameters (requires authentication).

## Blog Endpoints

### 1. GET `/blogs/`

Retrieve a paginated list of blogs.

**Method:** GET  
**URL:** `/blogs/`

**Query Parameters:**
| Parameter | Type   | Required | Description                  |
|-----------|--------|----------|------------------------------|
| page      | int    | No       | Page number (default: 1)     |
| limit     | int    | No       | Number of blogs per page     |
| sort      | string | No       | "recent" (default) or "popular" |
| authorID  | string | No       | Filter by author ID          |

**Response:**
Returns a JSON object with paginated blog data and metadata.

**Example:**
```json
{
  "data": [
    {
      "id": "12345",
      "title": "Latest Tech Trends",
      "content": "...",
      "tags": ["technology"],
      "createdAt": "2025-08-08T12:00:00Z",
      "view_count": 150
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100
  }
}
```

---

### 2. GET `/blogs/:id`

Retrieve a single blog by its ID.

**Method:** GET  
**URL:** `/blogs/:id`

**Path Parameter:**
| Parameter | Type   | Required | Description         |
|-----------|--------|----------|---------------------|
| id        | string | Yes      | Blog unique ID      |

**Response:**
Returns a JSON object with the blog details.

**Example:**
```json
{
  "id": "12345",
  "title": "Latest Tech Trends",
  "content": "...",
  "tags": ["technology"],
  "createdAt": "2025-08-08T12:00:00Z",
  "view_count": 150
}
```

---

### 3. POST `/blogs/`

Create a new blog post. **Requires authentication.**

**Method:** POST  
**URL:** `/blogs/`

**Headers:**
| Header         | Value                | Description                  |
|----------------|----------------------|------------------------------|
| Authorization  | Bearer {accessToken} | JWT access token (required)  |

**Request Body:**
| Field     | Type     | Required | Description         |
|-----------|----------|----------|---------------------|
| title     | string   | Yes      | Blog title          |
| content   | string   | Yes      | Blog content        |
| tags      | [string] | No       | Array of tags       |

**Example:**
```json
{
  "title": "Latest Tech Trends",
  "content": "...",
  "tags": ["technology", "innovation"]
}
```

**Response:**
Returns the created blog object.

**Example:**
```json
{
  "id": "12345",
  "title": "Latest Tech Trends",
  "content": "...",
  "tags": ["technology", "innovation"],
  "createdAt": "2025-08-08T12:00:00Z",
  "view_count": 0
}
```

---

### 4. PUT `/blogs/:id`

Update an existing blog post by its ID. **Requires authentication.**

**Method:** PUT  
**URL:** `/blogs/:id`

**Headers:**
| Header         | Value                | Description                  |
|----------------|----------------------|------------------------------|
| Authorization  | Bearer {accessToken} | JWT access token (required)  |

**Path Parameter:**
| Parameter | Type   | Required | Description         |
|-----------|--------|----------|---------------------|
| id        | string | Yes      | Blog unique ID      |

**Request Body:**
| Field     | Type     | Required | Description         |
|-----------|----------|----------|---------------------|
| title     | string   | No       | Blog title          |
| content   | string   | No       | Blog content        |
| tags      | [string] | No       | Array of tags       |

**Example:**
```json
{
  "title": "Updated Tech Trends",
  "content": "Updated content...",
  "tags": ["technology", "update"]
}
```

**Response:**
Returns the updated blog object.

**Example:**
```json
{
  "id": "12345",
  "title": "Updated Tech Trends",
  "content": "Updated content...",
  "tags": ["technology", "update"],
  "createdAt": "2025-08-08T12:00:00Z",
  "view_count": 151
}
```

---

### 5. DELETE `/blogs/:id`

Delete a blog post by its ID. **Requires authentication.**

**Method:** DELETE  
**URL:** `/blogs/:id`

**Headers:**
| Header         | Value                | Description                  |
|----------------|----------------------|------------------------------|
| Authorization  | Bearer {accessToken} | JWT access token (required)  |

**Path Parameter:**
| Parameter | Type   | Required | Description         |
|-----------|--------|----------|---------------------|
| id        | string | Yes      | Blog unique ID      |

**Response:**
Returns a success message.

**Example:**
```json
{
  "message": "Blog deleted successfully"
}
```

### Likes Sub-Routes (all require authentication)

- `POST /blogs/:id/likes/` — Like a blog post.
    
- `DELETE /blogs/:id/likes/` — Remove a like from a blog post.
    
- `GET /blogs/:id/likes/` — Get the total like count for a blog post.
    
- `GET /blogs/:id/likes/is-liked` — Check if the authenticated user has liked the blog post.


### **GET** `/blogs/filter`

Retrieve blogs filtered by **tag**, **date**, and **sort** parameters.

### Authorization

Requires Bearer Token (JWT) in the header.

Authorization: Bearer {{accessToken}}.

### Request

**Method:** GET  
**URL:** `/blogs/filter`

### Request Parameters

The request accepts the following **query parameters**:

### Query Parameters

| Parameter | Type | Required | Description | Example |
| --- | --- | --- | --- | --- |
| `tag` | string | No | Filter blogs by a single tag | `technology` |
| `date` | string | No | Filter blogs by creation date (YYYY-MM-DD) | `2025-08-08` |
| `sort` | string | No | Sort by "recent" (default) or "popular" | `popular` |

Request Example

``` json
GET /blogs/filter?tag=technology&date=2025-08-08&sort=popular 
Authorization: Bearer {{accessToken}}

 ```

## Response

On successful request, the server responds with status `200 OK` and a JSON object with the following structure:

| Field | Type | Description |
| --- | --- | --- |
| status | string | `"success"` indicating the request was successful. |
| data | array | An array of blog objects matching the filter criteria. |
| meta | object | Metadata about the response: |
| filter | object | The filter parameters used (`tag`, `date`, `sort`). |
| count | integer | Number of blogs returned. |

---

## Example Response

``` json
{
  "status": "success",
  "data": [
    {
      "id": "12345",
      "title": "Latest Tech Trends",
      "content": "Content of the blog...",
      "tags": ["technology", "innovation"],
      "createdAt": "2025-08-08T12:00:00Z",
      "view_count": 150
    }
  ],
  "meta": {
    "filter": {
      "tag": "technology",
      "date": "2025-08-08",
      "sort": "popular"
    },
    "count": 1
  }
}

 ```

If no blogs match the filter criteria, the `data` array will be empty and `count` will be :

``` json
{
  "status": "success",
  "data": [],
  "meta": {
    "filter": {
      "tag": "nonexistenttag",
      "date": "",
      "sort": "recent"
    },
    "count": 0
  }
}

 ```

- **401 Bad Request**
    

If the `sort` parameter is invalid, the server responds with a `400 Bad Request`:

``` json
{
  "error": "Invalid sort parameter. Allowed values: recent, popular"
}

 ```

- **401 Unauthorized**
    
    Returned if the request does not include a valid authorization token.

# AI Blog Suggestion

## Request

**Method:** POST  
**URL:** `/blogs/aisuggestion`  
**Authorization:** Bearer Token required

---

## Request Body

At least one of the following fields **must** be provided: `title`, `content`, or `tags`.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| title | string | Optional | The title of the blog draft. |
| content | string | Optional | The main content/body of the blog draft. |
| tags | \[\]string | Optional | An array of tags related to the blog content. |

---

## Example Request Body

``` json
{
  "title": "How to Write Better Go Code",
  "content": "Writing clean and efficient Go code is essential for building scalable applications...",
  "tags": ["go", "programming", "clean code"]
}

 ```

## Response

Success (200 OK)

Returns a JSON object containing the AI-generated suggestion text.

``` json
{
  "Suggestion": "To improve the blog, consider breaking down long paragraphs for better readability and use more examples..."
}

 ```

#### Error Responses

- **400 Bad Request**
    
    Returned if the request body is invalid JSON or if none of the required fields are provided.
    

``` json
{
  "error": "at least one of title, content, or tags must be provided"
}

 ```

or

``` json
{
  "error": "invalid request"
}

 ```

- **401 Unauthorized**
    
    Returned if the request does not include a valid authorization token.

## GET /blogs/search

Search blogs based on various optional filters such as tag, date, sort order, title keywords, and user ID.

---

### Request

- **Method:** GET
    
- **URL:** `/blogs/search`
    
- **Headers:**
    
    - `Authorization: Bearer`
        

---

### Query Parameters

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| tag | string | Optional | Filter blogs by a single tag. |
| date | string | Optional | Filter blogs created on a specific date (format: YYYY-MM-DD). |
| sort | string | Optional | Sort order: `popular` (by view count descending) or `recent` (by creation date descending). Defaults to `recent` if omitted or invalid. |
| title | string | Optional | Search for blogs containing this keyword in their title (case-insensitive). |
| userID | string | Optional | Filter blogs created by a specific user ID. |

---

### Response

#### Success (200 OK)

Returns an array of blog objects matching the search criteria.

``` json
[
  {
    "id": "abc123",
    "title": "Concurrency Patterns in Go",
    "content": "Go provides rich support for concurrency...",
    "tags": ["golang", "concurrency"],
    "user_id": "user123",
    "view_count": 250,
    "createdAt": "2025-08-07T14:22:00Z"
  },
  {
    "id": "def456",
    "title": "Understanding Channels in Go",
    "content": "Channels are powerful tools for goroutine communication...",
    "tags": ["golang"],
    "user_id": "user123",
    "view_count": 200,
    "createdAt": "2025-08-06T09:15:30Z"
  }
]

 ```

#### Error Responses

- **401 Unauthorized**
    
    Returned if the request lacks valid authentication.
    
- **500 Internal Server Error**
    
    Returns if there is an error fetching or processing blogs.
    

``` json
{
  "error": "failed to fetch blogs: <error details>"
}

 ```

### Notes

- If no query parameters are provided, all blogs are returned sorted by creation date descending (most recent first).
    
- Date filtering expects the date to be in the `YYYY-MM-DD` format.
    
- Sorting defaults to "recent" if an invalid value is provided.

# Like a Blog

## POST /blog/{id}/likes/

Like a specific blog post by its ID. Requires authentication.

---

### Request

- **Method:** POST
    
- **URL:** `/blog/{blog_id}/likes/`  
    Replace `{blog_id}` with the blog's unique identifier.
    
- **Headers:**
    

| Header | Value | Description |
| --- | --- | --- |
| Authorization | Bearer {{accessToken}} | Required. Bearer token for authentication. |

---

### Request Parameters

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| id | string | Yes | The unique ID of the blog to like (path parameter). |

---

### Response

#### Success (200 OK)

``` json
{
  "message": "blog liked"
}

 ```

- Indicates that the blog was successfully liked by the authenticated user.
    
- If the user has already liked the blog, the operation is idempotent and returns success without duplicating the like.
    

#### Unauthorized (401 Unauthorized)

``` json
{
  "error": "unauthorized"
}

 ```

- Returned when no valid authentication token is provided or the user is not authenticated.
    

### Notes

- Ensure the `id` path parameter is a valid blog identifier in hexadecimal format.

# Remove Like from Blog

## DELETE /blog/{blog_id}/likes/

Remove a like from a specific blog post by its ID. Requires authentication.

---

### Request

- **Method:** DELETE
    
- **URL:** `/blog/{blog_id}/likes/`  
    Replace `{blog_id}` with the blog's unique identifier.
    
- **Headers:**
    

| Header | Value | Description |
| --- | --- | --- |
| Authorization | Bearer {{accessToken}} | Required. Bearer token for authentication. |

---

### Request Parameters

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| id | string | Yes | The unique ID of the blog to remove the like from (path parameter). |

---

### Response

#### Success (200 OK)

``` json
{
  "message": "blog unliked"
}

 ```

- Indicates that the like was successfully removed for the authenticated user.
    

#### Unauthorized (401 Unauthorized)

``` json
{
  "error": "unauthorized"
}

 ```

- Returned when no valid authentication token is provided or the user is not authenticated.
    

#### Bad Request (400 Bad Request)

``` json
{
  "error": "blog id is required"
}

 ```

- Returned when the blog ID path parameter is missing or empty.
    

### Notes

- If the like does not exist, the operation is still considered successful.
    
- Ensure the `blog_id` path parameter is a valid blog identifier in hexadecimal format

# Get Like Count for a Blog

## GET /blog/{blog_id}/likes/

Retrieve the total number of likes for a specific blog post by its ID.

---

### Request

- **Method:** GET
    
- **URL:** `/blog/{blog_id}/likes/`  
    Replace `{blog_id}` with the blog's unique identifier.
    
- **Headers:**
    

| Header | Value | Description |
| --- | --- | --- |
| Authorization | Bearer {{accessToken}} | Required. Bearer token for authentication. |

---

### Request Parameters

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| id | string | Yes | The unique ID of the blog to get the like count for (path parameter). |

---

### Response

#### Success (200 OK)

``` json
{
  "like_count": 10
}

 ```

- Returns the total number of likes for the specified blog.
    

#### Bad Request (400 Bad Request)

``` json
{
  "error": "blog id is required"
}

 ```

- Returned when the blog ID path parameter is missing or empty.
    

---

### Notes

- The `id` path parameter must be a valid blog identifier in hexadecimal format.
    
- Authentication is required to access this endpoint.

# Check If Blog Is Liked

## GET /blog/{id}/likes/is-liked

Check whether the authenticated user has liked a specific blog post.

---

### Request

- **Method:** GET
    
- **URL:** `/blog/{blog_id}/likes/is-liked`  
    Replace `{blog_id}` with the blog's unique identifier.
    
- **Headers:**
    

| Header | Value | Description |
| --- | --- | --- |
| Authorization | Bearer {{accessToken}} | Required. Bearer token for authentication. |

---

### Request Parameters

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| id | string | Yes | The unique ID of the blog to check (path parameter). |

---

### Response

#### Success (200 OK)

``` json
{
  "is_liked": true
}

 ```

#### Bad Request (400 Bad Request)

``` json
{
  "error": "blog id is required"
}

 ```

- Returned when the blog ID path parameter is missing or empty.
    

#### Unauthorized (401 Unauthorized)

``` json
{
  "error": "unauthorized"
}

 ```

- Returned if the user is not authenticated.
    

### Notes

- The `id` path parameter must be a valid blog identifier in hexadecimal format.
    
- Authentication is required to access this endpoint.
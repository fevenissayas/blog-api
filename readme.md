# Blog API

A comprehensive blog management system built with Go using Clean Architecture principles. The API provides user authentication, blog management, commenting, and liking functionality with secure JWT-based authentication.

## Features

### Authentication & Authorization

- User registration with email verification (OTP-based)
- Admin registration capability
- Secure login with JWT access tokens
- HTTP-only cookie-based refresh tokens
- Role-based access control (User/Admin)
- Password reset functionality

### Blog Management

- Create, read, update, delete blog posts
- Pagination support for blog listing
- Advanced filtering (by tags, date, popularity)
- Search functionality
- View count tracking
- AI-powered content suggestions

### Social Features

- Blog post commenting system
- Like/unlike functionality
- User profile management

### Security Features

- Password strength validation
- Secure password hashing (bcrypt)
- JWT token validation
- HTTP-only cookies for refresh tokens
- Rate limiting and input validation

## Tech Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: MongoDB
- **Authentication**: JWT tokens
- **Email Service**: SMTP (Gmail)
- **AI Service**: Google Gemini API
- **Architecture**: Clean Architecture

## Project Structure

```
blog-api/
├── Delivery/
│   ├── Controllers/     # HTTP handlers
│   └── Router/         # Route definitions
├── Domain/             # Business entities and interfaces
├── Infrastructure/     # External services and utilities
├── Repositories/       # Data access layer
├── Usecases/          # Business logic
└── main.go           # Application entry point
```

## Environment Variables

Create a `.env` file in the root directory:

```env
# JWT Configuration
JWT_SECRET=your_jwt_secret_key

# MongoDB Configuration
MONGODB_URI=mongodb://localhost:27017
DB_NAME=blog_api

# SMTP Email Configuration
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=465
EMAIL_USERNAME=your_email@gmail.com
EMAIL_PASSWORD=your_app_password
EMAIL_FROM=your_email@gmail.com

# AI Service
API_Key=your_gemini_api_key

# Server Port (optional)
PORT=8080
```

## Installation & Setup

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd blog-api
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Set up environment variables**

   - Copy `.env.example` to `.env`
   - Fill in your actual configuration values

4. **Start MongoDB**

   - Ensure MongoDB is running locally or update `MONGODB_URI` for remote instance

5. **Run the application**
   ```bash
   go run main.go
   ```

The API will be available at `http://localhost:8080`

## API Endpoints

### Authentication

- `POST /auth/register` - User/Admin registration
- `POST /auth/verify-email` - Email verification with OTP
- `POST /auth/login` - User login
- `POST /auth/logout` - User logout
- `POST /auth/refresh` - Refresh access token
- `POST /auth/promote` - Promote user to admin (Admin only)
- `POST /auth/update` - Update user profile

### Password Management

- `POST /password/request-reset` - Request password reset
- `POST /password/reset` - Reset password with token

### Blogs

- `GET /blogs` - Get paginated blogs
- `GET /blogs/:id` - Get single blog (increments view count)
- `POST /blogs` - Create new blog (Auth required)
- `PUT /blogs/:id` - Update blog (Auth required)
- `DELETE /blogs/:id` - Delete blog (Auth required)
- `GET /blogs/filter` - Filter blogs by tags/date
- `GET /blogs/search` - Search blogs
- `POST /blogs/aisuggestion` - Get AI content suggestions

### Blog Interactions

- `POST /blogs/:id/likes` - Like a blog
- `DELETE /blogs/:id/likes` - Unlike a blog
- `GET /blogs/:id/likes` - Get like count
- `GET /blogs/:id/likes/is-liked` - Check if user liked blog

### Comments

- `POST /blogs/:id/comments` - Add comment to blog
- `GET /blogs/:id/comments` - Get blog comments
- `DELETE /comments/:commentID` - Delete comment

## Authentication Flow

1. **Registration**: User provides email, username, password
2. **Email Verification**: OTP sent to email, must be verified before account activation
3. **Login**: Returns access token and sets refresh token in HTTP-only cookie
4. **Token Refresh**: Automatic refresh using secure cookies
5. **Logout**: Clears all tokens and cookies

## Security Features

- Passwords hashed using bcrypt
- JWT tokens with configurable expiration
- HTTP-only cookies for refresh tokens
- Input validation and sanitization
- Role-based authorization
- Secure email verification process

## Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -o blog-api main.go
```

### Docker Support

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support or questions, please create an issue in the repository or contact the maintainers.

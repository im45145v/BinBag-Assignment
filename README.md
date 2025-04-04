# BinBag-Assignment

A RESTful API for user authentication and profile management built with Go, Gin, and MongoDB.
[Postman Documentation](https://www.postman.com/cryosat-candidate-5965552/binbag/collection/r8lr2x3/api-documentation-reference?origin=tab-menu)
## Overview

This project provides a simple API for user registration, login, and profile management. It uses the following technologies:

- **Go**: Programming language.
- **Gin**: Web framework for building RESTful APIs.
- **MongoDB**: NoSQL database for storing user data.
- **JWT**: JSON Web Tokens for authentication.

---

## Project Structure

```
/workspaces/BinBag-Assignment
├── config/                # Configuration files
│   └── config.go          # Application configuration (e.g., MongoDB URI, JWT secret)
├── controllers/           # API controllers
│   └── authController.go  # Handles user registration, login, and profile retrieval
├── middlewares/           # Middleware for request handling
│   └── authMiddleware.go  # JWT authentication middleware
├── models/                # Data models
│   └── user.go            # User model and utility methods
├── routes/                # API route definitions
│   └── routes.go          # Registers API routes
├── utils/                 # Utility functions
│   └── jwt.go             # JWT generation and validation
├── main.go                # Application entry point
├── go.mod                 # Go module dependencies
├── go.sum                 # Dependency checksums
└── README.md              # Documentation
```

---

## API Endpoints

### Authentication

#### Register User
- **URL**: `/register`
- **Method**: `POST`
- **Auth Required**: No
- **Request Body**:
  ```json
  {
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "bio": "Software developer",
    "profile_picture": "https://example.com/profile.jpg"
  }
  ```
- **Success Response**:
  - **Code**: 201 Created
  - **Content**:
    ```json
    {
      "message": "User registered successfully"
    }
    ```
- **Error Response**:
  - **Code**: 400 Bad Request
  - **Content**: `{"error": "Invalid input"}`
  - **Code**: 500 Internal Server Error
  - **Content**: `{"error": "Failed to register user"}`

#### Login User
- **URL**: `/login`
- **Method**: `POST`
- **Auth Required**: No
- **Request Body**:
  ```json
  {
    "email": "john@example.com",
    "password": "password123"
  }
  ```
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    ```
- **Error Response**:
  - **Code**: 400 Bad Request
  - **Content**: `{"error": "Invalid input"}`
  - **Code**: 401 Unauthorized
  - **Content**: `{"error": "Invalid credentials"}`

### User Profile

#### Get User Profile
- **URL**: `/profile`
- **Method**: `GET`
- **Auth Required**: Yes (Bearer Token)
- **Headers**:
  ```
  Authorization: Bearer <JWT_TOKEN>
  ```
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "id": "60d21b4967d0d8992e610c85",
      "name": "John Doe",
      "email": "john@example.com",
      "address": "123 Main St",
      "bio": "Software developer",
      "profile_picture": "https://example.com/profile.jpg"
    }
    ```
- **Error Response**:
  - **Code**: 401 Unauthorized
  - **Content**: `{"error": "Authorization header is required"}`
  - **Code**: 404 Not Found
  - **Content**: `{"error": "User not found"}`

#### Update User Profile
- **URL**: `/update-profile`
- **Method**: `POST`
- **Auth Required**: Yes (Bearer Token)
- **Headers**:
  ```
  Authorization: Bearer <JWT_TOKEN>
  ```
- **Request Body**:
  ```json
  {
    "name": "John Doe",
    "address": "123 Main St",
    "bio": "Software developer",
    "profile_picture": "https://example.com/profile.jpg"
  }
  ```
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "message": "Profile updated successfully"
    }
    ```
- **Error Response**:
  - **Code**: 400 Bad Request
  - **Content**: `{"error": "Invalid input format"}`
  - **Code**: 401 Unauthorized
  - **Content**: `{"error": "Authorization header is required"}`
  - **Code**: 500 Internal Server Error
  - **Content**: `{"error": "Failed to update profile"}`

---

## Code Documentation

### `/main.go`
The entry point of the application:
- Initializes the Gin router.
- Connects to MongoDB.
- Registers API routes.
- Starts the server on port `8080`.

---

### `/config/config.go`
Handles application configuration:
- Defines the MongoDB connection URI.
- Configures the JWT secret key and token expiration time.

---

### `/controllers/authController.go`
Contains the following handlers:
1. **Register**: Handles user registration by hashing the password and storing user data in MongoDB.
2. **Login**: Authenticates users by verifying their email and password, then generates a JWT token.
3. **GetProfile**: Retrieves the authenticated user's profile using their ID from the JWT token.

---

### `/controllers/userController.go`
Contains the following handler:
1. **UpdateProfile**: Handles user profile updates by validating and updating user fields like `name`, `address`, `bio`, and `profile_picture`.

---

### `/middlewares/authMiddleware.go`
Defines the `AuthMiddleware`:
- Validates the JWT token from the `Authorization` header.
- Extracts the user ID from the token claims and sets it in the request context.

---

### `/models/user.go`
Defines the `User` struct:
- Fields: `ID`, `Name`, `Email`, `Password`, `Address`, `Bio`, `ProfilePicture`.
- Methods:
  - `HashPassword`: Hashes the user's password using bcrypt.
  - `CheckPassword`: Verifies the password against the hashed password.

---

### `/routes/routes.go`
Registers API routes:
- Public routes: `/register`, `/login`.
- Protected routes: `/profile`, `/update-profile` (requires JWT authentication).

---

### `/utils/jwt.go`
Provides utility functions for JWT:
- `GenerateToken`: Creates a JWT token with user ID and email as claims.
- `ValidateToken`: Validates the JWT token and returns the parsed token.

---

## Running the Application

### Prerequisites
- Install [Go](https://golang.org/dl/).
- Install [MongoDB](https://www.mongodb.com/try/download/community).

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/im45145v/BinBag-Assignment.git
   cd BinBag-Assignment
   ```

2. Set up environment variables (optional):
   - `JWT_SECRET_KEY`: Secret key for signing JWT tokens.

3. Run the application:
   ```bash
   go run main.go
   ```

4. The server will start on `http://localhost:8080`.

---

## Testing the API

### Using `curl`
#### Register a User
```bash
curl -X POST http://localhost:8080/register \
-H "Content-Type: application/json" \
-d '{"name":"John Doe","email":"john@example.com","password":"password123", "bio": "Software developer", "profile_picture": "https://example.com/profile.jpg"}'
```

#### Login
```bash
curl -X POST http://localhost:8080/login \
-H "Content-Type: application/json" \
-d '{"email":"john@example.com","password":"password123"}'
```

#### Get Profile
```bash
curl -X GET http://localhost:8080/profile \
-H "Authorization: Bearer <JWT_TOKEN>"
```

#### Update Profile
```bash
curl -X POST http://localhost:8080/update-profile \
-H "Authorization: Bearer <JWT_TOKEN>" \
-H "Content-Type: application/json" \
-d '{"name":"John Doe","address":"123 Main St","bio":"Software developer","profile_picture":"https://example.com/profile.jpg"}'
```

---

## Error Handling

The API returns appropriate HTTP status codes and error messages in JSON format:
- `200 OK`: Request succeeded.
- `201 Created`: Resource successfully created.
- `400 Bad Request`: Invalid input data.
- `401 Unauthorized`: Authentication required or failed.
- `404 Not Found`: Resource not found.
- `500 Internal Server Error`: Server-side error.

---

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

---

## Docker Deployment

### Prerequisites
- Install [Docker](https://www.docker.com/get-started).
- Install [Docker Compose](https://docs.docker.com/compose/install/).

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/im45145v/BinBag-Assignment.git
   cd BinBag-Assignment
   ```

2. Build the Docker image:
   ```bash
   docker build -t binbag-assignment .
   ```

3. Run the Docker container:
   ```bash
   docker run -p 8080:8080 binbag-assignment
   ```

4. The server will start on `http://localhost:8080`.

### Using Docker Compose

1. Clone the repository:
   ```bash
   git clone https://github.com/im45145v/BinBag-Assignment.git
   cd BinBag-Assignment
   ```

2. Create a `.env` file in the root directory and set the environment variables:
   ```env
   MONGO_URI=mongodb://mongo:27017
   DB_NAME=binbag_db
   USERS_COLLECTION=users
   JWT_SECRET_KEY=your_jwt_secret_key
   ```

3. Start the services using Docker Compose:
   ```bash
   docker-compose up
   ```

4. The server will start on `http://localhost:8080` and MongoDB will be available on `mongodb://localhost:27017`.

5. To stop the services, press `Ctrl+C` and run:
   ```bash
   docker-compose down
   ```


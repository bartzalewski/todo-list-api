# Todo List API

Welcome to the Todo List API! This is a simple API for managing a todo list with user authentication. The API supports creating, reading, updating, and deleting todo items, along with user signup and login functionalities.

## Features

- User authentication with JWT tokens
- CRUD operations for todo items
- Secure password storage using bcrypt
- Concurrent safe in-memory storage for users and todos

## Getting Started

Follow these instructions to get a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

Ensure you have the following installed on your system:

- [Go](https://golang.org/doc/install) (version 1.22+)
- [Git](https://git-scm.com/)

### Installation

1. **Clone the repository**:

   ```sh
   git clone https://github.com/bartzalewski/todo-list-api.git
   cd todo-list-api
   ```

2. **Initialize Go modules**:

   ```sh
   go mod tidy
   ```

3. **Run the application**:

   ```sh
   go run main.go
   ```

The server will start on `http://localhost:8080`.

## API Endpoints

### User Authentication

- **Sign Up**

  `POST /signup`

  Request:

  ```json
  {
    "username": "your-username",
    "password": "your-password"
  }
  ```

  Response:

  ```json
  {
    "status": "User created successfully"
  }
  ```

- **Sign In**

  `POST /signin`

  Request:

  ```json
  {
    "username": "your-username",
    "password": "your-password"
  }
  ```

  Response:

  ```json
  {
    "status": "Signed in successfully"
  }
  ```

  Cookie: `token`

### Todo Management

- **Create Todo**

  `POST /todos`

  Request:

  ```json
  {
    "title": "First Todo"
  }
  ```

  Response:

  ```json
  {
    "id": 1,
    "title": "First Todo",
    "completed": false,
    "created_at": "2023-05-24T12:34:56Z"
  }
  ```

- **Get Todos**

  `GET /todos`

  Response:

  ```json
  [
    {
      "id": 1,
      "title": "First Todo",
      "completed": false,
      "created_at": "2023-05-24T12:34:56Z"
    }
  ]
  ```

- **Update Todo**

  `PUT /todos/{id}`

  Request:

  ```json
  {
    "id": 1,
    "title": "Updated Todo",
    "completed": true
  }
  ```

  Response:

  ```json
  {
    "id": 1,
    "title": "Updated Todo",
    "completed": true
  }
  ```

- **Delete Todo**

  `DELETE /todos/{id}`

  Response: `204 No Content`

## Built With

- [Go](https://golang.org/) - The Go programming language
- [Gorilla Mux](https://github.com/gorilla/mux) - A powerful URL router and dispatcher for Golang
- [JWT-Go](https://github.com/dgrijalva/jwt-go) - A Go implementation of JSON Web Tokens
- [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - A package for password hashing

## Contributing

Feel free to submit issues or pull requests. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to the Go community for their invaluable resources and support.
- [Gorilla Mux](https://github.com/gorilla/mux) and [JWT-Go](https://github.com/dgrijalva/jwt-go) for making development easier.

---

Happy coding! 🚀

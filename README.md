# Webapp
  
## Installation
1. Install dependencies
```sh
$ go mod tidy
```
2. Set up environment variables
   Create a `config.yaml` file in the project root and add the necessary environment variables
```yaml
app:
    name: 
    environment:

database:
    host:
    port:
    username:
    password:
    name:

    # testing
    test_host:
    test_port:
    test_username:
    test_password:
    test_name:

server:
    port:

auth:
    jwt_secret:
    token_expiration: # in second

user:
    password_min_length:
```
3. Run the application
   ```sh
   $ go run cmd/api/main.go
   ```

## App Structure
This template follows the Domain-Driven Design (DDD) principle

### Folder Structure
```
project-root/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── auth/
│   │   ├── domain/
│   │   │   ├── user.go
│   │   │   ├── auth.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── application/
│   │   │   └── service.go
│   │   ├── infrastructure/
│   │   │   ├── jwt_auth.go
│   │   │   └── postgres_repository.go
│   │   └── interfaces/
│   │   │   ├── http/
│   │   │   │   └── middleware/
│   │   │   │       └── auth.go
│   │   │   └── http_handler.go
│   │   └── module.go
│   └── shared/
│       ├── auth/   
│       │   └── interface.go
│       ├── infrastructure/
│       │   └── postgres/
│       │       └── db.go
│       └── interfaces/
│           └── http/
│               └── server.go
└── pkg/
    └── apperrors/
        └── errors.go
```

### Interface
The `interface` layer is responsible for handling the communication between the external world and the application. This includes HTTP handlers, middleware, and any other components that interact with external systems.
- HTTP Handlers: These are responsible for handling HTTP requests and responses. They parse incoming requests, call the appropriate services, and return the responses.
- Middleware: Middleware functions are used to process requests before they reach the handlers. Common uses include authentication, logging, and request validation.
  
### Infrastructure
The `infrastructure` layer contains the implementation details of external systems and services that the application depends on. This includes database connections, third-party services, and other external resources.

- Database Repositories: These are responsible for interacting with the database. They contain methods for querying and manipulating data.
- JWT Authentication: This handles the creation and validation of JWT tokens for authentication purposes.
  
### Domain
The `domain` layer contains the core business logic of the application. This is where the main entities, value objects, and business rules are defined.

- Entities: These are the core objects of the application, such as `User` and `Auth`.
- Repositories: Interfaces that define the methods for data access. The actual implementation is provided in the infrastructure layer.
- Services: These contain the business logic and use the repositories to interact with the data.

### Application
The `application` layer contains the application-specific logic that orchestrates the use of domain services and entities. This layer is responsible for implementing use cases and application workflows.

- Services: These are higher-level services that coordinate the use of domain services to fulfill application-specific requirements.

## License
This project is licensed under a proprietary license. All rights are reserved by the owner. Unauthorized copying, distribution, or modification of this code is strictly prohibited.
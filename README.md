# Write Gin Like gRPC

## Example
```go

func main() {
	svc := &Service{}

	engine := gin.New()

	engine.POST("/users", JSON(svc.CreateUser))

	engine.Run()
}

func (s *Service) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	// do something.
}

type CreateUserRequest struct {
	User *User  `json:"user"`
}

type CreateUserResponse struct {
	User *User  `json:"user"`
}

func JSON[Request, Response any](f func(context.Context, *Request) (*Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Request

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := f(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

```

## [Detail Source Code](main.go)

## Try it

### Run the HTTP server
```sh
make run
```

### List Users with query parameters
```sh
curl -X GET 'http://localhost:8080/users?name=min'

# Output: {"users":[{"id":"1","name":"min","email":"min@mail.example.com"}]}
```

### Create User with body and query parameters
```sh
curl -X POST 'http://localhost:8080/users?opt=ignore' -d '{"user": {"name": "min", "email": "min@example.com"}}'

# Output: {"user":{"id":"","name":"min","email":"min@example.com"},"opt":"ignore"}
```

### Get User with path parameters
```sh
curl -X GET 'http://localhost:8080/users/123'

# Output: {"user":{"id":"123","name":"min","email":"min@mail.example.com"}}


# Try not found case
# Use `-i` show the status code
curl -i -X GET 'http://localhost:8080/users/404'

# Output:
# HTTP/1.1 404 Not Found
# Content-Type: application/json; charset=utf-8
#
# {"error":"user not found"}

```

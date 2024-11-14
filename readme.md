# Training Exercise 🎯
## Setup 🛠
1. Start Docker containers:
```bash
make docker-up
```
2. Start API:
```bash
make api
```

## API Exercise - User CRUD Implementation 🚀

I have implemented a complete CRUD API for user management with:
- Full unit tests for the service layer
- Redis caching implementation
- Error handling
- Input validation
- Clean architecture separation (handler, service, repository)
- Kafka event publishing for user creation events

Key files:
- Handler: `app/user/handler.go`
- Service: `app/user/service.go` 
- Repository: `app/user/repository.go`
- Unit Tests: `app/user/service_test.go`

## Basic Exercises 📚

I have completed exercises 1-9 in the `exercise` folder:

1. Basic boolean return ✅
2. Product struct creation with JSON structure ✅
3. Order struct with nested products ✅
4. JSON string to Product slice conversion ✅
5. Map to Product conversion ✅
6. Product search by code ✅
7. Product deletion from slice ✅
8. Calculate total product price ✅
9. Sort products by name and code ✅
# AGENTS.md - Developer Guide for omiai-server

This document provides guidelines for agentic coding agents working in this repository.

## Project Overview

omiai-server is a Go-based matchmaking server using Gin web framework, GORM for database, and Redis for caching/queues. The project follows a layered architecture (controller → biz → data).

## Build, Lint, and Test Commands

### Running the Application
```bash
# Build the binary
make build

# Run the server
./bin/server
```

### Linting
```bash
# Run linting (includes fmt check)
make lint

# Run go fmt only
make fmt
```

### Testing
```bash
# Run all tests with coverage
make test

# Run a single test
go test -v ./internal/controller/match -run TestParseTime

# Run tests in a specific package
go test -v ./internal/biz/omiai/...

# Run tests with verbose output
go test -v -gcflags=all=-l ./...
```

Note: The `-gcflags=all=-l` flag disables inlining which is sometimes needed for certain test scenarios.

## Code Style Guidelines

### Project Structure
```
internal/
├── api/         # API definitions
├── biz/         # Business logic layer (entities, interfaces)
├── conf/        # Configuration
├── controller/  # HTTP handlers (Gin)
├── cron/        # Scheduled jobs
├── data/        # Data layer (GORM repositories)
├── middleware/  # HTTP middleware
├── queues/      # Message queues
├── server/      # Server initialization (Wire DI)
├── service/     # External service integrations
└── validates/   # Request validation structs
cmd/
├── script/      # CLI scripts
├── seeder/      # Database seeding
└── server/      # Main entry point
pkg/
└── response/    # HTTP response helpers
```

### Naming Conventions
- **Packages**: lowercase, singular (e.g., `client`, `match`, `biz_omiai`)
- **Files**: lowercase with underscores (e.g., `followup_test.go`, `client_list.go`)
- **Structs**: PascalCase (e.g., `Client`, `Controller`)
- **Interfaces**: PascalCase with "Interface" suffix (e.g., `ClientInterface`)
- **Methods**: PascalCase (e.g., `TableName()`, `RealAge()`)
- **Variables**: camelCase (e.g., `req`, `clause`, `respList`)
- **Constants**: PascalCase or SCREAMING_SNAKE_CASE for enum values

### Import Organization
Go imports are organized in three groups (standard library, third-party, internal):
```go
import (
    "context"
    "time"
    
    "omiai-server/internal/biz"
    "omiai-server/internal/validates"
    "omiai-server/pkg/response"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)
```

### GORM Model Conventions
- Use `gorm:"column:..."` tags to map struct fields to database columns
- Add `comment:` tags for database schema documentation
- Use `-` to exclude fields from database (e.g., `gorm:"-"`)
- Add indexes with `index`, `uniqueIndex` tags
- Define `TableName()` method for each model:
```go
func (t *Client) TableName() string {
    return "client"
}
```

### JSON Serialization
- Use `json:"field_name"` tags for API responses
- Keep JSON field names in snake_case (e.g., `partner_id`, `created_at`)

### Error Handling
- Return errors from functions rather than logging directly in business logic
- Use the response package for HTTP error responses:
```go
// Validation errors
response.ValidateError(ctx, err, response.ValidateCommonError)

// Database errors
response.ErrorResponse(ctx, response.DBSelectCommonError, "获取客户列表失败")

// Success responses
response.SuccessResponse(ctx, "ok", map[string]interface{}{
    "list": respList,
})
```

### Validation
- Use `validates/` package for request validation structs
- Use ` ShouldBind` methods from Gin for request parsing:
```go
var req validates.ClientListValidate
if err := ctx.ShouldBind(&req); err != nil {
    response.ValidateError(ctx, err, response.ValidateCommonError)
    return
}
```

### Testing Conventions
- Test files are named `*_test.go` in the same package
- Use `github.com/stretchr/testify/assert` for assertions
- Use table-driven tests for multiple test cases:
```go
func TestParseTime(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        // test cases
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

### Database Queries
- Use `biz.WhereClause` struct for building dynamic queries:
```go
clause := &biz.WhereClause{
    OrderBy: "created_at desc",
    Where:   "1=1",
    Args:    []interface{}{},
}
if req.Name != "" {
    clause.Where += " AND name LIKE ?"
    clause.Args = append(clause.Args, "%"+req.Name+"%")
}
```

### Configuration
- Configuration files are in `configs/` directory
- Use Viper for configuration management (see `internal/conf/`)
- Never commit sensitive credentials; use `.gitignore` and environment variables

### Dependency Injection
- Uses Google Wire for compile-time dependency injection
- Run `make wire` to regenerate wire_gen.go files after modifying injectors

### Linter Configuration
The project uses golangci-lint with these linters enabled:
- bodyclose, deadcode, dogsled, durationcheck, errcheck
- exportloopref, govet, gosimple, gofmt, goconst
- gomnd, gocyclo, ineffassign, prealloc, revive
- staticcheck, structcheck, typecheck, unused, unconvert
- varcheck, whitespace, wastedassign

Run `make lint` before committing to ensure code passes all checks.

### Common Error Codes
See `pkg/response/code_common.go` for common error codes used in this project.

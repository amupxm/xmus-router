# XMUS Router v2.0 - Complete Redesign Roadmap

## Current State Analysis

### Strengths ✅
- **Performance**: Excellent benchmark results (55,938 ops with 3,421 ns/op)
- **Simplicity**: Clean basic routing concept
- **Lightweight**: Minimal dependencies

### Critical Issues ❌
- **Panic-driven errors**: Routes panic on conflicts instead of returning errors
- **No thread safety**: Concurrent registration will cause race conditions
- **No middleware support**: Essential for production use
- **Poor error handling**: Limited debugging capabilities
- **No parameter extraction**: Manual parsing required
- **Security gaps**: Missing CORS, security headers, rate limiting
- **No observability**: No metrics, tracing, or structured logging

## V2.0 Complete Redesign Plan

### 1. Core Architecture (Week 1-2)

#### Generic-First Design
```go
// Core router with generics for type safety
type Router[T Context] struct {
    tree *radixTree[T]
    middleware []Middleware[T]
    errorHandler ErrorHandler[T]
    config *Config
}

// Type-safe context
type Context interface {
    Request() *http.Request
    Response() ResponseWriter
    Param(key string) string
    Query(key string) string
    Set(key string, value any)
    Get(key string) (any, bool)
}

// Generic handlers
type Handler[T Context] func(T) error
type Middleware[T Context] func(Handler[T]) Handler[T]
```

#### Thread-Safe Route Registration
```go
type RouteRegistry[T Context] struct {
    mu sync.RWMutex
    routes map[string]*Route[T]
    compiled bool
}

func (r *Router[T]) Register(method, path string, handler Handler[T]) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if r.compiled {
        return ErrRouterCompiled
    }
    
    return r.addRoute(method, path, handler)
}
```

### 2. Type-Safe Parameter System (Week 2)

#### Parameter Extraction with Generics
```go
// Type-safe parameter binding
type ParamBinder[T any] interface {
    Bind(ctx Context, dest *T) error
}

// Usage examples
type UserParams struct {
    ID   int64  `param:"id" validate:"required,min=1"`
    Name string `query:"name" validate:"max=100"`
}

func GetUser(ctx *XmusContext) error {
    var params UserParams
    if err := ctx.BindParams(&params); err != nil {
        return err
    }
    
    user := ctx.Get("user").(*User) // Type-safe context access
    return ctx.JSON(200, user)
}
```

### 3. Security-First Design (Week 3)

#### Built-in Security Middleware
```go
// Security configuration
type SecurityConfig struct {
    CORS *CORSConfig
    RateLimit *RateLimitConfig
    JWT *JWTConfig
    CSRF bool
    Headers map[string]string
}

// Rate limiting with generics
type RateLimiter[T Context] struct {
    store   Store
    keyFunc func(T) string
    limit   Rate
}

// JWT middleware with type safety
func JWTMiddleware[T Context](config JWTConfig) Middleware[T] {
    return func(next Handler[T]) Handler[T] {
        return func(ctx T) error {
            token, err := extractToken(ctx.Request())
            if err != nil {
                return NewHTTPError(401, "Invalid token")
            }
            
            claims, err := validateToken(token, config)
            if err != nil {
                return NewHTTPError(401, "Token validation failed")
            }
            
            ctx.Set("user", claims.User)
            return next(ctx)
        }
    }
}
```

### 4. Advanced Middleware System (Week 3-4)

#### Composable Middleware Chain
```go
// Middleware registry with ordering
type MiddlewareRegistry[T Context] struct {
    global []MiddlewareEntry[T]
    groups map[string][]MiddlewareEntry[T]
}

type MiddlewareEntry[T Context] struct {
    middleware Middleware[T]
    priority   int
    conditions []Condition
}

// Route groups with middleware
func (r *Router[T]) Group(prefix string, middleware ...Middleware[T]) *Group[T] {
    return &Group[T]{
        router:     r,
        prefix:     prefix,
        middleware: middleware,
    }
}
```

### 5. Integration with XMUS Logger (Week 4)

#### Native Logger Integration
```go
// Request logging middleware
func RequestLogger[T Context](logger *xmuslogger.Logger) Middleware[T] {
    return func(next Handler[T]) Handler[T] {
        return func(ctx T) error {
            start := time.Now()
            
            err := next(ctx)
            
            logger.Info().
                Str("method", ctx.Request().Method).
                Str("path", ctx.Request().URL.Path).
                Int("status", ctx.Response().Status()).
                Dur("duration", time.Since(start)).
                Err(err).
                Msg("Request completed")
                
            return err
        }
    }
}

// Error logging with context
func ErrorHandler[T Context](logger *xmuslogger.Logger) func(error, T) {
    return func(err error, ctx T) {
        logger.Error().
            Err(err).
            Str("method", ctx.Request().Method).
            Str("path", ctx.Request().URL.Path).
            Interface("headers", ctx.Request().Header).
            Msg("Request error")
    }
}
```

### 6. Performance Optimizations (Week 5)

#### High-Performance Routing Tree
```go
// Optimized radix tree with generics
type RadixNode[T Context] struct {
    path     string
    wildcard bool
    param    string
    handler  Handler[T]
    children []*RadixNode[T]
    indices  string
}

// Zero-allocation parameter extraction
type ParamExtractor struct {
    keys   []string
    values []string
}

func (pe *ParamExtractor) Extract(path string, pattern string) {
    // Implementation that reuses slices to avoid allocations
}
```

### 7. Observability & Monitoring (Week 5-6)

#### Built-in Metrics
```go
// Prometheus-compatible metrics
type Metrics struct {
    requests   *prometheus.CounterVec
    duration   *prometheus.HistogramVec
    errors     *prometheus.CounterVec
    concurrent *prometheus.GaugeVec
}

// Tracing middleware
func TracingMiddleware[T Context](tracer trace.Tracer) Middleware[T] {
    return func(next Handler[T]) Handler[T] {
        return func(ctx T) error {
            spanCtx, span := tracer.Start(ctx.Request().Context(), 
                fmt.Sprintf("%s %s", ctx.Request().Method, ctx.Request().URL.Path))
            defer span.End()
            
            ctx.Request().WithContext(spanCtx)
            return next(ctx)
        }
    }
}
```

### 8. Authentication Framework (Week 6-7)

#### Pluggable Auth System
```go
// Auth provider interface
type AuthProvider[T Context] interface {
    Authenticate(T) (*User, error)
    Authorize(T, ...Permission) error
}

// JWT provider
type JWTProvider[T Context] struct {
    secret    []byte
    algorithm string
    claims    ClaimsFunc[T]
}

// OAuth2 provider
type OAuth2Provider[T Context] struct {
    config oauth2.Config
    store  TokenStore
}

// Multi-provider auth
func MultiAuth[T Context](providers ...AuthProvider[T]) Middleware[T] {
    return func(next Handler[T]) Handler[T] {
        return func(ctx T) error {
            var lastErr error
            for _, provider := range providers {
                if user, err := provider.Authenticate(ctx); err == nil {
                    ctx.Set("user", user)
                    return next(ctx)
                } else {
                    lastErr = err
                }
            }
            return NewHTTPError(401, "Authentication failed")
        }
    }
}
```

### 9. Developer Experience (Week 7-8)

#### Code Generation & Validation
```go
//go:generate xmus-router-gen -input=routes.yaml -output=routes_gen.go

// Automatic route registration from struct tags
type UserController struct{}

func (uc *UserController) GetUser(ctx *XmusContext) error `route:"GET /users/:id" auth:"jwt" validate:"id:required,numeric"`
func (uc *UserController) CreateUser(ctx *XmusContext) error `route:"POST /users" auth:"jwt" roles:"admin"`

// Automatic OpenAPI generation
func (r *Router[T]) GenerateOpenAPI() *openapi.Document {
    // Implementation
}
```

#### Enhanced Error Handling
```go
// Rich error types
type HTTPError struct {
    Code    int           `json:"code"`
    Message string        `json:"message"`
    Details []ErrorDetail `json:"details,omitempty"`
    Cause   error         `json:"-"`
}

// Error middleware with context
func ErrorMiddleware[T Context](logger *xmuslogger.Logger) Middleware[T] {
    return func(next Handler[T]) Handler[T] {
        return func(ctx T) error {
            defer func() {
                if r := recover(); r != nil {
                    logger.Error().
                        Interface("panic", r).
                        Str("path", ctx.Request().URL.Path).
                        Msg("Panic recovered")
                    
                    ctx.Response().WriteHeader(500)
                    ctx.JSON(500, HTTPError{
                        Code:    500,
                        Message: "Internal server error",
                    })
                }
            }()
            
            return next(ctx)
        }
    }
}
```

## Implementation Timeline

### Phase 1: Foundation (Weeks 1-2)
- [ ] Generic router architecture
- [ ] Thread-safe route registration
- [ ] Basic middleware system
- [ ] Type-safe context
- [ ] Parameter extraction

### Phase 2: Security & Auth (Weeks 3-4)  
- [ ] Security middleware (CORS, Rate limiting)
- [ ] JWT authentication
- [ ] OAuth2 integration
- [ ] CSRF protection
- [ ] XMUS Logger integration

### Phase 3: Performance & Observability (Weeks 5-6)
- [ ] Optimized routing tree
- [ ] Zero-allocation optimizations
- [ ] Metrics and monitoring
- [ ] Distributed tracing
- [ ] Health checks

### Phase 4: Developer Experience (Weeks 7-8)
- [ ] Route groups and composition
- [ ] Code generation tools
- [ ] OpenAPI documentation
- [ ] Validation framework
- [ ] Testing utilities

## Usage Examples

### Simple API
```go
func main() {
    router := xmus.New[*xmus.Context]()
    
    // Middleware
    router.Use(
        xmus.RequestLogger(logger),
        xmus.CORS(),
        xmus.RateLimit(100, time.Minute),
    )
    
    // Routes with type safety
    router.GET("/users/:id", GetUser)
    router.POST("/users", CreateUser)
    
    // Route groups
    api := router.Group("/api/v1", xmus.JWT(jwtConfig))
    api.GET("/profile", GetProfile)
    
    log.Fatal(http.ListenAndServe(":8080", router))
}

func GetUser(ctx *xmus.Context) error {
    id := ctx.ParamInt64("id")
    user, err := userService.GetByID(id)
    if err != nil {
        return xmus.NewHTTPError(404, "User not found")
    }
    return ctx.JSON(200, user)
}
```

### Enterprise API with Full Features
```go
func main() {
    config := &xmus.Config{
        Security: &xmus.SecurityConfig{
            CORS: &xmus.CORSConfig{
                AllowOrigins: []string{"https://app.example.com"},
                AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
            },
            RateLimit: &xmus.RateLimitConfig{
                Requests: 1000,
                Window:   time.Hour,
            },
            JWT: &xmus.JWTConfig{
                Secret:    os.Getenv("JWT_SECRET"),
                Algorithm: "HS256",
            },
        },
        Observability: &xmus.ObservabilityConfig{
            Metrics: true,
            Tracing: true,
            Health:  true,
        },
    }
    
    router := xmus.NewWithConfig[*xmus.Context](config)
    
    // Multi-auth support
    auth := xmus.MultiAuth(
        xmus.NewJWTProvider(config.Security.JWT),
        xmus.NewOAuth2Provider(oauthConfig),
        xmus.NewAPIKeyProvider(apiKeyConfig),
    )
    
    // Public routes
    router.GET("/health", healthCheck)
    router.POST("/auth/login", login)
    
    // Protected API
    api := router.Group("/api/v1", auth, xmus.RequireRoles("user"))
    api.GET("/users/:id", GetUser)
    api.PUT("/users/:id", UpdateUser)
    
    // Admin only
    admin := api.Group("/admin", xmus.RequireRoles("admin"))
    admin.DELETE("/users/:id", DeleteUser)
    
    // Start server with graceful shutdown
    server := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }
    
    xmus.StartWithGracefulShutdown(server)
}
```

## Key Benefits

### 🚀 Performance
- Zero-allocation routing
- Optimized radix tree
- Efficient middleware chains
- Benchmark targets: >50k ops/sec

### 🔒 Security
- Built-in CORS, CSRF, rate limiting
- JWT/OAuth2 ready
- Security headers by default
- Input validation and sanitization

### 🛠 Developer Experience  
- Full type safety with generics
- Auto-generated documentation
- Rich error handling
- Comprehensive testing tools

### 📊 Production Ready
- Structured logging with XMUS Logger
- Metrics and tracing
- Health checks
- Graceful shutdown

### 🔄 Future Proof
- Plugin architecture
- Extensible auth system
- Cloud-native features
- GraphQL support ready

This redesign transforms XMUS Router from a basic routing library into a comprehensive, production-ready web framework while maintaining its performance advantages and simplicity philosophy.
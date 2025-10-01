# PACT Router Implementation

A high-performance HTTP router implementation based on the **Path-Aware Compression Tree (PACT)** algorithm, optimized for modern web APIs with cache-line aware data structures and predictive hot-path caching.

## Features

- **Cache-Line Optimized**: 64-byte aligned nodes for optimal CPU cache performance
- **Hot Path Caching**: Pre-caches frequently accessed routes for O(1) lookups
- **Build-Time Analysis**: Two-phase optimization with route pattern recognition
- **Adaptive Storage**: Different storage strategies based on child count
- **Concurrent Access**: Lock-free reads with copy-on-write updates
- **Performance Monitoring**: Built-in metrics and health checks
- **SIMD Optimizations**: Optional SIMD instructions for prefix matching
- **Path Compression**: Optional compression for memory efficiency

## Quick Start

```go
package main

import (
    "fmt"
    "log"
)

func main() {
    // Create router
    router := NewPACTRouter()
    
    // Define routes
    routes := []Route{
        {Path: "/api/v1/users", Method: "GET", Handler: "getUsers"},
        {Path: "/api/v1/users/:id", Method: "GET", Handler: "getUser"},
        {Path: "/api/v1/posts", Method: "GET", Handler: "getPosts"},
    }
    
    // Build router with optimization
    router.Build(routes)
    
    // Lookup routes
    handler := router.Lookup("/api/v1/users")
    fmt.Println(handler) // Output: getUsers
}
```

## Advanced Usage

```go
// Create advanced router with custom configuration
config := &RouterConfig{
    HotPathCacheSize: 64,
    HotPathThreshold: 60.0,
    MaxMemoryUsage:   1024 * 1024 * 5, // 5MB
    CompressionEnabled: true,
    SIMDEnabled:      true,
    ConcurrentAccess: true,
}

router := NewAdvancedPACTRouter(config)
router.Build(routes)

// Concurrent lookups
handler := router.ConcurrentLookup("/api/v1/users")

// Performance monitoring
stats := router.GetStats()
fmt.Printf("Cache hit rate: %.2f%%\n", router.GetCacheHitRate())
fmt.Printf("Average lookup time: %.2f ns\n", router.GetAverageLookupTime())
```

## Performance Characteristics

| Operation | Best Case | Average Case | Worst Case |
|-----------|-----------|--------------|------------|
| Lookup (cached) | O(1) | O(1) | O(1) |
| Lookup (uncached) | O(1) | O(log n) | O(k) |
| Insert | O(1) | O(log n) | O(k) |

Where:
- n = number of routes
- k = path length

## Benchmarking

Run the included benchmarks:

```bash
go test -bench=.
```

Example output:
```
BenchmarkPACT-8             1000000    150 ns/op
BenchmarkPACTLookup-8       2000000    100 ns/op
BenchmarkPACTBuild-8         100000   1500 ns/op
BenchmarkPACTMemory-8        100000   2000 ns/op
```

## Configuration Options

```go
type RouterConfig struct {
    // Cache settings
    HotPathCacheSize int     // Default: 32
    HotPathThreshold float64 // Default: 70.0
    
    // Memory settings
    MaxMemoryUsage    uint64 // Default: 10MB
    CompressionEnabled bool   // Default: true
    
    // Performance settings
    SIMDEnabled      bool    // Default: true
    ConcurrentAccess bool    // Default: true
}
```

## Memory Usage

| Component | Space |
|-----------|-------|
| Node (empty) | 64 bytes |
| Node (with children) | 64 + 16×8 = 192 bytes |
| Hot path cache | O(h) where h = number of hot paths |
| Total | O(n) where n = number of routes |

## When to Use PACT

✅ **Good fit**:
- Web APIs with predictable route patterns
- Applications with 100+ routes
- High-throughput services (>10k req/s)
- REST APIs with common prefixes
- Services where latency matters

❌ **Not ideal**:
- Very few routes (<10)
- Highly dynamic routing (routes change frequently)
- Non-HTTP routing scenarios
- Extremely memory-constrained environments

## Examples

See `example_usage.go` for comprehensive examples including:
- Basic usage patterns
- Advanced configuration
- Performance comparison
- Real-world scenarios
- Error handling
- Monitoring and observability

## Testing

Run all tests:

```bash
go test -v
```

Run specific test categories:

```bash
go test -v -run TestPACTCorrectness
go test -v -run TestPACTHotPathCaching
go test -v -run TestPACTMemoryUsage
```

## Algorithm Details

This implementation is based on the PACT (Path-Aware Compression Tree) algorithm described in the research paper. Key innovations include:

1. **Build-Time Analysis**: Analyzes all routes before tree construction
2. **Cache-Line Optimization**: 64-byte aligned nodes for optimal CPU performance
3. **Hot Path Prediction**: Identifies and pre-caches frequently accessed routes
4. **Adaptive Storage**: Different storage strategies based on child count
5. **First-Child Optimization**: Optimizes for the most common child at each level

## Contributing

When contributing to this implementation:

1. **Maintain cache-line alignment**: Keep nodes at 64 bytes
2. **Profile memory access patterns**: Use `perf` to measure cache misses
3. **Benchmark realistically**: Use actual web traffic patterns
4. **Document tradeoffs**: Explain memory vs speed decisions
5. **Test edge cases**: Empty routes, very long paths, many parameters

## License

This implementation is provided for educational and research purposes.

## References

- [PACT Paper](https://rs3lab.github.io/assets/papers/2021/kim:pactree.pdf)
- Radix trees: Knuth, TAOCP Vol. 3
- Cache-oblivious algorithms: Frigo et al. (1999)
- Adaptive data structures: Sleator & Tarjan (1985)

```go
// Logger provides a thin wrapper around a Logr instance. This is a struct instead of
// an interface so that there are no allocations on the heap each interface method
// invocation. Normally not something to be concerned about, but logging calls for
// disabled levels should have as little CPU and memory impact as possible. Most of
// these wrapper calls will be inlined as well.
type Logger struct {
	log        *logr.Logger
	lockConfig *int32
}
```


### Why `Logger` is a `struct` instead of an `interface` (Performance Optimization)

The comment explains that this is a performance-critical design choice to 
**minimize allocations and CPU overhead**, particularly for disabled log levels. 
Here’s a breakdown of the reasoning:

---

### **1. Avoiding Heap Allocations with Interfaces**
In Go, when you use an **interface**, the underlying concrete type (e.g., a struct) is 
**boxed** into an interface value, which may allocate memory on the heap. For example:
```go
type LoggerInterface interface {
    Info(msg string)
}

func logSomething(logger LoggerInterface) {
    logger.Info("hello")  // May allocate heap memory if `logger` is a struct.
}
```
- Each interface method call incurs a small overhead (dynamic dispatch + potential heap allocation).

#### **Why This Matters for Logging**
Logging libraries often check if a log level (e.g., `Debug`, `Info`) is enabled before 
processing the log. If the logger were an interface:
```go
if logger.Enabled() {  // Interface method call → potential heap allocation.
    logger.Info("expensive log")  // Another allocation.
}
```
- Even **disabled logs** would waste CPU/memory due to interface overhead.

---

### **2. Benefits of Using a `struct` Wrapper**
By wrapping `logr.Logger` in a `struct` (value type), the code ensures:
- **No heap allocations**: Structs passed by value avoid boxing.
- **Method calls may be inlined**: The compiler can optimize calls to `Logger`’s methods 
  (since it’s a concrete type).
- **Zero-cost disabled logs**: Checks like `logger.Enabled()` become virtually free.

#### **Example of Inlining**
For a struct method:
```go
func (l Logger) Info(msg string) {
    if atomic.LoadInt32(l.lockConfig) == 0 {  // Cheap atomic check.
        l.log.Info(msg)  // Direct call to underlying logger.
    }
}
```
The compiler might inline this, eliminating the function call overhead entirely.

---

### **3. Why Not an Interface?**
Interfaces are typically preferred for abstraction, but here **performance trumps flexibility** because:
- Logging is a **high-frequency operation** (called often in hot paths).
- Disabled logs should **do almost nothing** (no allocations, minimal CPU).
- The `Logger` struct still abstracts the underlying `logr.Logger` 
  (just without interface overhead).

---

### **4. Trade-offs**
- **Pros**: Faster, no allocations, better for low-level logging.
- **Cons**: Less flexible (can’t swap implementations at runtime without changing the struct).

---

### **Key Takeaway**
The choice of a `struct` wrapper is a deliberate optimization to ensure logging has 
**near-zero overhead when disabled**, which is critical for performance-sensitive 
applications. Interfaces would introduce unnecessary costs for this use case.


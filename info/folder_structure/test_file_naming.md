In Go, the package declaration in test files (`_test.go`) can be either:

1. **`package packagename`** (same as main package)
2. **`package packagename_test`** (separate test package)

Here's the key difference:

### **1. `package packagename` (Normal Tests)**
- Tests are in the **same package** as the code being tested
- Can test **unexported** (lowercase) functions/variables
- Typical structure:
  ```go
  // mypkg.go
  package mypkg
  
  func ExportedFunc() {}
  func unexportedFunc() {}
  
  // mypkg_test.go
  package mypkg
  
  func TestExportedFunc(t *testing.T) {
      ExportedFunc()  // Can call directly
      unexportedFunc() // Can test internal functions
  }
  ```

### **2. `package packagename_test` (External Tests)**
- Tests are in a **separate package**
- Can only test **exported** (uppercase) functions/variables
- Typical structure:
  ```go
  // mypkg.go
  package mypkg
  
  func ExportedFunc() {}
  func unexportedFunc() {}
  
  // mypkg_test.go
  package mypkg_test
  
  import (
      "testing"
      "path/to/mypkg"
  )
  
  func TestExportedFunc(t *testing.T) {
      mypkg.ExportedFunc()  // Must use package prefix
      // mypkg.unexportedFunc() // ERROR: can't access
  }
  ```

### **Key Differences**
| Feature                | `package packagename` | `package packagename_test` |
|------------------------|-----------------------|----------------------------|
| Package scope          | Same package          | Separate package           |
| Access to unexported   | Yes                   | No                         |
| Import needed          | No                    | Yes                        |
| Testing perspective    | White-box testing     | Black-box testing          |
| Cyclic dependencies    | Possible              | Avoided                    |

### **When to Use Each**
- **Use same package (`packagename`)** when:
  - You need to test internal implementation details
  - Writing unit tests that depend on package internals

- **Use test package (`packagename_test`)** when:
  - You want to test only the public API
  - Need to avoid import cycles
  - Writing integration/black-box tests
  - Testing how other packages would use your code

### **Special Case: `package main`**
For `main` packages, you must use `package main_test` because you can't import `main`.

### **Example Project Structure**
```
mypkg/
├── mypkg.go         # package mypkg
├── internal_test.go # package mypkg (white-box)
└── external_test.go # package mypkg_test (black-box)
```

The choice depends on whether you want to test internals (same package) or just the 
public interface (test package). Many projects use both approaches for different 
types of tests.


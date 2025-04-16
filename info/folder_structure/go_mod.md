Here’s a **`go.mod` cheat sheet** covering the most common operations and syntax 
for Go modules:

---

### **1. Initialize a Module**
```bash
go mod init <module-name>  # e.g., `go mod init github.com/yourname/project`
```

---

### **2. Add/Upgrade Dependencies**
```bash
go get <package>@<version>  # e.g., `go get github.com/gorilla/mux@v1.8.0`
```
| Command | Purpose |
|---------|---------|
| `go get -u` | Upgrade to latest minor/patch version |
| `go get -u=patch` | Upgrade only patch versions |
| `go get package@latest` | Latest version |
| `go get package@none` | Remove dependency |

---

### **3. Tidy Dependencies**
```bash
go mod tidy  # Add missing/unused dependencies
```

---

### **4. Vendor Dependencies**
```bash
go mod vendor  # Create a `vendor/` directory
go build -mod=vendor  # Build using vendored deps
```

---

### **5. Replace/Exclude Dependencies**
In `go.mod`:
```go
replace (
    github.com/old/pkg => ./local/pkg  // Local replacement
    github.com/old/pkg v1.2.3 => github.com/new/pkg v1.4.5
)

exclude (
    github.com/bad/pkg v1.0.0  // Exclude specific version
)
```

---

### **6. Verify Checksums**
```bash
go mod verify  # Check dependencies haven’t been modified
```

---

### **7. Download Dependencies**
```bash
go mod download  # Download deps to cache (no install)
```

---

### **8. Graph Dependencies**
```bash
go mod graph  # Print dependency graph
```

---

### **9. Why Is a Package Used?**
```bash
go mod why <package>  # e.g., `go mod why github.com/gorilla/mux`
```

---

### **10. Clean Cache**
```bash
go clean -modcache  # Delete module cache
```

---

### **Key `go.mod` Directives**
| Directive | Purpose | Example |
|-----------|---------|---------|
| `module` | Declares module name | `module github.com/your/project` |
| `go` | Go version | `go 1.21` |
| `require` | Direct dependencies | `require github.com/foo/bar v1.2.3` |
| `replace` | Override a dependency | `replace github.com/A => ./local/A` |
| `exclude` | Ignore a version | `exclude github.com/B v1.0.0` |

---

### **Example `go.mod` File**
```go
module github.com/your/project

go 1.21

require (
    github.com/gorilla/mux v1.8.0
    golang.org/x/text v0.3.8
)

replace (
    github.com/old/pkg => ./local/pkg
)
```

---

### **Common Workflow**
1. **Initialize**: `go mod init`
2. **Add Deps**: `go get`
3. **Tidy**: `go mod tidy`
4. **Build**: `go build`

---

### **FAQ**
- **How to upgrade all dependencies?**  
  ```bash
  go get -u ./...
  go mod tidy
  ```
- **How to use a local package?**  
  ```go
  replace github.com/remote/pkg => /path/to/local/pkg
  ```
- **How to pin a version?**  
  ```bash
  go get package@v1.2.3
  ```


---
---
---

Here’s a **detailed cheat sheet for the `go.mod` file itself**, covering all directives 
(including `exclude`), their syntax, and practical examples:

---

### **1. `module` Directive**
**Purpose**: Declares the module’s name (must match the repository path).  
**Syntax**:
```go
module github.com/yourusername/project
```
**Notes**:
- Used by `go get` to fetch dependencies.
- Must be the first line in `go.mod`.

---

### **2. `go` Directive**
**Purpose**: Declares the minimum Go version required.  
**Syntax**:
```go
go 1.21
```
**Notes**:
- Affects toolchain behavior (e.g., `go mod tidy` uses this version).
- Not enforced during compilation unless you use version-specific features.

---

### **3. `require` Directive**
**Purpose**: Lists direct dependencies and their versions.  
**Syntax**:
```go
require (
    github.com/gorilla/mux v1.8.0
    golang.org/x/text v0.3.8
)
```
**Notes**:
- Versions can be:
  - **Semantic**: `v1.2.3`
  - **Pseudo-versions**: `v0.0.0-20210709164813-5f38f98423a1` (for unreleased commits).
- Added automatically by `go get` or `go mod tidy`.

---

### **4. `exclude` Directive**
**Purpose**: Ignores specific dependency versions (e.g., broken or insecure versions).  
**Syntax**:
```go
exclude (
    github.com/bad/dependency v1.0.0
    golang.org/x/net v1.2.3
)
```
**When to use**:
- A dependency version has critical bugs.
- Security vulnerabilities exist in a specific version.  
**Example**:
```bash
go mod edit -exclude=golang.org/x/net@v1.2.3
```

---

### **5. `replace` Directive**
**Purpose**: Overrides a dependency with another version or local path.  
**Syntax**:
```go
replace (
    github.com/old/module => github.com/forked/module v1.2.3
    github.com/local/module => ./local/module
)
```
**Use cases**:
- Testing local changes without publishing.
- Using a fork temporarily.  
**Example**:
```bash
go mod edit -replace=github.com/old/module=./local/module
```

---

### **6. `retract` Directive**
**Purpose**: Marks specific versions of your own module as deprecated.  
**Syntax**:
```go
retract (
    v1.0.0 // Published accidentally
    v1.1.0 // Contains security flaw
)
```
**Notes**:
- Added to warn users upgrading to problematic versions.
- Appears in `go list -m -versions <module>`.

---

### **7. `indirect` Tag**
**Purpose**: Marks dependencies not directly imported by your module.  
**Example**:
```go
require (
    github.com/foo/bar v1.2.3 // indirect
)
```
**How it’s added**:
- Auto-generated by `go mod tidy` for transitive dependencies.
- Manually added if you `require` an unused package.

---

### **Full Example `go.mod` File**
```go
module github.com/your/project

go 1.21

require (
    github.com/gorilla/mux v1.8.0
    golang.org/x/text v0.3.8 // indirect
)

exclude (
    github.com/insecure/pkg v1.0.0
)

replace (
    github.com/old/pkg => github.com/new/pkg v1.2.3
)

retract (
    v0.9.0 // Buggy release
)
```

---

### **Key Commands to Edit `go.mod`**
| Command | Purpose |
|---------|---------|
| `go mod edit -require=<package>@<version>` | Add a dependency |
| `go mod edit -exclude=<package>@<version>` | Exclude a version |
| `go mod edit -replace=<old>=<new>` | Replace a dependency |
| `go mod edit -retract=<version>` | Deprecate a version |

---

### **FAQ**
**Q: How to remove a `replace` or `exclude`?**  
```bash
go mod edit -dropreplace=github.com/old/pkg
go mod edit -dropexclude=github.com/bad/pkg@v1.0.0
```

**Q: How to downgrade a dependency?**  
```bash
go get github.com/foo/bar@v1.2.3  # Updates go.mod
go mod tidy
```

**Q: What’s the difference between `exclude` and `retract`?**  
- `exclude`: Your module ignores a version of someone else’s dependency.  
- `retract`: Your module warns users not to use specific versions of **your own module**.

---

This covers everything you’d need to manage `go.mod` like a pro! 🚀  


## An Example of using **exclude** in go.mod file

```gomod
exclude (
	github.com/RoaringBitmap/roaring v0.7.0
	github.com/RoaringBitmap/roaring v0.7.1
	github.com/dyatlov/go-opengraph v0.0.0-20210112100619-dae8665a5b09
	github.com/willf/bitset v1.2.0
)
```

### **Explanation of the `go.mod` Hack**

This `exclude` block is a **defensive workaround** to prevent automatic upgrades to 
specific dependency versions that would break the build due to **module path changes** 
or **incompatible updates**. Here's what's happening:

---

### **1. The Core Issue**
- **`willf/bitset` v1.2.0** changed its module path from:  
  `github.com/willf/bitset` → `github.com/bits-and-blooms/bitset`  
  - This is a **breaking change** because Go treats different module paths as entirely 
    different packages.
- Some dependencies (e.g., `RoaringBitmap/roaring` or `dyatlov/go-opengraph`) still reference 
  the **old path** (`github.com/willf/bitset`).  
  - If `bitset` upgrades to `v1.2.0`, those dependencies will fail to resolve, causing compilation errors.

---

### **2. Why `exclude` is Used**
- The `exclude` directive **blocks specific versions** from being selected by Go's dependency resolver.  
- By excluding `willf/bitset v1.2.0`, the module ensures:  
  - The dependency tree stays on `v1.1.0` (or earlier), which uses the old path (`github.com/willf/bitset`).  
  - Dependent packages (like `RoaringBitmap/roaring`) continue to work without updates.

---

### **3. Why Other Packages Are Excluded**
The `exclude` block also lists:  
- `github.com/RoaringBitmap/roaring v0.7.0` and `v0.7.1`  
- `github.com/dyatlov/go-opengraph v0.0.0-20210112100619-dae8665a5b09`  

These are likely excluded because:  
1. They **depend on `willf/bitset`** but haven’t updated to the new path (`bits-and-blooms/bitset`).  
2. They might pull in `bitset v1.2.0` transitively, triggering the same issue.

---

### **4. How This Hack Works**
1. **Dependency Resolution**:  
   - When you run `go get` or `go mod tidy`, Go's resolver avoids the excluded versions.  
   - Instead, it picks the latest **non-excluded** version (e.g., `willf/bitset v1.1.0`).  

2. **Impact**:  
   - Your project and its dependencies continue using the old module path.  
   - No urgent need to fork/update transitive dependencies.

---

### **5. Better Long-Term Solutions**
1. **Upstream Fixes**:  
   - Wait for dependencies (e.g., `RoaringBitmap/roaring`) to update their imports to 
     `github.com/bits-and-blooms/bitset`.  

2. **Manual `replace`**:  
   ```go
   replace github.com/willf/bitset => github.com/bits-and-blooms/bitset v1.2.0
   ```
   - Forces all dependencies to use the new path.  
   - Requires testing (might break other dependencies).  

3. **Fork and Patch**:  
   - Fork affected dependencies and update their `go.mod` to use the new path.  

---

### **Key Takeaways**
- **`exclude` is a temporary fix** for version conflicts.  
- **Root cause**: A dependency (`bitset`) changed its module path, breaking backward compatibility.  
- **Alternatives**: `replace` or waiting for upstream fixes are cleaner but may require more work.  

This hack is a **band-aid** to keep the build working until dependencies catch up. 
Use `go mod why` to analyze dependency chains!  

```bash
go mod why github.com/willf/bitset
```


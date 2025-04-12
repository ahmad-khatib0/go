The following `fixedBigInt` struct represents a fixed-size multi-word integer. This is essentially 
a way to handle large integers that cannot fit into a single machine word (e.g., a 64-bit `uint64`). 

---

### **1. `fixedBigInt` Struct**
```go
type fixedBigInt struct {
	bits    []uint64
	n       uint
	msbMask uint64
}
```
- **`bits []uint64`**:
  - This is a slice of `uint64` values that stores the individual chunks (or "words") of the
    multi-word integer.
  - Each `uint64` in the slice represents a 64-bit chunk of the overall integer.
  - For example, if the integer is 128 bits long, `bits` will have 2 elements.

- **`n uint`**:
  - This is the total number of bits in the multi-word integer.
  - For example, if `n = 128`, the integer is 128 bits long.

- **`msbMask uint64`**:
  - This is a bitmask used to handle the most significant bit (MSB) of the integer.
  - It ensures that only the relevant bits in the most significant chunk are used, and the 
    rest are masked out.

---

### **2. `newFixedBigInt` Function**
```go
func newFixedBigInt(n uint) *fixedBigInt {
	chunkSize := (n + 63) / 64
	if chunkSize == 0 {
		chunkSize = 1
	}

	return &fixedBigInt{
		bits:    make([]uint64, chunkSize),
		n:       n,
		msbMask: (1 << (64 - n%64)) - 1,
	}
}
```
This function creates and initializes a new `fixedBigInt` instance.

#### **Step-by-Step Breakdown**
1. **Calculate `chunkSize`**:
   ```go
   chunkSize := (n + 63) / 64
   ```
   - This calculates the number of `uint64` chunks required to store `n` bits.
   - The formula `(n + 63) / 64` ensures that any remainder (bits that don’t fill a full 
     64-bit chunk) are accounted for by rounding up.
   - For example:
     - If `n = 64`, `chunkSize = 1`.
     - If `n = 65`, `chunkSize = 2`.

2. **Handle Edge Case**:
   ```go
   if chunkSize == 0 {
       chunkSize = 1
   }
   ```
   - If `n = 0`, `chunkSize` would be 0, which is invalid. This ensures that at least one 
     chunk is allocated.

3. **Create the `fixedBigInt` Instance**:
   ```go
   return &fixedBigInt{
       bits:    make([]uint64, chunkSize),
       n:       n,
       msbMask: (1 << (64 - n%64)) - 1,
   }
   ```
   - **`bits`**:
     - A slice of `uint64` is created with the calculated `chunkSize`.
     - This slice will store the bits of the multi-word integer.
   - **`n`**:
     - The total number of bits is stored in the `n` field.
   - **`msbMask`**:
     - This is a bitmask that ensures only the relevant bits in the most significant chunk 
       are used.
     - The formula `(1 << (64 - n%64)) - 1` calculates the mask:
       - `n % 64` gives the number of bits used in the most significant chunk.
       - `64 - (n % 64)` gives the number of unused bits in the most significant chunk.
       - `1 << (64 - n%64)` creates a mask with a `1` in the position of the first unused bit.
       - Subtracting `1` from this value creates a mask with `1`s in all the relevant bit positions.
     - For example:
       - If `n = 65`, the most significant chunk uses 1 bit (`65 % 64 = 1`).
       - `64 - 1 = 63`, so `1 << 63` is `0x8000000000000000`.
       - Subtracting `1` gives `0x7FFFFFFFFFFFFFFF`, which masks out the unused bits.

---

### **Example Usage**
Let’s say we want to create a `fixedBigInt` with `n = 130` bits.

1. **Calculate `chunkSize`**:
   ```go
   chunkSize := (130 + 63) / 64 = 193 / 64 = 3
   ```
   - We need 3 `uint64` chunks to store 130 bits.

2. **Calculate `msbMask`**:
   ```go
   msbMask := (1 << (64 - 130%64)) - 1
   ```
   - `130 % 64 = 2` (130 bits use 2 bits in the most significant chunk).
   - `64 - 2 = 62`.
   - `1 << 62` is `0x4000000000000000`.
   - Subtracting `1` gives `0x3FFFFFFFFFFFFFFF`.

3. **Create the `fixedBigInt`**:
   ```go
   fb := &fixedBigInt{
       bits:    make([]uint64, 3),
       n:       130,
       msbMask: 0x3FFFFFFFFFFFFFFF,
   }
   ```
   - The `bits` slice has 3 elements, each initialized to `0`.
   - The `msbMask` ensures that only the first 2 bits of the third chunk are used.

---

### **Purpose of `fixedBigInt`**
The `fixedBigInt` struct is designed to:
1. Handle large integers that cannot fit into a single machine word 
   (e.g., 128-bit, 256-bit, or larger integers).
2. Provide a fixed-size representation, ensuring that the integer always uses exactly `n` bits.
3. Use a bitmask (`msbMask`) to handle the most significant chunk efficiently, avoiding 
   unnecessary operations on unused bits.

---

### **Use Cases**
This kind of data structure is useful in:
- Cryptography (e.g., handling large keys or modular arithmetic).
- Arbitrary-precision arithmetic (e.g., big integer libraries).
- Low-level bit manipulation tasks where fixed-size integers are required.

---

### **Summary**
- The `fixedBigInt` struct represents a fixed-size multi-word integer.
- The `newFixedBigInt` function initializes the struct by:
  - Calculating the number of `uint64` chunks needed.
  - Creating a bitmask (`msbMask`) to handle the most significant chunk.
- This structure is useful for handling large integers efficiently and ensuring
  fixed-size behavior.

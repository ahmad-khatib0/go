### **Set Operations**
**Set operations** are mathematical operations performed on collections of elements 
   (sets) to combine, compare, or filter them. Common set operations include:

1. **Union (OR, `∪`)**  
   - Combines elements from both sets.  
   - Example: `{1, 2} ∪ {2, 3} = {1, 2, 3}`  

2. **Intersection (AND, `∩`)**  
   - Returns elements present in both sets.  
   - Example: `{1, 2} ∩ {2, 3} = {2}`  

3. **Difference (NOT, `\`)**  
   - Returns elements in the first set but not the second.  
   - Example: `{1, 2} \ {2, 3} = {1}`  

4. **Symmetric Difference (XOR, `∆`)**  
   - Returns elements in either set but not both.  
   - Example: `{1, 2} ∆ {2, 3} = {1, 3}`  

---

### **Bitmaps (Bit Arrays)**
A **bitmap** (or bit array/bitset) is a compact way to represent a set of integers using 
binary bits, where each bit indicates the presence (`1`) or absence (`0`) of an element.

#### **How It Works**
- Suppose we have a set `{1, 3, 5}` for numbers `0` to `7`:
  ```
  Position: 0 1 2 3 4 5 6 7  
  Value:    0 1 0 1 0 1 0 0  
  ```
  - The bitmap is `01010100` (binary) or `0x54` (hex).

#### **Advantages**
- **Fast logical operations**: Bitwise AND/OR/XOR compute intersections/unions in 
   constant time.
- **Memory-efficient**: For dense sets, it uses ~1 bit per element.

#### **Limitations**
- **Wasteful for sparse sets**: Storing `{1, 1000000}` requires 1 million bits.
- **No compression**: Traditional bitmaps don’t optimize for gaps.

---

### **Roaring Bitmaps: Optimized Bitmaps**
Roaring Bitmaps solve traditional bitmap limitations by:
1. **Hybrid Storage**:
   - **Array**: Stores sparse elements compactly (e.g., `[1, 1000000]`).
   - **Bitmap**: For dense ranges (e.g., `0-65535` as a 8192-byte block).
   - **Run-length encoding (RLE)**: Compresses long consecutive values.

2. **Performance**:
   - Faster than hash sets (e.g., Java `HashSet`) for large datasets.
   - Outperforms compressed bitmaps (like EWAH) in most cases.

---

### **Example: Set Operations with Bitmaps**
```python
# Python example using the `roaringbitmap` library
from roaringbitmap import RoaringBitmap

# Create two sets
set1 = RoaringBitmap([1, 2, 3, 1000])
set2 = RoaringBitmap([3, 4, 1000, 5000])

# Union
union = set1 | set2  # {1, 2, 3, 4, 1000, 5000}

# Intersection
intersection = set1 & set2  # {3, 1000}

# Difference
diff = set1 - set2  # {1, 2}

# Symmetric difference
sym_diff = set1 ^ set2  # {1, 2, 4, 5000}
```

---

### **Use Cases**
1. **Databases**  
   - Filtering rows efficiently (e.g., PostgreSQL, Apache Druid).  
2. **Search Engines**  
   - Inverted indexes for document retrieval (e.g., Elasticsearch).  
3. **Data Analytics**  
   - Accelerating OLAP queries (e.g., ClickHouse).  
4. **Networking**  
   - IP address filtering or firewall rules.  

---

### **Why Not Always Use Bitmaps?**
- **Hash sets** are simpler for small, non-integer data.  
- **Bloom filters** are better for probabilistic membership tests.  
- **Traditional bitmaps** waste memory for sparse data.  

Roaring Bitmaps strike a balance by dynamically choosing the best storage method.  


---
---
---
# Example Use Case 

### **How Databases Use Roaring Bitmaps for Efficient Row Filtering**

Databases leverage **Roaring Bitmaps** (or similar compressed bitmaps) to accelerate 
**filtering**, **indexing**, and **query execution** by representing row matches as 
compact bit vectors. Here’s a step-by-step breakdown with real-world examples:

---

## **1. Storing Row IDs as Bitmaps**
When a database executes a query like:
```sql
SELECT * FROM users WHERE age > 30 AND country = 'US';
```
It needs to quickly identify which rows match the conditions. Instead of scanning every 
row, it uses **bitmap indexing**:

### **Bitmap Index Structure**
- Each **distinct value** in a column (e.g., `country='US'`) gets a bitmap.
- Each **bit** in the bitmap corresponds to a **row ID**:
  - `1` = row matches the value.  
  - `0` = row does not match.

**Example:**
| Row ID | `country='US'` | `age > 30` |
|--------|----------------|------------|
| 1      | 1              | 0          |
| 2      | 0              | 1          |
| 3      | 1              | 1          |
| 4      | 1              | 0          |

- Bitmap for `country='US'`: `1011` (rows 1, 3, 4 match)  
- Bitmap for `age > 30`: `0110` (rows 2, 3 match)  

---

## **2. Fast Filtering with Bitwise Operations**
To combine conditions, databases use **bitwise operations**:
- **AND (`&`)** → Intersection (rows matching **both** conditions).  
- **OR (`|`)** → Union (rows matching **either** condition).  
- **NOT (`~`)** → Complement (rows **not** matching).  

**Example Query:**  
```sql
SELECT * FROM users WHERE country = 'US' AND age > 30;
```
- Compute `1011 (US) & 0110 (age>30) = 0010` (only **row 3** matches both).  

**Result:** Only row 3 is returned, avoiding full-table scans.

---

## **3. Roaring Bitmap Optimizations**
Traditional bitmaps waste space for sparse data (e.g., if row IDs are 1 and 1,000,000).
**Roaring Bitmaps solve this** by:
1. **Splitting into Containers**  
   - Divide row IDs into chunks (e.g., 16-bit blocks).  
   - For each chunk, choose the best storage:  
     - **Array**: Sparse values (e.g., `[1, 1000]`).  
     - **Bitmap**: Dense ranges (e.g., `0-65535`).  
     - **Run-length encoding (RLE)**: Long consecutive values (e.g., `1-10000`).  

2. **CPU Cache Efficiency**  
   - Bitwise operations are optimized for modern CPUs.  

**Before (Raw Bitmap):**  
Storing rows `{1, 1000000}` requires 1,000,000 bits (125 KB).  

**After (Roaring Bitmap):**  
Stores only `[1, 1000000]` in ~8 bytes (compressed).  

---

## **4. Real-World Database Implementations**
### **A. Apache Druid**
- Uses Roaring Bitmaps for **segment filtering**.  
- Each segment (data shard) stores bitmaps for column values.  
- Queries like `WHERE user_id IN (1, 2, 3)` are resolved via bitmap unions.  

### **B. ClickHouse**
- Implements **bitmap indexes** for low-cardinality columns.  
- Aggregations (e.g., `COUNT(DISTINCT user_id)`) use bitmap merges.  

### **C. PostgreSQL**
- Uses **Bitmap Index Scans** to combine multiple conditions.  
- Converts B-tree indexes to in-memory bitmaps for fast filtering.  

### **D. Elasticsearch/Lucene**
- **Inverted indexes** store document IDs as bitmaps.  
- Boolean queries (`AND/OR/NOT`) are resolved via bitwise ops.  

---

## **5. Performance Gains**
| Approach           | Time Complexity | Memory Usage | Example Query Speed |
|--------------------|-----------------|--------------|---------------------|
| Full Table Scan    | O(N)            | High         | Slow (reads all rows) |
| B-tree Index       | O(log N)        | Medium       | Fast for single-column |
| **Roaring Bitmap** | O(1)* per op    | Low          | **10-100x faster** for multi-column filters |

_(*Bitwise operations are effectively O(1) for fixed-size chunks.)_

---

## **6. Limitations**
- **Best for integer IDs**: Less efficient for string-heavy filters.  
- **Overhead for small datasets**: Hash sets may be faster for tiny tables.  
- **Updates can be costly**: Modifying bitmaps requires recompression.  

---

### **Summary**
Databases use **Roaring Bitmaps** to:
1. **Encode row matches** as compressed bit vectors.  
2. **Accelerate filtering** via bitwise ops (`AND/OR/NOT`).  
3. **Reduce I/O and CPU usage** by skipping unmatched rows.  

This is why systems like Druid, ClickHouse, and Elasticsearch handle 
**high-cardinality filters** so efficiently!  


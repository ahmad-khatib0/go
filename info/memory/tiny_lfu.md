### **TinyLFU (Tiny Least Frequently Used) Eviction Policy**

TinyLFU is a **probabilistic cache eviction algorithm** designed to balance high hit 
rates with low memory overhead. It’s a lightweight approximation of the classic 
**LFU (Least Frequently Used)** policy, optimized for modern caching workloads. 
Here’s how it works and why it’s used in Ristretto:

---

## **1. Core Idea**
- **Goal**: Evict the **least frequently accessed** items when the cache is full.
- **Trade-off**: Unlike exact LFU (which requires heavy metadata), TinyLFU 
  **approximates** frequency counts *efficiently*.

---

## **2. How TinyLFU Works**
### **(1) Frequency Sketch (Probabilistic Counting)**
- Instead of tracking exact access counts for every item, TinyLFU uses a 
  **Count-Min Sketch** (a compact, probabilistic data structure).
  - **Hash functions** map items to a small fixed number of counters.
  - Increments counters on each access.
  - Estimates frequency by taking the **minimum value** across hashed counters 
    (reduces hash collision errors).

### **(2) Doorkeeper (Freshness Filter)**
- A **Bloom filter** or small cache that:
  - Filters out **one-hit wonders** (items accessed only once).
  - Prevents them from polluting the frequency sketch.

### **(3) Eviction Decision**
- When the cache is full, TinyLFU compares the **estimated frequency** of the new item 
  vs. the victim (oldest item in a sample).
  - If the new item’s frequency > victim’s frequency, it replaces the victim.
  - Otherwise, the new item is discarded.

---

## **3. Why Use TinyLFU?**
| **Feature**               | **TinyLFU** vs. Other Policies                                  |
|---------------------------|----------------------------------------------------------------|
| **Memory Efficiency**     | Uses ~10 bits per item (vs. exact LFU’s 32+ bits).             |
| **Speed**                 | O(1) operations (thanks to hashing).                           |
| **Adaptability**          | Resists pollution by infrequent items (e.g., one-hit wonders). |
| **Scalability**           | Works well for large caches (unlike exact LFU).                |

**Compared to:**
- **LRU (Least Recently Used)**: Only cares about *recency*, not frequency.  
  → Bad for *scan-resistant* workloads (e.g., repeated access to old items).  
- **LFU (Exact)**: Accurate but memory-heavy (needs per-item counters).  

---

## **4. TinyLFU in Ristretto**
Ristretto’s config leverages TinyLFU via:
```go
NumCounters: 2500  // Size of the frequency sketch (10x MaxCost).
MaxCost:     250   // Max items (or total "cost") in cache.
```
- **Why `NumCounters = 10 * MaxCost`?**  
  - Ensures the frequency sketch has enough resolution to distinguish 
    frequently vs. rarely used items.  
  - Empirical testing shows this ratio balances accuracy/memory overhead.  

---

## **5. Example Workflow**
1. **Item `A` is accessed**:  
   - TinyLFU hashes `A` into the Count-Min Sketch, increments counters.  
2. **Cache fills up**:  
   - New item `B` arrives. TinyLFU compares `B`’s estimated frequency to a sampled victim (`C`).  
   - If `B` > `C`, `B` replaces `C`.  
3. **One-hit wonders**:  
   - Items with only one access are filtered by the Doorkeeper (avoiding cache pollution).  

---

## **6. When to Use TinyLFU?**
✅ **Read-heavy workloads** (e.g., caching database queries, API responses).  
✅ **Mixed access patterns** (some items hot, some cold).  
❌ **Write-heavy workloads** (TinyLFU’s overhead may not justify gains).  

---

### **Key Takeaways**
- TinyLFU approximates **"least frequently used"** eviction with minimal memory.  
- Ristretto uses it via `NumCounters` (frequency sketch size) and `MaxCost` (cache capacity).  
- Beats LRU in most real-world scenarios (better hit rates).  

For deep dives, see:  
- [TinyLFU Paper](https://arxiv.org/abs/1512.00727)  
- [Ristretto’s Design](https://dgraph.io/blog/post/introducing-ristretto/)


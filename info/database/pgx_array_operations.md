# **Using PostgreSQL Arrays (`TEXT[]`) with `pgx/v5` in Go**

Here’s a complete guide to working with **array operations**
(like permissions stored as `TEXT[]`) using the `pgx/v5` library.

---

## **1. Basic Setup**

### **(A) Import pgx**

```go
import (
    "context"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)
```

### **(B) Connect to PostgreSQL**

```go
pool, err := pgxpool.New(context.Background(), "postgres://user:password@localhost:5432/db")
if err != nil {
    log.Fatal(err)
}
defer pool.Close()
```

---

## **2. Inserting Arrays**

### **(A) Simple Insert**

```go
_, err = pool.Exec(context.Background(), `
    INSERT INTO roles (id, permissions)
    VALUES ($1, $2)`,
    "role123",
    []string{"create_post", "delete_post"}, // Automatically converted to `TEXT[]`
)
if err != nil {
    log.Fatal(err)
}
```

### **(B) Using `pgtype.TextArray` (Explicit Control)**

```go
import "github.com/jackc/pgx/v5/pgtype"

permissions := pgtype.TextArray{
    Elements: []pgtype.Text{
        {String: "create_post", Valid: true},
        {String: "delete_post", Valid: true},
    },
    Dimensions: []pgtype.ArrayDimension{{Length: 2, LowerBound: 1}},
    Status:     pgtype.Present,
}

_, err = pool.Exec(context.Background(), `
    INSERT INTO roles (id, permissions)
    VALUES ($1, $2)`,
    "role123",
    permissions,
)
```

---

## **3. Querying Arrays**

### **(A) Check if Array Contains a Value (`ANY`)**

```go
var roleID string
err = pool.QueryRow(context.Background(), `
    SELECT id FROM roles
    WHERE $1 = ANY(permissions)`,
    "create_post",
).Scan(&roleID)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Role with permission:", roleID)
```

### **(B) Check if Array Contains All Values (`@>`)**

```go
var hasAllPermissions bool
err = pool.QueryRow(context.Background(), `
    SELECT $1::TEXT[] <@ permissions
    FROM roles WHERE id = $2`,
    []string{"create_post", "delete_post"},
    "role123",
).Scan(&hasAllPermissions)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Has all permissions:", hasAllPermissions)
```

### **(C) Get All Rows with Matching Permissions**

```go
rows, err := pool.Query(context.Background(), `
    SELECT id, permissions FROM roles
    WHERE permissions && $1::TEXT[]`,  -- Overlap operator (`&&`)
    []string{"create_post", "admin"},
)
defer rows.Close()

for rows.Next() {
    var (
        id          string
        permissions []string
    )
    if err := rows.Scan(&id, &permissions); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Role %s has permissions: %v\n", id, permissions)
}
```

---

## **4. Updating Arrays**

### **(A) Append to an Array (`||`)**

```go
_, err = pool.Exec(context.Background(), `
    UPDATE roles
    SET permissions = permissions || $1
    WHERE id = $2`,
    []string{"new_permission"},
    "role123",
)
```

### **(B) Remove a Value (`array_remove`)**

```go
_, err = pool.Exec(context.Background(), `
    UPDATE roles
    SET permissions = array_remove(permissions, $1)
    WHERE id = $2`,
    "delete_post",
    "role123",
)
```

---

## **5. Using Custom Types (Optional)**

If you frequently work with arrays, define a **custom type**:

```go
type Permissions []string

// Scan implements sql.Scanner
func (p *Permissions) Scan(src interface{}) error {
    arr := pgtype.TextArray{}
    if err := arr.Scan(src); err != nil {
        return err
    }
    *p = make(Permissions, len(arr.Elements))
    for i, e := range arr.Elements {
        (*p)[i] = e.String
    }
    return nil
}

// Value implements driver.Valuer
func (p Permissions) Value() (driver.Value, error) {
    if p == nil {
        return nil, nil
    }
    arr := pgtype.TextArray{}
    arr.Elements = make([]pgtype.Text, len(p))
    for i, s := range p {
        arr.Elements[i] = pgtype.Text{String: s, Valid: true}
    }
    arr.Dimensions = []pgtype.ArrayDimension{{Length: int32(len(p)), LowerBound: 1}}
    arr.Status = pgtype.Present
    return arr.Value()
}
```

Now use it in queries:

```go
var perms Permissions
err = pool.QueryRow(context.Background(), "SELECT permissions FROM roles WHERE id = $1", "role123").Scan(&perms)
```

---

## **6. Performance Tips**

1. **Index Arrays for Faster Queries**:
   ```sql
   CREATE INDEX idx_roles_permissions ON roles USING GIN(permissions);
   ```
2. **Use `&&` (Overlap) for Multi-Permission Checks**:
   ```go
   rows, err := pool.Query(ctx, `
       SELECT id FROM roles
       WHERE permissions && $1`,  -- Returns roles with ANY matching permission
       []string{"perm1", "perm2"},
   )
   ```

---

## **7. Full Example (Mattermost-Style Permissions)**

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://user:password@localhost:5432/mattermost")
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Insert a role with permissions
	_, err = pool.Exec(ctx, `
		INSERT INTO roles (id, name, permissions)
		VALUES ($1, $2, $3)`,
		"role_1",
		"team_admin",
		[]string{"create_post", "delete_post", "invite_user"},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Query roles with a specific permission
	var roleID string
	err = pool.QueryRow(ctx, `
		SELECT id FROM roles
		WHERE 'delete_post' = ANY(permissions)
		LIMIT 1`,
	).Scan(&roleID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Role with 'delete_post':", roleID)
}
```

---

## **Summary**

| Operation           | `pgx/v5` Code                                     |
| ------------------- | ------------------------------------------------- | --- | --- |
| **Insert**          | `pool.Exec(..., []string{"a", "b"})`              |
| **Query (`ANY`)**   | `WHERE $1 = ANY(permissions)`                     |
| **Query (`@>`)**    | `WHERE permissions @> $1`                         |
| **Update (Append)** | `SET permissions = permissions                    |     | $1` |
| **Update (Remove)** | `SET permissions = array_remove(permissions, $1)` |

**Recommendation**:  
✅ Use `TEXT[]` for permissions (faster queries).  
✅ Use `pgx`’s built-in array support (no ORM needed).

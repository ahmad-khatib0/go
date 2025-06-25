# **Mockery & Testify/Mock Cheat Sheet**

_(Complete Guide to Using Generated Mocks in Go)_

Mockery generates mocks based on `testify/mock`, so understanding **Testify’s mock package**
is key. Below is a full breakdown of the most important methods, patterns, and how they
work together.

---

## **1. Mockery vs. Testify/Mock**

| **Component**      | **Role**                                                                       |
| ------------------ | ------------------------------------------------------------------------------ |
| **`mockery`**      | Code generator that creates mock structs from your interfaces.                 |
| **`testify/mock`** | Underlying library that provides mocking functionality (e.g., `On`, `Return`). |

---

## **2. Key Mock Methods (From `testify/mock`)**

### **(A) Setting Expectations**

| Method         | Purpose                                             | Example                                   |
| -------------- | --------------------------------------------------- | ----------------------------------------- |
| **`On()`**     | Defines a mock call expectation.                    | `mock.On("MethodName", arg1, arg2)`       |
| **`Return()`** | Sets return values for the mock.                    | `.Return(val1, err)`                      |
| **`Run()`**    | Executes a function when mock is called.            | `.Run(func(args mock.Arguments) { ... })` |
| **`Once()`**   | Expects the call **exactly once**.                  | `mock.On(...).Once()`                     |
| **`Times(n)`** | Expects the call **`n` times**.                     | `mock.On(...).Times(3)`                   |
| **`Maybe()`**  | Marks the call as optional (no strict requirement). | `mock.On(...).Maybe()`                    |

### **(B) Argument Matchers**

| Matcher                             | Purpose                                  | Example                                                                |
| ----------------------------------- | ---------------------------------------- | ---------------------------------------------------------------------- |
| **`mock.Anything`**                 | Matches any argument (type-safe).        | `mock.On("Method", mock.Anything)`                                     |
| **`mock.AnythingOfType("string")`** | Matches any argument of a specific type. | `mock.On("Method", mock.AnythingOfType("int"))`                        |
| **`mock.MatchedBy(func(x) bool)`**  | Custom matcher logic.                    | `mock.On("Method", mock.MatchedBy(func(x int) bool { return x > 0 }))` |
| **`mock.Eq(x)`**                    | Matches if `arg == x`.                   | `mock.On("Method", mock.Eq(42))`                                       |

### **(C) Assertions & Verification**

| Method                                    | Purpose                                | Example                                    |
| ----------------------------------------- | -------------------------------------- | ------------------------------------------ |
| **`AssertExpectations(t)`**               | Verifies all expected calls were made. | `mock.AssertExpectations(t)`               |
| **`AssertCalled(t, method, args...)`**    | Asserts a specific call happened.      | `mock.AssertCalled(t, "Method", 42)`       |
| **`AssertNotCalled(t, method, args...)`** | Asserts a call **did not** happen.     | `mock.AssertNotCalled(t, "Method", 99)`    |
| **`AssertNumberOfCalls(t, method, n)`**   | Asserts a method was called `n` times. | `mock.AssertNumberOfCalls(t, "Method", 2)` |

---

## **3. Full Example: Using a Generated Mock**

### **(1) Generated Mock (via Mockery)**

Given an interface:

```go
//go:generate mockery --name=UserRepository
type UserRepository interface {
    GetUser(id int) (*User, error)
    CreateUser(user *User) error
}
```

Mockery generates:

```go
// mock_UserRepository.go
type MockUserRepository struct {
    mock.Mock
}
func (m *MockUserRepository) GetUser(id int) (*User, error) {
    args := m.Called(id)
    return args.Get(0).(*User), args.Error(1)
}
func (m *MockUserRepository) CreateUser(user *User) error {
    args := m.Called(user)
    return args.Error(0)
}
```

### **(2) Using the Mock in Tests**

```go
func TestGetUser(t *testing.T) {
    // Setup
    repo := new(MockUserRepository)
    user := &User{ID: 1, Name: "Alice"}

    // Expectation: GetUser(1) -> user, nil
    repo.On("GetUser", 1).Return(user, nil)

    // Test
    result, err := repo.GetUser(1)
    assert.NoError(t, err)
    assert.Equal(t, user, result)

    // Verify expectations
    repo.AssertExpectations(t)
}
```

### **(3) Advanced Example (Argument Matchers)**

```go
func TestCreateUser(t *testing.T) {
    repo := new(MockUserRepository)

    // Expect CreateUser with any *User -> nil error
    repo.On("CreateUser", mock.AnythingOfType("*User")).Return(nil)

    // Test
    err := repo.CreateUser(&User{Name: "Bob"})
    assert.NoError(t, err)

    // Verify
    repo.AssertCalled(t, "CreateUser", mock.MatchedBy(func(u *User) bool {
        return u.Name == "Bob"
    }))
}
```

---

## **4. Common Patterns**

### **(A) Dynamic Return Values**

```go
repo.On("GetUser", mock.Anything).Return(func(id int) *User {
    return &User{ID: id}
}, nil)
```

### **(B) Run() for Side Effects**

```go
repo.On("CreateUser", mock.Anything).Run(func(args mock.Arguments) {
    user := args.Get(0).(*User)
    user.ID = 42 // Modify the user
}).Return(nil)
```

### **(C) Testing Errors**

```go
repo.On("GetUser", 99).Return(nil, errors.New("not found"))
```

### **(D) Mocking Chained Calls**

```go
dbMock := new(MockDatabase)
dbMock.On("Query").Return(new(MockQuery))
queryMock := new(MockQuery)
queryMock.On("Where", "id = ?", 1).Return(queryMock)
queryMock.On("First", mock.Anything).Run(func(args mock.Arguments) {
    result := args.Get(0).(*User)
    *result = User{ID: 1}
})
```

---

## **5. Key Takeaways**

1. **`On()` + `Return()`** define mock behavior.
2. **`mock.Anything`** is your friend for flexible argument matching.
3. **`AssertExpectations(t)`** ensures all expected calls happened.
4. **`Run()`** lets you execute custom logic during mocking.
5. **`Times(n)` / `Once()`** control call frequency expectations.

---

## **6. Troubleshooting**

| Issue                             | Fix                                                          |
| --------------------------------- | ------------------------------------------------------------ |
| **"Expected call not called"**    | Forgot `AssertExpectations(t)`?                              |
| **"No matching expected call"**   | Argument matchers (`mock.Anything`) too strict?              |
| **"Panic: interface conversion"** | Wrong `Return()` types? Use `args.Get(0).(YourType)` safely. |

---

### **Final Notes**

- Mockery just **generates** the mocks; `testify/mock` does the **heavy lifting**.
- Use **`go:generate mockery --name=YourInterface`** to update mocks automatically.
- Prefer **`mock.MatchedBy`** for complex argument validations.

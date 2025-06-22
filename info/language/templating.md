# Go Templating Engine: Comprehensive Cheat Sheet

## Table of Contents

1. [Basic Syntax](#basic-syntax)
2. [Variables](#variables)
3. [Control Structures](#control-structures)
4. [Functions](#functions)
5. [Pipelines](#pipelines)
6. [Template Composition](#template-composition)
7. [HTML Templates](#html-templates)
8. [Custom Functions](#custom-functions)
9. [Advanced Features](#advanced-features)
10. [Best Practices](#best-practices)

---

## Basic Syntax

### Template Delimiters

```go
{{ }}   // Standard action delimiter
{{- }}  // Trim preceding whitespace
{{ -}}  // Trim following whitespace
{{/* */}} // Comments
```

### Text vs Actions

```go
// Text is output as-is
Hello, World!

// Actions are processed
{{ "Hello, World!" }}
```

---

## Variables

### Definition

```go
{{ $name := "Alice" }}
```

### Usage

```go
{{ $name }}               // Output: Alice
{{ printf "Hello %s" $name }} // Output: Hello Alice
```

### Special Variables

```go
{{ . }}     // Current context (dot)
{{ .Field }} // Access struct field
{{ .Method }} // Call method
```

---

## Control Structures

### If-Else

```go
{{ if .Condition }}
  Show when true
{{ else if .OtherCondition }}
  Alternative
{{ else }}
  Default case
{{ end }}
```

### Range

```go
{{ range .Items }}
  {{ . }} // Current item
{{ end }}

// With index
{{ range $index, $item := .Items }}
  {{ $index }}: {{ $item }}
{{ end }}

// Empty case
{{ range .Items }}
  {{ . }}
{{ else }}
  No items found
{{ end }}
```

### With

```go
{{ with .User }}
  {{ .Name }} ({{ .Email }})
{{ else }}
  No user data
{{ end }}
```

---

## Functions

### Built-in Functions

```go
{{ eq .A .B }}        // Equality check
{{ lt .A .B }}        // Less than
{{ len .Slice }}      // Length
{{ index .Slice 0 }}  // Access element
{{ printf "%d" 42 }}  // Formatted output
{{ html "<b>" }}      // HTML escape
{{ js "alert" }}      // JavaScript escape
{{ urlquery "a=b" }}  // URL encoding
```

### String Functions

```go
{{ "hello" | title }}    // Hello
{{ "HELLO" | lower }}    // hello
{{ "hello" | upper }}    // HELLO
{{ "hello" | trunc 3 }}  // hel
{{ "a,b,c" | split "," }} // [a b c]
```

---

## Pipelines

### Basic Pipeline

```go
{{ .Name | upper | printf "%s!" }}
```

### Nested Pipelines

```go
{{ .Date | formatDate | trunc 10 }}
```

### Conditional Pipeline

```go
{{ if (gt (len .Items) 5 }}Many items{{ else }}Few items{{ end }}
```

---

## Template Composition

### Define Templates

```go
{{ define "header" }}
  <header>Site Header</header>
{{ end }}
```

### Execute Templates

```go
{{ template "header" }}
```

### Template Inheritance

```go
// base.html
{{ define "base" }}
<html>
  {{ template "content" . }}
</html>
{{ end }}

// page.html
{{ define "content" }}
  <body>Page Content</body>
{{ end }}

// Usage
tmpl.ExecuteTemplate(w, "base", data)
```

---

## HTML Templates

### Auto-escaping

```go
{{ .UserInput }} // Auto-escaped
{{ .SafeHTML | safe }} // Marked as safe
```

### URL Construction

```go
<a href="/user/{{ .ID }}">Profile</a>
```

### CSS/JS in Templates

```go
<style>
  .error { color: {{ .ErrorColor }}; }
</style>

<script>
  var config = {{ .Config | json }};
</script>
```

---

## Custom Functions

### Registering Functions

```go
funcMap := template.FuncMap{
    "formatDate": formatDate,
    "add": func(a, b int) int { return a + b },
}

t := template.New("").Funcs(funcMap)
```

### Example Function

```go
func formatDate(t time.Time) string {
    return t.Format("2006-01-02")
}

// In template
{{ .CreatedAt | formatDate }}
```

---

## Advanced Features

### Nested Templates

```go
{{ define "list" }}
  <ul>
    {{ range . }}
      {{ template "item" . }}
    {{ end }}
  </ul>
{{ end }}

{{ define "item" }}
  <li>{{ . }}</li>
{{ end }}
```

### Global Variables

```go
{{ $global := "value" }}
{{ define "partial" }}
  {{ $global }} // Accessible
{{ end }}
```

### Template Variables

```go
{{ $count := 0 }}
{{ range .Items }}
  {{ $count = add $count 1 }}
{{ end }}
Total: {{ $count }}
```

---

## Best Practices

1. **Separate Logic and Presentation**

   - Keep complex logic in Go code
   - Templates should focus on display

2. **Error Handling**

   ```go
   t := template.Must(template.ParseFiles("template.html"))
   ```

3. **Security**

   - Always escape user input
   - Use `html/template` instead of `text/template` for web

4. **Performance**

   - Parse templates once at startup
   - Use `ExecuteTemplate` for named templates

5. **Organization**

   - Group related templates in files
   - Use consistent naming conventions

6. **Testing**
   ```go
   func TestTemplate(t *testing.T) {
       tmpl := template.Must(template.ParseFiles("template.html"))
       var buf bytes.Buffer
       err := tmpl.Execute(&buf, testData)
       if err != nil {
           t.Fatal(err)
       }
       // Verify output
   }
   ```

---

## Complete Example

```go
package main

import (
    "os"
    "text/template"
    "time"
)

type User struct {
    Name    string
    Email   string
    Joined  time.Time
    Active  bool
    Friends []string
}

func main() {
    funcMap := template.FuncMap{
        "formatDate": func(t time.Time) string {
            return t.Format("2006-01-02")
        },
    }

    tmpl := template.Must(template.New("").Funcs(funcMap).Parse(`
        {{ define "user" }}
            User: {{ .Name }}
            Email: {{ .Email }}
            Joined: {{ .Joined | formatDate }}
            Status: {{ if .Active }}Active{{ else }}Inactive{{ end }}

            Friends:
            {{ range .Friends }}
                - {{ . }}
            {{ else }}
                No friends yet
            {{ end }}
        {{ end }}
    `))

    user := User{
        Name:    "Alice",
        Email:   "alice@example.com",
        Joined:  time.Now(),
        Active:  true,
        Friends: []string{"Bob", "Charlie"},
    }

    tmpl.ExecuteTemplate(os.Stdout, "user", user)
}
```

# Regex Cheat Sheet for Go (Golang) vs JavaScript

## Basic Regex Syntax

| Pattern  | Matches                          | Go | JS |
|----------|----------------------------------|----|----|
| `.`      | Any character except newline     | ✓  | ✓  |
| `\d`     | Digit (0-9)                      | ✓  | ✓  |
| `\D`     | Non-digit                        | ✓  | ✓  |
| `\w`     | Word character (a-z, A-Z, 0-9, _) | ✓  | ✓  |
| `\W`     | Non-word character               | ✓  | ✓  |
| `\s`     | Whitespace (space, tab, newline) | ✓  | ✓  |
| `\S`     | Non-whitespace                   | ✓  | ✓  |
| `\n`     | Newline                          | ✓  | ✓  |
| `\t`     | Tab                              | ✓  | ✓  |

## Quantifiers

| Pattern  | Meaning                          | Go | JS |
|----------|----------------------------------|----|----|
| `*`      | 0 or more                        | ✓  | ✓  |
| `+`      | 1 or more                        | ✓  | ✓  |
| `?`      | 0 or 1                           | ✓  | ✓  |
| `{n}`    | Exactly n times                  | ✓  | ✓  |
| `{n,}`   | n or more times                  | ✓  | ✓  |
| `{n,m}`  | Between n and m times            | ✓  | ✓  |

## Position Anchors

| Pattern  | Meaning                          | Go | JS |
|----------|----------------------------------|----|----|
| `^`      | Start of string/line             | ✓  | ✓  |
| `$`      | End of string/line               | ✓  | ✓  |
| `\A`     | Start of string                  | ✓  | ✗  |
| `\z`     | End of string                    | ✓  | ✗  |
| `\b`     | Word boundary                    | ✓  | ✓  |
| `\B`     | Not word boundary                | ✓  | ✓  |

## Character Classes

| Pattern      | Meaning                      | Go | JS |
|--------------|------------------------------|----|----|
| `[abc]`      | a, b, or c                   | ✓  | ✓  |
| `[^abc]`     | Not a, b, or c               | ✓  | ✓  |
| `[a-z]`      | Character between a to z     | ✓  | ✓  |
| `[[:alpha:]]`| Alphabetic character         | ✓  | ✗  |
| `[[:digit:]]`| Digit (same as \d)           | ✓  | ✗  |

## Groups and Capturing

| Pattern      | Meaning                      | Go | JS |
|--------------|------------------------------|----|----|
| `(expr)`     | Capturing group              | ✓  | ✓  |
| `(?:expr)`   | Non-capturing group          | ✓  | ✓  |
| `(?P<name>expr)` | Named capture group      | ✓  | ✓* |
| `\1`         | Backreference to group 1     | ✓  | ✓  |

*JavaScript uses `(?<name>expr)` syntax for named groups

## Flags/Modifiers

| Flag | Meaning                      | Go Syntax         | JS Syntax |
|------|------------------------------|-------------------|-----------|
| `i`  | Case insensitive             | `(?i)` in pattern | `/regex/i` |
| `m`  | Multi-line (^/$ match lines) | `(?m)` in pattern | `/regex/m` |
| `s`  | Dot matches newline          | `(?s)` in pattern | `/regex/s` |

## Key Differences

1. **Syntax**:
   - Go uses `regexp.MustCompile()` or `regexp.Compile()`
   - JavaScript uses `/regex/` literals or `new RegExp()`

2. **Named Groups**:
   - Go: `(?P<name>expr)`
   - JS: `(?<name>expr)`

3. **Flags**:
   - Go embeds flags in pattern: `(?i)caseinsensitive`
   - JS adds flags after regex: `/caseinsensitive/i`

4. **Unicode Support**:
   - Go has more extensive Unicode character classes
   - JavaScript's support varies by engine

5. **Lookaround Assertions**:
   - JavaScript supports lookaheads (`(?=...)`, `(?!...)`) and lookbehinds (`(?<=...)`, `(?<!...)`)
   - Go only supports lookaheads (as of Go 1.21)

## Common Go Regex Functions

```go
import "regexp"

// Compile regex (panics on error)
re := regexp.MustCompile(`pattern`) 

// Match string
matched := re.MatchString("input") 

// Find first match
match := re.FindString("input") 

// Find all matches
matches := re.FindAllString("input", -1) 

// Replace matches
result := re.ReplaceAllString("input", "replacement") 

// Split by regex
parts := re.Split("input", -1) 
```

## Common JavaScript Regex Methods

```javascript
const re = /pattern/flags; // or new RegExp('pattern', 'flags')

// Test for match
re.test('input'); 

// Find first match
'input'.match(re); 

// Find all matches
'input'.matchAll(re); 

// Replace matches
'input'.replace(re, 'replacement'); 

// Split by regex
'input'.split(re); 
```

Remember that Go's regex implementation is RE2, which is different from JavaScript's
(which typically uses PCRE-style regex). The main practical differences are that Go 
doesn't support backreferences in patterns and has more limited lookaround support.


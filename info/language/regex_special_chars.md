Complete the cheat sheet with all the special characters and clarify how they 
work in both Go and JavaScript regex.

## Special Character Cheat Sheet Addition

| Character | Meaning in Regex | Needs Escaping? | Go Example | JS Example |
|-----------|------------------|-----------------|------------|------------|
| `~`       | Literal tilde    | No              | `~`        | `~`        |
| `\`       | Escape character | Yes (`\\`)      | `\\`       | `\\`       |
| `[ ]`     | Character class  | Yes (`\[` `\]`) | `\[`       | `\[`       |
| `{ }`     | Quantifier       | Yes (`\{` `\}`) | `\{`       | `\{`       |
| `( )`     | Grouping         | Yes (`\(` `\)`) | `\(`       | `\(`       |
| `^`       | Anchor/negation  | Yes in `[]`     | `\^`       | `\^`       |
| `$`       | Anchor           | No              | `$`        | `$`        |
| `.`       | Any character    | Yes (`\.`)      | `\.`       | `\.`       |
| `*`       | Quantifier       | Yes (`\*`)      | `\*`       | `\*`       |
| `+`       | Quantifier       | Yes (`\+`)      | `\+`       | `\+`       |
| `?`       | Quantifier       | Yes (`\?`)      | `\?`       | `\?`       |
| `\|`      | Alternation      | Yes (`\\|`)     | `\\|`      | `\\|`      |

## Your Specific Example Explained

```go
regexp.MustCompile(`\B~[a-zA-Z0-9\-_]+`)
```

Breaking this down:
1. `\B` - Not a word boundary
2. `~` - Literal tilde character
3. `[a-zA-Z0-9\-_]` - Character class matching:
   - Letters a-z (case insensitive)
   - Numbers 0-9
   - Hyphen (escaped as `\-`)
   - Underscore
4. `+` - One or more of the preceding characters

## Complete Special Characters List

These characters have special meaning and need escaping with `\` when you want their literal meaning:
``` 
^ $ . * + ? ( ) [ ] { } \ | -
```

## When to Escape in Character Classes

Inside `[]` (character classes):
- Always escape: `]`, `\`, `^`, `-`
- Don't need to escape: `.`, `*`, `+`, `?`, `(`, `)`, `{`, `}`, `|`

Example:
```go
// Match either ~ or !
regexp.MustCompile(`[~!]`)  // Correct
regexp.MustCompile(`[\~\!]`) // Over-escaping (still works but unnecessary)
```

## Go vs JavaScript Differences in Character Escaping

1. **In string literals** (Go only):
   ```go
   // Go string literals need double escaping
   regexp.MustCompile("\\B~[a-zA-Z0-9\\-_]+")
   // Raw string literals (backticks) don't
   regexp.MustCompile(`\B~[a-zA-Z0-9\-_]+`)
   ```

2. **JavaScript**:
   ```javascript
   // In regex literals
   /\B~[a-zA-Z0-9\-_]+/
   // In RegExp constructor
   new RegExp("\\B~[a-zA-Z0-9\\-_]+")
   ```

---
---
---

Here’s a **comprehensive reference table** of all regex special characters and their 
meanings in Go and JavaScript, formatted for easy copying/saving. I’ve included every 
character with special regex meaning, along with escape rules and examples:

---

### **Complete Regex Special Characters Reference**
| Char | Meaning                  | Needs Escape? | Go Example            | JS Example           | Notes                               |
|------|--------------------------|---------------|-----------------------|----------------------|-------------------------------------|
| `\`  | Escape char              | `\\`          | `\\d`                 | `\\d`                | Required for `\d`, `\w`, etc.       |
| `^`  | Start of string/line     | `\^` in `[]`  | `^Start`              | `/^Start/`           | Negation inside `[^a]`              |
| `$`  | End of string/line       | No            | `End$`                | `/End$/`             |                                     |
| `.`  | Any char (except \n)     | `\.`          | `\.com`               | `/\.com/`            | Literal dot needs escaping          |
| `*`  | 0+ repetitions           | `\*`          | `a\*b` (literal `a*b`) | `/a\*b/`             | Quantifier otherwise                |
| `+`  | 1+ repetitions           | `\+`          | `1\+1=2`              | `/1\+1=2/`           | Quantifier otherwise                |
| `?`  | 0 or 1 (optional)        | `\?`          | `colou\?r`            | `/colou\?r/`         | Quantifier otherwise                |
| `\|` | OR (alternation)         | `\\|`         | `cat\\|dog`           | `/cat\|dog/`         |                                     |
| `()` | Grouping/capture         | `\( \)`       | `\(abc\)`             | `/\(abc\)/`          | Literal parentheses                 |
| `[]` | Character class          | `\[ \]`       | `\[abc\]`             | `/\[abc\]/`          | Literal brackets                    |
| `{}` | Quantifier bounds        | `\{ \}`       | `a\{2}`               | `/a\{2}/`            | `{n,m}` quantifier otherwise        |
| `-`  | Range in `[]`            | `\-`          | `[a\-z]` (literal `-`) | `/[a\-z]/`           | Unescaped: `[a-z]` (range)          |
| `~`  | Literal tilde            | No            | `~file`               | `/~file/`            |                                     |
| `!`  | Literal exclamation      | No            | `!important`          | `/!important/`       |                                     |
| `@`  | Literal at-sign          | No            | `@domain`             | `/@domain/`          |                                     |
| `#`  | Literal hash             | No            | `#tag`                | `/#tag/`             |                                     |
| `%`  | Literal percent          | No            | `100%`                | `/100%/`             |                                     |
| `&`  | Literal ampersand        | No            | `A&B`                 | `/A&B/`              |                                     |
| `=`  | Literal equals           | No            | `x=1`                 | `/x=1/`              |                                     |
| `:`  | Literal colon            | No            | `https:`              | `/https:/`           |                                     |
| `;`  | Literal semicolon        | No            | `alert();`            | `/alert();/`         |                                     |
| `'`  | Literal single quote     | No            | `'text'`              | `/'text'/`           |                                     |
| `"`  | Literal double quote     | No            | `"text"`              | `/"text"/`           |                                     |
| `/`  | Regex delimiter (JS)     | `\/`          | N/A (Go uses `"`/`` ` ``) | `/path\/file/`   | Only in JS regex literals           |

---

### **Key Rules for Escaping**
1. **Always escape these when literal**:  
   `^ $ . * + ? ( ) [ ] { } \ | -`

2. **Inside character classes `[]`**:
   - Must escape: `]`, `\`, `^`, `-`
   - Example:  
     - `[\^a-z]` → Matches `^` or lowercase letters  
     - `[a\-z]` → Matches `a`, `-`, or `z`

3. **Go-specific**:
   - Use raw strings (backticks) to avoid double-escaping:  
     ```go
     // Good (raw string)
     regexp.MustCompile(`\d+\.\d+`) 
     // Avoid (double escapes)
     regexp.MustCompile("\\d+\\.\\d+")
     ```

4. **JavaScript-specific**:
   - In regex literals (`/.../`), escape `/` as `\/`:  
     ```javascript
     /path\/file/.test("path/file") // true
     ```

---

### **Full ASCII Chart (Literal vs. Special)**
| Chars       | Treat as Literal | Treat as Special (Need Escape) |
|-------------|------------------|--------------------------------|
| `a-z` `A-Z` `0-9` | ✓          |                                |
| `_` ` ` (space)   | ✓          |                                |
| `~!@#$%&_=:;'"`   | ✓          |                                |
| `^$.*+?()[]{}|\-` |            | ✓                              |

---

### **Pro Tip**
When in doubt, escape any non-alphanumeric character. Modern regex engines will still 
work correctly (e.g., `\@` and `@` behave the same, but escaping makes your intent clear).



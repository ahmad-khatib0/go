## Backreferences

**What they are**: Backreferences allow you to match the same text that was matched 
by a capturing group earlier in the regex.

**Example**: Match repeated words
```regex
\b(\w+)\b\s+\1\b
```
- `(\w+)` captures a word into group 1
- `\1` matches the exact same text as group 1
- Matches: "the the", "cat cat", but not "the cat"

**Language Support**:
- ✅ JavaScript supports backreferences
- ❌ Go does NOT support backreferences in patterns

## Lookaround (Lookahead/Lookbehind)

These are "zero-width assertions" - they check for a pattern but don't consume characters.

### 1. Lookahead (`(?=...)` and `(?!...)`)

**Positive Lookahead (`?=`)**:
Matches a group AFTER the main expression without including it in the result.

**Example**: Find numbers followed by "px"
```regex
\d+(?=px)
```
- Matches "10" in "10px" but not "10pt"

**Negative Lookahead (`?!`)**:
Ensures a pattern does NOT follow.

**Example**: Find numbers NOT followed by "px"
```regex
\d+(?!px)
```
- Matches "10" in "10pt" but not "10px"

### 2. Lookbehind (`(?<=...)` and `(?<!...)`)

**Positive Lookbehind (`?<=`)**:
Matches a group BEFORE the main expression.

**Example**: Find numbers preceded by "$"
```regex
(?<=\$)\d+
```
- Matches "100" in "$100"

**Negative Lookbehind (`?<!`)**:
Ensures a pattern does NOT precede.

**Example**: Find numbers NOT preceded by "$"
```regex
(?<!\$)\d+
```
- Matches "100" in "€100" but not "$100"

**Language Support**:
- ✅ JavaScript supports both lookahead and lookbehind
- ✅ Go supports lookahead (since Go 1.0)
- ❌ Go does NOT support lookbehind (as of Go 1.21)

## Practical Examples

1. **Password Validation** (at least 1 digit, 1 letter, 8 chars):
   ```regex
   ^(?=.*\d)(?=.*[a-zA-Z]).{8,}$
   ```
   (Uses positive lookaheads)

2. **Find words not followed by "ing"**:
   ```regex
   \b\w+\b(?!ing\b)
   ```

3. **Match "USD" only when preceded by number**:
   ```regex
   (?<=\d)USD
   ```

The main difference is that Go's regex engine (RE2) prioritizes safety and performance
over advanced features, while JavaScript's regex is more feature-rich but can sometimes 
be slower or vulnerable to "regex denial of service" attacks with complex patterns.


---
---
---

## What is: Not word boundary, Non-word character, Named groups?

### 1. `\B` - Not Word Boundary
**What it is**: The opposite of `\b`. Matches positions where a word boundary does *not* exist.

**Examples**:
- `\Bcat\B` matches "cat" in "concatenate" but not "cat" at start/end of word
- `\d\B` matches "3" in "3D" (number followed by non-word char) but not "3" in "3 dogs"

**Visualization**:
```
"Hello-World" 
   ^   ^
   |   \B matches here (between 'l' and 'l')
   \b matches here (between 'o' and '-')
```

### 2. `\W` - Non-Word Character
**What it is**: Matches any character that is *not* a word character (opposite of `\w` = `[^a-zA-Z0-9_]`)

**Examples**:
- `\W+` matches " @$ " in "hello@$ world"
- `\w\W\w` matches "o W" in "Hello World"

**Common non-word chars**: Spaces, punctuation, symbols like @, #, $, etc.

### 3. Named Groups (`(?P<name>...)` in Go, `(?<name>...)` in JS)
**What it is**: Captures a group and gives it a name for later reference.

**Go Example**:
```go
re := regexp.MustCompile(`(?P<year>\d{4})-(?P<month>\d{2})`)
match := re.FindStringSubmatch("2023-10")
fmt.Println(match["year"])  // "2023"
fmt.Println(match["month"]) // "10"
```

**JavaScript Example**:
```javascript
const re = /(?<year>\d{4})-(?<month>\d{2})/;
const match = "2023-10".match(re);
console.log(match.groups.year);  // "2023"
console.log(match.groups.month); // "10"
```

### 4. `?` - Three Different Meanings
The question mark has multiple uses in regex:

| Usage       | Meaning | Example |
|-------------|---------|---------|
| `?` (quantifier) | Makes preceding item optional (0 or 1) | `colou?r` matches "color" or "colour" |
| `??` | Lazy/lazy quantifier (matches as few as possible) | `a.*?b` in "aabab" matches "aab" not "aabab" |
| `(?...)` | Special group syntax | `(?:...)` non-capturing group, `(?=...)` lookahead |

**Key Difference**:
- When used **after a character/group**: Quantifier (`a?`)
- When used **after `(`**: Special group syntax (`(?i)` for flags)


# Goldmark Subscript Extension

<!-- markdownlint-disable MD033 -->

A [Goldmark](https://github.com/yuin/goldmark) extension that adds subscript support using single-tilde syntax (`~text~`). This extension allows you to render subscripts in your Markdown documents while maintaining full compatibility with Goldmark's built-in strikethrough extension.

## Installation

```bash
go get github.com/zmtcreative/gm-subscript
```

## Configuration

### Basic Usage

```go
package main

import (
    "bytes"
    "fmt"

    "github.com/yuin/goldmark"
    "github.com/zmtcreative/gm-subscript"
)

func main() {
    md := goldmark.New(
        goldmark.WithExtensions(
            subscript.Subscript, // Use the pre-configured instance
        ),
    )

    var buf bytes.Buffer
    if err := md.Convert([]byte("H~2~O"), &buf); err != nil {
        panic(err)
    }
    fmt.Print(buf.String()) // Output: <p>H<sub>2</sub>O</p>
}
```

### Alternative Configuration

```go
md := goldmark.New(
    goldmark.WithExtensions(
        subscript.NewSubscript(), // Create a new instance
    ),
)
```

### With Other Extensions

This extension works seamlessly with other Goldmark extensions, including the built-in strikethrough:

```go
import (
    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/extension"
    "github.com/zmtcreative/gm-subscript"
)

md := goldmark.New(
    goldmark.WithExtensions(
        extension.GFM,              // Includes strikethrough
        extension.DefinitionList,
        extension.Footnote,
        subscript.Subscript,        // Add subscript support
    ),
)
```

## Examples

### Basic Chemical Formulas

```markdown
Input:  H~2~O
Output: H<sub>2</sub>O (water)

Input:  C~6~H~12~O~6~
Output: C<sub>6</sub>H<sub>12</sub>O<sub>6</sub> (glucose)
```

### Mathematical Expressions

```markdown
Input:  x~1~, x~2~, x~n~
Output: x<sub>1</sub>, x<sub>2</sub>, x<sub>n</sub>

Input:  log~10~(100) = 2
Output: log<sub>10</sub>(100) = 2
```

### Combined with Strikethrough

```markdown
Input:  H~2~O is ~~not~~ essential for life
Output: H<sub>2</sub>O is <del>not</del> essential for life

// subscripts cannot be preceeded by whitespace -- they must be part of other text
Input:  NH~4~ with ~subscript~ and ~~strikethrough~~
Output: NH<sub>4</sub> with <del>subscript</del> and <del>strikethrough</del>
```

### Special Characters and Unicode

```markdown
Input:  R~*~ x f~p~ x n~e~
Output: R<sub>*</sub> x f<sub>p</sub> x n<sub>e</sub>

Input:  Text~αβγ123~end
Output: Text<sub>αβγ123</sub>end
```

## Compatibility

This extension is designed to work alongside Goldmark's built-in `extension.Strikethrough` and other extensions that use the `~` character. The parsing rules ensure proper disambiguation:

### Syntax Rules

1. **No whitespace inside subscripts**: Content between tildes cannot contain spaces
   - ✅ `H~2~O` → H<sub>2</sub>O
   - ❌ `H~2 ~O` → H~2 ~O (not parsed as subscript)

2. **Must be preceded by non-whitespace**: Subscripts cannot start at the beginning of a line or after whitespace
   - ✅ `H~2~O` → H<sub>2</sub>O
   - ❌ `~2~O` → <del>2</del>O (parsed as strikethrough)
   - ❌ `H ~2~O` → H <del>2</del>O (parsed as strikethrough)

3. **Single tildes only**: Double tildes are reserved for strikethrough
   - ✅ `~text~` → subscript (when rules 1-2 are met)
   - ✅ `~~text~~` → <del>text</del> (strikethrough)

4. **No nested markdown**: Other markdown syntax is not processed inside subscripts
   - ✅ `Text~<em>word</em>~` → Text<sub>&lt;em&gt;word&lt;/em&gt;</sub>
   - For complex formatting, use HTML directly: `Text<sub><em>word</em></sub>`

5. **No nested tildes**: Content cannot contain additional tilde characters
   - ✅ `H~2~O` → H<sub>2</sub>O
   - ❌ `H~2~O~` → H<sub>2</sub>O~ (only first pair is parsed)

### Strikethrough Coexistence

The extension tries to intelligently handle the shared use of the `~` character:

- **Single tildes** (following the rules above) → subscripts
- **Double tildes** → strikethrough
- **Single tildes with spaces** → strikethrough
- **Single tildes at line start or after whitespace** → strikethrough

### Limitations

- **Simple content only**: Subscripts are best suited for simple text, numbers, and basic symbols
- **No complex formatting**: For complex subscripts with multiple formatting options, use HTML `<sub>` tags directly
- **Mathematics**: For complex mathematical expressions and equations, consider using **`KaTeX`** or **`MathJax`** instead

## Use Cases

This extension is ideal for:

- Scientific and chemical formulas
- Simple Mathematical notation
- Reference numbering
- Simple technical documentation

For more complex scenarios requiring nested formatting or advanced mathematical notation, consider using:

- Direct HTML `<sub>` tags for complex formatting
- [KaTeX](https://katex.org/) or [MathJax](https://www.mathjax.org/) for advanced mathematics

## License

This project is licensed under the MIT License. See the [LICENSE.md](LICENSE.md) file for details.

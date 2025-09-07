# Goldmark Subscript Extension

<!-- markdownlint-disable MD033 -->

[![Go Reference](https://pkg.go.dev/badge/github.com/zmtcreative/gm-subscript.svg)](https://pkg.go.dev/github.com/zmtcreative/gm-subscript)
[![Go version](https://img.shields.io/github/go-mod/go-version/zmtcreative/gm-subscript)](https://github.com/zmtcreative/gm-subscript)
[![License](https://img.shields.io/github/license/zmtcreative/gm-subscript)](./LICENSE.md)
![GitHub Tag](https://img.shields.io/github/v/tag/zmtcreative/gm-subscript?include_prereleases&sort=semver)

A [Goldmark](https://github.com/yuin/goldmark) extension that adds subscript support using single-tilde syntax (`H~2~O`). This extension allows you to render subscripts in your Markdown documents while maintaining full compatibility with Goldmark's built-in strikethrough extension.

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

This extension works with other Goldmark extensions, including the built-in strikethrough:

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

### Syntax Rules

> [!TIP]
>
> **Use double-tilde for strikethrough everywhere** &nbsp;&mdash;&nbsp;&nbsp;To reduce ambiguity, try to always use
> double-tilde delimiters for strikethrough. This reduces the ambiguity when subscript is enabled. It's not a perfect
> solution &mdash; the conflicting syntax is frustrating, but since no one entity is truly setting a definitive standard
> for markdown, this is what we have.

1. **No whitespace inside subscripts**: Content between tildes cannot contain spaces (*strikethrough **cannot** have leading or trailing spaces between tildes either, but it **can** have spaces between words*)
   - ✅ `H~2~O` → H<sub>2</sub>O
   - ❌ `H~2 ab~O` → H<del>2 ab</del>O (*not parsed as subscript but **will** be parsed as strikethrough*)
   - ❌ `H~2 ~O` → H~2 ~O (*not parsed as subscript OR strikethrough*)
   - ❌ `H~ 2~O` → H~ 2~O (*not parsed as subscript OR strikethrough*)
   - ❌ `H~ 2 ~O` → H~ 2~O (*not parsed as subscript OR strikethrough*)

2. **Must be preceded by non-whitespace**: Subscripts cannot start at the beginning of a line or after whitespace
   - ✅ `H~2~O` → H<sub>2</sub>O
   - ❌ `~2~O` → <del>2</del>O (*parsed as strikethrough -- tilde at beginning of line*)
   - ❌ `H ~2~O` → H <del>2</del>O (*parsed as strikethrough -- space before opening tilde*)

3. **Single tildes only**: Double tildes are reserved for strikethrough
   - ✅ `body~text~` → body<sub>text</sub> subscript (*when rules 1-2 are met*)
   - ✅ `body~~text~~` → body<del>text</del> (*strikethrough*)

4. **No nested markdown**: Other markdown syntax is not processed inside subscripts
   - ✅ `Text~<em>word</em>~` → Text<sub>&lt;em&gt;word&lt;/em&gt;</sub>
   - For complex formatting, use HTML directly: `Text<sub><em>word</em></sub>`
   - For **really** complex formatting and output use `KaTex` or `Mathjax` (*or similar LaTeX rendering*)

5. **No nested tildes**: Content cannot contain tilde characters (*opening and closing tilde are consumed during parsing*)
   - ✅ `H~2~O` → H<sub>2</sub>O
   - ✅ `H~2~~O~` → H<sub>2</sub><sub>O</sub>  &nbsp;&mdash;&nbsp;  `~2~` and `~O~` are parsed as separate subscripts!
   - ❌ `H~2~O~` → H<sub>2</sub>O~  &nbsp;&mdash;&nbsp;  only the tildes around the `2` are parsed &mdash; the last tilde is just a plain text character
   - ❌ `~H~2~O~` → <del>H<sub>2</sub>O</del>  &nbsp;&mdash;&nbsp;  here is how it's parsed:
     - The first tilde is at the beginning of the line (*which is not allowed*), so the parser skips it, and
     - Moves on to the next valid subscript `2`, which **is** rendered as a subscript (*and the two tilde delimiters are consumed*)
     - Leaving `~H<sub>2</sub>O~` for the strikethrough parser to process, so the entire H<sub>2</sub>O is struck through

> [!NOTE]
>
> The subscript parser has higher priority than the strikethrough parser, so it can find and consume valid subscript
> markdown before the strikethrough parser makes its pass. The subscript parser is a single-pass parser, so there is no
> backtracking to look for nested tildes. Both parsers consume their tildes. This means that whatever is left over
> after the subscript parses the markdown is what the strikethrough parser sees on its pass. This can sometimes lead to
> unexpected rendering results, but is the best we can do given the shared syntax.

## Examples

### Basic Chemical Formulas

```markdown
Input:  H~2~O
Output: H<sub>2</sub>O (water)
```

&nbsp;&nbsp;&nbsp;H<sub>2</sub>O (water)

```markdown
Input:  C~6~H~12~O~6~
Output: C<sub>6</sub>H<sub>12</sub>O<sub>6</sub> (glucose)
```

&nbsp;&nbsp;&nbsp;C<sub>6</sub>H<sub>12</sub>O<sub>6</sub> (glucose)

### Mathematical Expressions

```markdown
Input:  x~1~, x~(2)~, x~n+1~
Output: x<sub>1</sub>, x<sub>(2)</sub>, x<sub>n+1</sub>
```

&nbsp;&nbsp;&nbsp;x<sub>1</sub>, x<sub>(2)</sub>, x<sub>n+1</sub>

```markdown
Input:  log~10~(100) = 2
Output: log<sub>10</sub>(100) = 2
```

&nbsp;&nbsp;&nbsp;log<sub>10</sub>(100) = 2

### Combined with Strikethrough

```markdown
Input:  H~2~O is ~~not~~ essential for life
Output:
```

&nbsp;&nbsp;&nbsp;H<sub>2</sub>O is <del>not</del> essential for life

```markdown
// subscripts cannot be preceded by whitespace -- they must be part of other text
Input:  NH~4~ with ~subscript~ and ~~strikethrough~~
Output: NH<sub>4</sub> with <del>subscript</del> and <del>strikethrough</del>
```

&nbsp;&nbsp;&nbsp;NH<sub>4</sub> with <del>subscript</del> and <del>strikethrough</del>

### Special Characters, HTML Entities and Unicode

```markdown
Input:  N~&#x1f47d;~ = R~&#x1f7af;~ × _f_~p~ × n~e~ × _f_~l~ × _f_~i~ × _f_~c~ × L
Output: N<sub>👽</sub> = R<sub>🞯</sub> × <em>f</em><sub>p</sub> × n<sub>e</sub> × <em>f</em><sub>l</sub> × <em>f</em><sub>i</sub> × <em>f</em><sub>c</sub> × L
```

&nbsp;&nbsp;&nbsp;N<sub>👽</sub> = R<sub>🞯</sub> × <em>f</em><sub>p</sub> × n<sub>e</sub> × <em>f</em><sub>l</sub> × <em>f</em><sub>i</sub> × <em>f</em><sub>c</sub> × L

```markdown
Input:  Text~αβγ123~end
Output: Text<sub>αβγ123</sub>end
```

&nbsp;&nbsp;&nbsp;Text<sub>αβγ123</sub>end

> [!NOTE]
>
> Unicode and HTML Entities **should** render properly on most modern browsers, but a user's font selection **might**
> result in some Unicode and HTML Entity output rendering as unknown characters (*e.g., &#xFFFD;*).

## Compatibility

This extension is designed to work alongside Goldmark's built-in `extension.Strikethrough` extension. This requires some strict parsing rules regarding whitespace and embedding other markdown or HTML inside subscripts.

The parsing rules try to ensure proper disambiguation:

### Strikethrough Coexistence

The extension tries to intelligently handle the shared use of the `~` character:

- **Single tildes** (*following the rules above*) → subscripts
- **Double tildes** → strikethrough
- **Single tildes with spaces between words** → strikethrough
- **Single tildes at line start or after whitespace** → strikethrough
- **Leading and/or trailing spaces on content text between tildes** → plain text

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

# Notex

A lightweight markup language that compiles to HTML. Notex uses a brace-based block structure and an extensible function system where functions are plain executable scripts.

## Build

```sh
go build -o notex
```

Set the `NOTEX_PATH` environment variable to a colon-separated list of directories containing your function scripts:

```sh
export NOTEX_PATH="/path/to/NotexFuncs"
```

## Usage

```sh
notex input.ntx              # prints HTML to stdout
notex input.ntx output.html  # writes to file
cat input.ntx | notex        # reads from stdin
```

## Syntax

### Document header

An optional `[...]` block at the start of a file sets document metadata:

```
[
#Page Title
A short description of the page.
]
```

- `#text` — sets the `<title>`
- Plain lines — set the `<meta description>`
- `<raw>` — injected verbatim into `<head>`

### Groups

Curly braces create block-level containers (`<div>`):

```
{
    Some content here.
    {
        #heading
        Nested content.
    }
}
```

### Headings

A line beginning with `#` inside a group becomes a heading. The heading level is determined by nesting depth.

### Paragraphs

Plain text lines become `<p>` tags. Prefix a line with `\` to force paragraph treatment and strip the backslash.

### Lists

- Lines starting with `-` create unordered list items (`<ul><li>`).
- Lines starting with `.` create ordered list items (`<ol><li>`).

A list item followed by a group uses that group as the `<li>` content:

```
{
    - {
        - nested item one
        - nested item two
    }
}
```

### Raw HTML

Lines beginning with `<` are passed through as raw HTML.

### Comments

`%` starts a comment to the end of the line. Use `\%` to escape.

### Functions

Functions are called with `/@name(params){object}`:

```
/@bold(){some text}
/@link(https://example.com){click here}
/@i(){/@b(){nested functions}}
```

Functions are external executables discovered via `NOTEX_PATH`. Each script must accept `--supplies` as its first argument and print a space-separated list of the function names it handles.

### Attributes

Attributes apply a CSS class to a group and inject the corresponding CSS into `<head>`. They are written on their own line inside a group:

```
{
    !@dark()
    Content styled with the dark theme.
}
```

Attribute names are prefixed with `!` in the `--supplies` output. When called with `--style` as the last argument, the script returns the CSS for that class.

## Writing Function Scripts

A function script receives arguments as: `name param1 param2 ... object`.

```bash
#!/bin/bash
if [[ "$1" == "--supplies" ]]; then
    echo "myfunc"
else
    case "$1" in
        myfunc) echo "<span>${!#}</span>" ;;
    esac
fi
```

To provide an attribute, prefix its name with `!` in `--supplies` and handle the `--style` flag:

```bash
#!/bin/bash
if [[ "$1" == "--supplies" ]]; then
    echo "!mytheme"
else
    if [[ "$1" == "!mytheme" ]]; then
        if [[ "${!#}" == "--style" ]]; then
            echo ".mytheme { color: red; }"
        else
            echo "mytheme"
        fi
    fi
fi
```

## Included Scripts

### `NotexFuncs/core`

| Function | Output |
|---|---|
| `b`, `bold` | `<b>` |
| `i`, `italics` | `<i>` |
| `ul`, `u` | `<u>` |
| `a`, `link` | `<a href="param">` |
| `img` | `<img src="...">` |
| `!dark` | Dark theme attribute |
| `!light` | Light theme attribute |

### `NotexFuncs/now`

`now` — inserts the current date/time via `date`.

### `NotexFuncs/green-vibes`

`!greenVibes` — a dark green terminal-style theme.

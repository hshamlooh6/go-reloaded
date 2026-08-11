# go-reloaded

A small command-line tool written in Go that cleans up and reformats text files. It takes an input file and an output file, applies a series of transformations, and writes the result.

## What it does

- Converts `(hex)` and `(bin)` markers into their decimal equivalents
- Applies `(up)`, `(low)`, and `(cap)` formatting tags, including counted forms like `(up, 2)`
- Turns `a` into `an` when the following word starts with a vowel or `h`
- Fixes spacing around punctuation (`.`, `,`, `!`, `?`, `:`, `;`), while keeping groups like `...` or `!?` together
- Cleans up spacing around single-quote pairs

## Usage

```bash
go run . input.txt output.txt
```

The program reads `input.txt`, applies all transformations, and writes the result to `output.txt`.

## Project structure

| File | Responsibility |
|---|---|
| `main.go` | Entry point, file I/O |
| `transform.go` | Orchestrates all transformations |
| `hexBin.go` | Hex/bin to decimal conversion |
| `article.go` | `a` → `an` logic |
| `case.go` | `(up)`, `(low)`, `(cap)` handling |
| `punctuations.go` | Punctuation spacing |
| `quotes.go` | Quote mark spacing |

## Requirements

Standard library only — no external dependencies.

## Notes

This project was built as part of the Reboot01 curriculum. Test files are recommended for validating behavior before submitting to auditors.

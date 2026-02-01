# 🦜 parrot

A Go linter that detects comments which merely parrot what the code already says because we are Gophers, not parrots!

> "Comments should explain *why*, not *what*."
> "A common fallacy is to assume authors of incomprehensible code will be able to express themselves clearly in comments." – Kevlin Henney
> "When you feel the need to write a comment, first try to refactor the code so that any comment becomes superfluous." - Martin Fowler

Also https://www.youtube.com/watch?v=0lKjFLYkXTE

## Examples

```go
// ❌ Bad: parrots the code
// connect to database
db := connectToDatabase()

// ❌ Bad: restates the obvious
// if error is not nil, return
if err != nil {
    return err
}

// ✅ Good: explains WHY
// Retry connection because k8s networking can be flaky during pod startup
db := connectToDatabase()
```

## Installation

```bash
go install github.com/oliveiraethales/parrot/cmd/parrot@latest
```

## Usage

### Standalone

```bash
parrot ./...
```

### With golangci-lint

Add to `.golangci.yml`:

```yaml
linters-settings:
  custom:
    parrot:
      path: /path/to/parrot.so
      description: Detects comments that parrot the code
      original-url: github.com/oliveiraethales/parrot

linters:
  enable:
    - parrot
```

Build the plugin:

```bash
go build -buildmode=plugin -o parrot.so ./plugin
```

## How It Works

Parrot uses heuristic analysis:

1. Extracts identifiers from code (function names, variables, etc.)
2. Tokenizes adjacent comments
3. Measures word overlap between comment and code
4. Flags comments where >60% of meaningful words match code identifiers

## Configuration

Currently uses a fixed 60% overlap threshold.

## License

MIT

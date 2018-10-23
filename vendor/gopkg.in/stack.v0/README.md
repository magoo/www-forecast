[![GoDoc](https://godoc.org/gopkg.in/stack.v0
?status.svg)](https://godoc.org/gopkg.in/stack.v0
) [![Build Status](https://travis-ci.org/go-stack/stack.svg?branch=master)](https://travis-ci.org/go-stack/stack)

# stack

Package stack implements utilities to capture, manipulate, and format call stacks. It provides a simpler API than package runtime.

The implementation takes care of the minutia and special cases of interpreting the program counter (pc) values returned by runtime.Callers.

## Versioning
Package stack publishes stable APIs via gopkg.in. The most recent is v0 (the API is still baking, please file an issue if you have any suggestions), which is imported like so:

    import "gopkg.in/stack.v0"

## Formatting
Package stack's types implement fmt.Formatter, which provides a simple and flexible way to declaratively configure formatting when used with logging or error tracking packages.

```go
func DoTheThing() {
    c := stack.Caller(0)
    log.Print(c)          // might log "source.go:10"
    log.Printf("%+v", c)  // might log "pkg/path/source.go:10"
    log.Printf("%n", c)   // might log "DoTheThing"

    s := stack.Trace().TrimRuntime()
    log.Print(s)          // might log "[source.go:15 caller.go:42 main.go:14]"
}
```

See the docs for all of the supported formatting options.

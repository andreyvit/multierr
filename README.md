# multierr

[![Go Reference](https://pkg.go.dev/badge/github.com/andreyvit/multierr.svg)](https://pkg.go.dev/github.com/andreyvit/multierr)

Provides an idiomatic Append function for errors, allowing to combine multiple errors into one.

Basically the whole API is:

```go
err = multierr.Append(err, someFuncThatFails())
err = multierr.Append(err, anotherFuncThatFails())
```

Printing the error results in something like:

```text
2 errors occurred:
(1) some failure
(2) another failure
```

Among other cases, this is especially useful in failable deferred funcs
to avoid silently ignoring cleanup errors:

```go
func sprinkleMagicDust() (err error) {
    d := &Dummy{}
    defer func() {
        err = multierr.Append(err, d.Close())
    }()

return d.DoSomeMagic()
}
```

Unlike other overly complicated multierror packages, this one does not even
expose its multierror type, and will only use it when you actually end up with
more than a single error to return.

The returned multierror type supports errors.Is and errors.As by delegating
to each of the suberrors it contains.


## License

© 2020-2021, Andrey Tarantsov. Published under the MIT License.

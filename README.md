# serrors
simple error handling for golang

## usage
### wrap errors
```golang
func otherMethod() error {
	return errors.New("unimplemented")
}

func someMethod() error {
	// wrap error for displaying stack trace
	return serrors.Wrap(otherMethod())
}

func main() {
	if err := someMethod(); err != nil {
		// show stack trace by Printf
		log.Fatalf("%+v", err)
	}
}
```

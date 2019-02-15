The deriveToError function is useful for composing bool-returning functions along with traditional error-returning functions.

Given the following input pseudo-code:

```go
expectKey := func(k some_key) (some_value, bool) {
    value, keyExists = some_map[k]
    return value, keyExists
}
transformed := deriveToError(fmt.Errorf("eFalse"), expectKey)
```

goderive will generate the following code:

```go
func deriveToError(err error, f func(i some_key) (a some_value, b bool)) func(i some_key) (some_value, error) {
	return func(i some_key) (some_value, error) {
		out, success := f(i)
		if success {
			return out, nil
		}
		return out, err
	}
}
```
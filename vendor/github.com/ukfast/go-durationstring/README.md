# durationstring

A simple Go package for working with string format durations, e.g. `1d4h`

## Usage

#### Parsing

```go
y, mo, d, h, m, s, ms, us, ns, err :=  durationstring.Parse("1d4h5ns")

assert.Equal(t, 1, d)
assert.Equal(t, 4, h)
assert.Equal(t, 5, ns)
```

#### String formatting

```go
s := String(1, 0, 0, 4, 0, 0, 0, 0, 0)

assert.Equal(t, "1y4h", s)
```

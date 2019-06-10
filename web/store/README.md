> Forked from [Humble/locstor](https://github.com/go-humble/locstor) 
> Forked for add wasm support and *include* to [gas](https://github.com/gascore/gas) ecosystem

gas-web/store
=============

Version 0.3

locstor provides gopherjs and wasm bindings for the localStorage API. It allows you to
store and retrieve any arbitrary go data structure. locstor works great as a stand-alone package or in
combination with other [Humble](https://github.com/go-humble) packages.

locstor is written in pure go. It feels like go, follows go idioms when
possible, and compiles with the go tools.

Example Usage
-------------

### Accessing the localStorage API Directly

Use [`SetItem`](https://godoc.org/github.com/go-humble/locstor#SetItem) to
store an item in localStorage:

```go
if err := locstor.SetItem('foo', 'bar'); err != nil {
	// Handle err
}
```

Use [`GetItem`](https://godoc.org/github.com/go-humble/locstor#GetItem) to get
an item from localStorage:

```go
item, err := locstor.GetItem('foo')
if err != nil {
	// Handle err
}
fmt.Println(item)
// Output:
//   bar
```

Use [`Key`](https://godoc.org/github.com/go-humble/locstor#Key) to get the key
for a specific item:

```go
key, err := locstor.Key('bar')
if err != nil {
	// Handle err
}
fmt.Println(key)
// Output:
//   foo
```

Use [`RemoveItem`](https://godoc.org/github.com/go-humble/locstor#RemoveItem)
to remove an existing item from localStorage:

```go
if err := locstor.RemoveItem('foo'); err != nil {
	// Handle err
}
_, err := locstor.GetItem('foo')
fmt.Println(err)
// Output:
//   Could not find an item with the given key: foo
```

Use [`Length`](https://godoc.org/github.com/go-humble/locstor#Length) to get
the number of items currently in localStorage:

```go
count, err := locstor.Length()
if err != nil {
	// Handle err
}
```

Use [`Clear`](https://godoc.org/github.com/go-humble/locstor#Clear) to remove
all items from localStorage:

```go
if err := locstor.Clear(); err != nil {
	// Handle err
}
```

### Using a DataStore

You can also use a
[`DataStore`](https://godoc.org/github.com/go-humble/locstor#DataStore), which
is an abstraction layer built on top of localStorage capable of storing and
retrieving arbitrary go data structures, not just strings.

Use
[`NewDataStore`](https://godoc.org/github.com/go-humble/locstor#NewDataStore) to
create a new `DataStore`. It accepts an `EncoderDecoder` as an argument. There
are two encodings provided out-of-the-box: `JSONEncoding` and `BinaryEncoding`.
You should choose `JSONEncoding` if you want the data stored in localStorage to
be more readable and `BinaryEncoding` if you want the data to take up less
space. You can also provide a custom encoding by implementing the
[`EncoderDecoder`](https://godoc.org/github.com/go-humble/locstor#EncoderDecoder)
interface.

```go
store := locstor.NewDataStore(JSONEncoding)
```

Use [`Save`](https://godoc.org/github.com/go-humble/locstor#DataStore.Save) to
save data structures in localStorage:

```go
if err := store.Save("numbers", []int{1, 2, 3}); err != nil {
	// Handle err
}
```

Use [`Find`](https://godoc.org/github.com/go-humble/locstor#DataStore.Find) to
get existing data structures out of localStorage. Find works similarly to
[`json.Unmarshal`](http://golang.org/pkg/encoding/json/#Unmarshal) from the
standard library. The second argument to `Find`, called `holder`, is a pointer
to a variable that is capable of holding the decoded data structure. Since in
this case we stored a slice of ints, the type of holder should be `*[]int`.

```go
gotNumbers := []int{}
if err := store.Find("numbers", &gotNumbers); err != nil {
	// Handle err
}
fmt.Println(gotNumbers)
// Output:
//   [1 2 3]
```

Use [`Delete`](https://godoc.org/github.com/go-humble/locstor#DataStore.Delete)
to delete an existing data structure from localStorage:

```go
if err := store.Delete("numbers"); err != nil {
	// Handle err
}
gotNumbers := []int{}
_, err := locstor.Find('numbers', &gotNumbers)
fmt.Println(err)
// Output:
//   Could not find an item with the given key: numbers
```

### Handling Errors

[`ErrLocalStorageNotSupported`](https://godoc.org/github.com/go-humble/locstor#pkg-variables)
will be returned by any function or method if localStorage is not supported in
the current browser. `ErrLocalStorageNotSupported` is just a variable, so
you can do direct comparisons:

```go
if err := locstor.GetItem("foo"); err != nil {
	if err == locstor.ErrLocalStorageNotSupported {
		// Handle an ErrLocalStorageNotSupported error 
	} else {
		// Handle some other type of error
	}
}
```

[`ItemNotFoundError`](https://godoc.org/github.com/go-humble/locstor#ItemNotFoundError)
is the type of error returned if the item you were looking for does not exist in
localStorage. It is a type that implements the `error` interface and tells you
which key or item was not found. To check if an error is an `ItemNotFoundError`,
you can use a type assertion:

```go
if err := locstor.GetItem("foo"); err != nil {
	if _, ok := err.(locstor.ItemNotFoundError); ok {
		// Handle an ItemNotFoundError error 
	} else {
		// Handle some other type of error
	}
}
```

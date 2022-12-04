# patch

This package aims to answer the question asked when implementing **PATCH** requests of RESTful API in Go:

<div align="center">
    <p style="font-weight: bold;">
        Q: How to tell if a field is missing in the payload of a PATCH request?
    </p>
</div>

Since we are using generics, **Go 1.8+ is required**.

---

Here we only talk about JSON payloads as it's the most frequently used format when developing a RESTful API.

```go
type UserPatch struct {
	Name   string
	Age    int
	Gender string
}

func PatchUser(rw http.ResponseWriter, r *http.Request) {
	var payload UserPatch
	json.NewDecoder(r.Body).Decode(&payload)

    // Both the requests `{"Name":"","Age":18}` and `{"Age":18}` can
    // cause `payload.Name == ""`.
    // ** How can we distinguish the two? **
    // "A field is missing" and "a field is empty" are semantically different.
	if payload.Name == "" {
		// do sth...
	}
}
```

## Solution: add a sentinel to each field

Using `patch.Field` to define/wrap your fields in a struct.

```go
import "github.com/ggicci/patch"

type UserPatch struct {
	Name   patch.Field[string]
	Age    patch.Field[int]
	Gender patch.Field[string]
}

func PatchUser(rw http.ResponseWriter, r *http.Request) {
	var payload UserPatch
	json.NewDecoder(r.Body).Decode(&payload)

	if !payload.Name.Valid {
		// error: field "Name" is missing
	}
}
```

Now we can tell **a field is missing** from **a field is empty** by consulting the sentinel `Field.Valid`:

- when `Field.Valid == false`, then _field is missing_
- when `Field.Valid == true`, then _field is provided_

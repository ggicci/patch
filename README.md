# patch

This package aims to answer the question asked when implementing **PATCH** requests of RESTful API in Go:

<p align="center">
    <b>Q: How to tell if a field is missing in the payload of a PATCH request?</b>
</p>

Since we are using generics, **Go 1.18+ is required**.

---

## Definition

Here we only talk about JSON payloads as it's the most frequently used format when developing a RESTful API.

Before implementaion, we need to define what is a missing field?

In this package, a field is defined as a missing field when:

- if the name/key of the field is **not found** in the JSON object
- **or** the name/key of the field is **present but its value is `null`**, null is interpreted as having no value

For example, in the following two JSON objects, field `Name` is missing:

```json
{ "Age": 18 }
{ "Name": null, "Age": 18 }
```

A Go snippet example:

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
		// error: field "Name" is missing (not found or null)
	}
}
```

Now we can tell **a field is missing** from **a field is empty** by consulting the sentinel `Field.Valid`:

- when `Field.Valid == false`, then _field is missing_ or **is null**
- when `Field.Valid == true`, then _field is provided_ and **is not null**

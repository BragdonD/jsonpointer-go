#  jsonpointer-go

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

An implementation of RFC 6901 - JavaScript Object Notation (JSON) Pointer in Golang.

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Introduction

`jsonpointer-go` is a Go library for working with JSON Pointers as defined in [RFC 6901](https://tools.ietf.org/html/rfc6901). JSON Pointers provide a way to reference specific parts of a JSON document.

## Installation

To install the library, simply run:

```sh
go get github.com/bragdond/jsonpointer-go
```

## Usage

Here is a simple example of how to use `jsonpointer-go`:

```go
package main

import (
    "fmt"
    "github.com/bragdond/jsonpointer-go"
)

func main() {   
	jsonStr := `{
        "foo": ["bar", "baz"],
        "qux": {"corge": "grault"}
    }`

	var jsonData map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		panic(fmt.Sprintf("failed to unmarshal json string: %v", err))
	}

    pointer, err := jsonpointer.Parse("/foo/0")
    if err != nil {
        panic(err)
    }

    result, err := pointer.GetValue(jsonData)
    if err != nil {
        panic(err)
    }

    fmt.Println(result) // Output: bar
}
```

## Contributing

We welcome contributions! Please see the [CONTRIBUTING.md](CONTRIBUTING.md) file for guidelines on how to contribute.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Happy coding!
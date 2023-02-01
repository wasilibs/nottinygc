# nottinygc

nottinygc is a replacement memory allocator for TinyGo targetting WASI. The default allocator
is built for small code size which can cause performance issues in higher-scale use cases.
nottinygc replaces it with [bdwgc][1] for garbage collection and [mimalloc][2] for standard
malloc-based allocation. These mature libraries can dramatically improve the performance of
TinyGo WASI applications at the expense of several hundred KB of additional code footprint.

Note that this library currently only works when the scheduler is disabled.

## Usage

Using the library requires both importing it and modifying flags when invoking TinyGo.

Add a blank import to the library to your code's `main` package.

```go
import _ "github.com/wasilibs/nottinygc"
```

Additionally, add `-gc=custom` and `-tags=custommalloc` to your TinyGo build flags.

```bash
tinygo build -o main.wasm -gc=custom -tags=custommalloc -target=wasi -scheduler=none main.go
```

[1]: https://github.com/ivmai/bdwgc
[2]: https://github.com/microsoft/mimalloc

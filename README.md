# nottinygc

nottinygc requires TinyGo 0.28+

_nottinygc has reached end-of-life. The Go compiler [will support][6] wasmexport soon so should compile most binaries, including proxy-wasm. As soon as gotip includes it, this repository will be archived, and until then, it is recommended to pause rather than gamble on this bandaid._

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

### Using with Envoy

This library relies on WASI to implement certain functionality, for which Envoy's implementation
is incomplete, and in addition, it's ABI, proxy-wasm has design issues that prevent working with
defaults appropriate for normal TinyGo applications. As this project is commonly used with Envoy,
we provide a build tag, to work around these issues. If building an Envoy plugin, add
`-tags=nottinygc_envoy` (or combine it with additional tags) to your TinyGo build flags. This
will disable export of `malloc`/`free`/ and define a no-op `sched_yield` function.

Other hosts that implement WASI fully, such as [wazero][3] and ABIs with correct memory semantics,
such as [http-wasm][5], will not have any issue. Implementations of proxy-wasm other than Envoy
may also work fine without the build tag.

## Performance

Benchmarks are run against every commit in the [bench][4] workflow. GitHub action runners are highly
virtualized and do not have stable performance across runs, but the relative numbers within a run
should still be somewhat, though not precisely, informative.

One run looks like this

```
BenchmarkGC/bench.wasm-2         	      52	 220294559 ns/op
BenchmarkGC/benchref.wasm-2      	       6	2000167805 ns/op
```

The benchmark is very simple, allocating some large strings in a loop. We see nottinygc perform almost
10x better in this benchmark. Note that just allocation / collection time is not the only aspect of
GC performance, allocation locality and fragmentation can also affect performance of real-world
applications. We have found that the default allocator can cause applications to run out of memory,
possibly due to fragmentation, whereas this library will continue to run indefinitely.

[1]: https://github.com/ivmai/bdwgc
[2]: https://github.com/microsoft/mimalloc
[3]: https://github.com/tetratelabs/wazero
[4]: https://github.com/wasilibs/nottinygc/actions/workflows/bench.yaml
[5]: https://http-wasm.io/
[6]: https://github.com/golang/go/issues/65199

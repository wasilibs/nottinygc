module github.com/wasilibs/nottinygc/e2e/envoy-dispatch-call

go 1.20

replace github.com/wasilibs/nottinygc => ../..

require (
	github.com/magefile/mage v1.14.0
	github.com/tetratelabs/proxy-wasm-go-sdk v0.22.0
	github.com/tidwall/gjson v1.14.4
	github.com/wasilibs/nottinygc v0.6.0
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d
)

require (
	github.com/stretchr/testify v1.8.3 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
)

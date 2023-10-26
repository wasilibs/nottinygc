module bench

go 1.19

require github.com/wasilibs/nottinygc v0.0.0-00010101000000-000000000000

require (
	github.com/magefile/mage v1.14.0 // indirect
	github.com/tetratelabs/wazero v1.5.0
)

replace github.com/wasilibs/nottinygc => ../

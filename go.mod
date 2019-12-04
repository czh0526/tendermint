module github.com/czh0526/tendermint

go 1.13

require (
	github.com/spf13/cobra v0.0.1
	github.com/status-im/keycard-go v0.0.0-20191119114148-6dd40a46baa0
	github.com/tendermint/go-amino v0.14.1
	github.com/tendermint/tendermint v0.32.8
)

replace github.com/tendermint/tendermint v0.32.8 => ./tendermint

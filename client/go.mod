module github.com/decred/vspd/client/v2

go 1.19

require (
	decred.org/dcrwallet/v3 v3.0.0-20230519033517-96817277627d
	github.com/decred/dcrd/blockchain/stake/v5 v5.0.0
	github.com/decred/dcrd/chaincfg/chainhash v1.0.4
	github.com/decred/dcrd/chaincfg/v3 v3.1.2
	github.com/decred/dcrd/dcrutil/v4 v4.0.0
	github.com/decred/dcrd/txscript/v4 v4.0.0
	github.com/decred/dcrd/wire v1.5.0
	github.com/decred/slog v1.2.0
	github.com/decred/vspd/types/v2 v2.0.0
)

replace (
	github.com/decred/dcrd/blockchain/stake/v5 => github.com/decred/dcrd/blockchain/stake/v5 v5.0.0-20221022042529-0a0cc3b3bf92
	github.com/decred/dcrd/blockchain/standalone/v2 => github.com/decred/dcrd/blockchain/standalone/v2 v2.1.1-0.20230430213532-f95870f9c6af
	github.com/decred/dcrd/chaincfg/v3 => github.com/decred/dcrd/chaincfg/v3 v3.1.2-0.20230412145739-9aa79ec168f6
	github.com/decred/dcrd/gcs/v4 => github.com/decred/dcrd/gcs/v4 v4.0.0-20221022042529-0a0cc3b3bf92
)

require (
	decred.org/cspp/v2 v2.0.0 // indirect
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412 // indirect
	github.com/companyzero/sntrup4591761 v0.0.0-20200131011700-2b0d299dbd22 // indirect
	github.com/dchest/siphash v1.2.3 // indirect
	github.com/decred/base58 v1.0.4 // indirect
	github.com/decred/dcrd/blockchain/standalone/v2 v2.1.1-0.20230430213532-f95870f9c6af // indirect
	github.com/decred/dcrd/crypto/blake256 v1.0.1 // indirect
	github.com/decred/dcrd/crypto/ripemd160 v1.0.1 // indirect
	github.com/decred/dcrd/database/v3 v3.0.0 // indirect
	github.com/decred/dcrd/dcrec v1.0.0 // indirect
	github.com/decred/dcrd/dcrec/edwards/v2 v2.0.2 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.1.0 // indirect
	github.com/decred/dcrd/dcrjson/v4 v4.0.1 // indirect
	github.com/decred/dcrd/gcs/v4 v4.0.0 // indirect
	github.com/decred/dcrd/hdkeychain/v3 v3.1.0 // indirect
	github.com/decred/dcrd/rpc/jsonrpc/types/v4 v4.0.0 // indirect
	github.com/decred/go-socks v1.1.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/jrick/bitset v1.0.0 // indirect
	github.com/jrick/wsrpc/v2 v2.3.5 // indirect
	go.etcd.io/bbolt v1.3.7 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
)

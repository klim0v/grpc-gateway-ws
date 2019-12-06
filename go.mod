module github.com/klim0v/grpc-gateway-ws

go 1.13

require (
	github.com/MinterTeam/minter-go-node v1.0.5
	github.com/danil-lashin/iavl v0.11.2 // indirect
	github.com/danil-lashin/tendermint v0.31.4 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.2
	github.com/gorilla/websocket v1.4.1
	github.com/grpc-ecosystem/grpc-gateway v1.12.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.32.8
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	google.golang.org/genproto v0.0.0-20191205163323-51378566eb59
	google.golang.org/grpc v1.25.1
)

replace github.com/MinterTeam/minter-go-node v1.0.5 => github.com/MinterTeam/minter-go-node v1.0.5-0.20191113165918-fa18116d6a26

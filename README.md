# simple go grpc pubsub server

## Description
イベントをpushすると購読しているクライアントにイベントが通知されるだけのシンプルなサーバーです。  
gRPCを使用しているのでクライアント側は任意の言語で実装することができます。  

イベント名は任意のstring、push時にイベント名とmessageを指定。そのイベント名を購読している  
クライアントのみにmessageが通知される(詳細はprotocol buffersを参照)


## Requirement
- go version >=v1.11 (use modules)


## Installation
1. git clone ${GOPATH}/github.com/sakokazuki/simplegrpc or any directory
`git clone ${remoteurl} ${GOPATH}/github.com/sakokazuki/simplegrpc` or  
`git clone ${remoteurl} simplegrpc`
2. cd
`cd simplegrpc`
3. build and automatically update go.mod and download dependencies if needed. 
`make build` or `go build -tags=release`

## Development
`make dev` or `go run main.go`

## Useage
`./simplegrpc` and server start at `localhost:10151`

## Test Client

## go
サーバーに接続して定期的にイベントを発行するだけのシンプルなクライアントを一応用意しました。  
1. `make client`
or  
1. `cd client`
2. `go run main.go`　　


## unity
https://github.com/sakokazuki/SimplegrpcClientForUnity

## TODO
- サーバーのポートを起動時に設定できるように
- デバッグ用webサーバー用意
- テスト書いてみる？

## reference
コードの殆どはこれを参考にした  
plasma (https://github.com/openfresh/plasma)

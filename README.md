# simple go grpc pubsub server

## Description
イベントをpushすると購読しているクライアントにイベントが通知されるだけのシンプルなサーバーです。  
gRPCを使用しているのでクライアント側は任意の言語で実装することができます。  


## Requirement
- go version >=v1.11 (use modules)


## Installation
1. git clone any folder
`git clone ${remoteurl} simplegrpc`
2. cd
`cd simplegrpc`
3. build and automatically update go.mod and download dependencies if needed. 
`go build`

## Development
`go run main.go`

## Useage
`./simplegrpc` and server start at `localhost:10151`

## Test Client
サーバーに接続して定期的にイベントを発行するだけのシンプルなクライアントを一応用意しました。  

1. `cd clident`
2. `go run main.go`

## TODO
- サーバーのポートを起動時に設定できるように
- デバッグ用webサーバー用意
- テスト書いてみる？

## reference
plasma (https://github.com/openfresh/plasma)

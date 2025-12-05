生成API服务：
	goctl api new demo

根据指定的api文件生成API服务：
	goctl api go --api demo.api --dir .

生成RPC服务：
	goctl rpc new user




goctl model mysql ddl --src user.sql --dir .

goctl model mongo --type user --dir .

goctl api format --dir demo.api



# 单个 rpc 服务生成示例指令
$ goctl rpc protoc greet.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=. --client=true 
# 多个 rpc 服务生成示例指令
$ goctl rpc protoc greet.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=. --client=true -m


goctl rpc protoc test.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.


goctl rpc protoc test2.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=. --client=true -m


goctl api swagger --api swagger_test.api --dir docs --filename swaggertest


protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative hello.proto



protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/hello.proto





================================================================================================================================

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative wsserver.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative wscallback.proto


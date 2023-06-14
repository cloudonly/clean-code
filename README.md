## Connect-go

https://connect.build/docs/go/getting-started

```shell
go install github.com/bufbuild/buf/cmd/buf@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest
```

```shell
buf link
buf generate
```

```shell
curl \
    --header "Content-Type: application/json" \
    --data '{"name": "Jane"}' \
    http://localhost:8080/apis.luomu.greet.v1.GreetService/Greet
```

```shell
grpcurl \
    -protoset <(buf build -o -) -plaintext \
    -d '{"name": "Jane"}' \
    localhost:8080 apis.luomu.greet.v1.GreetService/Greet
```

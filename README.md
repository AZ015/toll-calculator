# toll-calculator

```
docker-compose up
```

## Install protobuf compiler
```
brew install protobuff
```

## Installing GRPC and Protobuffer plugins for Golang
1. Protobuffers
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
```

2. GRPC
```
go install google.golang.org/grpc/cmd/protoc-gen-go@v1.28
```

3. NOTE that you need to set the /go/bin dir in your path
```
PATH="${PATH}:${HOME}/go/bin"
```

4. Install the package dependencies
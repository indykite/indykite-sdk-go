# Manually Update gRPC API

To keep up the SDK with the latest API changes it requires frequent update of
the generated client code.

## Steps to generate new client

To make the flow simpler all the API definitions are published  to
[buf.build/indykite/indykiteapis](https://buf.build/indykite/indykiteapis) and
will be publicly available soon.

### Install the `buf` CLI

The buf CLI enables you to create consistent Protobuf APIs that preserve
compatibility and comply with best practices.

Follow the [Installation](https://docs.buf.build/installation) instructions.

#### With Homebrew

```shell
brew install bufbuild/buf/buf
```

### Generate the client code

```shell
buf generate buf.build/indykite/indykiteapis
```

or

```shell
make generate-proto:
```

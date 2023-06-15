<div align="left">
<a href="https://indykite.com">
<img src="https://github.com/indykite/.github/blob/master/assets/IndyKITE_Rough_red.png" alt="IndyKite Red Logo" width="100px" height="183px" align="right">
</a>
</div>

## Documentation
[IndyKite documentation](https://docs.indykite.com)

## Terminology

[IndyKite glossary](https://docs.indykite.com/docs/resources/glossary)

# IndyKite Client Libraries for Go

IndyKite is a cloud identity platform built to secure and manage
human & non-person (IoT) identities and their data. This repository contains the
Golang Library packages for [IndyKite Platform](https://indykite.com) Client SDK.

[![Build](https://github.com/indykite/indykite-sdk-go/actions/workflows/pr-test.yaml/badge.svg)](https://github.com/indykite/indykite-sdk-go/actions/workflows/pr-test.yaml)
[![codecov](https://codecov.io/gh/indykite/indykite-sdk-go/branch/master/graph/badge.svg?token=TFCDLXbnsh)](https://codecov.io/gh/indykite/indykite-sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/indykite/indykite-sdk-go)](https://goreportcard.com/report/github.com/indykite/indykite-sdk-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/indykite/indykite-sdk-go.svg)](https://pkg.go.dev/github.com/indykite/indykite-sdk-go)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Findykite%2Findykite-sdk-go.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Findykite%2Findykite-sdk-go?ref=badge_shield)

```go
import "github.com/indykite/indykite-sdk-go"
```

In order to access to the platform you must first obtain AppAgent credentials or Service account credentials.
These credentials can be obtained either from the [Admin Console](https://console.indykite.id)
or request one from your point of contact at IndyKite.

## Requirements

- Go version 1.20 or later
- An AppSpace, App, AppAgent, and AppAgent credentials for all APIs except for Config API
- Service account credentials for Config API

# Getting Started

## Initial settings

1. You need to have a configuration json file to be able to use the IndyKite Proto SDK. You can get it from your
   Indykite contact or from Indykite console if you have access to it.

   Example configuration file:

```json
{
  "appSpaceId": "gid:AAAAAimlBsA9X0Gnlv2UTituW4Q",
  "baseUrl": "https://jarvis.indykite.com",
  "applicationId": "gid:AAAABKiMFRKkFEiEn3c-izbz9sk",
  "appAgentId": "gid:AAAABUqz8GXJB0cDikdAH8ywz8I",
  "endpoint": "jarvis.indykite.com",
  "privateKeyJWK": {
    "kty": "EC",
    "d": "aa",
    "use": "sig",
    "crv": "P-256",
    "kid": "2e5lIxxb6obIwpok",
    "x": "6d83se2Eg",
    "y": "lshzMo",
    "alg": "ES256"
  },
  "privateKeyPKCS8Base64": "LS0tLS==",
  "privateKeyPKCS8": "-----BEGIN PRIVATE KEY-----\nM\n-----END PRIVATE KEY-----"
}
```

Conditionally optional parameters:

- baseUrl
- defaultTenantId
- endpoint

1. You have two choices to set up the necessary credentials. You either pass the json to the `INDYKITE_APPLICATION_CREDENTIALS`
   environment variable or set the `INDYKITE_APPLICATION_CREDENTIALS_FILE` environment variable to the configuration file's path.

   1. on Linux and OSX

      ```shell
      export INDYKITE_APPLICATION_CREDENTIALS='{"appSpaceId":"00000000-0000-4000-a000-000000000000","appAgentId":"00000000-0000-4000-a000-000000000001","endpoint": "application.indykite.com","privateKeyJWK":{"kty":"EC","d": "abcdef","use": "sig","crv": "P-256","kid":"efghij","x":"klmnop","y":"qrstvw","alg":"ES256"}}'`
      ```

      or

      ```shell
      export INDYKITE_APPLICATION_CREDENTIALS_FILE=/Users/xx/configuration.json
      ```

   1. on Windows command line

      ```shell
      setex INDYKITE_APPLICATION_CREDENTIALS='{"appSpaceId":"00000000-0000-4000-a000-000000000000","appAgentId":"00000000-0000-4000-a000-000000000001","endpoint": "application.indykite.com","privateKeyJWK":{"kty":"EC","d": "abcdef","use": "sig","crv": "P-256","kid":"efghij","x":"klmnop","y":"qrstvw","alg":"ES256"}}'`
      ```

      or

      ```shell
      setex INDYKITE_APPLICATION_CREDENTIALS_FILE "C:\Users\xx\Documents\configuration.json"
      ```

## SDK Development

Commit message follows
[commit guidelines](./doc/guides/commit-message.md#commit-message-guidelines)

## Roadmap

Checkout our roadmap on our
[issues page](https://github.com/indykite/indykite-sdk-go/issues)

## Contributing

[Contribution guidelines for this project](contributing.md)

## Support, Feedback, Connect with other developers

Feel free to file a bug, submit an issue or give us feedback on our
[issues page](https://github.com/indykite/indykite-sdk-go/issues)

## Vulnerability Reporting

[Responsible Disclosure](responsible_disclosure.md)

## Changelog

[Changelog](CHANGELOG.md)

## Contributers / Acknowledgements

Coming Soon!

## What is IndyKite

IndyKite is a cloud identity platform built to secure and manage human & non-person (IoT) identities and their data. Based on open source standards, the cloud platform gives developers the ability to secure data and embed identity controls into their Web 3.0 applications. Empowering the world’s 23 million developers without the need to involve security and identity specialists.

## License

[This project is licensed under the terms of the Apache 2.0 license.](LICENSE)

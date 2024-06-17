# Keeper Security's Linux Keyring Utility

This utility interacts with the native Linux APIs to store and retrieve secrets from the Keyring using [Secret Service](https://specifications.freedesktop.org/secret-service/latest/).

While initially developed to help Keeper secure KSM configs, this utility can be used by any integration, plugin, or code base, to store and retrieve credentials, secrets, and passwords in the Linux Keyring simply and natively.

This utility can be used using the pre-build binary from the releases page, or by importing it into your code base. Both use cases are covered below.

## Using the Executable 

Download the latest version from the releases page and optionally add it to PATH to get started.

### Usage

The executable supports two commands:

1. `set`
2. `get`

Both commands require an application `name` (i.e. the name of the secret in / to be stored in the Keyring) as the first argument.

### `set`

`set` requires a second argument of the secret to be stored. This can be either a:

1. BASE64 string
2. JSON string
3. Path to an existing JSON file

When the secret is saved to the Keyring it is first encoded into a BASE64 format (if not already a BASE64 string). This standardizes the format for both consistent storage and to make it easier to consume by Keeper integrations and products. 

> If you need a support for a different format, please submit a feature request. We'd be happy to extend this to support other use cases.

### `get`

`get` returns the stored BASE64 encoded config to `stdout` and exits with a `0` exit code. The requesting integration can capture the output for consumption. Any errors encountered retrieving the config will return a `non-zero` exit code and write to `stderr`.

### Example

```shell
# Save a secret
lku set APPNAME eyJ1c2VybmFtZSI6ICJnb2xsdW0iLCAicGFzc3dvcmQiOiAiTXlQcmVjaW91cyJ9
# or
lku set APPNAME config.json

# Retrieve a secret
lku get APPNAME
```

## Using in Your Code

You can install this utility into your code base using standard `go` commands:

```bash
go get -u github.com/Keeper-Security/linux-keyring-utility@latest
```

You can now include the `keyring` package in your application for easy keyring management:

```go
import (
    //...
    "github.com/Keeper-Security/linux-keyring-utility/pkg/keyring"
)
```

### Usage

### `set`

The `Set()` function of the `SecretProvider` takes two arguments. The first is the name of the secret to be stored, which is usually an application name. The second is the secret itself. This should be either:

1. A BASE64 string
2. A JSON string
3. A path to an existing JSON file

When the secret is saved to the Keyring it is first encoded into a BASE64 format (if not already a BASE64 string). This standardizes the format for both consistent storage and to make it easier to consume by Keeper integrations and products.

```go
provider := keyring.SecretProvider{}

err := provider.Set("MY_APP_NAME", "eyJ1c2VybmFtZSI6InVzZXIiLCAicGFzc3dvcmQiOiJwYXNzIn0=")
if err != nil {
    fmt.Println("Unable to set secret:", err)
}
```

### `get`

The `Get()` function of the `SecretProvider` returns the stored BASE64 encoded secret. Pass the name of the secret/application to retrieve the stored value.

```go
provider := keyring.SecretProvider{}

secret, err := provider.Get("MY_APP_NAME")
if err != nil {
    fmt.Println("Unable to get secret:", err)
}

fmt.Println(secret) // Prints the BASE64 encoded secret
```

## Contributing

Please read and refer to the contribution guide before making your first PR.

For bugs, feature requests, etc., please submit an issue!
# Linux Keyring Utility

This utility interacts with the native Linux APIs to store and retrieve secrets from the Keyring using [Secret Service](https://specifications.freedesktop.org/secret-service/latest/). It can be used by any integration, plugin, or code base to store and retrieve credentials, secrets, and passwords in any Linux Keyring simply and natively.

To use this utility, you can deploy the pre-built binary from the releases page, or by importing it into your code base. Both use cases are covered below.

For Windows implementations, see the [Windows Credential Utility](https://github.com/Keeper-Security/windows-credential-utility).

## Details

The Linux Keyring Utility gets and sets _secrets_ in a Linux
[Keyring](http://man7.org/linux/man-pages/man7/keyrings.7.html) using the
[D-Bus](https://dbus.freedesktop.org/doc/dbus-tutorial.html)
[Secret Service](https://specifications.freedesktop.org/secret-service/latest/).

It has been tested with
[GNOME Keyring](https://wiki.gnome.org/Projects/GnomeKeyring/) and
[KDE Wallet Manager](https://userbase.kde.org/KDE_Wallet_Manager).
It _should_ work with any implementation of the D-Bus Secrets Service.

It includes a simple get/set/del(ete) CLI implemented with
[Cobra](https://cobra.dev).

### Interface

There are two packages, `dbus_secrets` and `secret_collection`.
The `secret_collection` object uses the functions in `dbus_secrets`.
It unifies the D-Bus _Connection_, _Session_ and _Collection Service_ objects.
It offers the simple get/set/delete interface that the CLI uses.

## Usage

The Go Language API has offers `Get()`, `Set()` and `Delete()` methods.
The first two accept and return `string` data.

### Example (get)

```go
package main

import (
    "os"
    sc "github.com/Keeper-Security/linux-keyring-utility/pkg/secret_collection"
)

func doit() {
    if collection, err := sc.DefaultCollection(); err == nil {
        if err := collection.Unlock(); err == nil {
            if secret, err := collection.Get("myapp", "mysecret"); err == nil {
                print(string(secret))
                os.Exit(0)
            }
        }
    }
    os.Exit(1)
}
```

The `.DefaultCollection()` returns whatever collection the _default_ _alias_ refers to.
It will generate an error if the _default_ alias is not set.
It usually points to the _login_ keyring.
Most Linux Keyring interfaces allow the user to set it.

The `.NamedCollection(string)` method provides access to collections by name.

### Example (set)

Set takes the data as a parameter and only returns an error.

```go
if err := collection.Set("myapp", "mysecret", "mysecretdata"); err == nil {
    // success
}
```

Set accepts _any_ string as secret data.

## Binary Interface (CLI)

The Linux binary supports three subcommands:

1. `get`
2. `set`
3. `del`

_Get_ and _del_ require one parameter; name, which is the secret _Label_ in D-Bus API terms.

_Del_ accepts one or more secret labels and deletes all of them.
If it generates an error it will stop.

_Set_ also requires the data as a _single_ string in the second parameter.
For example, `set foo bar baz` will generate an error but `set foo 'bar baz'` will work.
If the string is `-` then the string is read from standard input.

### Base64 encoding

_Get_ and _set_ take a `-b` or `--base64` flag that handles base64 automatically.
If used, _Set_ will encode the input before storing it and/or _get_ will decode it before printing.

Note that calling `get -b` on a secret that is _not_ base64 encoded secret will generate an error.

### Examples

```shell
# set has no output
lkru set root_cred '{
    "username": "root"
    "password": "rand0m."
}'
# get prints (to stdout) whatever was set
lku get root_cred
{
    "username": "root"
    "password": "rand0m."
}
lkru set -b root_cred2 '{"username": "gollum", "password": "MyPrecious"}'
lkru get root_cred2
eyJ1c2VybmFtZSI6ICJnb2xsdW0iLCAicGFzc3dvcmQiOiAiTXlQcmVjaW91cyJ9
lkru get -b root_cred2
{"username": "gollum", "password": "MyPrecious"}
cat ./good_cred.json | lkru set -b root_cred3 -
lkru get root_cred3
ewogICJ1c2VybmFtZSI6ICJhZGFtIiwKICAicGFzc3dvcmQiOiAicGFzc3dvcmQxMjMuIgp9
# errors go to stderr
lkru get root_cred4 2>/dev/null
lkru get root_cred4
Unable to get secret 'root_cred4': Unable to retrieve secret 'root_cred4' for application 'lkru' from collection '/org/freedesktop/secrets/aliases/default': org.freedesktop.Secret.Collection.SearchItems returned nothing
# a missing wallet error
lkru -c missing_wallet get root_cred
Error unlocking the keyring: Unable to unlock collection '/org/freedesktop/secrets/collection/missing_wallet': Object /org/freedesktop/secrets/collection/missing_wallet does not exist
# this error happens when there is no D-Bus running
lkru get test_cred
Unable to get the default keyring: Unable to connect to the D-Bus Session Bus: exec: "dbus-launch": executable file not found in $PATH
# this happens when D-Bus has no Secrets Service, i.e., there are no Keyring(s)
lkru get test_cred
Unable to get the default keyring: Unable to open a D-Bus session: The name org.freedesktop.secrets was not provided by any .service files
```

## Contributing

Please read and refer to the contribution guide before making your first PR.

For bugs, feature requests, etc., please submit an issue!

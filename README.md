# Keeper Security's Linux Keyring Utility

A utility that gets and sets _secrets_ in a Linux
[Keyring](http://man7.org/linux/man-pages/man7/keyrings.7.html) using the
[D-Bus](https://dbus.freedesktop.org/doc/dbus-tutorial.html)
[Secret Service](https://specifications.freedesktop.org/secret-service/latest/).
It tested with using
[GNOME Keyring](https://wiki.gnome.org/Projects/GnomeKeyring/) and
[KDE Wallet Manager](https://userbase.kde.org/KDE_Wallet_Manager).
It _should_ work with any implementation of the D-Bus Secrets Service.

It includes a simple get/set/del(ete) CLI implemented with
[Cobra](https://cobra.dev).

## Usage

The Go Language API has offers `Get()`, `Set()` and `Delete()` methods.
The first two accept and return `string` data.

### Go API

The `secret_collection` API is a wrapper object for the function in the `dbus_secrets`.
It unifies the D-Bus _Connection_, _Session_ and _Collection Service_ objects.

#### Example (get)

Complete `get` example:

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

#### Example (set)

Set takes the data as a parameter and only returns an error

```go
if err := collection.Set("myapp", "mysecret", "mysecretdata"); err == nil {
    // success
}
```

### Binary Interface (CLI)

The Linux binary supports three subcommands:

1. `get`
2. `set`
3. `del`

_Get_ and _del_ require one parameter; name, which is the secret _Label_ in D-Bus API terms.

_Set_ also requires the data as a string as the second parameter.

#### Base64 encoding

_Get_ and _set_ take a `-b` or `--base64` flag that handles base64 automatically.
If used, _Set_ will encode the input before storing it and/or _get_ will decode it before printing.

Note that calling `get -b` on a secret that is _not_ base64 encoded secret will generate an error.

### CLI Examples

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
lkru set -b root_cred2 $(echo '{"username": "gollum", "password": "MyPrecious"}')
lkru get root_cred2
eyJ1c2VybmFtZSI6ICJnb2xsdW0iLCAicGFzc3dvcmQiOiAiTXlQcmVjaW91cyJ9
lkru get -b root_cred2
{"username": "gollum", "password": "MyPrecious"}
# errors go to stderr
lkru get root_cred3 2>/dev/null
lkru get root_cred3
Unable to get secret 'root_cred3': Unable to retrieve secret 'root_cred3' for application 'lkru' from collection '/org/freedesktop/secrets/aliases/default': org.freedesktop.Secret.Collection.SearchItems returned nothing
# most errors are obvious
lkru -c missing_wallet get root_cred
Error unlocking the keyring: Unable to unlock collection '/org/freedesktop/secrets/collection/missing_wallet': Object /org/freedesktop/secrets/collection/missing_wallet does not exist
```

## Contributing

Please read and refer to the contribution guide before making your first PR.

For bugs, feature requests, etc., please submit an issue!

# Keeper Security's Linux Keyring Utility

A utility that stores and retrieves secrets from the Linux Keyring using the
[D-Bus](https://dbus.freedesktop.org/doc/dbus-tutorial.html)
[Secret Service](https://specifications.freedesktop.org/secret-service/latest/).
It tested with using
[GNOME Keyring](https://wiki.gnome.org/Projects/GnomeKeyring/) and
[KDE Walltet Manager](https://userbase.kde.org/KDE_Wallet_Manager).
It _should_ work with any implementation of the D-Bus Secrets Service.

It includes a simple get/set/del(ete) CLI implemented with
[Cobra](https://cobra.dev).

## Download

A static binary is available via [Releases](releases).

## Usage (API)

The `secret_collection` API is a wrapper object for the function in the `dbus_secrets`.
It unifies the D-Bus _Connection_, _Session_ and _Collection Service_ objects.

### Example (get)

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

### Example (set)

Set takes the data as a parameter and only returns an error

```go
if err := collection.Set("myapp", "mysecret", "mysecretdata"); err == nil {
    // success
}
```

## Usage (CLI)

The executable supports two subcommands:

1. `set`
2. `get`

Get requires one parameter; name, which is the secret _Label_ in D-Bus API terms.

Set requires the name and the data as a string as the second parameter.

### Example

```shell
# set has no output
lkru set root_cred '{
    "username": "root"
    "password": "rand0m."
}'
# get prints (to stdout) whatever was set
lku get root_cred
{
"foo": "bar"
}
lkru set root_cred2 $(echo '{"username": "gollum", "password": "MyPrecious"}' | base64 -w0 -)
lku get root_cred2
eyJ1c2VybmFtZSI6ICJnb2xsdW0iLCAicGFzc3dvcmQiOiAiTXlQcmVjaW91cyJ9
# errors go to stderr
get root_cred3 2>/dev/null
get root_cred3
Unable to get secret 'root_cred3': Unable to retrieve secret 'root_cred3' for application 'lku' from collection '/org/freedesktop/secrets/aliases/default': org.freedesktop.Secret.Collection.SearchItems returned nothing
```

## Contributing

Please read and refer to the contribution guide before making your first PR.

For bugs, feature requests, etc., please submit an issue!

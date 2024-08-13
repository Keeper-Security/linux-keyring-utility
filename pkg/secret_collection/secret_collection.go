package secret_collection

import (
	"fmt"

	secrets "github.com/Keeper-Security/linux-keyring-utility/pkg/dbus_secrets"
	dbus "github.com/godbus/dbus/v5"
)

type SecretCollection struct {
	Conn    *dbus.Conn
	Session dbus.ObjectPath
	Path    dbus.ObjectPath
}

func Collection(name string) (*SecretCollection, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to the D-Bus Session Bus: %v", err)
	}
	session, err := secrets.Open(conn)
	if err != nil {
		return nil, fmt.Errorf("Unable to open a D-Bus session: %v", err)
	}
	if name != "" && name != "default" {
		return &SecretCollection{conn, session, secrets.NamedCollection(conn, name)}, nil
	} else {
		return &SecretCollection{conn, session, secrets.DefaultCollection(conn)}, nil
	}
}

func DefaultCollection() (*SecretCollection, error) {
	return Collection("default")
}

func (sc *SecretCollection) Get(applicationName, label string) ([]byte, error) {
	if secret, err := secrets.GetItem(sc.Conn, sc.Path, sc.Session, applicationName, label); err == nil {
		return secret.Value, nil
	} else {
		return nil, fmt.Errorf("Unable to retrieve the secret '%s' for application '%s' from collection '%s': %s: %v", label, applicationName, sc.Path, label, err)
	}
}

func (sc *SecretCollection) Set(applicationName, label string, data []byte) error {
	if _, err := secrets.CreateItem(sc.Conn, sc.Path, applicationName, label, secrets.Secret(sc.Session, data)); err != nil {
		return fmt.Errorf("Unable to create a secret in collection '%s' for application '%s' with label '%s': %v", sc.Path, applicationName, label, err)
	}
	return nil
}

func (sc *SecretCollection) Unlock() error {
	errorf := func(err error) error {
		return fmt.Errorf("Unable to unlock collection '%s': %v", sc.Path, err)
	}

	if unlocked, err := secrets.Unlock(sc.Conn, []dbus.ObjectPath{sc.Path}); err == nil {
		for _, secret := range unlocked {
			if secret == sc.Path {
				return nil
			}
		}
		if len(unlocked) > 0 {
			return errorf(fmt.Errorf("unlocked %d collections not including it", len(unlocked)))
		} else {
			return errorf(fmt.Errorf("no collections were unlocked"))
		}
	} else {
		return errorf(err)
	}
}

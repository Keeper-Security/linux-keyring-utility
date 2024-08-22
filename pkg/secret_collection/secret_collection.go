package secret_collection

import (
	"fmt"

	secrets "github.com/Keeper-Security/linux-keyring-utility/pkg/dbus_secrets"
	dbus "github.com/godbus/dbus/v5"
)

// SecretCollection represents a D-BUS Secrets API Collection which notably requires a Session.
type SecretCollection struct {
	Conn    *dbus.Conn
	Session dbus.ObjectPath
	Path    dbus.ObjectPath
}

// Collection returns a collection by its name.
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

// DefaultCollection returns the "default" collection, i.e., whatever ".../aliases/default" refers to.
func DefaultCollection() (*SecretCollection, error) {
	return Collection("default")
}

// Delete removes a secret from the collection
func (sc *SecretCollection) Delete(applicationName, label string) error {
	errorf := func(err error) error {
		return fmt.Errorf("Unable to delete secret '%s' for application '%s' from collection '%s': %v", label, applicationName, sc.Path, err)
	}
	if item, err := secrets.GetItem(sc.Conn, sc.Path, sc.Session, applicationName, label); err == nil {
		if err := secrets.DeleteItem(sc.Conn, item.Object); err == nil {
			return nil
		} else {
			return errorf(err)
		}
	} else {
		return errorf(err)
	}
}

// Get retrieves a secret from the collection
func (sc *SecretCollection) Get(applicationName, label string) ([]byte, error) {
	if secret, err := secrets.GetItem(sc.Conn, sc.Path, sc.Session, applicationName, label); err == nil {
		return secret.Value, nil
	} else {
		return nil, fmt.Errorf("Unable to retrieve secret '%s' for application '%s' from collection '%s': %v", label, applicationName, sc.Path, err)
	}
}

// Set stores a secret in the collection
func (sc *SecretCollection) Set(applicationName, label string, data []byte) error {
	if _, err := secrets.SetItem(sc.Conn, sc.Path, sc.Session, applicationName, label, data, secrets.StringContentType); err != nil {
		return fmt.Errorf("Unable to create secret in collection '%s' for application '%s' with label '%s': %v", sc.Path, applicationName, label, err)
	} else {
		return nil
	}
}

// Unlock unlocks the collection which is required to access the secrets in it.
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

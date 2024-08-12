package dbus_secrets

import (
	"fmt"

	dbus "github.com/godbus/dbus/v5"
)

const (
	applicationName       = "Keeper Keyring Utility"
	completedSignal       = "org.freedesktop.Secret.Prompt.Completed"
	createItemMethod      = "org.freedesktop.Secret.Collection.CreateItem"
	defaultCollectionPath = "/org/freedesktop/secrets/aliases/default"
	getSecretMethod       = "org.freedesktop.Secret.Item.GetSecret"
	itemAttributesVariant = "org.freedesktop.Secret.Item.Attributes"
	itemDeletetMethod     = "org.freedesktop.Secret.Item.Delete"
	itemLabelVariant      = "org.freedesktop.Secret.Item.Label"
	openSessionMethod     = "org.freedesktop.Secret.Service.OpenSession"
	promptMethod          = "org.freedesktop.Secret.Prompt.Prompt"
	searchItemsMethod     = "org.freedesktop.Secret.Collection.SearchItems"
	serviceName           = "org.freedesktop.secrets"
	servicePath           = "/org/freedesktop/secrets"
	unlockMethod          = "org.freedesktop.Secret.Service.Unlock"
)

// secretObject is a DBus Secret Server secret object.
type secretObject struct {
	Session     dbus.ObjectPath
	Parameters  []byte
	Value       []byte
	ContentType string `dbus:"content_type"`
}

// attributes returns a map of attributes to be attached to the secret.
func attributes(application string) map[string]string {
	return map[string]string{
		"Agent":       "keeper-keyring-utility",
		"Application": application,
	}
}

// busObject returns a Secret Service specific dbus.BusObject from its path.
func busObject(conn *dbus.Conn, path dbus.ObjectPath) dbus.BusObject {
	return conn.Object(serviceName, path)
}

// Secret returns a secret object with the given session and secret data.
func Secret(session dbus.ObjectPath, secretData []byte) secretObject {
	return secretObject{
		ContentType: "text/plain; charset=utf8",
		Session:     session,
		Value:       secretData,
	}
}

// DefaultCollection returns the default collection path.
func DefaultCollection(conn *dbus.Conn) dbus.ObjectPath {
	return dbus.ObjectPath(defaultCollectionPath)
}

// NamedCollection returns an object representing a collection with the given name.
func NamedCollection(conn *dbus.Conn, name string) dbus.ObjectPath {
	return dbus.ObjectPath(servicePath + "/collection/" + name)
}

// Open creates a new session with the secret service.
func Open(conn *dbus.Conn) (dbus.ObjectPath, error) {
	var output dbus.Variant
	var sessionPath dbus.ObjectPath

	// The D-Bus API requires a lot of empty strings and unused variables. For example:
	if err := busObject(conn, servicePath).Call(openSessionMethod, 0, "plain", dbus.MakeVariant("")).Store(&output, &sessionPath); err != nil {
		return "", err
	}
	return conn.Object(serviceName, sessionPath).Path(), nil
}

// Close closes the session with the secret service.
func Close(conn *dbus.Conn, session dbus.ObjectPath) error {
	return busObject(conn, session).Call("org.freedesktop.Secret.Session.Close", 0).Err
}

// Unlock unlocks the given objects.
func Unlock(conn *dbus.Conn, objects []dbus.ObjectPath) ([]dbus.ObjectPath, error) {
	var unlocked []dbus.ObjectPath
	var prompt dbus.ObjectPath

	for _, object := range objects {
		if err := busObject(conn, servicePath).Call(unlockMethod, 0, []dbus.ObjectPath{object}).Store(&unlocked, &prompt); err != nil {
			return nil, err
		}
		// Anything other than "/" means prompt the user per the D-Bus Secrets API docs.
		if prompt != dbus.ObjectPath("/") {
			promptSignal := make(chan *dbus.Signal, 1)
			conn.Signal(promptSignal)
			if err := busObject(conn, prompt).Call(promptMethod, 0, "").Err; err != nil {
				return nil, err
			}
			signal := <-promptSignal
			switch signal.Name {
			case completedSignal:
				dismissed := signal.Body[0].(bool)
				result := signal.Body[1].(dbus.Variant)

				if dismissed {
					return nil, fmt.Errorf("%s: prompt was dismissed", promptMethod)
				} else if result.Value() != nil {
					collections := result.Value()
					switch c := collections.(type) {
					case []dbus.ObjectPath:
						unlocked = append(unlocked, c...)
					}
				}
			}
		}
	}
	return unlocked, nil
}

// CreateItem creates a new item in the given collection with the given secret.
func CreateItem(conn *dbus.Conn, collection dbus.ObjectPath, applicationName string,
	secretLabel string, secret secretObject) (
	dbus.BusObject, error) {
	var item, prompt dbus.ObjectPath

	if err := busObject(conn, collection).Call(createItemMethod, 0, map[string]dbus.Variant{
		itemLabelVariant:      dbus.MakeVariant(secretLabel),
		itemAttributesVariant: dbus.MakeVariant(attributes(applicationName)),
	}, secret, true).Store(&item, &prompt); err != nil {
		return nil, err
	} else if prompt != dbus.ObjectPath("/") {
		return nil, nil
	}
	return busObject(conn, item), nil
}

// DeleteItem deletes the given item.
func DeleteItem(conn *dbus.Conn, item dbus.BusObject) error {
	var prompt dbus.ObjectPath
	if err := item.Call(itemDeletetMethod, 0).Store(&prompt); err != nil {
		return err
	}
	return nil
}

// GetItem returns the secret object with the corresponding application and label.
func GetItem(conn *dbus.Conn, collection dbus.ObjectPath, session dbus.ObjectPath, applicationName string, label string) (*secretObject, error) {
	var items []dbus.ObjectPath

	if err := busObject(conn, collection).Call(searchItemsMethod, 0, attributes(applicationName)).Store(&items); err == nil {
		if len(items) == 1 {
			var secret secretObject
			if err := busObject(conn, items[0]).Call(getSecretMethod, 0, session).Store(&secret); err != nil {
				return nil, err
			}
			return &secret, nil
		} else if len(items) > 1 {
			return nil, fmt.Errorf("%s returned %d items", searchItemsMethod, len(items))
		} else {
			return nil, fmt.Errorf("%s returned nothing", searchItemsMethod)
		}
	} else {
		return nil, err
	}
}

package dbus_secrets

import (
	"testing"

	dbus "github.com/godbus/dbus/v5"
	"github.com/stretchr/testify/assert"
)

const (
	testApplicationName = "dbus_secrets_test"
	testSecretData      = "ewogICAgInVzZXJuYW1lIjogImNhcHRhaW5fdW5kZXJwYW50cyIsCiAgICAicGFzc3dvcmQiOiAiWWFtMXczdD8hIgp9Cg=="
	testSecretData2     = "ewogICJ1c2VybmFtZSI6ICJndW1uYmFsbCIsCiAgInBhc3N3b3JkIjogIkxVTkNIQk9YMSEhIgp9Cg=="
	testSecretLabel     = "test-secret"
	testSecretLabel2    = "test-secret2"
	testNamedCollection = "login"
)

func TestUnlockDefaultCollection(t *testing.T) {
	conn, err := dbus.SessionBus()
	if err != nil {
		t.Fatalf("Failed to connect to session bus: %v", err)
	}
	defer conn.Close()

	session, err := Open(conn)
	if err != nil {
		t.Fatalf("Failed to open session: %v", err)
	}
	defer Close(conn, session)

	collection := DefaultCollection(conn)

	Unlock(conn, []dbus.ObjectPath{collection})
}

func TestUnlockNamedCollection(t *testing.T) {
	conn, err := dbus.SessionBus()
	if err != nil {
		t.Fatalf("Failed to connect to session bus: %v", err)
	}
	defer conn.Close()

	session, err := Open(conn)
	if err != nil {
		t.Fatalf("Failed to open session: %v", err)
	}
	defer Close(conn, session)

	collection := NamedCollection(conn, testNamedCollection)

	Unlock(conn, []dbus.ObjectPath{collection})
}

func TestSetItem(t *testing.T) {
	conn, err := dbus.SessionBus()
	if err != nil {
		t.Fatalf("Failed to connect to session bus: %v", err)
	}
	defer conn.Close()

	session, err := Open(conn)
	if err != nil {
		t.Fatalf("Failed to open session: %v", err)
	}
	defer Close(conn, session)

	collection := DefaultCollection(conn)

	Unlock(conn, []dbus.ObjectPath{collection})

	if _, err := SetItem(conn, collection, session, testApplicationName, testSecretLabel, []byte(testSecretData), StringContentType); err != nil {
		t.Fatalf("Failed to create item: %v", err)
	}
}

func TestSetEmptyItem(t *testing.T) {
	conn, err := dbus.SessionBus()
	if err != nil {
		t.Fatalf("Failed to connect to session bus: %v", err)
	}
	defer conn.Close()

	session, err := Open(conn)
	if err != nil {
		t.Fatalf("Failed to open session: %v", err)
	}
	defer Close(conn, session)

	collection := DefaultCollection(conn)

	Unlock(conn, []dbus.ObjectPath{collection})

	if _, err := SetItem(conn, collection, session, testApplicationName, "", []byte(""), StringContentType); err == nil {
		t.Fatalf("Created item with empty secret")
	}
}

func TestSetAndGetItem(t *testing.T) {
	conn, err := dbus.SessionBus()
	if err != nil {
		t.Fatalf("Failed to connect to session bus: %v", err)
	}
	defer conn.Close()

	session, err := Open(conn)
	if err != nil {
		t.Fatalf("Failed to open session: %v", err)
	}
	defer Close(conn, session)

	collection := DefaultCollection(conn)

	Unlock(conn, []dbus.ObjectPath{collection})

	if _, err := SetItem(conn, collection, session, testApplicationName, testSecretLabel, []byte(testSecretData), StringContentType); err == nil {
		if item, err := GetItem(conn, collection, session, testApplicationName, testSecretLabel); err == nil {
			assert.EqualValues(t, item.Value, []byte(testSecretData))
		} else {
			t.Fatalf("Failed to get item: %v", err)
		}
	} else {
		t.Fatalf("Failed to create item: %v", err)
	}
}

func TestSetAndDeleteItem(t *testing.T) {
	conn, err := dbus.SessionBus()
	if err != nil {
		t.Fatalf("Failed to connect to session bus: %v", err)
	}
	defer conn.Close()

	session, err := Open(conn)
	if err != nil {
		t.Fatalf("Failed to open session: %v", err)
	}
	defer Close(conn, session)

	collection := DefaultCollection(conn)

	Unlock(conn, []dbus.ObjectPath{collection})

	if item, err := SetItem(conn, collection, session, testApplicationName, testSecretLabel, []byte(testSecretData), StringContentType); err == nil {
		if err := DeleteItem(conn, item); err != nil {
			t.Fatalf("Failed to delete item: %v", err)
		}
	} else {
		t.Fatalf("Failed to create item: %v", err)
	}
}

func TestDeleteNonExistentItem(t *testing.T) {
	conn, err := dbus.SessionBus()
	if err != nil {
		t.Fatalf("Failed to connect to session bus: %v", err)
	}
	defer conn.Close()

	session, err := Open(conn)
	if err != nil {
		t.Fatalf("Failed to open session: %v", err)
	}
	defer Close(conn, session)

	collection := DefaultCollection(conn)

	Unlock(conn, []dbus.ObjectPath{collection})

	if item, err := SetItem(conn, collection, session, testApplicationName, testSecretLabel, []byte(testSecretData), StringContentType); err == nil {
		if err := DeleteItem(conn, item); err != nil {
			t.Fatalf("Failed to delete item: %v", err)
		}
		if err := DeleteItem(conn, item); err == nil {
			t.Fatalf("Deleted non-existent item")
		}
	} else {
		t.Fatalf("Failed to create item: %v", err)
	}
}

func TestOverwriteItem(t *testing.T) {
	conn, err := dbus.SessionBus()
	if err != nil {
		t.Fatalf("Failed to connect to session bus: %v", err)
	}
	defer conn.Close()

	session, err := Open(conn)
	if err != nil {
		t.Fatalf("Failed to open session: %v", err)
	}
	defer Close(conn, session)

	collection := DefaultCollection(conn)

	Unlock(conn, []dbus.ObjectPath{collection})

	if _, err := SetItem(conn, collection, session, testApplicationName, testSecretLabel, []byte(""), StringContentType); err == nil {
		if _, err := SetItem(conn, collection, session, testApplicationName, testSecretLabel, []byte(testSecretData), StringContentType); err == nil {
			if item, err := GetItem(conn, collection, session, testApplicationName, testSecretLabel); err == nil {
				assert.EqualValues(t, testSecretData, item.Value)
			} else {
				t.Fatalf("Failed to get item: %v", err)
			}
		}
	} else {
		t.Fatalf("Failed to create item: %v", err)
	}
}

func TestMultiSetGetDelete(t *testing.T) {
	conn, err := dbus.SessionBus()
	if err != nil {
		t.Fatalf("Failed to connect to session bus: %v", err)
	}
	defer conn.Close()

	session, err := Open(conn)
	if err != nil {
		t.Fatalf("Failed to open session: %v", err)
	}
	defer Close(conn, session)

	collection := DefaultCollection(conn)

	Unlock(conn, []dbus.ObjectPath{collection})

	if item, err := SetItem(conn, collection, session, testApplicationName, testSecretLabel, []byte(testSecretData), StringContentType); err == nil {
		if item, err := SetItem(conn, collection, session, testApplicationName, testSecretLabel2, []byte(testSecretData2), StringContentType); err == nil {
			if item, err := GetItem(conn, collection, session, testApplicationName, testSecretLabel); err == nil {
				assert.EqualValues(t, testSecretData, string(item.Value))
			} else {
				t.Fatalf("Failed to get secret with label '%s': %v", testSecretLabel, err)
			}
			if item, err := GetItem(conn, collection, session, testApplicationName, testSecretLabel2); err == nil {
				assert.EqualValues(t, testSecretData2, string(item.Value))
			} else {
				t.Fatalf("Failed to get secret with label '%s': %v", testSecretLabel2, err)
			}
			if err := DeleteItem(conn, item); err != nil {
				t.Fatalf("Failed to delete item: %v", err)
			}
		}
		if err := DeleteItem(conn, item); err != nil {
			t.Fatalf("Failed to delete item: %v", err)
		}
	} else {
		t.Fatalf("Failed to create item: %v", err)
	}
}

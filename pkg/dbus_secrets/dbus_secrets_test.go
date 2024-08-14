package dbus_secrets

import (
	"testing"

	dbus "github.com/godbus/dbus/v5"
)

const (
	testApplicationName = "tests"
	testSecretData      = "ewogICAgInVzZXJuYW1lIjogImNhcHRhaW5fdW5kZXJwYW50cyIsCiAgICAicGFzc3dvcmQiOiAiWWFtMXczdD8hIgp9Cg=="
	testSecretLabel     = "test-secret"
	testNamedCollection = "kdewallet"
)

// TestT1 tests the t1 function.
func TestT1(t *testing.T) {
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

// TestT2 tests the t2 function.
func TestT2(t *testing.T) {
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

// TestT3 tests the t3 function.
func TestT3(t *testing.T) {
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
	// Add additional assertions or checks as needed
	if _, err := CreateItem(conn, collection, session, testApplicationName, testSecretLabel, []byte(testSecretData)); err != nil {
		t.Fatalf("Failed to create item: %v", err)
	}
}

func TestT4(t *testing.T) {
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
	// Add additional assertions or checks as needed
	if _, err := CreateItem(conn, collection, session, testApplicationName, testSecretLabel, []byte(testSecretData)); err == nil {
		if _, err := GetItem(conn, collection, session, testApplicationName, testSecretLabel); err != nil {
			t.Fatalf("Failed to get item: %v", err)
		}
	} else {
		t.Fatalf("Failed to create item: %v", err)
	}
}

func TestT5(t *testing.T) {
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
	// Add additional assertions or checks as needed
	if item, err := CreateItem(conn, collection, session, testApplicationName, testSecretLabel,[]byte(testSecretData)); err == nil {
		if err := DeleteItem(conn, item); err != nil {
			t.Fatalf("Failed to delete item: %v", err)
		}
	} else {
		t.Fatalf("Failed to create item: %v", err)
	}
}

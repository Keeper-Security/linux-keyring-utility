package secret_collection

import (
	"testing"
)

const (
	testApplication = "test-app"
	testCollection  = "kdewallet"
	testData        = "ewogICJ1c2VybmFtZSI6ICJteV91c2VybmFtZSIsCiAgInBhc3N3b3JkIjogIm15X3Bhc3N3b3JkIgp9Cg=="
	testLabel       = "test-label"
)

// TestNamedCollection tests the NamedCollection function.
func TestNamedCollection(t *testing.T) {
	collection, _ := Collection(testCollection)
	if collection == nil {
		t.Fatalf("Expected a valid SecretCollection, got nil")
	}
	if collection.Conn == nil {
		t.Fatalf("Expected a valid dbus.Conn, got nil")
	}
	if collection.Session == "" {
		t.Fatalf("Expected a valid dbus.ObjectPath for session, got empty")
	}
	if collection.Path == "" {
		t.Fatalf("Expected a valid dbus.ObjectPath for path, got empty")
	}
}

// TestDefaultCollection tests the DefaultCollection function.
func TestDefaultCollection(t *testing.T) {
	collection, _ := DefaultCollection()
	if collection == nil {
		t.Fatalf("Expected a valid SecretCollection, got nil")
	}
	if collection.Conn == nil {
		t.Fatalf("Expected a valid dbus.Conn, got nil")
	}
	if collection.Session == "" {
		t.Fatalf("Expected a valid dbus.ObjectPath for session, got empty")
	}
	if collection.Path == "" {
		t.Fatalf("Expected a valid dbus.ObjectPath for path, got empty")
	}
}

// TestSecretCollection_Unlock tests the Unlock method of SecretCollection.
func TestSecretCollection_Unlock(t *testing.T) {
	if collection, err := Collection(testCollection); err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	} else {
		if err := collection.Unlock(); err != nil {
			t.Fatalf("Expected no error but got: %v", err)
		}
	}
}

// TestSecretCollection_Set tests the Set method of SecretCollection.
func TestSecretCollection_Set(t *testing.T) {
	if collection, err := Collection(testCollection); err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	} else {
		if err := collection.Unlock(); err != nil {
			t.Fatalf("Expected no error but got: %v", err)
		}
		if err := collection.Set(testApplication, testLabel, []byte(testData)); err != nil {
			t.Fatalf("Expected no error but got: %v", err)
		}
	}
}

// TestSecretCollection_Get tests the Get method of SecretCollection.
func TestSecretCollection_Get(t *testing.T) {
	if collection, err := Collection(testCollection); err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	} else {
		if err := collection.Unlock(); err != nil {
			t.Fatalf("Expected no error but got: %v", err)
		}
		if item, err := collection.Get(testApplication, testLabel); err != nil {
			t.Fatalf("Expected no error but got: %v", err)
		} else {
			t.Log(string(item))
		}
	}
}

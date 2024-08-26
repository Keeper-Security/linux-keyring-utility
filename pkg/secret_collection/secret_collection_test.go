package secret_collection

import (
	"testing"
)

const (
	testApplication = "secret_collection_test"
	testCollection  = "login"
	testData        = "one thing"
	testData2       = "versus another"
	testLabel       = "arbitrary-text"
	testLabel2      = "intentionally-different"
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

func TestSecretCollection_Delete(t *testing.T) {
	if collection, err := Collection(testCollection); err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	} else {
		if err := collection.Unlock(); err != nil {
			t.Fatalf("Expected no error but got: %v", err)
		}
		if err := collection.Delete(testApplication, testLabel); err != nil {
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
		if err := collection.Set(testApplication, testLabel, []byte(testData)); err != nil {
			t.Fatalf("Expected no error but got: %v", err)
		}
		if _, err := collection.Get(testApplication, testLabel); err != nil {
			t.Fatalf("Expected no error but got: %v", err)
		}
	}
}

func TestSecretCollection_WithTwoItems(t *testing.T) {
	if collection, err := Collection(testCollection); err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	} else {
		if err := collection.Unlock(); err != nil {
			t.Fatalf("Expected no error but got: %v", err)
		}
		if err := collection.Set(testApplication, testLabel, []byte(testData)); err != nil {
			t.Fatalf("Expected no error but got: %v", err)
		}
		if err := collection.Set(testApplication, testLabel2, []byte(testData2)); err != nil {

			t.Fatalf("Expected no error but got: %v", err)
		}

		if item, err := collection.Get(testApplication, testLabel); err != nil {
			t.Fatalf("Expected no error but got: %v", err)
		} else {
			if string(item) != testData {
				t.Fatalf("Expected %s but got %s", testData, string(item))
			}
			collection.Delete(testApplication, testLabel)
		}

		if item, err := collection.Get(testApplication, testLabel2); err != nil {
			t.Fatalf("Expected no error but got: %v", err)
		} else {
			if string(item) != testData2 {
				t.Fatalf("Expected %s but got %s", testData2, string(item))
			}
			collection.Delete(testApplication, testLabel2)
		}
	}
}

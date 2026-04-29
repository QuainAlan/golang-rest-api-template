package database

import "testing"

func TestMongoConnectionURI(t *testing.T) {
	t.Run("default when empty", func(t *testing.T) {
		t.Setenv("MONGODB_URI", "")
		if got, want := mongoConnectionURI(), "mongodb://localhost:27017"; got != want {
			t.Fatalf("mongoConnectionURI() = %q, want %q", got, want)
		}
	})
	t.Run("custom", func(t *testing.T) {
		t.Setenv("MONGODB_URI", "mongodb://mongo:27017")
		if got, want := mongoConnectionURI(), "mongodb://mongo:27017"; got != want {
			t.Fatalf("mongoConnectionURI() = %q, want %q", got, want)
		}
	})
}

package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	messenger "github.com/slidebolt/sb-messenger-sdk"
	storage "github.com/slidebolt/sb-storage-sdk"
)

type keyed string

func (k keyed) Key() string { return string(k) }

func TestHelloManifest(t *testing.T) {
	h := New(DefaultConfig()).Hello()
	if h.ID != "storage" {
		t.Fatalf("id: got %q want %q", h.ID, "storage")
	}
	if len(h.DependsOn) != 1 || h.DependsOn[0] != "messenger" {
		t.Fatalf("dependsOn: got %v want [messenger]", h.DependsOn)
	}
}

func TestOnStartRegistersWorkingStorageService(t *testing.T) {
	msg, payload, err := messenger.MockWithPayload()
	if err != nil {
		t.Fatal(err)
	}
	defer msg.Close()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd)

	s := New(DefaultConfig())
	if _, err := s.OnStart(map[string]json.RawMessage{"messenger": payload}); err != nil {
		t.Fatal(err)
	}
	defer s.OnShutdown()

	client, err := storage.Connect(map[string]json.RawMessage{"messenger": payload})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	if err := client.Save(rawValue{key: "plugin.dev.ent", data: json.RawMessage(`{"power":true}`)}); err != nil {
		t.Fatal(err)
	}
	got, err := client.Get(keyed("plugin.dev.ent"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"power":true}` {
		t.Fatalf("data: got %s", string(got))
	}
}

func TestOnStartLoadsDiskStateAndWatchesExternalUpdates(t *testing.T) {
	msg, payload, err := messenger.MockWithPayload()
	if err != nil {
		t.Fatal(err)
	}
	defer msg.Close()

	tmp := t.TempDir()
	key := "plugin.frigate.camera01"
	filePath := filepath.Join(tmp, "plugin", "frigate", "camera01", "camera01.json")
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filePath, []byte(`{"power":false}`), 0644); err != nil {
		t.Fatal(err)
	}

	s := New(Config{DataDir: tmp})
	if _, err := s.OnStart(map[string]json.RawMessage{"messenger": payload}); err != nil {
		t.Fatal(err)
	}
	defer s.OnShutdown()

	client, err := storage.Connect(map[string]json.RawMessage{"messenger": payload})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	assertEventuallyData(t, client, keyed(key), `{"power":false}`)

	if err := os.WriteFile(filePath, []byte(`{"power":true}`), 0644); err != nil {
		t.Fatal(err)
	}

	assertEventuallyData(t, client, keyed(key), `{"power":true}`)
}

func assertEventuallyData(t *testing.T, client storage.Storage, key keyed, want string) {
	t.Helper()

	deadline := time.Now().Add(3 * time.Second)
	for {
		got, err := client.Get(key)
		if err == nil && string(got) == want {
			return
		}
		if time.Now().After(deadline) {
			if err != nil {
				t.Fatalf("get %q: %v", key, err)
			}
			t.Fatalf("data for %q: got %s want %s", key, string(got), want)
		}
		time.Sleep(25 * time.Millisecond)
	}
}

type rawValue struct {
	key  string
	data json.RawMessage
}

func (r rawValue) Key() string                  { return r.key }
func (r rawValue) MarshalJSON() ([]byte, error) { return r.data, nil }

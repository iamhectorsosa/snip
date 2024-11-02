package database_test

import (
	"testing"

	"github.com/iamhectorsosa/snip/internal/database"
	"github.com/iamhectorsosa/snip/internal/store"
)

func TestCreateAndRead(t *testing.T) {
	db, cleanup, err := database.NewInMemory()
	if err != nil {
		t.Fatalf("database.NewInMemory, err=%v", err)
	}
	defer cleanup()

	want := store.Snippet{
		Key:   "hello",
		Value: "world",
	}

	if err := db.Create(want.Key, want.Value); err != nil {
		t.Errorf("db.Create, err=%v", err)
	}

	got, err := db.Read(want.Key)
	if err != nil {
		t.Errorf("db.Read, err=%v", err)
	}

	if got.Key != want.Key || got.Value != want.Value {
		t.Errorf("unexpected key values, wanted %s=%q got %s=%q", want.Key, want.Value, got.Key, got.Value)
	}
}

func TestCreateAndReadAll(t *testing.T) {
	db, cleanup, err := database.NewInMemory()
	if err != nil {
		t.Fatalf("database.NewInMemory, err=%v", err)
	}
	defer cleanup()

	wantSnippets := []store.Snippet{
		store.Snippet{
			Key:   "hello",
			Value: "world",
		},
		store.Snippet{
			Key:   "ahoj",
			Value: "ciao",
		},
	}

	for _, s := range wantSnippets {
		if err := db.Create(s.Key, s.Value); err != nil {
			t.Errorf("db.Create, err=%v", err)
		}
	}

	snippets, err := db.ReadAll()
	if err != nil {
		t.Errorf("db.Read, err=%v", err)
	}

	if len(snippets) != 2 {
		t.Errorf("unexpected number of snippets, wanted 2 got %d", len(snippets))
	}

	for i, got := range snippets {
		want := wantSnippets[i]
		if got.Key != want.Key || got.Value != want.Value {
			t.Errorf("unexpected key values, wanted %s=\"%q\" got %s=%q", want.Key, want.Value, got.Key, got.Value)
		}
	}
}

func TestCreateUpdateAndRead(t *testing.T) {
	db, cleanup, err := database.NewInMemory()
	if err != nil {
		t.Fatalf("database.NewInMemory, err=%v", err)
	}
	defer cleanup()

	want := store.Snippet{
		Key:   "hello",
		Value: "world",
	}

	if err := db.Create(want.Key, want.Value); err != nil {
		t.Errorf("db.Create, err=%v", err)
	}

	got, err := db.Read(want.Key)
	if err != nil {
		t.Errorf("db.Read, err=%v", err)
	}

	if got.Key != want.Key || got.Value != want.Value {
		t.Errorf("unexpected key values, wanted %s=%q got %s=%q", want.Key, want.Value, got.Key, got.Value)
	}

	updateWant := store.Snippet{
		Key:   "hello",
		Value: "goodbye",
	}

	err = db.Update(updateWant)
	if err != nil {
		t.Errorf("db.Update, err=%v", err)
	}

	got, err = db.Read(updateWant.Key)
	if err != nil {
		t.Errorf("db.Read, err=%v", err)
	}

	if got.Key != updateWant.Key || got.Value != updateWant.Value {
		t.Errorf("unexpected key values, wanted %s=%q got %s=%q", updateWant.Key, updateWant.Value, got.Key, got.Value)
	}
}

func TestCreateDeleteAndRead(t *testing.T) {
	db, cleanup, err := database.NewInMemory()
	if err != nil {
		t.Fatalf("database.NewInMemory, err=%v", err)
	}
	defer cleanup()

	want := store.Snippet{
		Key:   "hello",
		Value: "world",
	}

	if err := db.Create(want.Key, want.Value); err != nil {
		t.Errorf("db.Create, err=%v", err)
	}

	got, err := db.Read(want.Key)
	if err != nil {
		t.Errorf("db.Read, err=%v", err)
	}

	if got.Key != want.Key || got.Value != want.Value {
		t.Errorf("unexpected key values, wanted %s=%q got %s=%q", want.Key, want.Value, got.Key, got.Value)
	}

	err = db.Delete(want.Key)
	if err != nil {
		t.Errorf("db.Delete, err=%v", err)
	}

	if _, err := db.Read(want.Key); err == nil {
		t.Errorf("expected error err=sql: no rows in result set")
	}
}

func TestCreateResetAndReadAll(t *testing.T) {
	db, cleanup, err := database.NewInMemory()
	if err != nil {
		t.Fatalf("database.NewInMemory, err=%v", err)
	}
	defer cleanup()

	wantSnippets := []store.Snippet{
		store.Snippet{
			Key:   "hello",
			Value: "world",
		},
		store.Snippet{
			Key:   "ahoj",
			Value: "ciao",
		},
	}

	for _, s := range wantSnippets {
		if err := db.Create(s.Key, s.Value); err != nil {
			t.Errorf("db.Create, err=%v", err)
		}
	}

	snippets, err := db.ReadAll()
	if err != nil {
		t.Errorf("db.Read, err=%v", err)
	}

	if len(snippets) != 2 {
		t.Errorf("unexpected number of snippets, wanted 2 got %d", len(snippets))
	}

	for i, got := range snippets {
		want := wantSnippets[i]
		if got.Key != want.Key || got.Value != want.Value {
			t.Errorf("unexpected key values, wanted %s=\"%q\" got %s=%q", want.Key, want.Value, got.Key, got.Value)
		}
	}

	if err := db.Reset(); err != nil {
		t.Errorf("db.Reset, err=%v", err)
	}

	snippets, err = db.ReadAll()
	if err != nil {
		t.Errorf("db.Read, err=%v", err)
	}

	if len(snippets) != 0 {
		t.Errorf("unexpected number of snippets, wanted 0 got %d", len(snippets))
	}
}

func TestImportAndReadAll(t *testing.T) {
	db, cleanup, err := database.NewInMemory()
	if err != nil {
		t.Fatalf("database.NewInMemory, err=%v", err)
	}
	defer cleanup()

	wantSnippets := []store.Snippet{
		store.Snippet{
			Key:   "hello",
			Value: "world",
		},
		store.Snippet{
			Key:   "ahoj",
			Value: "ciao",
		},
	}

	if err := db.Import(wantSnippets); err != nil {
		t.Errorf("db.Create, err=%v", err)
	}

	snippets, err := db.ReadAll()
	if err != nil {
		t.Errorf("db.Read, err=%v", err)
	}

	if len(snippets) != 2 {
		t.Errorf("unexpected number of snippets, wanted 2 got %d", len(snippets))
	}

	for i, got := range snippets {
		want := wantSnippets[i]
		if got.Key != want.Key || got.Value != want.Value {
			t.Errorf("unexpected key values, wanted %s=\"%q\" got %s=%q", want.Key, want.Value, got.Key, got.Value)
		}
	}
}

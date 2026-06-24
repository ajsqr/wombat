package cmap

import (
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	m := New[string, int]()

	if m == nil {
		t.Fatal("expected non-nil map")
	}

	if m.data == nil {
		t.Fatal("expected underlying map to be initialized")
	}

	if m.mx == nil {
		t.Fatal("expected mutex to be initialized")
	}
}

func TestSetAndGet(t *testing.T) {
	m := New[string, int]()

	m.Set("foo", 42)

	v, ok := m.Get("foo")

	if !ok {
		t.Fatal("expected key to exist")
	}

	if v != 42 {
		t.Fatalf("expected 42 got %d", v)
	}
}

func TestGetMissingKey(t *testing.T) {
	m := New[string, int]()

	_, ok := m.Get("missing")

	if ok {
		t.Fatal("expected key to not exist")
	}
}

func TestOverwriteExistingKey(t *testing.T) {
	m := New[string, int]()

	m.Set("foo", 1)
	m.Set("foo", 2)

	v, ok := m.Get("foo")

	if !ok {
		t.Fatal("expected key to exist")
	}

	if v != 2 {
		t.Fatalf("expected 2 got %d", v)
	}
}

func TestDelete(t *testing.T) {
	m := New[string, int]()

	m.Set("foo", 42)

	m.Delete("foo")

	_, ok := m.Get("foo")

	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestDeleteMissingKey(t *testing.T) {
	m := New[string, int]()

	// should not panic
	m.Delete("missing")

	_, ok := m.Get("missing")

	if ok {
		t.Fatal("did not expect key to exist")
	}
}

func TestConcurrentAccess(t *testing.T) {
	m := New[int, int]()

	const workers = 1000

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			m.Set(i, i)

			v, ok := m.Get(i)

			if !ok {
				t.Errorf("expected key %d to exist", i)
				return
			}

			if v != i {
				t.Errorf("expected %d got %d", i, v)
			}
		}(i)
	}

	wg.Wait()
}

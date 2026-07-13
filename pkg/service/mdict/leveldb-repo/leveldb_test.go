//
// Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package leveldb_repo

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"
)

func TestLvDB_PutGetPrefix(t *testing.T) {
	db, err := NewLvDB(filepath.Join(t.TempDir(), "kv"))
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Put("PFKW#_apple", []byte("a")); err != nil {
		t.Fatal(err)
	}
	if err := db.Put("PFKW#_apply", []byte("b")); err != nil {
		t.Fatal(err)
	}
	if err := db.Put("PFKW#_banana", []byte("c")); err != nil {
		t.Fatal(err)
	}

	got, err := db.Get("PFKW#_apple")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if string(got) != "a" {
		t.Fatalf("Get = %q, want %q", got, "a")
	}

	kvs, err := db.Prefix("PFKW#_app")
	if err != nil {
		t.Fatalf("Prefix: %v", err)
	}
	if len(kvs) != 2 {
		t.Fatalf("Prefix returned %d pairs, want 2: %v", len(kvs), kvs)
	}
}

// TestLvDB_UseAfterCloseErrors verifies Close releases the handle and later
// ops return an error (leveldb's ErrClosed) instead of panicking.
func TestLvDB_UseAfterCloseErrors(t *testing.T) {
	db, err := NewLvDB(filepath.Join(t.TempDir(), "close"))
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if _, err := db.Get("anything"); err == nil {
		t.Fatal("Get after Close: want error, got nil")
	}
	if err := db.Put("k", []byte("v")); err == nil {
		t.Fatal("Put after Close: want error, got nil")
	}
}

// TestLvDB_ConcurrentAccess guards the single-handle design (issue #722): the
// old silenceper/pool wrapper could return a nil connection and panic in
// acquire under concurrent load. The shared *leveldb.DB must tolerate
// concurrent Put/Get/Prefix without panicking or losing writes.
func TestLvDB_ConcurrentAccess(t *testing.T) {
	db, err := NewLvDB(filepath.Join(t.TempDir(), "concurrent"))
	if err != nil {
		t.Fatal(err)
	}

	const writers = 4
	const perWriter = 200

	var wg sync.WaitGroup
	for w := 0; w < writers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < perWriter; i++ {
				key := fmt.Sprintf("k%d-%04d", id, i)
				if err := db.Put(key, []byte(key)); err != nil {
					t.Errorf("Put(%s): %v", key, err)
					return
				}
				if v, err := db.Get(key); err != nil {
					t.Errorf("Get(%s): %v", key, err)
					return
				} else if string(v) != key {
					t.Errorf("Get(%s) = %q", key, v)
					return
				}
			}
		}(w)
	}

	// Concurrent prefix scans while writes are in flight.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			if _, err := db.Prefix("k"); err != nil {
				t.Errorf("Prefix: %v", err)
				return
			}
		}
	}()

	wg.Wait()

	kvs, err := db.Prefix("k")
	if err != nil {
		t.Fatalf("final Prefix: %v", err)
	}
	if got, want := len(kvs), writers*perWriter; got != want {
		t.Fatalf("after concurrent writes got %d pairs, want %d", got, want)
	}
}

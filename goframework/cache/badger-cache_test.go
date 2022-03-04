package cache

import (
	"testing"
)

func TestBadgerCache_Has(t *testing.T) {
	err := testBadgerCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("foo")
	if err != nil {
		t.Error(err)
	}
	if inCache {
		t.Error("foo found in cache, and it shouldn't be there")
	}

	_ = testBadgerCache.Set("foo", "bar")
	inCache, err = testBadgerCache.Has("foo")
	if err != nil {
		t.Error(err)
	}
	if !inCache {
		t.Error("foo not found in cache")
	}

	err = testBadgerCache.Forget("foo")
}

func TestBadgerCache_Get(t *testing.T) {
	err := testBadgerCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	x, err := testBadgerCache.Get("foo")
	if err != nil {
		t.Error(err)
	}
	if x != "bar" {
		t.Error("Get should return bar")
	}
}

func TestBadgerCache_Forget(t *testing.T) {
	err := testBadgerCache.Set("alpha", "beta")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Forget("alpha")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("aplha")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("test badger forget(), alhpa should have been deleted")
	}
}

func TestBadgerCache_Empty(t *testing.T) {
	err := testBadgerCache.Set("alpha", "beta")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Empty()
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("aplha")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("test badger Empty(), alhpa should have been deleted")
	}

}

func TestBadgerCache_EmptyByMatch(t *testing.T) {
	err := testBadgerCache.Set("alpha", "beta")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Set("allo", "papa")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Set("me", "you")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.EmptyByMatch("a")

	inCache, err := testBadgerCache.Has("aplha")
	if err != nil {
		t.Error(err)
	}
	if inCache {
		t.Error("test badger EmptyByMatch(), alpha should have been deleted")
	}

	inCache, err = testBadgerCache.Has("allo")
	if err != nil {
		t.Error(err)
	}
	if inCache {
		t.Error("test badger EmptyByMatch(), allo should have been deleted")
	}

	inCache, err = testBadgerCache.Has("me")
	if err != nil {
		t.Error(err)
	}
	if !inCache {
		t.Error("test badger EmptyByMatch(), you should still be here")
	}

}

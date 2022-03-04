package cache

import (
	"testing"
)

func TestRedisCache_Has(t *testing.T) {
	// test the Cache
	err := testRedisCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}
	if inCache {
		t.Error("foo found in cache, but should not be there")
	}

	// test Set()
	err = testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	inCache, err = testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}
	if !inCache {
		t.Error("foo is not found in cache, but should be")
	}
}

func TestRedisCache_Get(t *testing.T) {
	// set a value in the cache
	err := testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	// get the value
	value, err := testRedisCache.Get("foo")
	if err != nil {
		t.Error(err)
	}

	if value != "bar" {
		t.Error("Get do not get the value from the cache")
	}

}

func TestRedisCache_Forget(t *testing.T) {
	// set a value in the cache
	err := testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}

	value, err := testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if value {
		t.Error("The value from the cache has been Forget(), should not be there")
	}

}

func TestRedisCache_Empty(t *testing.T) {
	err := testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Empty()
	if err != nil {
		t.Error(err)
	}

	value, err := testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if value {
		t.Error("The cache has been Empty(), should not have any values")
	}
}

func TestRedisCache_EmptyByMatch(t *testing.T) {
	err := testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Set("alpha", "beta")
	if err != nil {
		t.Error(err)
	}

	// emptyByMatch foo
	err = testRedisCache.EmptyByMatch("foo")
	if err != nil {
		t.Error(err)
	}

	// foo should be gone
	value, err := testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if value {
		t.Error("The cache has been EmptyByMatch(), should not have foo value")
	}

	// alpha should be there
	value, err = testRedisCache.Has("alpha")
	if err != nil {
		t.Error(err)
	}

	if !value {
		t.Error("The cache has been EmptyByMatch(), should have apha values")
	}

}

func TestEncodeDecode(t *testing.T) {
	entry := Entry{}

	entry["foo"] = "bar"

	bytes, err := encode(entry)
	if err != nil {
		t.Error(err)
	}

	_, err = decode(string(bytes))
	if err != nil {
		t.Error(err)
	}
}

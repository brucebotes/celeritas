package cache

import "testing"

func TestRedisCache_Has(t *testing.T) {
	err := testRedisCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo found in cache, and it should not be there")
	}

	err = testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	inCache, err = testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("foo not found in cache, and it should be there")
	}
}

func TestRedisCache_Get(t *testing.T) {
	err := testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	x, err := testRedisCache.Get("foo")
	if err != nil {
		t.Error(err)
	}

	if x != "bar" {
		t.Error("did not get correct value from cache")
	}
}

func testRedisCache_Forget(t *testing.T) {
	err := testRedisCache.Set("alpha", "beta")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Forget("alpha")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("alpha")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("alpha found in cache and it should not be there")
	}
}

func TestRedisCache_Empty(t *testing.T) {
	err := testRedisCache.Set("alpha", "beta")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Empty()
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("alpha")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("alpha found in cache and it should not be there")
	}
}

func TestRedisCache_EmptyByMatch(t *testing.T) {
	err := testRedisCache.Set("alpha", "foo")
	if err != nil {
		t.Error(err)
	}
	err = testRedisCache.Set("alpha1", "beta")
	if err != nil {
		t.Error(err)
	}
	err = testRedisCache.Set("beta", "foo")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.EmptyByMatch("alph")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("alpha")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("alpha found in cache and it should not be there")
	}

	inCache, err = testRedisCache.Has("alpha1")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("alpha1 found in cache and it should not be there")
	}

	inCache, err = testRedisCache.Has("beta")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("beta not found in cache and it should be there")
	}
}

func TestEncoceDecode(t *testing.T) {
	entry := Entry{}
	entry["foo"] = "bar"
	bytes, err := encode(entry)
	if err != nil {
		t.Error(err)
	}
	new_entry, err := decode(string(bytes))
	if err != nil {
		t.Error(err)
	}

	if new_entry["foo"] != entry["foo"] {
		t.Error("expected new entry to be the same as original entry")
	}
}

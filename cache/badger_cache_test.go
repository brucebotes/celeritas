package cache

import "testing"

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
		t.Error("foo found in cache and it should not be there")
	}

	_ = testBadgerCache.Set("foo", "bar")

	inCache, err = testBadgerCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("foo not found in cache and it should be there")
	}

	_ = testBadgerCache.Forget("foo")
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
		t.Error("did not get correct value from cache")
	}
}

func TestBadgerCache_Forget(t *testing.T) {
	err := testRedisCache.Set("foo", 123)
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo found in cache and it should not be there")
	}
}

func TestBadgerCache_Empty(t *testing.T) {
	err := testRedisCache.Set("alpha", 123)
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Empty()
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("alpha")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("alpha found in cache and it should not be there")
	}
}

func TestBadgerCache_EmptyByMatch(t *testing.T) {
	err := testRedisCache.Set("alpha", "beta")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Set("alpha1", "beta")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.EmptyByMatch("a")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("alpha")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("alpha found in cache and it should not be there")
	}

	inCache, err = testBadgerCache.Has("alpha1")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("alpha1 found in cache and it should not be there")
	}

	inCache, err = testBadgerCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo not found in cache and it should be there")
	}
}

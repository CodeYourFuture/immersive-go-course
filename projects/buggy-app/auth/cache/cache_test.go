package cache

import "testing"

type TestValue string

func TestNew(t *testing.T) {
	New[TestValue]()
}

func TestPut(t *testing.T) {
	c := New[TestValue]()
	k := c.Key("foo")
	v := TestValue("entry")
	c.Put(k, &v)
}

func TestGet(t *testing.T) {
	c := New[TestValue]()
	k := c.Key("foo")
	v := TestValue("entry")
	c.Put(k, &v)
	gV, ok := c.Get(k)
	if !ok {
		t.Fatalf("cache: get not ok")
	}
	if v != *gV {
		t.Fatalf("cache: expected %s, got %s", v, *gV)
	}
}

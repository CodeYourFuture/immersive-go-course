package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestBytes(t *testing.T) {
	t.Run("buffer with no bytes", func(t *testing.T) {
		b := NewBufferString("")
		got := b.Bytes()
		want := []byte("")

		assertBuffer(got, want, t)

	})

	t.Run("can write and get bytes using .Bytes()", func(t *testing.T) {
		want := []byte("Hello")

		b := NewBufferString("Hello")

		got := b.Bytes()

		assertBuffer(got, want, t)
	})

	t.Run("writing more bytes to buffer should add more bytes to the buffer", func(t *testing.T) {
		want := []byte("Hello wolrd")

		b := NewBufferString("Hello")

		_, err := b.Write([]byte(" wolrd"))
		if err != nil {
			t.Error("did not want an error but got one")
		}

		got := b.Bytes()

		assertBuffer(got, want, t)
	})

	t.Run("can read some bytes using bytes.Read with slice big enough for whole bytes", func(t *testing.T) {
		b := bytes.NewBufferString("Hello world")

		slice := make([]byte, 50)

		n, err := b.Read(slice)
		if err != nil {
			t.Error("did not want an error but got one")
		}

		if n != 11 {
			t.Errorf("got %d, want %d", n, 11)
		}

		assertBuffer([]byte("Hello world"), slice[:n], t)
	})

	t.Run("if you call bytes.Read with smaller slice, only that size is read and if you call it again the rest are read", func(t *testing.T) {
		b := NewBufferString("Hello world")

		slice := make([]byte, 6)
		n, err := b.Read(slice)
		if err != nil {
			t.Error("did not want an error but got one")
		}

		if n != 6 {
			t.Errorf("got %d, want %d", n, 6)
		}

		assertBuffer([]byte("Hello "), slice[:n], t)

		n, ererr := b.Read(slice)
		if ererr != nil {
			t.Error("did not want an error but got one")
		}

		if n != 5 {
			t.Errorf("got %d, want %d", n, 5)
		}

		assertBuffer([]byte("world"), slice[:n], t)

	})

}

func assertBuffer(got []byte, want []byte, t testing.TB) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

}

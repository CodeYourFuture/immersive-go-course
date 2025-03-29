package interfaces

import (
	"reflect"
	"testing"
    "log"
)
func TestBytesBuffer(t *testing.T) {
    t.Run("testing b.Bytes()", func(t *testing.T) {
        //b := bytes.Buffer{}
        b:= OurByteBuffer{}
        b.Write([]byte("Hello"))

        got := b.Bytes()
        want := []byte("Hello")

        if !reflect.DeepEqual(got, want){
            t.Errorf("got %v want %v", got, want)
        }

        b.Write([]byte(" World"))

        got = b.Bytes()
        want = []byte("Hello World")

        if !reflect.DeepEqual(got, want) {
            t.Errorf("got %v want %v", got, want)
        }
    })
    t.Run("testing b.Read() with a slice big enough", func(t *testing.T) {
        // b := bytes.Buffer{}
        b := OurByteBuffer{}
        b.Write([]byte("12345678"))

        eightBytes := make([]byte, 8)

        _, err := b.Read(eightBytes)
        if err != nil {
            log.Fatal(err)
        }
       
        got := string(eightBytes)
        want := "12345678"

        if got != want {
            t.Errorf("got %s want %s", got, want)
        }
    })
    t.Run("testing b.Read() with a limited slice", func(t *testing.T) {
        //b := bytes.Buffer{}
        b := OurByteBuffer{}
        b.Write([]byte("123456789"))
        
        fourBytes := make([]byte, b.Len() - 5)
        
        // Some of the bytes are read.
        _, err := b.Read(fourBytes)
        if err != nil {
            log.Fatal(err)
        }

        got := string(fourBytes)
        want := "1234"

        if got != want {
            t.Errorf("got %s want %s", got, want)
        }
        
        // The rest of the bytes are read.
        _, err = b.Read(fourBytes)
        if err != nil {
            log.Fatal(err)
        }

        got = string(fourBytes)
        want = "5678"

        if got != want {
            t.Errorf("got %s want %s", got, want)
        }
    })
    
}

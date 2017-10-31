package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/willf/bloom"
)

func TestAdd(t *testing.T) {
	b := bloom.New(100, 2)
	email1 := []byte("apple@uci.edu")
	email2 := []byte("orange@uci.edu")
	email3 := []byte("potato@gmail.com")
	email4 := []byte("sandwich@yahoo.com")
	b.Add(email1)
	b.Add(email2)
	b.Add(email3)
	assert.Equal(t, false, b.Test(email4))
	assert.Equal(t, true, b.Test(email3))
	assert.Equal(t, true, b.Test(email2))
	assert.Equal(t, true, b.Test(email1))
}

func TestAddString(t *testing.T) {
	b := bloom.New(1000, 2)
	email1 := "selena@gmail.com"
	email2 := "alice@uci.edu"
	email3 := "sam@gmail.com"
	b.AddString(email1)
	assert.Equal(t, false, b.TestString(email3))
	b.AddString(email3)
	assert.Equal(t, false, b.TestString(email2))
	assert.Equal(t, true, b.TestString(email3))
	assert.Equal(t, true, b.TestString(email1))
}

func TestTestAndAddString(t *testing.T) {
	b := bloom.New(1000, 2)
	email1 := "jake@gmail.com"
	assert.Equal(t, false, b.TestString(email1))
	assert.Equal(t, false, b.TestAndAddString(email1))
	// should not add the string until after TestAndAddString() is called
	assert.Equal(t, true, b.TestString(email1))
}

func TestCopy(t *testing.T) {
	b := bloom.New(1000, 2)
	b2 := b.Copy()
	b3 := b2.Copy()
	assert.Equal(t, b, b2)
	assert.Equal(t, b3, b2)
	assert.Equal(t, b, b3)
}

func TestEqual(t *testing.T) {
	b := bloom.New(1000, 4)
	b2 := b.Copy()
	assert.Equal(t, true, b.Equal(b2))
}

func TestClearAll(t *testing.T) {
	b := bloom.New(1000, 4)
	email1 := "james@gmail.com"
	email2 := "jessica@uci.edu"
	b.AddString(email1)
	b.AddString(email2)
	b2 := b.Copy()
	assert.Equal(t, true, b.TestString(email1))
	assert.Equal(t, true, b.TestString(email2))
	b.ClearAll()
	assert.Equal(t, false, b.TestString(email1))
	assert.Equal(t, false, b.TestString(email2))
	assert.NotEqual(t, b, b2)
	b2.ClearAll()
	assert.Equal(t, b, b2)
}

func TestMerge(t *testing.T) {
	b := bloom.New(1000, 4)
	bEmail1 := "stella@gmail.com"
	b2Email1 := "padthai@gmail.com"
	b2Email2 := "spagetti@uci.edu"
	b.AddString(bEmail1)
	assert.Equal(t, true, b.TestString(bEmail1))
	assert.Equal(t, false, b.TestString(b2Email1))
	assert.Equal(t, false, b.TestString(b2Email2))
	b2 := bloom.New(1000, 4)
	b2.AddString(b2Email1)
	b2.AddString(b2Email2)
	assert.Equal(t, true, b2.TestString(b2Email1))
	assert.Equal(t, true, b2.TestString(b2Email2))
	b.Merge(b2)
	assert.Equal(t, false, b2.TestString(bEmail1))
	assert.Equal(t, true, b.TestString(bEmail1))
	assert.Equal(t, true, b.TestString(b2Email1))
	assert.Equal(t, true, b.TestString(b2Email2))
}

func TestInvalidMerge(t *testing.T) {
	b := bloom.New(1000, 4)
	bEmail1 := "stella@gmail.com"
	b2Email1 := "padthai@gmail.com"
	b.AddString(bEmail1)
	b2 := bloom.New(500, 2)
	b2.AddString(b2Email1)
	fmt.Println(b2.Merge(b))
}

package main

import (
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

package main

import (
	"testing"

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
	if b.Test(email4) {
		t.Errorf("%v should not be in bloom filter", email4)
	}
	if !b.Test(email3) {
		t.Errorf("%v should be in bloom filter", email3)
	}
	if !b.Test(email2) {
		t.Errorf("%v should be in bloom filter", email2)
	}
	if !b.Test(email1) {
		t.Errorf("%v should be in bloom filter", email1)
	}
}

func TestAddString(t *testing.T) {
	b := bloom.New(1000, 2)
	email1 := "selena@gmail.com"
	email2 := "alice@uci.edu"
	email3 := "sam@gmail.com"
	b.AddString(email1)
	if b.TestString(email3) {
		t.Errorf("%v should not be in bloom filter", email3)
	}
	b.AddString(email3)
	if b.TestString(email2) {
		t.Errorf("%v should not be in bloom filter", email2)
	}
	if !b.TestString(email3) {
		t.Errorf("%v should be in bloom filter", email3)
	}
	if !b.TestString(email1) {
		t.Errorf("%v should be in bloom filter", email1)
	}
}

func TestTestAndAddString(t *testing.T) {
	b := bloom.New(1000, 2)
	email1 := "jake@gmail.com"
	if b.TestString(email1) {
		t.Errorf("%v should not be in bloom filter", email1)
	}
	if b.TestAndAddString(email1) {
		t.Errorf("%v should not be in bloom filter the first time we look", email1)
	}
	if !b.TestString(email1) {
		t.Errorf("%v should be in bloom filter", email1)
	}
}

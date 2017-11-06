package main

import (
	"fmt"
	"testing"
"github.com/stretchr/testify/assert"
"github.com/willf/bloom"
)

func TestBasic(t *testing.T) {
	//fmt.Println("Creating Bloom Filter")
	var bList [10]string
	b := bloom.New(100, 3)

	bList[0] = "1@e.com"
	bList[1] = "2@e.com"
	bList[2] = "3@e.com"
	bList[3] = "4@e.com"
	bList[4] = "5@e.com"
	bList[5] = "6@e.com"
	bList[6] = "7@e.com"
	bList[7] = "8@e.com"
	bList[8] = "9@e.com"
	bList[9] = "10@e.com"

	//fmt.Println("Populating Bloom Filter")
	for i := 0; i < 10; i++ {
		//	fmt.Printf("Adding %v to the  Bloom Filter\n", bList[i])
		b.AddString(bList[i])
	}

	//Test that all 10 items have been added to the bloom filter.
	if len(bList) != 10 {
		t.Error("Emials in bloom filter should be 10.")

	}

	//Testing certain strings
	test0 := b.TestString(bList[0])
	if !test0 {
		t.Errorf("%v should be in.", bList[0])
	}

	test3 := b.TestString(bList[3])
	if !test3 {
		t.Errorf("%v should be in.", bList[3])
	}

	test8 := b.TestString(bList[8])
	if !test8 {
		t.Errorf("%v should be in.", bList[8])
	}

	testnone := b.TestString("11@e.com")
	if testnone {
		t.Errorf("%v should NOT be in.", "11@e.com")
	}

	test6 := b.TestString(bList[6])
	if !test6 {
		t.Errorf("%v should be in.", bList[6])
	}

}

func TestCap(t *testing.T) {
	f := bloom.New(1000, 4)
	if f.Cap() != 1000 {
		t.Error("not accessing Cap() correctly")
	}
}

func TestMerge(t *testing.T) {
	f := bloom.New(1000, 4)
	n1 := []byte("f")
	f.Add(n1)

	g := bloom.New(1000, 4)
	n2 := []byte("g")
	g.Add(n2)

	h := bloom.New(999, 4)
	n3 := []byte("h")
	h.Add(n3)

	j := bloom.New(1000, 5)
	n4 := []byte("j")
	j.Add(n4)

	err := f.Merge(g)
	if err != nil {
		t.Errorf("There should be no error when merging two similar filters")
	}

	err = f.Merge(h)
	if err == nil {
		t.Errorf("There should be an error when merging filters with mismatched m")
	}

	err = f.Merge(j)
	if err == nil {
		t.Errorf("There should be an error when merging filters with mismatched k")
	}

	n2b := f.Test(n2)
	if !n2b {
		t.Errorf("The value doesn't exist after a valid merge")
	}

	n3b := f.Test(n3)
	if n3b {
		t.Errorf("The value exists after an invalid merge")
	}

	n4b := f.Test(n4)
	if n4b {
		t.Errorf("The value exists after an invalid merge")
	}
}

func TestCopy(t *testing.T) {
	f := bloom.New(1000, 4)
	n1 := []byte("f")
	f.Add(n1)

	// copy here instead of New
	g := f.Copy()
	n2 := []byte("g")
	g.Add(n2)

	n1fb := f.Test(n1)
	if !n1fb {
		t.Errorf("The value doesn't exist in original after making a copy")
	}

	n1gb := g.Test(n1)
	if !n1gb {
		t.Errorf("The value doesn't exist in the copy")
	}

	n2fb := f.Test(n2)
	if n2fb {
		t.Errorf("The value exists in the original, it should only exist in copy")
	}

	n2gb := g.Test(n2)
	if !n2gb {
		t.Errorf("The value doesn't exist in copy after Add()")
	}
}

func TestFrom(t *testing.T) {
	var (
		k    = uint(5)
		data = make([]uint64, 10)
		test = []byte("test")
	)

	bf := bloom.From(data, k)
	if bf.K() != k {
		t.Errorf("Constant k does not match the expected value")
	}

	if bf.Cap() != uint(len(data)*64) {
		t.Errorf("Capacity does not match the expected value")
	}

	if bf.Test(test) {
		t.Errorf("Bloom filter should not contain the value")
	}

	bf.Add(test)
	if !bf.Test(test) {
		t.Errorf("Bloom filter should contain the value")
	}

	// create a new Bloom filter from an existing (populated) data slice.
	bf = bloom.From(data, k)
	if !bf.Test(test) {
		t.Errorf("Bloom filter should contain the value")
	}

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
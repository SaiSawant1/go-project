package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathTransformFunc(t *testing.T) {
	key := "momsbetspictures"
	actualPathKey := PathKey{
		Pathname: "74aaf/3b3ae/d3e58/e3754/7d349/4b17f/ef871/7d276",
		Filename: "74aaf3b3aed3e58e37547d3494b17fef8717d276",
	}
	pathKey := CASPathTransformFunc(key)
	assert.Equal(t, actualPathKey.Pathname, pathKey.Pathname)
	assert.Equal(t, actualPathKey.Filename, pathKey.Filename)
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	key := "test"

	s := NewStore(opts)
	data := []byte("some jpg bytes updated\n this is new line")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, err := io.ReadAll(r)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, b, data)

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

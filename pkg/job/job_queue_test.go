package main

import (
	"testing"
	"time"
)

func TestLocalIdempotentHelper(t *testing.T) {
	ih := NewLocalIdempotentHelper(10)

	if ih.IsProduced([]byte("d1")) {
		t.Fatal("Expected no produced when ttl is 10")
	}

	if !ih.IsProduced([]byte("d1")) {
		t.Fatal("Expected produced when ttl is 10")
	}

	ih = NewLocalIdempotentHelper(0)

	if ih.IsProduced([]byte("d1")) {
		t.Fatal("Expected no produced when ttl is 0")
	}

	time.Sleep(time.Millisecond)

	if ih.IsProduced([]byte("d1")) {
		t.Fatal("Expected no produced when ttl is 0 after call IsProduced")
	}
}

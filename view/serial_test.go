package view

import (
	"testing"
)

func TestSerialPrinter_NewSerialPrinter(t *testing.T) {
	sp := NewSerialPrinter()
	if sp == nil {
		t.Error("Expected *SerialPrinter, got nil")
	}

	sp.pinstats.Add(1, 1.2)
	if sp.pinstats.pins[1].stats[3] != 1.2 {
		t.Error("Expected 1.2, got 0")
	}

	sp.pinstats.Add(1, 1.7)
	if sp.pinstats.pins[1].stats[3] != 1.7 {
		t.Error("Expected 1.7, got 0")
	}
	if sp.pinstats.pins[1].stats[2] != 1.2 {
		t.Error("Expected 1.2, got 0")
	}

	sp.pinstats.Add(1, 1.9)
	if sp.pinstats.pins[1].stats[3] != 1.9 {
		t.Error("Expected 1.9, got 0")
	}
	if sp.pinstats.pins[1].stats[2] != 1.7 {
		t.Error("Expected 1.7, got 0")
	}
	if sp.pinstats.pins[1].stats[1] != 1.2 {
		t.Error("Expected 1.2, got 0")
	}

	sp.pinstats.Add(1, 0.9)
	if sp.pinstats.pins[1].stats[3] != 0.9 {
		t.Error("Expected 0.9, got 0")
	}
	if sp.pinstats.pins[1].stats[2] != 1.9 {
		t.Error("Expected 1.9, got 0")
	}
	if sp.pinstats.pins[1].stats[1] != 1.7 {
		t.Error("Expected 1.7, got 0")
	}
	if sp.pinstats.pins[1].stats[0] != 1.2 {
		t.Error("Expected 1.2, got 0")
	}
}

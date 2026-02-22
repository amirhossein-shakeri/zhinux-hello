package greeting

import (
	"errors"
	"strings"
	"testing"
)

func TestNewName_NormalizesAndValidates(t *testing.T) {
	name, err := NewName("  Alice   Bob  ", 64)
	if err != nil {
		t.Fatalf("expected name to be valid, got error: %v", err)
	}
	if got, want := name.String(), "Alice Bob"; got != want {
		t.Fatalf("normalized name mismatch: got %q want %q", got, want)
	}
}

func TestNewName_RejectsEmptyName(t *testing.T) {
	_, err := NewName(" \n\t ", 32)
	if !errors.Is(err, ErrNameEmpty) {
		t.Fatalf("expected ErrNameEmpty, got: %v", err)
	}
}

func TestNewName_RejectsInvalidUTF8(t *testing.T) {
	_, err := NewName(string([]byte{0xff, 0xfe}), 32)
	if !errors.Is(err, ErrInvalidUTF8) {
		t.Fatalf("expected ErrInvalidUTF8, got: %v", err)
	}
}

func TestNewName_RespectsConfigurableMaxLength(t *testing.T) {
	_, err := NewName("Alicia", 5)
	if !errors.Is(err, ErrNameTooLong) {
		t.Fatalf("expected ErrNameTooLong, got: %v", err)
	}
}

func TestNewStreamCount_EnforcesBounds(t *testing.T) {
	_, err := NewStreamCount(0, 1, 5)
	if !errors.Is(err, ErrStreamCountOutOfRange) {
		t.Fatalf("expected ErrStreamCountOutOfRange, got: %v", err)
	}
}

func TestPolicy_ValidateRejectsInvalidBounds(t *testing.T) {
	policy := DefaultPolicy()
	policy.MinStreamCount = 5
	policy.MaxStreamCount = 3

	err := policy.Validate()
	if !errors.Is(err, ErrInvalidPolicy) {
		t.Fatalf("expected ErrInvalidPolicy, got: %v", err)
	}
}

func TestAggregate_SayHelloDeterministic(t *testing.T) {
	policy := DefaultPolicy()
	policy.GreetingTemplate = "Welcome, %s."
	aggregate, err := NewAggregate(policy)
	if err != nil {
		t.Fatalf("unexpected policy validation error: %v", err)
	}

	first, err := aggregate.SayHello("  Alice ")
	if err != nil {
		t.Fatalf("unexpected first greeting error: %v", err)
	}
	second, err := aggregate.SayHello("Alice")
	if err != nil {
		t.Fatalf("unexpected second greeting error: %v", err)
	}

	if first != second {
		t.Fatalf("expected deterministic output, got %q and %q", first, second)
	}
	if got, want := first, "Welcome, Alice."; got != want {
		t.Fatalf("greeting mismatch: got %q want %q", got, want)
	}
}

func TestAggregate_StreamGreetings_RespectsPolicyBounds(t *testing.T) {
	policy := DefaultPolicy()
	policy.MinStreamCount = 1
	policy.MaxStreamCount = 3
	aggregate, err := NewAggregate(policy)
	if err != nil {
		t.Fatalf("unexpected policy validation error: %v", err)
	}

	greetings, err := aggregate.StreamGreetings("Alice", 3)
	if err != nil {
		t.Fatalf("unexpected stream greeting error: %v", err)
	}
	if got, want := len(greetings), 3; got != want {
		t.Fatalf("greetings length mismatch: got %d want %d", got, want)
	}

	_, err = aggregate.StreamGreetings("Alice", 4)
	if !errors.Is(err, ErrStreamCountOutOfRange) {
		t.Fatalf("expected ErrStreamCountOutOfRange, got: %v", err)
	}
}

func TestAggregate_NormalizeMessage(t *testing.T) {
	policy := DefaultPolicy()
	policy.MaxMessageLength = 10
	aggregate, err := NewAggregate(policy)
	if err != nil {
		t.Fatalf("unexpected policy validation error: %v", err)
	}

	message, err := aggregate.NormalizeMessage("  hi   team  ")
	if err != nil {
		t.Fatalf("unexpected normalize error: %v", err)
	}
	if got, want := message.String(), "hi team"; got != want {
		t.Fatalf("normalized message mismatch: got %q want %q", got, want)
	}

	_, err = aggregate.NormalizeMessage(strings.Repeat("a", 11))
	if !errors.Is(err, ErrMessageTooLong) {
		t.Fatalf("expected ErrMessageTooLong, got: %v", err)
	}
}

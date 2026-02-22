package greeting

// Aggregate encapsulates greeting-related business rules.
type Aggregate struct {
	policy Policy
}

// NewAggregate builds an aggregate with validated policy.
func NewAggregate(policy Policy) (Aggregate, error) {
	normalized := policy.withDefaults()
	if err := normalized.Validate(); err != nil {
		return Aggregate{}, err
	}

	return Aggregate{policy: normalized}, nil
}

// NewDefaultAggregate creates an aggregate with baseline policy values.
func NewDefaultAggregate() Aggregate {
	aggregate, err := NewAggregate(DefaultPolicy())
	if err != nil {
		// DefaultPolicy is internally controlled; panic keeps bad defaults obvious.
		panic(err)
	}

	return aggregate
}

func (a Aggregate) Policy() Policy {
	return a.policy
}

// SayHello validates and normalizes input, then composes a deterministic greeting.
func (a Aggregate) SayHello(rawName string) (string, error) {
	name, err := NewName(rawName, a.policy.MaxNameLength)
	if err != nil {
		return "", err
	}

	return a.policy.ComposeGreeting(name), nil
}

// StreamGreetings enforces stream bounds and returns deterministic greeting payloads.
func (a Aggregate) StreamGreetings(rawName string, rawCount int) ([]string, error) {
	name, err := NewName(rawName, a.policy.MaxNameLength)
	if err != nil {
		return nil, err
	}

	count, err := NewStreamCount(rawCount, a.policy.MinStreamCount, a.policy.MaxStreamCount)
	if err != nil {
		return nil, err
	}

	message := a.policy.ComposeGreeting(name)
	greetings := make([]string, count.Int())
	for i := 0; i < len(greetings); i++ {
		greetings[i] = message
	}

	return greetings, nil
}

// NormalizeMessage enforces message normalization and policy-driven limits.
func (a Aggregate) NormalizeMessage(rawMessage string) (Message, error) {
	return NewMessage(rawMessage, a.policy.MaxMessageLength)
}

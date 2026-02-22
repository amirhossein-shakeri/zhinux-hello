package greeting

import (
	"fmt"
	"strings"
)

const (
	defaultMaxNameLength    = 1 << 7 // 128
	defaultMaxMessageLength = 1 << 9 // 512
	defaultMinStreamCount   = 1
	defaultMaxStreamCount   = 1 << 10 // 1024
	defaultGreetingTemplate = "Hello, %s!"
)

// Policy centralizes configurable business limits and deterministic output rules.
type Policy struct {
	MaxNameLength    int
	MaxMessageLength int
	MinStreamCount   int
	MaxStreamCount   int
	GreetingTemplate string
}

func DefaultPolicy() Policy {
	return Policy{
		MaxNameLength:    defaultMaxNameLength,
		MaxMessageLength: defaultMaxMessageLength,
		MinStreamCount:   defaultMinStreamCount,
		MaxStreamCount:   defaultMaxStreamCount,
		GreetingTemplate: defaultGreetingTemplate,
	}
}

func (p Policy) withDefaults() Policy {
	normalized := p
	if normalized.MaxNameLength <= 0 {
		normalized.MaxNameLength = defaultMaxNameLength
	}
	if normalized.MaxMessageLength <= 0 {
		normalized.MaxMessageLength = defaultMaxMessageLength
	}
	if normalized.MinStreamCount <= 0 {
		normalized.MinStreamCount = defaultMinStreamCount
	}
	if normalized.MaxStreamCount <= 0 {
		normalized.MaxStreamCount = defaultMaxStreamCount
	}
	if normalized.GreetingTemplate == "" {
		normalized.GreetingTemplate = defaultGreetingTemplate
	}

	return normalized
}

func (p Policy) Validate() error {
	normalized := p.withDefaults()

	if normalized.MaxNameLength <= 0 {
		return fmt.Errorf("%w: max name length must be greater than zero", ErrInvalidPolicy)
	}
	if normalized.MaxMessageLength <= 0 {
		return fmt.Errorf("%w: max message length must be greater than zero", ErrInvalidPolicy)
	}
	if normalized.MinStreamCount <= 0 {
		return fmt.Errorf("%w: min stream count must be greater than zero", ErrInvalidPolicy)
	}
	if normalized.MaxStreamCount < normalized.MinStreamCount {
		return fmt.Errorf("%w: max stream count must be greater than or equal to min stream count", ErrInvalidPolicy)
	}
	if !strings.Contains(normalized.GreetingTemplate, "%s") {
		return fmt.Errorf("%w: greeting template must include %%s placeholder", ErrInvalidPolicy)
	}

	return nil
}

// ComposeGreeting is deterministic for identical name + policy input.
func (p Policy) ComposeGreeting(name Name) string {
	normalized := p.withDefaults()
	return fmt.Sprintf(normalized.GreetingTemplate, name.String())
}

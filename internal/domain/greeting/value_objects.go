package greeting

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Name struct {
	value string
}

func NewName(rawName string, maxLength int) (Name, error) {
	if maxLength <= 0 {
		return Name{}, fmt.Errorf("%w: max name length must be greater than zero", ErrInvalidPolicy)
	}

	normalized, err := normalizeText(rawName)
	if err != nil {
		return Name{}, err
	}
	if normalized == "" {
		return Name{}, ErrNameEmpty
	}
	if utf8.RuneCountInString(normalized) > maxLength {
		return Name{}, fmt.Errorf("%w: got=%d max=%d", ErrNameTooLong, utf8.RuneCountInString(normalized), maxLength)
	}

	return Name{value: normalized}, nil
}

func (n Name) String() string {
	return n.value
}

type Message struct {
	value string
}

func NewMessage(rawMessage string, maxLength int) (Message, error) {
	if maxLength <= 0 {
		return Message{}, fmt.Errorf("%w: max message length must be greater than zero", ErrInvalidPolicy)
	}

	normalized, err := normalizeText(rawMessage)
	if err != nil {
		return Message{}, err
	}
	if normalized == "" {
		return Message{}, ErrMessageEmpty
	}
	if utf8.RuneCountInString(normalized) > maxLength {
		return Message{}, fmt.Errorf("%w: got=%d max=%d", ErrMessageTooLong, utf8.RuneCountInString(normalized), maxLength)
	}

	return Message{value: normalized}, nil
}

func (m Message) String() string {
	return m.value
}

type StreamCount struct {
	value int
}

func NewStreamCount(rawCount, minCount, maxCount int) (StreamCount, error) {
	if rawCount < minCount || rawCount > maxCount {
		return StreamCount{}, fmt.Errorf(
			"%w: got=%d allowed=[%d,%d]",
			ErrStreamCountOutOfRange,
			rawCount,
			minCount,
			maxCount,
		)
	}

	return StreamCount{value: rawCount}, nil
}

func (c StreamCount) Int() int {
	return c.value
}

func normalizeText(raw string) (string, error) {
	if !utf8.ValidString(raw) {
		return "", ErrInvalidUTF8
	}

	normalized := strings.Join(strings.Fields(strings.TrimSpace(raw)), " ")
	return normalized, nil
}

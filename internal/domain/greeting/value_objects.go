package greeting

import (
	"errors"
	"fmt"

	platformvalidation "github.com/amirhossein-shakeri/zhinux-platform/validation"
)

type Name struct {
	value string
}

func NewName(rawName string, maxLength int) (Name, error) {
	if maxLength <= 0 {
		return Name{}, fmt.Errorf("%w: max name length must be greater than zero", ErrInvalidPolicy)
	}

	normalized, err := platformvalidation.NormalizeText(rawName)
	if err != nil {
		if errors.Is(err, platformvalidation.ErrInvalidUTF8) {
			return Name{}, ErrInvalidUTF8
		}
		return Name{}, err
	}
	if normalized == "" {
		return Name{}, ErrNameEmpty
	}
	if !platformvalidation.RuneCountAtMost(normalized, maxLength) {
		return Name{}, fmt.Errorf("%w: got=%d max=%d", ErrNameTooLong, len([]rune(normalized)), maxLength)
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

	normalized, err := platformvalidation.NormalizeText(rawMessage)
	if err != nil {
		if errors.Is(err, platformvalidation.ErrInvalidUTF8) {
			return Message{}, ErrInvalidUTF8
		}
		return Message{}, err
	}
	if normalized == "" {
		return Message{}, ErrMessageEmpty
	}
	if !platformvalidation.RuneCountAtMost(normalized, maxLength) {
		return Message{}, fmt.Errorf("%w: got=%d max=%d", ErrMessageTooLong, len([]rune(normalized)), maxLength)
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
	if !platformvalidation.IntInRange(rawCount, minCount, maxCount) {
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

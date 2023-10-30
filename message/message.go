package message

import (
	"errors"
	"fmt"
)

type Command byte

const MAGIC_BYTE = 0x27
const MSG_SEPARATOR = byte(44)
const MSG_END = byte(0)
const (
	Get Command = 1 << iota
	Set
	Delete
)

type Message struct {
	kind            Command
	numberArguments int
	arguments       [][]byte
}

// / [magic byte] [command kind one bytes] [number arguments] [arg 1] [seperator] [arg n]\0
func New(rawBytes []byte) (*Message, error) {

	if len(rawBytes) > 1024 {
		return nil, errors.New("Message too large")
	}

	if !(len(rawBytes) >= 3) {
		return nil, errors.New("Message too small")
	}

	if rawBytes[0] != MAGIC_BYTE {
		return nil, errors.New("Unknown message received")
	}

	command := Command(rawBytes[1])

	switch command {
	case Get:
	case Set:
	case Delete:
	default:
		return nil, errors.New(fmt.Sprintf("Unknown command %d", command))
	}

	numberArguments := int(rawBytes[2])

	arguments := make([][]byte, numberArguments)

	byteIndex := 3

	for i := 0; i < numberArguments; i++ {
		currentArgument := make([]byte, 0)

		for byteIndex < len(rawBytes) {
			if rawBytes[byteIndex] == MSG_SEPARATOR || rawBytes[byteIndex] == MSG_END {
				break
			}
			currentArgument = append(currentArgument, rawBytes[byteIndex])
			byteIndex++
		}

		arguments[i] = currentArgument
		byteIndex++
	}

	return &Message{
		kind:            command,
		numberArguments: int(numberArguments),
		arguments:       arguments,
	}, nil
}

func (m Message) GetKind() Command {
	return m.kind
}

func (m Message) GetArgument(index int) ([]byte, error) {
	if index-1 >= len(m.arguments) {
		return nil, errors.New("Index out of range ")
	}

	return m.arguments[index-1], nil

}

package memory

import (
	"errors"
)

var ErrMemoryType error = errors.New("Memory type does not exist")

// Memory manages memory for all users.
type Memory interface {
	// For returns memory of corresponding user
	For(string) Store
	Delete(string)
}

// Store is a key/value data store
type Store interface {
	Set(string, string) // Save data for the user by key, value
	Get(string) string  // Get data of the user that specified by key
	Delete(string)      // Delete data of the user that specified by key
}

func New(name string) Memory {
	switch name {
	case "ephemeral":
		return newEphemeralMemory()
	default:
		return newEphemeralMemory()
	}
}

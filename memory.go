package fbbot

// Memory is a key/value data store
type Memory interface {
	Set(User, string, string) // Save data for the user bye key, value
	Get(User, string) (string, bool) // Get data of the user that specified by key
	Delete(User, string) // Delete data of the user that specified by key
	DeleteAll(User) // Delete all data of the user
}

// EphemeralMemory is a memory that stores data in RAM
// Just use it for development
type EphemeralMemory struct {
	store map[string](*map[string]string) // user_id -> (key->value)
}

func NewEphemeralMemory() EphemeralMemory {
	s := make(map[string](*map[string]string))
	return EphemeralMemory{store: s}
}

func (m EphemeralMemory) Set(u User, key string, value string) {
	if _, ok := m.store[u.ID]; !ok {
		u_store := make(map[string]string)
		m.store[u.ID] = &u_store
	}
	u_store := *(m.store[u.ID])
	u_store[key] = value
}

func (m EphemeralMemory) Get(u User, key string) (string, bool) {
	if _, ok := m.store[u.ID]; !ok {
		return "", false
	}
	u_store := *(m.store[u.ID])
	value, ok := u_store[key]
	return value, ok
}

func (m EphemeralMemory) Delete(u User, key string) {
	if _, ok := m.store[u.ID]; !ok {
		return
	}
	u_store := *(m.store[u.ID])
	delete(u_store, key)
}

func (m EphemeralMemory) DeleteAll(u User) {
	delete(m.store, u.ID)
}
package storage

type Storage interface {
	// Length returns an integer representing the number of data items stored in the Storage object.
	Length() int
	// Key will return the name of the nth key in the storage.
	Key(ind int) string
	// GetItem will return that key's value.
	GetItem(key string) (string, bool)
	// SetItem will add that key to the storage, or update that key's value if it already exists.
	SetItem(key, val string)
	// RemoveItem will remove that key from the storage.
	RemoveItem(key string)
	// Clear will empty all keys out of the storage.
	Clear()
}

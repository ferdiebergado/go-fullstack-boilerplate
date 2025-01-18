package session

type Manager interface {
	// Saves the session key with its data.
	Save(sessionKey, sessionData string) error

	// Retrieve the session specified by the key
	Session(string) (string, error)

	// Delete a session with the given key.
	DeleteSession(string) error
}

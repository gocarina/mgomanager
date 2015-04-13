package mgomanager

import "gopkg.in/mgo.v2"

// Manager represents the mgo.Session manager.
type Manager struct {
	masterSession    *mgo.Session
	sessionAvailable chan *mgo.Session
}

// SessionAvailable returns the channel of sessions available.
func (m *Manager) SessionAvailable() chan *mgo.Session {
	return m.sessionAvailable
}

// CreateManager creates a new Manager with a given session and a maximum number of sessions that
// can be copied from the passed session.
func CreateManager(masterSession *mgo.Session, maxSessions int) *Manager {
	return &Manager{
		masterSession:    masterSession,
		sessionAvailable: make(chan *mgo.Session, maxSessions),
	}
}

// GetSession gets an available session or copies the master session if none is available.
func (m *Manager) GetSession() (session *mgo.Session) {
	select {
	case session = <-m.sessionAvailable:
		// We got a session, so we reuse one.
	default:
		// We create a new session from the master session.
		session = m.masterSession.Copy()
	}
	return
}

// Recycle the session. If the pool is full, it discards it.
func (m *Manager) Recycle(session *mgo.Session) {
	session.Refresh()
	select {
	case m.sessionAvailable <- session:
	// We put it in the session
	default:
		// The pool is full, we closes the session.
		session.Close()
	}
}

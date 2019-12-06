package log

// KV is a helper type to
type KV map[string]interface{}

// Logger knows how to log.
type Logger interface {
	WithValues(d map[string]interface{}) Logger
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type dummy int

// Dummy is a dummy logger.
const Dummy = dummy(0)

func (d dummy) WithValues(map[string]interface{}) Logger { return d }
func (dummy) Infof(string, ...interface{})               {}
func (dummy) Warningf(string, ...interface{})            {}
func (dummy) Errorf(string, ...interface{})              {}
func (dummy) Debugf(format string, args ...interface{})  {}

package stream_window

import (
	"backend/internal/common/dto"
	"strings"
)

type iStreamWindowOptions interface {
	String() string
}

// StreamWindowType represents stream window types
type StreamWindowType struct {
	Name    string
	Options iStreamWindowOptions
}

var allStreamWindowTypes = map[string]StreamWindowType{
	Session.Name:     Session,
	Interval.Name:    Interval,
	EventWindow.Name: EventWindow,
	StateWindow.Name: StateWindow,
	CountWindow.Name: CountWindow,
}

var (
	Session StreamWindowType = StreamWindowType{
		Name: "SESSION", Options: &dto.StreamWindowOptionsSession{},
	}
	Interval StreamWindowType = StreamWindowType{
		Name: "INTERVAL", Options: &dto.StreamWindowOptionsInterval{},
	}
	EventWindow StreamWindowType = StreamWindowType{
		Name: "EVENT_WINDOW", Options: &dto.StreamWindowOptionsEventWindow{},
	}
	StateWindow StreamWindowType = StreamWindowType{
		Name: "STATE_WINDOW", Options: &dto.StreamWindowOptionsStateWindow{},
	}
	CountWindow StreamWindowType = StreamWindowType{
		Name: "COUNT_WINDOW", Options: &dto.StreamWindowOptionsCountWindow{},
	}
)

// String returns the string representation
func (swt StreamWindowType) String() string {
	return swt.Name
}

// Of converts string to StreamWindowType
func Of(name string) (StreamWindowType, bool) {
	upperName := strings.ToUpper(name)
	swt, ok := allStreamWindowTypes[upperName]
	if ok {
		return swt, true
	}
	return StreamWindowType{}, false
}

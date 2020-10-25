package k8s

import (
	"time"
)

// Msg message struct
type Msg struct {
	Type    string
	Content map[string]string
}

// ExceptMsg expection struct
type ExceptMsg struct {
	Start   time.Time
	Message string
	End     time.Time
	Service string
}

package shimutil

import (
	"encoding/json"
	"errors"
	"time"
)

type TerminationStatus struct {
	raw attrs
}

type attrs struct {
	At       time.Time `json:"at"`
	ExitCode int32     `json:"exitCode"`
	Signal   int32     `json:"signal"`
	Reason   string    `json:"reason"`
}

const (
	reasonExited   string = "exited"
	reasonSignaled string = "signaled"
)

func ParseExitFile(bytes []byte) (*TerminationStatus, error) {
	raw := attrs{}
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return nil, err
	}
	if raw.Reason != reasonExited && raw.Reason != reasonSignaled {
		return nil, errors.New("Unexpected termination reason")
	}
	if raw.Reason == reasonExited && (raw.ExitCode < 0 || raw.ExitCode > 127) {
		return nil, errors.New("Unexpected exit code")
	}
	if raw.Reason == reasonSignaled && raw.Signal <= 0 {
		return nil, errors.New("Unexpected signal")
	}
	return &TerminationStatus{raw}, nil
}

func (t *TerminationStatus) IsSignaled() bool {
	return t.raw.Reason == reasonSignaled
}

func (t *TerminationStatus) At() time.Time {
	return t.raw.At
}

func (t *TerminationStatus) ExitCode() int32 {
	if t.IsSignaled() {
		panic("ExitCode() should not be used when " +
			"container terminated has been killed")
	}
	return t.raw.ExitCode
}

func (t *TerminationStatus) Signal() int32 {
	if !t.IsSignaled() {
		panic("Signal() should not be used when " +
			"container exited normally")
	}
	return t.raw.Signal
}

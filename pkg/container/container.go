package container

import (
	"encoding/json"
	"errors"
	"time"
)

type Container struct {
	impl
}

type impl struct {
	ID_     ID     `json:"id"`
	Name_   string `json:"name"`
	Status_ Status `json:"status"`

	CreatedAt_  string `json:"createdAt"`
	StartedAt_  string `json:"startedAt,omitempty"`
	FinishedAt_ string `json:"finishedAt,omitempty"`

	Command_ []string `json:"command"`
	Args_    []string `json:"args,omitempty"`

	Rootfs_ string `json:"rootfs"`
}

func New(
	id ID,
	name string,
) (*Container, error) {
	if !isValidName(name) {
		return nil, errors.New("Invalid container name")
	}

	return &Container{
		impl{
			ID_:        id,
			Name_:      name,
			CreatedAt_: time.Now().Format(time.RFC3339),
		},
	}, nil
}

func (c *Container) ID() ID {
	return c.ID_
}

func (c *Container) Name() string {
	return c.Name_
}

func (c *Container) CreatedAt() string {
	return c.CreatedAt_
}

func (c *Container) CreatedAtNano() int64 {
	return unixNanoTime(c.CreatedAt())
}

func (c *Container) StartedAt() string {
	return c.StartedAt_
}

func (c *Container) StartedAtNano() int64 {
	if c.StartedAt_ == "" {
		return 0
	}
	return unixNanoTime(c.StartedAt())
}

func (c *Container) FinishedAt() string {
	return c.FinishedAt_
}

func (c *Container) FinishedAtNano() int64 {
	if c.FinishedAt_ == "" {
		return 0
	}
	return unixNanoTime(c.FinishedAt())
}

func (c *Container) Status() Status {
	return c.Status_
}

func (c *Container) SetStatus(s Status) {
	c.Status_ = s
}

func (c *Container) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.impl)
}

func (c *Container) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &c.impl)
}

func isValidName(name string) bool {
	for _, c := range name {
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') && (c != '_') {
			return false
		}
	}
	return len(name) > 0 && len(name) <= 32
}

func unixNanoTime(s string) int64 {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t.UnixNano()
}

package flags

import (
	"errors"
	"fmt"
	"strconv"
)

// CustomInt is a custom flag type that tracks whether it was set.
type CustomInt struct {
	Value int
	IsSet bool
}

func (c *CustomInt) Set(s string) error {
	var err error
	c.Value, err = strconv.Atoi(s)
	if err == nil {
		c.IsSet = true
	}
	return err
}

func (c *CustomInt) String() string {
	return fmt.Sprintf("%d", c.Value)
}

// CustomString is a custom flag type that tracks whether it was set.
type CustomString struct {
	Value string
	IsSet bool
}

func (c *CustomString) Set(s string) error {
	c.Value = s
	c.IsSet = true
	return nil
}

func (c *CustomString) String() string {
	return c.Value
}

// CustomBool is a custom flag type that tracks whether it was set.
type CustomBool struct {
	Value bool
	IsSet bool
}

func (c *CustomBool) Set(s string) error {
	switch s {
	case "true":
		c.IsSet = true
		c.Value = true
	case "false":
		c.Value = false
		c.IsSet = true
	default:
		return errors.New("invalid boolean flag")
	}

	return nil
}

func (c *CustomBool) String() string {
	return fmt.Sprintf("%v", c.Value)
}

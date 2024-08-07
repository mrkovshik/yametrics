// Package flags provides custom flag types that track whether they were set.
package flags

import (
	"errors"
	"fmt"
	"strconv"
)

// CustomInt is a custom flag type that tracks whether it was set.
type CustomInt struct {
	Value int  // The integer value of the flag.
	IsSet bool // Indicates if the flag was set.
}

// Set assigns an integer value to the CustomInt flag from a string.
func (c *CustomInt) Set(s string) error {
	var err error
	c.Value, err = strconv.Atoi(s)
	if err == nil {
		c.IsSet = true
	}
	return err
}

// String returns the string representation of the CustomInt value.
func (c *CustomInt) String() string {
	return fmt.Sprintf("%d", c.Value)
}

// CustomString is a custom flag type that tracks whether it was set.
type CustomString struct {
	Value string // The string value of the flag.
	IsSet bool   // Indicates if the flag was set.
}

// Set assigns a string value to the CustomString flag.
func (c *CustomString) Set(s string) error {
	c.Value = s
	c.IsSet = true
	return nil
}

// String returns the string representation of the CustomString value.
func (c *CustomString) String() string {
	return c.Value
}

// CustomBool is a custom flag type that tracks whether it was set.
type CustomBool struct {
	Value bool // The boolean value of the flag.
	IsSet bool // Indicates if the flag was set.
}

// Set assigns a boolean value to the CustomBool flag from a string.
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

// String returns the string representation of the CustomBool value.
func (c *CustomBool) String() string {
	return fmt.Sprintf("%v", c.Value)
}

package util

import (
	"fmt"
	"strings"
)

type Panic2Err struct {
	err error
}

func (p *Panic2Err) Recover() {

}

func (p *Panic2Err) Check() error {
	return nil
}

type ChainError struct {
	errors []string
}

func NewChainError() *ChainError {
	return &ChainError{
		errors: nil,
	}
}

func (c *ChainError) Append(err error) *ChainError {
	c.errors = append(c.errors, err.Error())
	return c
}

func (c *ChainError) Recover() {
	if err := recover(); err != nil {
		c.Append(fmt.Errorf("recovery: %v", err))
	}
}

func (c *ChainError) Result() error {
	if c.errors != nil {
		return fmt.Errorf(strings.Join(c.errors, ": "))
	}

	return nil
}

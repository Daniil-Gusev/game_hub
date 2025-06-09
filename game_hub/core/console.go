package core

import (
	"fmt"
	"github.com/chzyer/readline"
	"io"
	"os"
	"strings"
)

type Console interface {
	Read() (string, error)
	Write(string) error
	Close() error
}

type ReadlineConsole struct {
	rl *readline.Instance
}

func NewReadlineConsole() (*ReadlineConsole, error) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "> ",
		HistoryLimit: 200,
	})
	if err != nil {
		return nil, err
	}
	return &ReadlineConsole{rl: rl}, nil
}

func (c *ReadlineConsole) Read() (string, error) {
	line, err := c.rl.Readline()
	if err == readline.ErrInterrupt {
		return "", NewAppError(ErrEOF, "interrupt", nil)
	}
	if err == io.EOF {
		return "", NewAppError(ErrEOF, "interrupt", nil)
	}
	if err != nil {
		return "", NewAppError(ErrInternal, "read_error", map[string]any{
			"error": fmt.Sprintf("%v", err),
		})
	}
	return strings.TrimSpace(line), nil
}

func (c *ReadlineConsole) Write(s string) error {
	if _, err := os.Stdout.WriteString(s); err != nil {
		return err
	}
	return nil
}

func (c *ReadlineConsole) Close() error {
	return c.rl.Close()
}

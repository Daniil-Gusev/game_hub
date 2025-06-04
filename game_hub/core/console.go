package core

import (
	"bufio"
	"fmt"
	"github.com/chzyer/readline"
	"io"
	"os"
	"strings"
)

// интерфейс ввода-вывода
type Console interface {
	Read() (string, error)
	Write(string) error
	Close() error
}

// стандартный терминал
type StdConsole struct {
	Reader *bufio.Scanner
}

func NewStdConsole() (*StdConsole, error) {
	return &StdConsole{Reader: bufio.NewScanner(os.Stdin)}, nil
}

func (c *StdConsole) Read() (string, error) {
	if c.Reader.Scan() {
		return strings.TrimSpace(c.Reader.Text()), nil
	}
	err := c.Reader.Err()
	if err == io.EOF {
		return "", NewAppError(ErrEOF, "interrupt", nil)
	}
	if err != nil {
		return "", NewAppError(ErrInternal, "read_error", map[string]any{
			"error": fmt.Sprintf("%v", err),
		})
	}
	return "", NewAppError(ErrEOF, "interrupt", nil)
}

func (c *StdConsole) Write(s string) error {
	if _, err := os.Stdout.WriteString(s); err != nil {
		return err
	}
	return nil
}

type StdReadlineConsole struct {
	rl *readline.Instance
}

func NewStdReadlineConsole() (*StdReadlineConsole, error) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "> ",
		HistoryLimit: 200,
	})
	if err != nil {
		return nil, NewAppError(ErrInit, "console_initialization_error", map[string]any{
			"error": fmt.Sprintf("%v", err),
		})
	}
	return &StdReadlineConsole{rl: rl}, nil
}

func (c *StdReadlineConsole) Read() (string, error) {
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

func (c *StdReadlineConsole) Write(s string) error {
	if _, err := os.Stdout.WriteString(s); err != nil {
		return err
	}
	return nil
}

func (c *StdReadlineConsole) Close() error {
	return c.rl.Close()
}

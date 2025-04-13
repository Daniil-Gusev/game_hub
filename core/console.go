package core

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// интерфейс ввода-вывода
type Console interface {
	Read() (string, error)
	Write(string)
}

// стандартный терминал
type StdConsole struct {
	Reader *bufio.Scanner
}

func (c StdConsole) Read() (string, error) {
	if c.Reader.Scan() {
		return strings.TrimSpace(c.Reader.Text()), nil
	}
	err := c.Reader.Err()
	if err == io.EOF {
		return "", NewAppError(ErrEOF, "eof", nil)
	}
	if err != nil {
		return "", NewAppError(ErrInternal, "read_error", map[string]any{
			"error": fmt.Sprintf("%v", err),
		})
	}
	return "", NewAppError(ErrEOF, "eof", nil)
}
func (c StdConsole) Write(s string) {
	os.Stdout.WriteString(s)
}
func NewStdConsole() *StdConsole {
	return &StdConsole{Reader: bufio.NewScanner(os.Stdin)}
}

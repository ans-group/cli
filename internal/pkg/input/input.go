package input

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ans-group/cli/internal/pkg/output"
)

var InputReader = func() io.Reader {
	return os.Stdin
}

func ReadInput(name string) (string, error) {
	buf := bytes.Buffer{}

	scanner := bufio.NewScanner(InputReader())
	output.Errorf("enter %s (single dot '.' on new line indicates EOF):", name)
	for scanner.Scan() {
		if scanner.Text() == "." {
			break
		}
		buf.WriteString(scanner.Text() + "\n")
	}

	err := scanner.Err()
	if err != nil {
		return "", fmt.Errorf("error reading %s from stdin input: %s", name, err)
	}

	return strings.TrimSuffix(buf.String(), "\n"), nil
}

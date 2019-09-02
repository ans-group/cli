package input

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/ukfast/cli/internal/pkg/output"
)

func ReadInput(name string) (string, error) {
	buf := bytes.Buffer{}

	scanner := bufio.NewScanner(os.Stdin)
	output.Errorf("Enter %s (single dot '.' on new line indicates EOF):", name)
	for scanner.Scan() {
		if scanner.Text() == "." {
			break
		}
		buf.WriteString(scanner.Text() + "\n")
	}

	err := scanner.Err()
	if err != nil {
		return "", fmt.Errorf("Error reading %s from stdin input: %s", name, err)
	}

	return strings.TrimSuffix(buf.String(), "\n"), nil
}

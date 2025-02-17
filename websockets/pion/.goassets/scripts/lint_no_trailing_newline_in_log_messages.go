package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	parentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//- \.:    Matches a literal dot (.).
	//         This is used to match the dot in function calls like log.Trace, log.Debug, etc.
	//- \(":   Matches a literal opening parenthesis ( followed by a double quote ".
	//- [^"]*: Matches zero or more characters that are not a double quote (").
	//         This is used to match the content of the log message (everything inside the quotes).
	//- \\n:   Matches a literal newline character (\n).
	//- "\)?:  Matches a closing double quote (") followed by an optional closing parenthesis ()).
	//         The ? makes the closing parenthesis optional, which accounts for cases where
	//         the log function might not have additional arguments.
	logReg, err := regexp.Compile(`\.(Trace|Debug|Info|Warn|Error)f?\("[^"]*\\n"\)?`)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = filepath.Walk(parentDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) != ".go" {
			return err
		}

		fileContent, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		for _, match := range logReg.FindAll(fileContent, -1) {
			if !strings.Contains(string(match), "nolint") {
				return fmt.Errorf("Log format strings should not have trailing new-line: %s", match)
			}
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

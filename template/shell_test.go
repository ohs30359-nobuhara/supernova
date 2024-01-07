package template

import (
	"testing"
)

func TestShellTemplate_Run(t *testing.T) {
	script := `#!/bin/sh
echo "Hello, Shell!"`

	shellTemplate := ShellTemplate{
		Script: script,
		Dir:    nil, // Optional: Value directory if needed
	}

	out, err := shellTemplate.execScript()
	if err != nil {
		t.Errorf("Error executing ShellTemplate: %v", err)
	}

	if string(out) != "Hello, Shell!\n" {
		t.Errorf("output error")
	}
}

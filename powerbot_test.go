package main

import "testing"

type testcase struct {
	input    string
	command  string
	argument string
}

var tests = []testcase{
	{"powerbot: code 1234", "code", "1234"},
}

func TestExtractCommandAndArgument(t *testing.T) {
	for _, tc := range tests {
		command, argument := ExtractCommandAndArgument(tc.input)
		if command != tc.command {
			t.Error("Command parsing failed, expected: ", tc.command, "  got: ", command)
		} else if argument != tc.argument {
			t.Error("Argument parsing failed, expected: ", tc.argument, "  got: ", argument)
		}
	}

}

package migrate

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var lastMigrationCall string

type mockMigration struct{}

func (m *mockMigration) Up() error {
	lastMigrationCall = "Up"

	return nil
}

func (m *mockMigration) Down() error {
	lastMigrationCall = "Down"

	return nil
}

func (m *mockMigration) Drop() error {
	lastMigrationCall = "Drop"

	return nil
}

func (m *mockMigration) Force(v int) error {
	lastMigrationCall = "Force"

	return nil
}

func (m *mockMigration) Version() (uint, bool, error) {
	lastMigrationCall = "Version"

	return 0, false, nil
}

func (m *mockMigration) Close() (error, error) {
	lastMigrationCall = "Close"

	return nil, nil
}

func TestExecuteOption(t *testing.T) {
	tests := []struct {
		option           string
		userInput        string
		instance         *mockMigration
		expectedFuncCall string
		expectedError    error
	}{
		{
			option:           optionUp,
			userInput:        "",
			instance:         &mockMigration{},
			expectedFuncCall: "Up",
			expectedError:    nil,
		},
		{
			option:           optionDown,
			userInput:        "",
			instance:         &mockMigration{},
			expectedFuncCall: "Down",
			expectedError:    nil,
		},
		{
			option:           optionDrop,
			userInput:        "",
			instance:         &mockMigration{},
			expectedFuncCall: "Drop",
			expectedError:    nil,
		},
		{
			option:           optionForce,
			userInput:        "1",
			instance:         &mockMigration{},
			expectedFuncCall: "Force",
			expectedError:    nil,
		},
		{
			option:           optionFullReset,
			userInput:        "",
			instance:         &mockMigration{},
			expectedFuncCall: "Drop",
			expectedError:    nil,
		},
		{
			option:           optionNothing,
			userInput:        "",
			instance:         &mockMigration{},
			expectedFuncCall: "",
			expectedError:    nil,
		},
	}

	for _, test := range tests {
		r := strings.NewReader(test.userInput)
		err := executeOption(r, test.instance, test.option)

		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		}

		if test.expectedError == nil {
			assert.NoError(t, err)
		}

		assert.Equal(t, test.expectedFuncCall, lastMigrationCall)

		lastMigrationCall = ""
	}
}

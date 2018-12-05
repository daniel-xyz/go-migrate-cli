package migrate

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

func TestExecuteOption(t *testing.T) {
	testScenarios := []struct {
		option            string
		userInput         string
		migrationInstance *mockMigration
		expectedFuncCall  string
		expectedError     error
	}{
		{
			option:            optionUp,
			userInput:         "",
			migrationInstance: &mockMigration{},
			expectedFuncCall:  "Up",
			expectedError:     nil,
		},
		{
			option:    optionUp,
			userInput: "",
			migrationInstance: &mockMigration{
				db: mockDB{
					migrationVersion: 1,
				},
				lastMigrationCall: "",
			},
			expectedFuncCall: "Up",
			expectedError:    errors.New("already up-to-date"),
		},
		{
			option:    optionDown,
			userInput: "",
			migrationInstance: &mockMigration{
				db: mockDB{
					migrationVersion: 1,
				},
				lastMigrationCall: "",
			}, expectedFuncCall: "Down",
			expectedError: nil,
		},
		{
			option:    optionDown,
			userInput: "",
			migrationInstance: &mockMigration{
				db: mockDB{
					migrationVersion: 0,
				},
				lastMigrationCall: "",
			},
			expectedFuncCall: "Down",
			expectedError:    errors.New("already on lowest possible version"),
		},
		{
			option:            optionDrop,
			userInput:         "",
			migrationInstance: &mockMigration{},
			expectedFuncCall:  "Drop",
			expectedError:     nil,
		},
		{
			option:            optionForce,
			userInput:         "1",
			migrationInstance: &mockMigration{},
			expectedFuncCall:  "Force(1)",
			expectedError:     nil,
		},
		{
			option:            optionForce,
			userInput:         "some string",
			migrationInstance: &mockMigration{},
			expectedFuncCall:  "",
			expectedError:     errors.New("expected integer"),
		},
		{
			option:            optionFullReset,
			userInput:         "",
			migrationInstance: &mockMigration{},
			expectedFuncCall:  "Drop",
			expectedError:     nil,
		},
		{
			option:            optionNothing,
			userInput:         "",
			migrationInstance: &mockMigration{},
			expectedFuncCall:  "",
			expectedError:     nil,
		},
	}

	for i, scenario := range testScenarios {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			r := strings.NewReader(scenario.userInput)
			err := executeOption(r, scenario.migrationInstance, scenario.option)

			if scenario.expectedError != nil {
				assert.EqualError(t, err, scenario.expectedError.Error())
			}

			if scenario.expectedError == nil {
				assert.NoError(t, err)
			}

			assert.Equal(t, scenario.expectedFuncCall, scenario.migrationInstance.lastMigrationCall)

			scenario.migrationInstance.lastMigrationCall = ""
		})
	}
}

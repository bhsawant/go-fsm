package fsm

import (
	"testing"
)

// --- Helper Transition Functions ---

func mod3TranFunc(state State, inputBit rune) State {
	switch state {
	case "S0":
		if inputBit == '0' {
			return "S0"
		}
		return "S1"
	case "S1":
		if inputBit == '0' {
			return "S2"
		}
		return "S0"
	case "S2":
		if inputBit == '0' {
			return "S1"
		}
		return "S2"
	}
	return state
}

// this is added to test the generic nature of the module
func mod2TranFunc(state State, inputBit rune) State {
	switch inputBit {
	case '0':
		return "S0"
	case '1':
		return "S1"
	default:
		return state
	}
}

func invalidTranFunc(state State, inputBit rune) State {
	return "invalid_state"
}

// --- Combined Table-Driven Tests ---

func TestEngine_Process(t *testing.T) {

	type engineSetup struct {
		states        []State
		alphabet      []rune
		initial       State
		accepting     []State
		stateTranFunc StateTransitionFunc
	}

	tests := []struct {
		name          string
		setup         engineSetup
		input         string
		expectedState State
		expectedBool  bool
		expectError   bool
	}{
		// --- Mod 3 Cases ---
		{"Mod3: input:0 (decimal 0)", engineSetup{[]State{"S0", "S1", "S2"}, []rune{'0', '1'}, "S0", []State{"S0", "S1", "S2"}, mod3TranFunc}, "0", "S0", true, false},
		{"Mod3: input:11 (decimal 3)", engineSetup{[]State{"S0", "S1", "S2"}, []rune{'0', '1'}, "S0", []State{"S0", "S1", "S2"}, mod3TranFunc}, "11", "S0", true, false},
		{"Mod3: input:10 (decimal 2)", engineSetup{[]State{"S0", "S1", "S2"}, []rune{'0', '1'}, "S0", []State{"S0", "S1", "S2"}, mod3TranFunc}, "10", "S2", true, false},
		{"Mod3 (test boolean): input:10 (decimal 2): expectedBool=false", engineSetup{[]State{"S0", "S1", "S2"}, []rune{'0', '1'}, "S0", []State{"S0"}, mod3TranFunc}, "10", "S2", false, false},

		// --- Mod 2 Cases ---
		{"Mod2: 10 (Even)", engineSetup{[]State{"S0", "S1"}, []rune{'0', '1'}, "S0", []State{"S0"}, mod2TranFunc}, "10", "S0", true, false},
		{"Mod2: 11 (Odd)", engineSetup{[]State{"S0", "S1"}, []rune{'0', '1'}, "S0", []State{"S0"}, mod2TranFunc}, "11", "S1", false, false},

		// --- Error Cases ---
		{"Error: Invalid Char", engineSetup{[]State{"S0"}, []rune{'0'}, "S0", []State{"S0"}, mod3TranFunc}, "2", "", false, true},
		{"Error: Bad Transition", engineSetup{[]State{"S0"}, []rune{'0'}, "S0", []State{"S0"}, invalidTranFunc}, "0", "", false, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize the engine for this specific case
			engine, err := NewEngine(tc.setup.states, tc.setup.alphabet, tc.setup.initial, tc.setup.accepting, tc.setup.stateTranFunc)
			if err != nil {
				t.Fatalf("Engine setup failed: %v", err)
			}

			finalState, isAccepted, err := engine.Process(tc.input)

			if (err != nil) != tc.expectError {
				t.Errorf("expected error: %v, got: %v", tc.expectError, err)
				return
			}

			if !tc.expectError {
				if finalState != tc.expectedState {
					t.Errorf("expected state %v, got %v", tc.expectedState, finalState)
				}
				if isAccepted != tc.expectedBool {
					t.Errorf("expected acceptance %v, got %v", tc.expectedBool, isAccepted)
				}
			}
		})
	}
}

func TestNewEngine_Validation(t *testing.T) {
	states := []State{"S0"}
	alphabet := []rune{'0'}

	t.Run("Invalid Initial", func(t *testing.T) {
		_, err := NewEngine(states, alphabet, "MISSING", states, mod3TranFunc)
		if err == nil {
			t.Error("expected error for missing initial state")
		}
	})

	t.Run("Invalid Accepting", func(t *testing.T) {
		_, err := NewEngine(states, alphabet, "S0", []State{"MISSING"}, mod3TranFunc)
		if err == nil {
			t.Error("expected error for missing accepting state")
		}
	})

	t.Run("Nil Logic", func(t *testing.T) {
		_, err := NewEngine(states, alphabet, "S0", states, nil)
		if err == nil {
			t.Error("expected error for nil transition function")
		}
	})
}

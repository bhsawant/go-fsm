package fsm

import (
	"fmt"
)

// State represents a system state
type State string

// StateTransitionFunc is the core function of the FSM engine. User needs to provide it.
type StateTransitionFunc func(State, rune) State

// Engine is to store the user provided 5 tuples (Q,Σ,q0,F,δ) input
type Engine struct {
	stateMap         map[State]bool      // Q - finite set of states
	inputAlphabetMap map[rune]bool       // Σ - finite input alphabet
	initialState     State               // q0 - the initial state
	acceptStateMap   map[State]bool      // F - set of accepting/final states
	nextState        StateTransitionFunc // δ - the transition function
}

// NewEngine is a constructor that initializes and validates the FSM Engine using the formal 5-tuple definition.
// It maps the provided arguments to the mathematical components of a Deterministic Finite Automaton (DFA).
//
// Parameters:
//   - states ([]State): Represents Q, the finite set of all possible states in the system.
//   - inputAlphabets ([]rune): Represents Σ, the finite set of input symbols (alphabet) the machine accepts.
//   - initialState (State): Represents q₀, the starting state of the machine (must be an element of Q).
//   - acceptingStates ([]State): Represents F, the set of states that result in an "Accepted" status (must be a subset of Q).
//   - nextState (StateTransitionFunc): Represents δ, the transition function that defines the stateTranFunc (Q × Σ → Q).
//
// Returns:
//   - *Engine: A pointer to the initialized FSM Engine if all validations pass.
//   - error: Returns an error if the 5-tuple is mathematically inconsistent (e.g., q₀ or F are not within Q).
func NewEngine(states []State,
	inputAlphabets []rune,
	initialState State,
	acceptingStates []State,
	nextState StateTransitionFunc) (*Engine, error) {

	// populating the states
	stateMap := make(map[State]bool)
	for _, state := range states {
		stateMap[state] = true
	}

	// populating the input alphabets map
	inputMap := make(map[rune]bool)
	for _, alphabet := range inputAlphabets {
		inputMap[alphabet] = true
	}

	// populating the map for valid accept states
	acceptMap := make(map[State]bool)
	for _, acceptingState := range acceptingStates {
		// if the accepting state is not in the initial states that user provided, then throw an error
		if !stateMap[acceptingState] {
			return nil, fmt.Errorf("accepting state %s is not in the valid states", acceptingState)
		}
		acceptMap[acceptingState] = true
	}

	// check if the user provided initialState is also valid
	if !stateMap[initialState] {
		return nil, fmt.Errorf("initial state %s is not in the valid states", initialState)
	}

	// check if the user provided the expected StateTransitionFunction, if not throw an error
	if nextState == nil {
		return nil, fmt.Errorf("nextState stateTransitionFunction cannot be nil")
	}

	// if all validations pass, return the Engine so that the 'process' method can be called.
	return &Engine{
		stateMap:         stateMap,
		inputAlphabetMap: inputMap,
		initialState:     initialState,
		acceptStateMap:   acceptMap,
		nextState:        nextState,
	}, nil
}

// Process is a generic implementation of the FSM engine. It iterates through each character
// and applies the StateTransitionFunction stateTranFunc (nextState) to move from the current state
// to the next state.
//
// Parameters:
//   - input: a string of characters to be processed by the Engine from left to right,
//     by transitioning to the next state. For smooth transitioning, each character must
//     be present into the list of inputAlphabets.
//
// Returns:
//
//   - State: It is the final state of the Engine after the processing.
//     In a mod3 Engine, it represents actual reminders S0, S1, and S2.
//
//   - bool: Indicates "Acceptance." This allows the user to define a success condition
//     independent of the final state.
//     Example: In a Mod-3 engine, all final states (S0, S1, S2) are valid mathematical
//     results. However, if the user only considers a number "valid" if it is
//     divisible by 3, they would define S0 as the only 'Accepting State.'
//     The bool will then return true only when the input satisfies that specific criteria.
//
//   - error: Returns a non-nil error if an invalid character is provided or the
//     machine transitions to an undefined state.
func (e *Engine) Process(input string) (State, bool, error) {
	// before starting the process, make the initialState as currentState
	currState := e.initialState

	for _, currRune := range input {
		// check if the user input is correct, else throw an error
		if !e.inputAlphabetMap[currRune] {
			return "", false, fmt.Errorf("invalid character %c in the provided input", currRune)
		}
		// if no error, then transition to the next state and make it as a current state for further processing
		currState = e.nextState(currState, currRune)

		// before returning, check if the new currentState is also a valid state, else throw an error
		if !e.stateMap[currState] {
			return "", false, fmt.Errorf("invalid transition state: %s", currState)
		}
	}
	return currState, e.acceptStateMap[currState], nil
}

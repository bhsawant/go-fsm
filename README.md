# Go FSM - Generic Finite State Machine Library

This library provides a generic, extensible implementation of a Finite State Machine (FSM) in Go, based on the formal mathematical definition of a Finite Automaton.

## Advanced Exercise: Finite Automaton (FA)

A finite automaton (FA) is defined as a 5-tuple (Q, Σ, q0, F, δ) where:

*   **Q**: A finite set of states.
*   **Σ**: A finite input alphabet.
*   **q0**: The initial state (q0 ∈ Q).
*   **F**: The set of accepting/final states (F ⊆ Q).
*   **δ**: The transition function (δ: Q × Σ → Q).

## Implementation Overview

The `fsm` package allows developers to create custom FSMs by defining their own states, alphabets, and transition logic.

### Core Components

1.  **`State`**: A string-based type representing a state in the machine.
2.  **`StateTransitionFunc`**: A function signature `func(State, rune) State` that implements the transition logic (δ).
3.  **`Engine`**: The core component that holds the FSM configuration and processes input strings.

## API Usage

### 1. Initialize the Engine

Use `fsm.NewEngine` to create a new FSM instance. It validates that the initial and accepting states are part of the defined state set.

```go
engine, err := fsm.NewEngine(states, alphabet, initial, accepting, transitionFunc)
```

### 2. Process Input

The `Process` method takes an input string and returns the final state, whether it's an accepting state, and any error encountered (e.g., invalid input characters).

```go
finalState, isAccepted, err := engine.Process("1101")
```

## Example: Modulo-Three FA

As a practical application of this library, we implement a modulo-three function. This function computes the remainder when an unsigned binary integer (provided as a string) is divided by three.

### Mod-Three Configuration

*   **Q**: {S0, S1, S2} (representing remainders 0, 1, and 2)
*   **Σ**: {0, 1}
*   **q0**: S0
*   **F**: {S0, S1, S2}
*   **δ (Transition Function)**:
    *   δ(S0, 0) = S0
    *   δ(S0, 1) = S1
    *   δ(S1, 0) = S2
    *   δ(S1, 1) = S0
    *   δ(S2, 0) = S1
    *   δ(S2, 1) = S2

### Code Snippet

```go
import "github.com/bhsawant/go-fsm/fsm"

// Define the transition logic
mod3Transition := func(state fsm.State, input rune) fsm.State {
    switch state {
    case "S0": 
        if input == '0' { return "S0" }
        return "S1"
    case "S1": 
        if input == '0' { return "S2" }
        return "S0"
    case "S2": 
        if input == '0' { return "S1" }
        return "S2"
    default: return state
    }
}

// Initialize and use the engine
engine, _ := fsm.NewEngine(
    []fsm.State{"S0", "S1", "S2"},
    []rune{'0', '1'},
    "S0",
    []fsm.State{"S0", "S1", "S2"},
    mod3Transition,
)

finalState, _, err := engine.Process("1101") // Returns "S1" (Remainder 1)
```

## Running Tests

The library is thoroughly tested with table-driven tests to ensure correctness of both the generic engine and the Modulo-Three logic.

To run tests and see detailed output every time, execute the following command from the **project root folder**:

```bash
go test -v -count=1 ./...
```

package main

import (
	"log"

	"github.com/bhsawant/go-fsm/fsm"
)

func main() {
	userInput := "1011" // binary for 11

	//FSM definitions ---

	//Q = (S0, S1, S2)
	states := []fsm.State{"S0", "S1", "S2"}

	//Σ = (0, 1)
	inputAlphabets := []rune{'0', '1'}

	//q0 = S0
	initialState := fsm.State("S0")

	//F = (S0, S1, S2)
	acceptingStates := []fsm.State{"S0"}

	//δ(S0,0) = S0; δ(S0,1) = S1; δ(S1,0) = S2; δ(S1,1) = S0; δ(S2,0) = S1; δ(S2,1) = S2
	mod3TranFunc := func(state fsm.State, inputBit rune) fsm.State {
		//Formula for modulo 3 : ((reminder * 2) + inputBit ) % 3
		switch state {

		case "S0": // reminder 0
			if inputBit == '0' {
				return "S0" // state = ((0 * 2) + 0) % 3 = 0 ==> S0
			}
			return "S1" // state = ((0 * 2) + 1) % 3 = 1 ==> S1

		case "S1": // reminder 1
			if inputBit == '0' {
				return "S2" // state = ((1 * 2) + 0) % 3 = 2 ==> S2
			}
			return "S0" // state = ((1 * 2) + 1) % 3 = 0 ==> S0

		case "S2": // reminder 2
			if inputBit == '0' {
				return "S1" // state = ((2 * 2) + 0) % 3 = 1 ==> S1
			}
			return "S2" // state = ((2 * 2) + 1) % 3 = 2 ==> S2
		}
		return state
	}

	//Create the engine based on the above arguments
	mod3Engine, err := fsm.NewEngine(states, inputAlphabets, initialState, acceptingStates, mod3TranFunc)
	if err != nil {
		log.Printf("Error while creating the mod3 FSM Engine, %v", err)
		return
	}

	result, _, err := mod3Engine.Process(userInput)
	if err != nil {
		log.Printf("Error while processing the mod3 FSM, %v", err)
		return
	}

	log.Printf("The final state is : %v", result)
}

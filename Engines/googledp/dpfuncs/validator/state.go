package validator

import "googledp/entities"

// DFA to check query format

func NewSMValidator() *StateMachine {
	sm := newStateMachine(0, false)

	sm.AddState(1, false)
	sm.AddState(2, false)
	sm.AddState(3, true)
	sm.AddState(4, false)

	sm.AddTransition(0, entities.FILTER, 1)
	sm.AddTransition(0, entities.BIN, 2)
	sm.AddTransition(0, entities.MEASUREMENT, 3)

	sm.AddTransition(1, entities.FILTER, 4)
	sm.AddTransition(1, entities.BIN, 2)
	sm.AddTransition(1, entities.MEASUREMENT, 3)

	sm.AddTransition(2, entities.FILTER, 4)
	sm.AddTransition(2, entities.BIN, 4)
	sm.AddTransition(2, entities.MEASUREMENT, 3)

	sm.AddTransition(3, entities.FILTER, 4)
	sm.AddTransition(3, entities.BIN, 4)
	sm.AddTransition(3, entities.MEASUREMENT, 4)

	sm.AddTransition(4, entities.FILTER, 4)
	sm.AddTransition(4, entities.BIN, 4)
	sm.AddTransition(4, entities.MEASUREMENT, 4)

	return sm
}

func newStateMachine(initstate int, isFinal bool) *StateMachine {
	sm := &StateMachine{
		transitions:  make(map[transitionInput]int),
		initState:    initstate,
		currentState: initstate,
	}

	sm.AddState(initstate, isFinal)

	return sm
}

type transitionInput struct {
	fromState int
	input     string
}

type StateMachine struct {
	initState    int
	currentState int
	states       []int
	finalStates  []int
	transitions  map[transitionInput]int
}

func (sm *StateMachine) AddState(state int, isFinal bool) {
	sm.states = append(sm.states, state)
	if isFinal {
		sm.finalStates = append(sm.finalStates, state)
	}
}

func (sm *StateMachine) AddTransition(fromState int, input string, toState int) {
	find := false

	for _, s := range sm.states {
		if s == fromState {
			find = true
		}
	}

	if !find {
		return
	}

	if input == "" {
		return
	}

	transition := transitionInput{fromState: fromState, input: input}

	sm.transitions[transition] = toState
}

func (sm *StateMachine) Input(testInput string) {
	currentState := sm.currentState
	transition := transitionInput{fromState: currentState, input: testInput}

	nextState, ok := sm.transitions[transition]

	if ok {
		sm.currentState = nextState
	}
}

func (sm *StateMachine) Verify() bool {
	for _, s := range sm.finalStates {
		if s == sm.currentState {
			return true
		}
	}

	return false
}

func (sm *StateMachine) Reset() {
	sm.currentState = sm.initState
}

func (sm *StateMachine) VerifyInputs(inputs []string) bool {
	for _, inp := range inputs {
		sm.Input(inp)
	}
	return sm.Verify()
}

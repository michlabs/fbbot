package fbbot

import (
	// log "github.com/Sirupsen/logrus"
)

// Event represents event triggered by state
type Event string

type State interface {
	Enter(*Bot, *Message) Event
	Process(*Bot, *Message) Event
	Leave(*Bot, *Message) Event
}

// BaseState is base struct for State
type BaseState struct {}
func (s BaseState) Enter(bot *Bot, msg *Message) (e Event)   { return e } // Do nothing
func (s BaseState) Process(bot *Bot, msg *Message) (e Event) { return e } // Do nothing
func (s BaseState) Leave(bot *Bot, msg *Message) (e Event)   { return e } // Do nothing

// manager represents a dialog manager
type manager struct {
	// maps an user's ID to pointer of his current dialog
	currentDialog map[string]*Dialog
}

// NewDialogManager returns pointer to a new dialog manager
func NewDialogManager() *manager {
	return &manager{currentDialog: make(map[string]*Dialog)}
}

// SetDialog sets current dialog for user
func (m *manager) SetDialog(user_id string, d *Dialog) {
	m.currentDialog[user_id] = d
}

// GetDialog returns pointer to current dialog of the user
// Returns nil if user is not in any dialog
func (m *manager) GetDialog(user_id string) *Dialog {
	return m.currentDialog[user_id]
}

type Dialog struct {
	BeginState State
	EndState   State

	stateMap map[string]State // maps an user ID to his current state
	transMap map[State]map[Event]State
}

func NewDialog() *Dialog {
	var d Dialog
	d.stateMap = make(map[string]State)
	d.transMap = make(map[State]map[Event]State)

	return &d
}

func (d *Dialog) AddTransition(src State, event Event, dst State) {
	_, exist := d.transMap[src]
	if !exist {
		d.transMap[src] = make(map[Event]State)
	}
	d.transMap[src][event] = dst
}

func (d *Dialog) Handle(bot *Bot, msg *Message) {
	var event Event
	state := d.getState(msg.Sender.ID)
	if state == nil { // Start dialog
		d.setState(msg.Sender.ID, d.BeginState)
		state = d.getState(msg.Sender.ID)
		event = state.Enter(bot, msg)
	} else {
		event = state.Process(bot, msg)
	}
	d.transition(bot, msg, state, event)
}

func (d *Dialog) transition(bot *Bot, msg *Message, src State, event Event) {
	dst, exist := d.transMap[src][event]
	if !exist {
		return
	}
	src.Leave(bot, msg)
	d.setState(msg.Sender.ID, dst)
	event = d.getState(msg.Sender.ID).Enter(bot, msg)
	if dst == d.EndState {
		d.resetState(msg.Sender.ID)
		return
	}
	d.transition(bot, msg, dst, event)
}

func (d *Dialog) setState(user_id string, state State) {
	d.stateMap[user_id] = state
}

func (d *Dialog) getState(user_id string) State {
	return d.stateMap[user_id]
}

func (d *Dialog) resetState(user_id string) {
	delete(d.stateMap, user_id)
}
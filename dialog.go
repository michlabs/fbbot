package fbbot

import (
	"fmt"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"
)

// Event represents event triggered by state
type Event string

const ResetEvent Event = "reset"
const NilEvent Event = ""

type State interface {
	Enter(*Bot, *Message) Event
	Process(*Bot, *Message) Event
	Leave(*Bot, *Message) Event
}

// BaseState is base struct for State
type BaseState struct {
	Name string
}
func (s BaseState) Enter(bot *Bot, msg *Message) (e Event)   { return e } // Do nothing
func (s BaseState) Process(bot *Bot, msg *Message) (e Event) { return e } // Do nothing
func (s BaseState) Leave(bot *Bot, msg *Message) (e Event)   { return e } // Do nothing

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

func (d *Dialog) Render() {
	nodes := make(map[State]bool)
	var edgesHTML string = ""
	for src, values := range d.transMap {
		nodes[src] = true
		for event, dst := range values {
			edgesHTML = edgesHTML + fmt.Sprintf("g.setEdge(\"%s\", \"%s\", {label: \"%s\" });\n", src, dst, event)
			nodes[dst] = true
		}
	}

	var nodesHTML string = `var states = [`
	for node, _ := range nodes {
		nodesHTML = nodesHTML + fmt.Sprintf(`"%s", `, node)
	}
	nodesHTML = nodesHTML + `];`

	html := fmt.Sprintf(TEMPLATE, nodesHTML, edgesHTML, d.BeginState, d.EndState)
	err := ioutil.WriteFile("dialog.html", []byte(html), 0644)
    if err != nil {
        log.Error("Could not write dialog.html file")
    }
}

func (d *Dialog) AddTransition(src State, event Event, dst State) {
	_, exist := d.transMap[src]
	if !exist {
		d.transMap[src] = make(map[Event]State)
	}
	d.transMap[src][event] = dst
}

func (d *Dialog) Handle(bot *Bot, msg *Message) {
	if d.BeginState == nil || d.EndState == nil {
		log.Fatal("BeginState and EndState are not set.")
	}

	var event Event
	state := d.getState(msg.Sender.ID)
	if state == nil || state == d.EndState {
		d.setState(msg.Sender.ID, d.BeginState)
		state = d.getState(msg.Sender.ID)
		event = state.Enter(bot, msg)
	} else {
		event = state.Process(bot, msg)
	}
	d.transition(bot, msg, state, event)
}

func (d *Dialog) transition(bot *Bot, msg *Message, src State, event Event) {
	if event == ResetEvent {
		d.resetState(msg.Sender.ID)
		return
	}
	
	dst, exist := d.transMap[src][event]
	if !exist {
		return
	}
	src.Leave(bot, msg)
	d.setState(msg.Sender.ID, dst)
	event = d.getState(msg.Sender.ID).Enter(bot, msg)
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
package fbbot

import (
	log "github.com/Sirupsen/logrus"
)

// Event represents event triggered in a step
type Event string

const ResetEvent Event = "reset"
const NilEvent Event = ""

type Step interface {
	Name() string
	Enter(*Bot, *Message) Event
	Process(*Bot, *Message) Event
	Leave(*Bot, *Message) Event
}

// BaseStep is base struct for steps
type BaseStep struct{}

func (s BaseStep) Name() string                             { return "unnamed step" }
func (s BaseStep) Enter(bot *Bot, msg *Message) (e Event)   { return e } // Do nothing
func (s BaseStep) Process(bot *Bot, msg *Message) (e Event) { return e } // Do nothing
func (s BaseStep) Leave(bot *Bot, msg *Message) (e Event)   { return e } // Do nothing

type Dialog struct {
	beginStep Step
	endStep   Step

	currentStepMap map[string]Step // maps an user ID to his current step
	p2pTransMap    map[Step]map[Event]Step
	globalTransMap map[Event]Step

	// Hooks
	PreHandleMessageHook   func(*Bot, *Message) bool
	PostHandleMessageHook  func(*Bot, *Message)
	PreHandlePostbackHook  func(*Bot, *Postback) bool
	PostHandlePostbackHook func(*Bot, *Postback)
}

func NewDialog() *Dialog {
	var d Dialog
	d.currentStepMap = make(map[string]Step)
	d.p2pTransMap = make(map[Step]map[Event]Step)
	d.globalTransMap = make(map[Event]Step)

	return &d
}

func (d *Dialog) SetBeginStep(s Step) {
	d.beginStep = s
}

func (d *Dialog) SetEndStep(s Step) {
	d.endStep = s
}

func (d *Dialog) AddTransition(event Event, steps ...Step) {
	n := len(steps)
	if n == 0 {
		return
	}

	if n == 1 { // global transition
		d.globalTransMap[event] = steps[0]
		return
	}

	// point-to-point transition
	srcs := steps[:n-1]
	dst := steps[n-1]
	for _, src := range srcs {
		d.addP2PTransition(src, event, dst)
	}
}

// Add point-to-point transition
func (d *Dialog) addP2PTransition(src Step, event Event, dst Step) {
	_, exist := d.p2pTransMap[src]
	if !exist {
		d.p2pTransMap[src] = make(map[Event]Step)
	}
	d.p2pTransMap[src][event] = dst
}

func (d *Dialog) HandleMessage(bot *Bot, msg *Message) {
	if d.PreHandleMessageHook != nil {
		isBreak := d.PreHandleMessageHook(bot, msg)
		if isBreak {
			return
		}
	}

	if d.beginStep == nil || d.endStep == nil {
		log.Fatal("BeginStep and EndStep are not set.")
	}

	var event Event
	step := d.getStep(msg.Sender.ID)
	if step == nil || step == d.endStep {
		bot.STMemory.Delete(msg.Sender.ID)
		d.setStep(msg.Sender.ID, d.beginStep)
		step = d.getStep(msg.Sender.ID)
		event = step.Enter(bot, msg)
	} else {
		event = step.Process(bot, msg)
	}
	d.transition(bot, msg, step, event)

	if d.PostHandleMessageHook != nil {
		d.PostHandleMessageHook(bot, msg)
	}
}

func (d *Dialog) HandlePostback(bot *Bot, pbk *Postback) {
	if d.PreHandlePostbackHook != nil {
		isBreak := d.PreHandlePostbackHook(bot, pbk)
		if isBreak {
			return
		}
	}

	msg := &Message{Sender: pbk.Sender, Text: pbk.Payload}
	d.HandleMessage(bot, msg)

	if d.PostHandlePostbackHook != nil {
		d.PostHandlePostbackHook(bot, pbk)
	}
}

func (d *Dialog) transition(bot *Bot, msg *Message, src Step, event Event) {
	if event == ResetEvent {
		d.Reset(msg.Sender.ID)
		return
	}

	var dst Step
	var exist bool
	// check point-to-point transition first
	dst, exist = d.p2pTransMap[src][event]
	if !exist { // if doesn't have point-to-point transition
		// then check global transition
		dst, exist = d.globalTransMap[event]
		if !exist { // if doesn't have global transition too
			return // then do nothing
		}
	}

	// if have destination step
	src.Leave(bot, msg)
	d.setStep(msg.Sender.ID, dst)
	event = d.getStep(msg.Sender.ID).Enter(bot, msg)
	d.transition(bot, msg, dst, event)
}

func (d *Dialog) setStep(user_id string, step Step) {
	d.currentStepMap[user_id] = step
}

func (d *Dialog) getStep(user_id string) Step {
	return d.currentStepMap[user_id]
}

func (d *Dialog) Reset(user_id string) {
	delete(d.currentStepMap, user_id)
}

// This function used for moving dialog to any step.
// It should be used with caution for adhoc cases only, since it breaks already defined dialog flow.
func (d *Dialog) Move(msg *Message, dst Step) {
	// Get out of current step nicely
	currentStep := d.getStep(msg.Sender.ID)
	if currentStep != nil {
		currentStep.Leave(bot, msg)
	}

	if dst == nil || dst == d.endStep {
		// Follow current logic for end step
		bot.STMemory.Delete(msg.Sender.ID)
		dst = d.beginStep
	}
	d.setStep(msg.Sender.ID, dst)
	event := dst.Enter(bot, msg)
	d.transition(bot, msg, dst, event)
}

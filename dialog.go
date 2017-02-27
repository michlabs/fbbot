package fbbot

import (
	log "github.com/Sirupsen/logrus"
)

// Event represents event triggered in a step
type Event string

const ResetEvent Event = "reset"
const NilEvent Event = ""

type Step interface {
	Enter(*Bot, *Message) Event
	Process(*Bot, *Message) Event
	Leave(*Bot, *Message) Event
}

// BaseStep is base struct for steps
type BaseStep struct {
	Name string
}

func (s BaseStep) Enter(bot *Bot, msg *Message) (e Event)   { return e } // Do nothing
func (s BaseStep) Process(bot *Bot, msg *Message) (e Event) { return e } // Do nothing
func (s BaseStep) Leave(bot *Bot, msg *Message) (e Event)   { return e } // Do nothing

type Dialog struct {
	beginStep Step
	endStep   Step

	steps   map[Step]bool   // stores all steps of this dialog
	stepMap map[string]Step // maps an user ID to his current step
	transMap map[Step]map[Event]Step
}

func NewDialog() *Dialog {
	var d Dialog
	d.steps = make(map[Step]bool)
	d.stepMap = make(map[string]Step)
	d.transMap = make(map[Step]map[Event]Step)

	return &d
}

func (d *Dialog) AddSteps(steps ...Step) {
	for _, step := range steps {
		d.steps[step] = true
	}
}

func (d *Dialog) SetBeginStep(s Step) {
	d.beginStep = s
}

func (d *Dialog) SetEndStep(s Step) {
	d.endStep = s
}

func (d *Dialog) AddTransition(src Step, event Event, dst Step) {
	_, exist := d.transMap[src]
	if !exist {
		d.transMap[src] = make(map[Event]Step)
	}
	d.transMap[src][event] = dst
}

func (d *Dialog) AddGlobalTransition(event Event, dst Step) {
	for step := range d.steps {
		d.AddTransition(step, event, dst)
	}
}

func (d *Dialog) Handle(bot *Bot, msg *Message) {
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
}

func (d *Dialog) transition(bot *Bot, msg *Message, src Step, event Event) {
	if event == ResetEvent {
		d.resetStep(msg.Sender.ID)
		return
	}

	dst, exist := d.transMap[src][event]
	if !exist {
		return
	}
	src.Leave(bot, msg)
	d.setStep(msg.Sender.ID, dst)
	event = d.getStep(msg.Sender.ID).Enter(bot, msg)
	d.transition(bot, msg, dst, event)
}

func (d *Dialog) setStep(user_id string, step Step) {
	d.stepMap[user_id] = step
}

func (d *Dialog) getStep(user_id string) Step {
	return d.stepMap[user_id]
}

func (d *Dialog) resetStep(user_id string) {
	delete(d.stepMap, user_id)
}

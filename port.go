package gossipswitch

import (
	"github.com/DSiSc/craft/log"
	"sync"
)

// state is used to record switch port state. e.g., message statistics
type state struct {
	InCount  uint
	OutCount uint
}

// InPort is switch in port. Message will be send to InPort, and then switch read the message from InPort
type InPort struct {
	State   state
	channel chan interface{}
}

// create a new in port instance
func newInPort() *InPort {
	return &InPort{
		State:   state{},
		channel: make(chan interface{}),
	}
}

// Channel return the input channel of this InPort
func (inPort *InPort) Channel() chan<- interface{} {
	return inPort.channel
}

// read message from this InPort, will be blocked until the message arrives.
func (inPort *InPort) read() <-chan interface{} {
	return inPort.channel
}

// OutPutFunc is binded to switch out port, and OutPort will call OutPutFunc when receive a message from switch
type OutPutFunc func(msg interface{}) error

// InPort is switch out port. Switch will broadcast message to out port
type OutPort struct {
	outPortMtx  sync.Mutex
	State       state
	outPutFuncs []OutPutFunc
}

// create a new out port instance
func newOutPort() *OutPort {
	return &OutPort{
		State: state{},
	}
}

// BindToPort bind a new OutPutFunc to this OutPort. Return error if bind failed
func (outPort *OutPort) BindToPort(outPutFunc OutPutFunc) error {
	log.Info("Bind OutPutFunc to OutPort %d")
	outPort.outPortMtx.Lock()
	defer outPort.outPortMtx.Unlock()
	outPort.outPutFuncs = append(outPort.outPutFuncs, outPutFunc)
	return nil
}

// write message to this OutPort
func (outPort *OutPort) write(msg interface{}) error {
	outPort.outPortMtx.Lock()
	defer outPort.outPortMtx.Unlock()
	for _, outPutFunc := range outPort.outPutFuncs {
		go outPutFunc(msg)
	}
	return nil
}

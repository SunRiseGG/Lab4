package engine

import (
  "sync"
)

type Command interface {
  Execute(handler Handler)
}
type Handler interface {
  Post(cmd Command)
}

type messageQueue struct {
  sync.Mutex
  data []Command
  waiting bool
}

var receiveSignal = make(chan struct{})

func (mq *messageQueue) push(cmd Command) {
  mq.Lock()
  defer mq.Unlock()
  mq.data = append(mq.data, cmd)
  if mq.waiting {
    mq.waiting = false
    receiveSignal <- struct{}{}
  }
}

func (mq *messageQueue) pull() Command {
  mq.Lock()
  defer mq.Unlock()
  if len(mq.data) == 0 {
    mq.waiting = true;
    mq.Unlock()
    <- receiveSignal
    mq.Lock()
  }
  res := mq.data[0]
  mq.data[0] = nil
  mq.data = mq.data[1:]
  return res
}

func (mq *messageQueue) size() int {
  return len(mq.data)
}

type Loop struct {
  queue *messageQueue
  terminateReceived bool
  stopSignal chan struct{}
}

func (l *Loop) Start() {
  l.queue = new(messageQueue)
  l.stopSignal = make(chan struct{})
  go func() {
    for (!l.terminateReceived) || (l.queue.size() != 0) {
      cmd := l.queue.pull()
      cmd.Execute(l)
    }
    l.stopSignal <- struct{}{}
  }()
}

type CommandFunc func (handler Handler)

func (c CommandFunc) Execute(handler Handler) {
  c(handler)
}

func (l *Loop) AwaitFinish() {
  l.Post(CommandFunc(func (h Handler) {
    h.(*Loop).terminateReceived = true
  }))
  <- l.stopSignal
}

func (l * Loop) Post(cmd Command) {
  l.queue.push(cmd)
}

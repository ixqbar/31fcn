package spider

import (
	"encoding/json"
)

const (
	IsSock5 = iota
)

type TProxyItem struct {
	Category string `json:"category"`
	Country  string `json:"country"`
	Address  string `json:"address"`
	Port     string `json:"port"`
	Time     int64  `json:"time"`
}

type TQueueItem struct {
	Type  int
	Value interface{}
}

func (obj *TQueueItem) Content() string {
	content, err := json.Marshal(obj.Value)
	if err != nil {
		Logger.Print(err)
		return ""
	}

	return string(content)
}

type TQueue struct {
	stopped chan bool
	values  chan TQueueItem
}

func NewTQueue(size int) *TQueue {
	return &TQueue{
		stopped: make(chan bool),
		values:  make(chan TQueueItem, size),
	}
}

func (obj *TQueue) Push(value TQueueItem) {
	obj.values <- value
}

func (obj *TQueue) Pull() TQueueItem {
	return <-obj.values
}

func (obj *TQueue) Stop() {
	obj.stopped <- true
}

func (obj *TQueue) PullToAsync(f func(TQueueItem)) {
E:
	for {
		select {
		case <-obj.stopped:
			break E
		case value := <-obj.values:
			go f(value)
			break
		}
	}

	Logger.Print("tqueue exit")
}

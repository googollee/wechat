package wechat

type EventType string

const (
	EventSubscribe   EventType = "subscribe"
	EventUnsubscribe           = "unsubscribe"
	EventSCAN                  = "SCAN"
	EventLocation              = "LOCATION"
	EventClick                 = "CLICK"
	EventView                  = "VIEW"
)

type EventHandler func(resp Response, req *Request)

type eventNode struct {
	event   EventType
	handler EventHandler
}

func (wc *Wechat) Event(event EventType, handler EventHandler) error {
	wc.eventNodes = append(wc.eventNodes, eventNode{
		event:   event,
		handler: handler,
	})
	return nil
}

func (wc *Wechat) handleEventMessage(resp Response, req *Request) {
	for _, node := range wc.eventNodes {
		if node.event == req.Message.Event {
			node.handler(resp, req)
			return
		}
	}
}

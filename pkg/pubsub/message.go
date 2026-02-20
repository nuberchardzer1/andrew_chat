package pubsub

import "encoding/json"

const (
	SubscribeAction = "subscribe"
	UnsubscribeAction = "unsubscribe"
	MessageAction = "publish"
)

type Message struct {
	Action string `json:"action"`
	Topic string `json:"topic"`
	Payload []byte `json:"payload"`
}

func (m *Message) Marshal()([]byte, error){
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil{
		return nil, err
	}
	return b, nil
}

func (m *Message) Unmarshal(b []byte)error{
	if err := json.Unmarshal(b, m); err != nil{
		return err
	}
	return nil
}
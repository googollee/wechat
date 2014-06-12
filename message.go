package wechat

import (
	"encoding/xml"
	"io"
)

type MsgType string

const (
	MsgText     MsgType = "text"
	MsgImage            = "image"
	MsgVoice            = "voice"
	MsgVideo            = "video"
	MsgLocation         = "location"
	MsgLink             = "link"
	MsgEvent            = "event"
)

type MessageHeader struct {
	MsgId        string
	MsgType      MsgType
	ToUserName   string
	FromUserName string
	CreateTime   int64
}

func (h MessageHeader) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	helper := encodeElementHelper{e, nil}
	start.Name.Local = "ToUserName"
	helper.Encode(start, xml.CharData(h.ToUserName))
	start.Name.Local = "FromUserName"
	helper.Encode(start, xml.CharData(h.FromUserName))
	start.Name.Local = "CreateTime"
	helper.Encode(start, h.CreateTime)
	start.Name.Local = "MsgType"
	helper.Encode(start, xml.CharData(h.MsgType))
	return helper.Error()
}

type Message struct {
	MessageHeader

	Content string

	PicUrl       string
	MediaId      string
	Format       string
	ThumbMediaId string

	LocationX float64 `xml:"Location_X"`
	LocationY float64 `xml:"Location_Y"`
	Scale     float64
	Label     string

	Event    string
	EventKey string
	Ticket   string

	Latitude  float64
	Longitude float64
	Precision float64
}

func ParseMessage(r io.Reader) (Message, error) {
	decoder := xml.NewDecoder(r)
	var ret Message
	err := decoder.Decode(&ret)
	return ret, err
}

func (msg Message) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := msg.MessageHeader.MarshalXML(e, start); err != nil {
		return err
	}
	helper := encodeElementHelper{e, nil}
	switch msg.MsgType {
	case MsgText:
		start.Name.Local = "Content"
		helper.Encode(start, xml.CharData(msg.Content))
	}
	return helper.Error()
}

type encodeElementHelper struct {
	e   *xml.Encoder
	err error
}

func (e *encodeElementHelper) Encode(start xml.StartElement, v interface{}) {
	if e.err != nil {
		return
	}
	e.err = e.e.EncodeElement(v, start)
}

func (e *encodeElementHelper) Error() error {
	return e.err
}

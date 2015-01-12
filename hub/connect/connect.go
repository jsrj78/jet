package connect

import (
	"fmt"
	"time"

	"github.com/dataence/glog"
	"github.com/surge/surgemq/message"
	"github.com/surge/surgemq/service"
	"github.com/ugorji/go/codec"
)

type Connection struct {
	clt *service.Client
	mh  *codec.MsgpackHandle
}

func NewConnection(ha string) (*Connection, error) {
	c := &Connection{}
	c.clt = &service.Client{}
	c.mh = &codec.MsgpackHandle{}

	// Creates a new MQTT CONNECT message and sets the proper parameters
	msg := message.NewConnectMessage()
	msg.SetVersion(4)
	msg.SetClientId([]byte(fmt.Sprintf("hubclient%d", time.Now().Unix())))
	msg.SetKeepAlive(300)

	// Connects to the remote server at 127.0.0.1 port 1883
	if err := c.clt.Connect("tcp://:1883", msg); err != nil {
		glog.Fatal(err)
	}

	// Creates a new SUBSCRIBE message to subscribe to topic "abc"
	submsg := message.NewSubscribeMessage()
	submsg.AddTopic([]byte("#"), 0)

	c.clt.Subscribe(submsg, nil, onPublish)

	c.Send("/test", "hello")

	return c, nil
}

func onPublish(msg *message.PublishMessage) error {
	glog.Infof("t: %q p: %02X", string(msg.Topic()), msg.Payload())
	return nil
}

func (c *Connection) Send(key string, val interface{}) {
	pubmsg := message.NewPublishMessage()
	pubmsg.SetQoS(0)
	pubmsg.SetTopic([]byte(key))

	var b []byte
	enc := codec.NewEncoderBytes(&b, c.mh)
	err := enc.Encode(val)
	if err != nil {
		glog.Error(err)
		return
	}
	pubmsg.SetPayload(b)

	c.clt.Publish(pubmsg, nil)
}

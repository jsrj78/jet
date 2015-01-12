package connect

import (
	"fmt"
	"time"

	"github.com/dataence/glog"
	"github.com/surge/surgemq/message"
	"github.com/surge/surgemq/service"
	"github.com/ugorji/go/codec"
)

var mh = &codec.MsgpackHandle{RawToString: true}

type Connection struct {
	clt *service.Client
}

func NewConnection(ha string) (*Connection, error) {
	c := &Connection{}
	c.clt = &service.Client{}

	// connect to the remote server
	msg := message.NewConnectMessage()
	msg.SetVersion(4)
	msg.SetClientId([]byte(fmt.Sprintf("hubclient%d", time.Now().Unix())))
	msg.SetKeepAlive(300)
	if err := c.clt.Connect("tcp://:1883", msg); err != nil {
		return nil, err
	}

	// subscribe to all topics
	submsg := message.NewSubscribeMessage()
	submsg.AddTopic([]byte("#"), 0)
	c.clt.Subscribe(submsg, nil, onPublish)

	// send a test data structure
	c.Send("/test", []interface{}{"hello", 123, nil, 45.67})

	return c, nil
}

func onPublish(msg *message.PublishMessage) error {
	b := msg.Payload()
	dec := codec.NewDecoderBytes(b, mh)
	var v interface{}
	if err := dec.Decode(&v); err != nil {
		return fmt.Errorf("cannot decode %v (%v)", msg, err)
	}

	glog.Infof("t: %q p: %v", string(msg.Topic()), v)
	return nil
}

func (c *Connection) Send(key string, val interface{}) {
	// don't return errors, only report them
	var err error
	defer func() {
		if err != nil {
			glog.Error(err)
		}
	}()

	// encode the payload as bytes
	var b []byte
	enc := codec.NewEncoderBytes(&b, mh)
	err = enc.Encode(val)
	if err != nil {
		return
	}

	// publish the message
	pubmsg := message.NewPublishMessage()
	pubmsg.SetTopic([]byte(key))
	pubmsg.SetPayload(b)
	err = c.clt.Publish(pubmsg, nil)
}

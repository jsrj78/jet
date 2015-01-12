package connect

import (
	"fmt"
	"os"
	"time"

	"github.com/dataence/glog"
	"github.com/surge/surgemq/message"
	"github.com/surge/surgemq/service"
	"github.com/ugorji/go/codec"
)

const protocolVersion = 1

var mh = &codec.MsgpackHandle{RawToString: true}

// Connection represents an open connection to the hub.
type Connection struct {
	Done chan struct{}
	clt  *service.Client
}

// NewConnection connects to the hub and announces this pack to it.
func NewConnection(name string) (*Connection, error) {
	c := &Connection{}
	c.Done = make(chan struct{})
	c.clt = &service.Client{}

	// collect some information
	now := time.Now().Unix()
	pid := os.Getpid()
	host, _ := os.Hostname()

	// connect to the remote server
	msg := message.NewConnectMessage()
	msg.SetVersion(4)
	msg.SetClientId([]byte(fmt.Sprintf("%s_%d", name, now)))
	msg.SetKeepAlive(300)
	if err := c.clt.Connect("tcp://:1883", msg); err != nil {
		return nil, err
	}

	// send out a greeting to sign up
	c.Send("hub/hello", []interface{}{protocolVersion, name, now, pid, host})

	// TODO add last-will message to broadcast when the connection goes away
	return c, nil
}

// Listen to a topic (may have wildcards) and call the provided callback
func (c *Connection) Listen(pat string, cb func(key string, val interface{})) {
	submsg := message.NewSubscribeMessage()
	submsg.AddTopic([]byte(pat), 0)

	// inject a function to perform the decoding before doing the callback
	c.clt.Subscribe(submsg, nil, func(msg *message.PublishMessage) error {

		// decode the payload
		b := msg.Payload()
		dec := codec.NewDecoderBytes(b, mh)
		var v interface{}
		if err := dec.Decode(&v); err != nil {
			return fmt.Errorf("cannot decode %v (%v)", msg, err)
		}

		// launch the supplied callback function
		cb(string(msg.Topic()), v)
		return nil
	})
}

// Send a key/value message to the hub.
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
	if err != nil {
		return
	}
}

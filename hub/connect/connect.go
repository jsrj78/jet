package connect

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/surge/glog"
	"github.com/surgemq/surgemq/message"
	"github.com/surgemq/surgemq/service"
	"github.com/ugorji/go/codec"
)

const protocolVersion = 1

var mh = &codec.MsgpackHandle{RawToString: true}

// Connection represents an open connection to the hub.
type Connection struct {
	Done chan struct{}
	clt  *service.Client
}

// NewConnection connects to the hub and registers this pack.
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
	msg.SetKeepAlive(3600)

	// set up a "will" which gets published on connection loss
	msg.SetWillFlag(true)
	msg.SetWillTopic([]byte("hub/goodbye"))
	val, err := encode([]interface{}{protocolVersion, name, now, pid, host})
	if err != nil {
		return nil, err
	}
	msg.SetWillMessage(val)

	if err := c.clt.Connect("tcp://:1883", msg); err != nil {
		return nil, err
	}

	// send out a greeting to sign up
	c.Send("hub/hello", val)

	return c, nil
}

// encode the payload as bytes, but only if val is not yet a []byte
func encode(val interface{}) ([]byte, error) {
	var e error
	v, ok := val.([]byte)
	if !ok {
		enc := codec.NewEncoderBytes(&v, mh)
		e = enc.Encode(val)
		if e != nil {
			v = nil
		}
	}
	return v, e
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

	// encode the payload as bytes, but only if val is not yet a []byte
	v, err := encode(val)
	if err != nil {
		return
	}

	// publish the message
	pubmsg := message.NewPublishMessage()
	pubmsg.SetTopic([]byte(key))
	pubmsg.SetPayload(v)
	pubmsg.SetRetain(strings.HasPrefix(key, "/") && val != nil)
	err = c.clt.Publish(pubmsg, nil)
	if err != nil {
		return
	}
}

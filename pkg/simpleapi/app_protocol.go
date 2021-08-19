package simpleapi

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"simple/pkg/simplenet"
	"time"
)

func (app *App) newClientCodec(rw io.ReadWriter) (simplenet.Codec, error) {
	return app.newCodec(rw, app.newResponse), nil
}

func (app *App) newServerCodec(rw io.ReadWriter) (simplenet.Codec, error) {
	return app.newCodec(rw, app.newRequest), nil
}

func (app *App) newCodec(rw io.ReadWriter, newMessage func(byte, byte) (Message, error)) simplenet.Codec {
	c := &codec{
		app:        app,
		conn:       rw.(net.Conn),
		reader:     bufio.NewReaderSize(rw, app.ReadBufSize),
		newMessage: newMessage,
	}
	c.headBuf = c.headData[:]
	return c
}

func (app *App) newRequest(serviceID, messageID byte) (Message, error) {
	if service := app.services[serviceID]; service != nil {
		if msg := service.(Service).NewRequest(messageID); msg != nil {
			return msg, nil
		}
		return nil, fmt.Errorf("unsupported message type: [%d, %d]", serviceID, messageID)
	}
	return nil, fmt.Errorf("unsupported service: [%d, %d]", serviceID, messageID)
}

func (app *App) newResponse(serviceID, messageID byte) (Message, error) {
	if service := app.services[serviceID]; service != nil {
		if msg := service.(Service).NewResponse(messageID); msg != nil {
			return msg, nil
		}
		return nil, fmt.Errorf("unsupported message type: [%d, %d]", serviceID, messageID)
	}
	return nil, fmt.Errorf("unsupported service: [%d, %d]", serviceID, messageID)
}

const packetHeadSize = 4 + 2

type codec struct {
	app        *App
	headBuf    []byte
	headData   [packetHeadSize]byte
	conn       net.Conn
	reader     *bufio.Reader
	newMessage func(byte, byte) (Message, error)
}

func (c *codec) Conn() net.Conn {
	return c.conn
}

func (c *codec) Receive() (msg interface{}, err error) {
	if c.app.RecvTimeout > 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.app.RecvTimeout))
		defer c.conn.SetReadDeadline(time.Time{})
	}

	if _, err = io.ReadFull(c.reader, c.headBuf); err != nil {
		return
	}

	packetSize := int(binary.LittleEndian.Uint32(c.headBuf))

	if packetSize > c.app.MaxRecvSize {
		return nil, fmt.Errorf("too large receive packet size: %d", packetSize)
	}

	packet := make([]byte, packetSize)

	if _, err = io.ReadFull(c.reader, packet); err == nil {
		msg1, err1 := c.newMessage(c.headData[4], c.headData[5])
		if err1 == nil {
			func() {
				defer func() {
					if panicErr := recover(); panicErr != nil {
						err = fmt.Errorf("%v", panicErr)
					}
				}()
				msg1.Unmarshal(packet)
			}()
			msg = msg1
		} else {
			err = err1
		}
	}

	return
}

func (c *codec) Send(m interface{}) (err error) {
	msg := m.(Message)

	pb, err := msg.Marshal()
	if err != nil {
		return err
	}

	packetSize := len(pb)

	if packetSize > c.app.MaxSendSize {
		panic(fmt.Sprintf("too large send packet size: %d", packetSize))
	}

	packet := make([]byte, packetHeadSize+packetSize)
	binary.LittleEndian.PutUint32(packet, uint32(packetSize))
	packet[4] = msg.ServiceID()
	packet[5] = msg.MessageID()
	copy(packet[packetHeadSize:], pb)

	if c.app.SendTimeout > 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.app.SendTimeout))
		defer c.conn.SetWriteDeadline(time.Time{})
	}

	_, err = c.conn.Write(packet)
	return
}

func (c *codec) Close() error {
	return c.conn.Close()
}

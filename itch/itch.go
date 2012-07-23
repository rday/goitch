package itch

import (
	"io"
	"encoding/binary"
)

type ITCHMessage interface {
	GetType() byte
}

type ITCHMessageReader struct {
	R     io.ByteReader // Source of data
	data  []byte        // Unread data
	state int           // State of current block
}

func (msg *ITCHMessageReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	msg.state = 0
	if len(msg.data) == 0 {
		for {
			switch msg.state {
			case 0:
				for {
					b, err := msg.R.ReadByte()
					if err != nil {
						return 0, err
					}

					if b == 0 {
						msg.state = 1
						break
					}
				}
			case 1:
				size, err := msg.R.ReadByte()
				if err != nil {
					return 0, err
				}
				for i := 0; i < int(size); i++ {
					b, err := msg.R.ReadByte()
					if err != nil {
						return 0, err
					}
					p[i] = b
				}

				return int(size), nil
				msg.state = 0
			}
		}
	}

	return 1, nil
}

type TimestampMessage struct {
	Type        byte
	Description string
	Timestamp   string
}

func (msg *TimestampMessage) GetType() byte {
	return 'T'
}

type SystemEventMessage struct {
	Type      byte
	Timestamp []byte
	EventCode byte
}

func (msg *SystemEventMessage) GetType() byte {
	return 'S'
}

type AddOrderMessage struct {
	Type        byte
	Timestamp   uint32
	OrderRefNum uint32
	Indicator   byte
	Shares      uint32
	Ticker      []byte
	Price       uint32
}

func (msg *AddOrderMessage) GetType() byte {
	return 'A'
}

type ITCHProcessor struct {
	Inp      io.ByteReader
	handlers map[byte]func(ITCHMessage)
}

func New(reader io.ByteReader) *ITCHProcessor {
	itch := &ITCHProcessor{Inp: reader,
		handlers: make(map[byte]func(ITCHMessage))}
	return itch
}

func (itch *ITCHProcessor) AddHandler(handler func(ITCHMessage), msgtype byte) bool {
	itch.handlers[msgtype] = handler
	return true
}

func (itch *ITCHProcessor) Process() error {
	var msg ITCHMessage
	itchReader := &ITCHMessageReader{R: itch.Inp}
	buf := make([]byte, 255)
	for {
		_, err := itchReader.Read(buf)
		if err != nil {
			return err
		}

		switch buf[0] {
		case 'S':
			msg = &SystemEventMessage{Type: 'S', Timestamp: buf[1:5], EventCode: buf[5]}
		case 'A':
			msg = &AddOrderMessage{Type: 'A', Timestamp: binary.BigEndian.Uint32(buf[1:5]), OrderRefNum: binary.BigEndian.Uint32(buf[5:13]), Indicator: buf[13], Shares: binary.BigEndian.Uint32(buf[14:18]), Ticker: buf[18:26], Price: binary.BigEndian.Uint32(buf[26:30])}
		}

		if msg != nil {
			itch.handlers[msg.GetType()](msg)
			msg = nil
		}
	}

	return nil
}

package rocketmq

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

type messageClientIDSetter struct {
	counter       int
	basePos       int
	startTime     int64
	nextStartTime int64
	stringBuilder *bytes.Buffer // ip  + pid + classloaderid + counter + time
	buffer        *bytes.Buffer
}

var stringBuilder = bytes.NewBuffer([]byte{})

func init() {
	binary.Write(stringBuilder, binary.BigEndian, GetLocalIp4()) // 4
	binary.Write(stringBuilder, binary.BigEndian, os.Getpid())   // 2
	binary.Write(stringBuilder, binary.BigEndian, hashCode())    // 4
	MessageClientIDSetter.stringBuilder = stringBuilder
	MessageClientIDSetter.setStartTime()
}

var MessageClientIDSetter = messageClientIDSetter{
	stringBuilder: bytes.NewBuffer([]byte{}), // length := 4 + 2 + 4 + 4 + 2
	basePos:       stringBuilder.Len() * 2,
	counter:       0,
}

func hashCode() []byte {
	tmpByte := []byte{1, 1, 1, 1}
	return tmpByte
}

func (m messageClientIDSetter) setUniqID(msg *Message) {
	if msg.Properties[MessageConst.PropertyUniqClientMessageIdKeyidx] == "" {
		msg.Properties[MessageConst.PropertyUniqClientMessageIdKeyidx] = m.createUniqID()
	}
}

func (m messageClientIDSetter) getUniqID(msg *Message) string {
	return msg.Properties[MessageConst.PropertyUniqClientMessageIdKeyidx]
}

func (m messageClientIDSetter) createUniqID() string {
	current := time.Now().UnixNano()
	if current > m.nextStartTime {
		m.setStartTime()
	}
	fmt.Printf("UnixNano:%s\n", time.Now().UnixNano()-m.startTime)
	binary.Write(m.stringBuilder, binary.BigEndian, time.Now().UnixNano()-m.startTime)
	m.counter++
	binary.Write(m.stringBuilder, binary.BigEndian, int32(m.counter))
	fmt.Printf("UniqKey:%s\n",  m.stringBuilder.String())
	return "C0A82B8D64C7355DA25450AF4DB50002"
	//return m.stringBuilder.String()
}

func (m messageClientIDSetter) setStartTime() {
	m.startTime = time.Now().UnixNano()
	m.nextStartTime = time.Now().UnixNano() + 2592000000000000 // next 30 days, 3600 * 24 * 30 * 1000 * 1000 *1000
}

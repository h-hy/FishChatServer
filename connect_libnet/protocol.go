package connect_libnet

import (
	"encoding/binary"
	"io"
	"fmt"
	"bytes"
	"errors"
)

var (
	BigEndian    = ByteOrder(binary.BigEndian)
	LittleEndian = ByteOrder(binary.LittleEndian)

	packet1BE = newSimpleProtocol(1, BigEndian)
	packet1LE = newSimpleProtocol(1, LittleEndian)
	packet2BE = newSimpleProtocol(2, BigEndian)
	packet2LE = newSimpleProtocol(2, LittleEndian)
	packet4BE = newSimpleProtocol(4, BigEndian)
	packet4LE = newSimpleProtocol(4, LittleEndian)
	packet8BE = newSimpleProtocol(8, BigEndian)
	packet8LE = newSimpleProtocol(8, LittleEndian)
)

type ByteOrder binary.ByteOrder

// Packet protocol.
type Protocol interface {
	// Create protocol state.
	// New(*Session) for session protocol state.
	// New(*Server) for server protocol state.
	// New(*Channel) for channel protocol state.
	New(interface{}) ProtocolState
}

// Protocol state.
type ProtocolState interface {
	// Packet a message.
	PrepareOutBuffer(buffer *OutBuffer, size int)

	// Write a packet.
	Write(writer io.Writer, buffer *OutBuffer) error

	// Read a packet.
	Read(reader io.Reader, buffer *InBuffer) error
}

// Create a {packet, N} protocol.
// The n means how many bytes of the packet header.
// n must is 1、2、4 or 8.
func PacketN(n int, byteOrder ByteOrder) Protocol {
	switch n {
	case 1:
		switch byteOrder {
		case BigEndian:
			return packet1BE
		case LittleEndian:
			return packet1LE
		}
	case 2:
		switch byteOrder {
		case BigEndian:
			return packet2BE
		case LittleEndian:
			return packet2LE
		}
	case 4:
		switch byteOrder {
		case BigEndian:
			return packet4BE
		case LittleEndian:
			return packet4LE
		}
	case 8:
		switch byteOrder {
		case BigEndian:
			return packet8BE
		case LittleEndian:
			return packet8LE
		}
	}
	panic("unsupported packet head size")
}

// The packet spliting protocol like Erlang's {packet, N}.
// Each packet has a fix length packet header to present packet length.
type simpleProtocol struct {
	n             int
	bo            binary.ByteOrder
	encodeHead    func([]byte)
	decodeHead    func([]byte) int
	MaxPacketSize int
}

func newSimpleProtocol(n int, byteOrder binary.ByteOrder) *simpleProtocol {
	protocol := &simpleProtocol{
		n:  n,
		bo: byteOrder,
	}

	switch n {
	case 1:
		protocol.encodeHead = func(buffer []byte) {
			buffer[0] = byte(len(buffer) - n)
		}
		protocol.decodeHead = func(buffer []byte) int {
			return int(buffer[0])
		}
	case 2:
		protocol.encodeHead = func(buffer []byte) {
			byteOrder.PutUint16(buffer, uint16(len(buffer)-n))
		}
		protocol.decodeHead = func(buffer []byte) int {
			return int(byteOrder.Uint16(buffer))
		}
	case 4:
		protocol.encodeHead = func(buffer []byte) {
			byteOrder.PutUint32(buffer, uint32(len(buffer)-n))
		}
		protocol.decodeHead = func(buffer []byte) int {
			return int(byteOrder.Uint32(buffer))
		}
	case 8:
		protocol.encodeHead = func(buffer []byte) {
			byteOrder.PutUint64(buffer, uint64(len(buffer)-n))
		}
		protocol.decodeHead = func(buffer []byte) int {
			return int(byteOrder.Uint64(buffer))
		}
	default:
		panic("unsupported packet head size")
	}

	return protocol
}

func (p *simpleProtocol) New(v interface{}) ProtocolState {
	return p
}

func (p *simpleProtocol) PrepareOutBuffer(buffer *OutBuffer, size int) {
	buffer.Prepare(size)
	buffer.Data = buffer.Data[:p.n]
}

func (p *simpleProtocol) Write(writer io.Writer, packet *OutBuffer) error {
	if p.MaxPacketSize > 0 && len(packet.Data) > p.MaxPacketSize {
		return PacketTooLargeError
	}
	p.encodeHead(packet.Data)
	if _, err := writer.Write(packet.Data); err != nil {
		return err
	}
	return nil
}

var SLC []byte = []byte{'J','H','U'}
var data_spilt []byte  = []byte{'#'}
var data_begin []byte  = []byte{'{','<'}
var data_end []byte  = []byte{'\r','\n'}
func checkPacket(buffer []byte,length int) (errcode int,err error) {
	var index,nowindex int

	index = bytes.Index(buffer,SLC)
	if (index==-1){
		return 2,nil //没找到协议头
	}else if (index > 0){
		buffer=buffer[index:]
	}
	//协议头定位完毕，开始查找数据头

	index = bytes.Index(buffer,data_begin)
	if (index==-1){
		return 1,nil //没找到数据头
	}
	nowindex=index+1
	//数据头定位完毕，开始检查协议完整性

	//检查协议完整性：命令编码
	index = bytes.Index(buffer[nowindex:],data_spilt)
	if (index==-1){
		return 1,nil
	}
	// var command_id string
	// command_id = string(buffer[nowindex+1:nowindex+index])
	nowindex=index+1
	
	index = bytes.Index(buffer[nowindex:],data_end)
	if (index==-1){
		return 1,nil //没找到数据尾
	}else{
		return 0,nil //没找到数据尾
	}
}

func (p *simpleProtocol) Read(reader io.Reader, buffer *InBuffer) error {
	buffer.Prepare(2048)//开辟2048字节数组
	var readedBuffer = make([]byte, 2048)
	var readedBufferLength int = 0
	for {
		length, err := reader.Read(buffer.Data) //读到了一部分长度
		if err != nil {
			return err
		}
		if (length + readedBufferLength > 2048){
			return errors.New("BUFFER TO LONG")
		}
		// readedBuffer[readedBufferLength:readedBufferLength+length]=buffer.Data[0:length]

		for i := 0; i < length; i++ {
			readedBuffer[i+readedBufferLength]=buffer.Data[i]
		}
		readedBufferLength+=length
		fmt.Printf("readedBuffer=%s",readedBuffer[0:readedBufferLength])

		errcode, err := checkPacket(readedBuffer,length) //checkPacket方法判断数据包是否完整
		fmt.Printf("errcode=%d",errcode)
		fmt.Printf("length=%d",length)
		fmt.Printf("\n\n")
        switch errcode {
            case 0:
            	buffer.Data = readedBuffer[0:length]
            	readedBufferLength=0
            	return nil //完整了，返回
            case 2:
            	readedBufferLength=0
            	break
            case 3:
            	return err //协议出错了
        }
	}
	return nil
}


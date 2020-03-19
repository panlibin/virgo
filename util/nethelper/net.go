package nethelper

import (
	"encoding/binary"
	"errors"
	"io"
)

// EndianType 大小端
type EndianType int32

// 大小端
const (
	LittleEndian EndianType = iota
	BigEndian
)

// MessageLengthSize 消息长度的长度
const MessageLengthSize int = 4

// MessageIDSize 消息ID的长度
const MessageIDSize int = 4

// MessageMaxLength 消息长度上限
const MessageMaxLength uint32 = 100000

// DefaultTCPRead 读取tcp socket默认方法
func DefaultTCPRead(r io.Reader, endianType EndianType) ([]byte, error) {
	bufMsgLen := make([]byte, MessageLengthSize)
	if _, err := io.ReadFull(r, bufMsgLen); err != nil {
		return nil, err
	}
	var msgLen uint32
	if endianType == LittleEndian {
		msgLen = binary.LittleEndian.Uint32(bufMsgLen)
	} else {
		msgLen = binary.BigEndian.Uint32(bufMsgLen)
	}

	if msgLen > MessageMaxLength {
		return nil, errors.New("message length out of range")
	}

	buf := make([]byte, msgLen)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// DefaultTCPWrite 写入tcp socket默认方法
func DefaultTCPWrite(msgID uint32, msgBuf []byte, endianType EndianType) ([]byte, error) {
	msgLen := len(msgBuf)
	writeBuf := make([]byte, msgLen+MessageIDSize+MessageLengthSize)

	if endianType == LittleEndian {
		binary.LittleEndian.PutUint32(writeBuf, uint32(msgLen+MessageIDSize))
		binary.LittleEndian.PutUint32(writeBuf[MessageLengthSize:], msgID)
	} else {
		binary.BigEndian.PutUint32(writeBuf, uint32(msgLen+MessageIDSize))
		binary.BigEndian.PutUint32(writeBuf[MessageLengthSize:], msgID)
	}

	copy(writeBuf[MessageIDSize+MessageLengthSize:], msgBuf)
	return writeBuf, nil
}

// DefaultWsWrite 写入websocket默认方法
func DefaultWsWrite(msgID uint32, msgBuf []byte, endianType EndianType) ([]byte, error) {
	msgLen := len(msgBuf)
	writeBuf := make([]byte, msgLen+MessageIDSize)

	if endianType == LittleEndian {
		binary.LittleEndian.PutUint32(writeBuf, msgID)
	} else {
		binary.BigEndian.PutUint32(writeBuf, msgID)
	}

	copy(writeBuf[MessageIDSize:], msgBuf)
	return writeBuf, nil
}

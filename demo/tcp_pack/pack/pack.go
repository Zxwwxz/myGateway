package pack

import (
	"encoding/binary"
	"errors"
	"io"
)

const MSG_HEADER = "12345678"

func Encode(writer io.Writer,msg []byte)(err error)  {
	headBuf := []byte(MSG_HEADER)
	if err = binary.Write(writer,binary.BigEndian,headBuf);err != nil {
		return err
	}
	msgLength := uint32(len(msg))
	if err = binary.Write(writer,binary.BigEndian,msgLength);err != nil {
		return err
	}
	if err = binary.Write(writer,binary.BigEndian,msg);err != nil {
		return err
	}
	return nil
}

func Decode(reader io.Reader)(msg []byte,err error) {
	headBuf := make([]byte,len(MSG_HEADER))
	if err := binary.Read(reader,binary.BigEndian,headBuf);err != nil {
		return nil,err
	}
	headMsg := string(headBuf)
	if headMsg != MSG_HEADER{
		return nil,errors.New("head err")
	}
	lenBuf := make([]byte,4)
	if err := binary.Read(reader,binary.BigEndian,lenBuf);err != nil {
		return nil,err
	}
	msgLength := binary.BigEndian.Uint32(lenBuf)
	msgBuf := make([]byte,msgLength)
	if err := binary.Read(reader,binary.BigEndian,msgBuf);err != nil {
		return nil,err
	}
	return msgBuf, nil
}


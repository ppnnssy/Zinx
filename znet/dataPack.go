package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinxProject/utils"
	"zinxProject/ziface"
)

type DataPack struct {
}

//拆包封包实例的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包的头的长度
func (dp *DataPack) GetHeadLen() uint32 {
	//DataLen uint32(4字节)+ID uint32(4字节)
	return 8
}

//封包方法
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//func NewBuffer(buf []byte) *Buffer
	//NewBuffer使用buf作为初始内容创建并初始化一个Buffer。本函数用于创建一个用于读取已存在数据的buffer；
	//也用于指定用于写入的内部缓冲的大小，此时，buf应为一个具有指定容量但长度为0的切片。buf会被作为返回值的底层缓冲切片。
	//大多数情况下，new(Buffer)（或只是声明一个Buffer类型变量）就足以初始化一个Buffer了。
	dataBuff := bytes.NewBuffer([]byte{})

	//将长度写进databuff中
	err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen())
	if err != nil {
		return nil, err
	}

	//将MsgID写进databuff中
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}

	//将data数据写进databuff中
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgData())
	if err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil

}

//拆包方法
//传进来的是二进制数据
//将包的Head信息读取出来，然后再根据Head信息里的data长度，进行一次读取
func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	dataBuff := bytes.NewReader(binaryData)

	msg := &Message{}

	//读取数据包长度
	err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		return nil, err
	}

	//读取msgID
	err = binary.Read(dataBuff, binary.LittleEndian, &msg.Id)
	if err != nil {
		return nil, err
	}

	//判断数据包datalen是否超出了最大长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too Large msg recv!")
	}

	return msg, err

}

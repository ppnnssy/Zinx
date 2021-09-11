package ziface

type IConnManager interface {
	//添加链接
	Add(conn IConnection)
	//删除链接
	Remove(conn IConnection)
	//根据链接ID获取链接
	Get(connId uint32) (IConnection, error)
	//得到当前的链接总数
	Len() int
	//清楚并终止所有链接。用于回收资源
	ClearConn()
}

package vgevent

// EventType 事件类型
type EventType int32

// IEvent 事件接口
type IEvent interface {
	GetType() EventType
}

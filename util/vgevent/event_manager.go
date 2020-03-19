package vgevent

// EventManager 事件管理器
type EventManager struct {
	pool       map[EventType]map[uint32]func(IEvent)
	eventIDMax uint32
}

// NewEventManager 创建事件管理器
func NewEventManager() *EventManager {
	pObj := new(EventManager)
	pObj.pool = make(map[EventType]map[uint32]func(IEvent))
	return pObj
}

// Register 注册事件回调
func (em *EventManager) Register(eventType EventType, f func(IEvent)) (eventID uint32) {
	mapListener, exist := em.pool[eventType]
	if !exist {
		mapListener = make(map[uint32]func(IEvent))
		em.pool[eventType] = mapListener
	}
	eventID = em.genEventID()
	mapListener[eventID] = f
	return
}

// Unregister 注销事件
func (em *EventManager) Unregister(eventType EventType, eventID uint32) {
	mapListener, exist := em.pool[eventType]
	if !exist {
		return
	}
	delete(mapListener, eventID)
}

// Clear 清空
func (em *EventManager) Clear() {
	em.pool = make(map[EventType]map[uint32]func(IEvent))
	em.eventIDMax = 0
}

// Dispatch 派发事件
func (em *EventManager) Dispatch(event IEvent) {
	mapListener, exist := em.pool[event.GetType()]
	if !exist {
		return
	}
	for _, f := range mapListener {
		f(event)
	}
}

func (em *EventManager) genEventID() uint32 {
	em.eventIDMax++
	return em.eventIDMax
}

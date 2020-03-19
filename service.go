package virgo

// IService 服务接口
type IService interface {
	OnInit(*Procedure)
	OnRelease()
}

// Launch 启动服务
func Launch(s IService) {
	p := NewProcedure(s)
	p.Start()
	p.waitQuit()
}

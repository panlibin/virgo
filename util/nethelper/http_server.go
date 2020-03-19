package nethelper

import (
	"net"
	"net/http"

	logger "github.com/panlibin/vglog"
)

type handlerWrapper struct {
	f func(w http.ResponseWriter, pReq *http.Request)
}

func (h *handlerWrapper) ServeHTTP(w http.ResponseWriter, pReq *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	pReq.ParseForm()
	h.f(w, pReq)
}

// HTTPServer http服务器
type HTTPServer struct {
	server *http.Server
	router *http.ServeMux
}

// NewHTTPServer 新建http服务器
func NewHTTPServer() *HTTPServer {
	pObj := new(HTTPServer)
	pObj.server = new(http.Server)
	pObj.router = http.NewServeMux()

	return pObj
}

// Start 启动
func (s *HTTPServer) Start(addr string, certFile string, keyFile string) error {
	logger.Infof("start http server")
	s.server.Addr = addr
	s.server.Handler = s.router

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Errorf("start http server error: %v", err)
		return err
	}

	go func() {
		if certFile == "" || keyFile == "" {
			s.server.Serve(ln)
		} else {
			s.server.ServeTLS(ln, certFile, keyFile)
		}
	}()

	logger.Infof("http server listen on %s", addr)

	return err
}

// Stop 停止
func (s *HTTPServer) Stop() {
	s.server.Close()
}

// Handle 注册
func (s *HTTPServer) Handle(pattern string, f func(w http.ResponseWriter, pReq *http.Request)) {
	s.router.Handle(pattern, &handlerWrapper{f})
}

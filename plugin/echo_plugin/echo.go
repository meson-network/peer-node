package echo_plugin

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/coreservice-io/echo_middleware"
	"github.com/coreservice-io/echo_middleware/tool"
	"github.com/coreservice-io/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type EchoServer struct {
	*echo.Echo
	Http_port int
	Logger    log.Logger
}

var instanceMap = map[string]*EchoServer{}

func GetInstance() *EchoServer {
	return instanceMap["default"]
}

func GetInstance_(name string) *EchoServer {
	return instanceMap[name]
}

/*
http_port
*/
type Config struct {
	Port int
}

func Init(serverConfig Config, OnPanicHanlder func(panic_err interface{}), logger log.Logger) error {
	return Init_("default", serverConfig, OnPanicHanlder, logger)
}

// Init a new instance.
//  If only need one instance, use empty name "". Use GetDefaultInstance() to get.
//  If you need several instance, run Init() with different <name>. Use GetInstance(<name>) to get.
func Init_(name string, serverConfig Config, OnPanicHanlder func(panic_err interface{}), logger log.Logger) error {
	if name == "" {
		name = "default"
	}

	_, exist := instanceMap[name]
	if exist {
		return fmt.Errorf("echo server instance <%s> has already been initialized", name)
	}

	if serverConfig.Port == 0 {
		serverConfig.Port = 8080
	}

	echoServer := &EchoServer{
		echo.New(),
		serverConfig.Port,
		logger,
	}

	//cros
	echoServer.Use(middleware.CORS())
	//logger
	echoServer.Use(echo_middleware.LoggerWithConfig(echo_middleware.LoggerConfig{
		Logger:            logger,
		RecordFailRequest: false,
	}))
	//recover and panicHandler
	echoServer.Use(echo_middleware.RecoverWithConfig(echo_middleware.RecoverConfig{
		OnPanic: OnPanicHanlder,
	}))
	echoServer.JSONSerializer = tool.NewJsoniter()

	instanceMap[name] = echoServer
	return nil
}

func (s *EchoServer) Start() error {
	s.Logger.Infoln("http server started on port :" + strconv.Itoa(s.Http_port))
	return s.Echo.Start(":" + strconv.Itoa(s.Http_port))
}

func (s *EchoServer) Close() {
	s.Echo.Close()
}

//check the server is indeed up
func (s *EchoServer) CheckStarted() bool {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			addr_tls := s.Echo.TLSListenerAddr()
			if addr_tls != nil && strings.Contains(addr_tls.String(), ":") {
				return true
			}
			addr := s.Echo.ListenerAddr()
			if addr != nil && strings.Contains(addr.String(), ":") {
				return true
			}
		}
	}
}

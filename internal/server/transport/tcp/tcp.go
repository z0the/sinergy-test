package tcp

import (
	"net"

	"go.uber.org/zap"

	"sinergy-test/internal/server/service"
)

func NewController(logger *zap.SugaredLogger, service service.Service) Controller {
	return &controller{
		logger:  logger,
		service: service,
	}
}

type controller struct {
	listener net.Listener
	service  service.Service
	logger   *zap.SugaredLogger
}

func (s *controller) Run(port string) error {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Warn("Recovered in controller: ", r)
		}
		s.logger.Info("Server has stopped...")
	}()

	err := s.loadListener(port)
	if err != nil {
		return err
	}

	s.logger.Infof("Starting controller on port %s...", port)
	for {
		conn, err := s.listener.Accept()

		if err != nil {
			switch typedErr := err.(type) {
			case *net.OpError:
				if typedErr.Timeout() {
					continue
				}
			default:
				s.logger.Error("Error during client conn attempt: ", err)
			}
		}
		s.logger.Info("Starting new client...")
		go s.service.ServeConn(conn)
	}
}

func (s *controller) loadListener(port string) error {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	s.listener = ln
	return nil
}

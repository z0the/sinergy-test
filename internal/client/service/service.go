package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"go.uber.org/zap"
)

func NewService(logger *zap.SugaredLogger, resourceServerHost, resourceServerPort string) Service {
	s := &service{
		logger:             logger,
		resourceServerHost: resourceServerHost,
		resourceServerPort: resourceServerPort,
	}
	err := s.listenUpdateServer()
	if err != nil {
		panic(err)
	}
	return s
}

type service struct {
	logger               *zap.SugaredLogger
	resourceServerHost   string
	resourceServerPort   string
	latestResourceUpdate *ResourceUpdate
	sync.RWMutex
	curConn net.Conn
}

func (s *service) GetLastAction(_ context.Context) (string, error) {
	s.RLock()
	defer s.RUnlock()
	if s.latestResourceUpdate.Action == "" {
		return "", errors.New("no active action available")
	}
	return s.latestResourceUpdate.Action, nil
}

func (s *service) listenUpdateServer() error {
	var decoder *json.Decoder
	var err error
	decoder, err = s.establishNewConnection()
	if err != nil {
		s.logger.Errorw("failed to establish new connection", "err", err)
		return err
	}
	go func() {
		for {
			update := new(ResourceUpdate)
			err := decoder.Decode(update)
			if err != nil {
				s.logger.Errorw("failed to decode resource update", "err", err)
				if errors.Is(err, io.EOF) {
					break
				}
				decoder, err = s.establishNewConnection()
				if err != nil {
					s.logger.Errorw("failed to reestablish new connection", "err", err)
					return
				}
				time.Sleep(time.Second * 2)
				continue
			}
			time.Sleep(time.Second * 2)
			s.Lock()
			s.latestResourceUpdate = update
			s.logger.Infow("updated resource", "action", update.Action)
			s.Unlock()
		}
	}()
	return nil
}

func (s *service) establishNewConnection() (*json.Decoder, error) {
	if s.curConn != nil {
		err := s.curConn.Close()
		if err != nil {
			return nil, err
		}
	}

	address := fmt.Sprintf("%s:%s", s.resourceServerHost, s.resourceServerPort)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	s.logger.Warn("establish connection successfully")
	s.curConn = conn
	dec := json.NewDecoder(conn)

	return dec, nil
}

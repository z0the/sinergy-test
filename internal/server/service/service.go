package service

import (
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func NewService(logger *zap.SugaredLogger, resourceUrls []string, maxSleepTimeSeconds int) Service {
	s := &service{
		logger:              logger,
		resourceUrls:        resourceUrls,
		maxSleepTimeSeconds: maxSleepTimeSeconds,
		connPool:            make(map[net.Conn]struct{}),
	}
	s.startResourcePooling()
	return s
}

type service struct {
	sync.Mutex
	logger              *zap.SugaredLogger
	resourceUrls        []string
	curURLIndex         int
	maxSleepTimeSeconds int
	connPool            map[net.Conn]struct{}
}

func (s *service) ServeConn(conn net.Conn) {
	s.Lock()
	defer s.Unlock()
	s.connPool[conn] = struct{}{}
}

func (s *service) startResourcePooling() {
	if len(s.resourceUrls) == 0 {
		panic("no resourceUrls provided")
	}

	go func() {
		for {
			url := s.resourceUrls[s.curURLIndex]
			s.rotateResourceURLIndex()

			resp, err := http.Get(url) //nolint:gosec
			if err != nil {
				s.logger.Errorw("failed to get url", "url", url, "err", err)
				continue
			}

			s.handleResp(resp)
			time.Sleep(time.Second * time.Duration(1+rand.Intn(s.maxSleepTimeSeconds))) //nolint:gosec
		}
	}()
}
func (s *service) handleResp(resp *http.Response) {
	defer s.closeBody(resp.Body)

	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.logger.Errorw("failed to read response body", "err", err)
		return
	}

	if len(rawData) == 0 {
		return
	}

	s.logger.Infof("updating resource: %s", string(rawData))

	wg := new(sync.WaitGroup)
	for conn := range s.connPool {
		wg.Add(1)
		go s.sendUpdate(wg, conn, rawData)
	}
	wg.Wait()
}

func (s *service) sendUpdate(wg *sync.WaitGroup, conn net.Conn, data []byte) {
	defer wg.Done()
	_, err := conn.Write(data)
	if err != nil {
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) || errors.Is(err, syscall.EPIPE) {
			delete(s.connPool, conn)
			return
		}
		s.logger.Errorw("failed to write to conn", "err", err)
	}
}

func (s *service) rotateResourceURLIndex() {
	if s.curURLIndex == len(s.resourceUrls)-1 {
		s.curURLIndex = 0
	} else {
		s.curURLIndex++
	}
}

func (s *service) closeBody(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		s.logger.Errorw("failed to close body", "err", err)
	}
}

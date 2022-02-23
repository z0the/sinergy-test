package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"sinergy-test/internal/client/service"
)

func NewController(logger *zap.SugaredLogger, srv service.Service) Controller {
	return &controller{
		logger:  logger,
		service: srv,
	}
}

type controller struct {
	logger  *zap.SugaredLogger
	service service.Service
}

func (c *controller) Run(port string) error {
	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/getLastAction").HandlerFunc(c.GetLastAction)
	c.logger.Infow("Starting controller", "port", port)
	return http.ListenAndServe(":"+port, router)
}

func (c *controller) GetLastAction(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	action, err := c.service.GetLastAction(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	resp := &GetLastActionResponse{Action: action}
	err = enc.Encode(resp)
	if err != nil {
		c.logger.Errorw("failed to encode GetLastActionResponse", "err", err)
	}
}

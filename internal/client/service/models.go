package service

type ResourceUpdate struct {
	Action string                 `json:"action"`
	Type   string                 `json:"type"`
	Data   map[string]interface{} `json:"data"`
}

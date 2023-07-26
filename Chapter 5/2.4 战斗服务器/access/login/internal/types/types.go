package types

type ConnectLogicRequest struct {
	Token    string `json:"token"`
	Platform string `json:"platform"`
}

type ConnectLogicResponse struct {
	Code         int32  `json:"code"`
	GateEndpoint string `json:"gate_endpoint"`
	GateToken    string `json:"gate_token"`
}

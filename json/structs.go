package json

type ScardRequest struct {
	Method string `json:"method"`
}

type ScardResponse struct {
	Error string `json:"error"`
}

type ScardVersionResponse struct {
	ScardResponse
	Version string `json:"version"`
}

type ScardContextResponse struct {
	ScardResponse
	Ctx Context `json:"ctx"`
}

type ScardCtxRequest struct {
	ScardRequest
	Ctx Context `json:"ctx"`
}

type ScardListReadersResponse struct {
	ScardResponse
	Readers []string `json:"readers"`
}

type ScardConnectRequest struct {
	ScardRequest
	Ctx Context					`json:"ctx"`
	Reader    string    `json:"reader"`
	ShareMode ShareMode `json:"shareMode"`
	Protocol  Protocol  `json:"protocol"`
}

type ScardConnectResponse struct {
	ScardResponse
	Card Card `json:"card"`
}

type ScardStatusRequest struct {
	ScardRequest
	Card Card `json:"card"`
}

type ScardStatusResponse struct {
	ScardResponse
	Card           Card     `json:"card"`
	Reader         string   `json:"reader"`
	State          uint32    `json:"state"`
	ActiveProtocol Protocol `json:"activeProtocol"`
	ATR            string   `json:"atr"`
}

type ScardDisconnectRequest struct {
	ScardRequest
	Card Card `json:"card"`
	Disposition Disposition `json:"disposition"`
}

type ScardTransmitRequest struct {
	ScardRequest
	Card Card `json:"card"`
	Data string `json:"data"`
}

type ScardTransmitResponse struct {
	ScardResponse
	Card           Card     `json:"card"`
	Data string `json:"data"`
	}

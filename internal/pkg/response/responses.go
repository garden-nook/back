package response

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type CreateUpdateIntId struct {
	Id int32 `json:"id"`
}

type CreateUpdateUuidId struct {
	Id string `json:"id"`
}

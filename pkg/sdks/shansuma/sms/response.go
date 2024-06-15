package sms

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Result  struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			MessageID string `json:"message_id"`
			Total     int    `json:"total"`
		} `json:"data"`
	} `json:"result"`
}

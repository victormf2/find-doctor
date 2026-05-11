package send_user_message

type Input struct {
	UserID  string `json:"userId"`
	Message string `json:"message"`
}
type Output struct {
	MessageID string `json:"messageId"`
}

type Operation struct{}

func NewOperation() *Operation {
	return &Operation{}
}

func (h *Operation) Do(input *Input) (*Output, error) {
	return nil, nil
}

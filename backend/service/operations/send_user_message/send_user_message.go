package send_user_message

type Input struct {
	UserID  string
	Message string
}
type Output struct {
	MessageID string
}

type Operation struct{}

func NewOperation() *Operation {
	return &Operation{}
}

func (h *Operation) Do(input *Input) (*Output, error) {
	return nil, nil
}

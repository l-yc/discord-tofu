package brain

// interface with tofu-ai
const (
	AIMessageTypeGroup = "group message"
	AIMessageTypeDirect = "private message"
	AIMessageTypeStatus = "status"
)
type AIMessage struct {
	Type string `json:"type"`
	Contents string `json:"contents,omitempty"`
}

type AIReply struct {
	StatusMessage string `json:"statusMessage"`
	PrimaryMood float32 `json:"primaryMood"`
	MoodStability float32 `json:"moodStability"`
	ExposedPositivity float32 `json:"exposedPositivity"`
	PositivityOverload bool `json:"positivityOverload"`
	Response string `json:"response"`
}

type AIError struct {
	ErrorMessage string `json:"error"`
	Response *string `json:"response"` // can be null
}


// brain io interface
const (
	BrainInputTypeMessageGroup = iota
	BrainInputTypeMessageDirect
	BrainInputTypeStatus
)
type BrainInput struct {
	Type int
	Content string
}

type BrainOutput struct {
	Error error
	Content string
}

const (
	BrainStateStatusOnline = iota
	BrainStateStatusIdle
	BrainStateStatusBusy
)
type BrainState struct {
	Mood string
	Status int
}

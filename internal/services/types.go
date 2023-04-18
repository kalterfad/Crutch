package services

type ChannelID string

type DirData struct {
	ChannelID `json:"channel_id"`
	Path      string            `json:"path"`
	Rules     map[string]string `json:"rules"`
}

type DirPathData struct {
	ChannelID
	OldPath  string
	NewPath  string
	FileName string
}

type SocketMsg struct {
	Command   string `json:"command"`
	ChannelID `json:"channel_id"`
	Path      string            `json:"path"`
	Rules     map[string]string `json:"rules"`
}

const NewWorker = "NEW_WORKER"
const RmWorker = "RM_WORKER"

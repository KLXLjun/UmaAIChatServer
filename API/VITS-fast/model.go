package vitsfast

type VitsPostData struct {
	Data        []interface{} `json:"data"`
	Index       int           `json:"fn_index"`
	SessionHash string        `json:"session_hash"`
}

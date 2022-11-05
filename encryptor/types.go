package encryptor

type Contents struct {
	Content []Content `json:"content"`
}

type Content struct {
	Path string `json:"path"`
	Key  string `json:"key"`
}

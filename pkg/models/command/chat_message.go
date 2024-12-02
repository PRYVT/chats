package command

type AddChatMessage struct {
	Text        string `json:"text"`
	ImageBase64 string `json:"image_base64"`
}

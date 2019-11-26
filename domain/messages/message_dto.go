package messages

type Message struct {
	Id int64 `json:"id"`
	Title string `json:"title"`
	Body string `json:"body"`
	CreatedAt string `json:"created_at"`
}
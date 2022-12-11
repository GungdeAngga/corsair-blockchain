package blockdata

type Block struct {
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Amount     int64  `json:"amount"`
	Timestamp  int64  `json:"timestamp"`
	Hash       string `json:"hash"`
	PrevHash   string `json:"prev_hash"`
}

type Blockchain []Block

type SearchByHash struct {
	Hash string `json:"hash"`
}

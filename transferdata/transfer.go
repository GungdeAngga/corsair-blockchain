package transferdata

type TransferData struct {
	From   int   `json:"from"`
	To     int   `json:"to"`
	Amount int64 `json:"amount"`
}

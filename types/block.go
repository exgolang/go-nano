package types

type Block struct {
	Index        int           `json:"index"`
	Transactions []interface{} `json:"transactions"`
	Prev         string        `json:"prev"`
	Current      string        `json:"current"`
	Fees         int           `json:"fees,omitempty"`
	Timestamp    int64         `json:"timestamp"`
}

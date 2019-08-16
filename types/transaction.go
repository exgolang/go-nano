package types

type (
	Transfer struct {
		Root  string `json:"root"`
		Input struct {
			From    string `json:"from,omitempty"`
			To      string `json:"to"`
			Value   int    `json:"value"`
			Message string `json:"message,omitempty"`
		} `json:"input"`
		Signature string `json:"signature,omitempty"`
		Timestamp int64  `json:"timestamp"`
	}

	Account struct {
		Master   string `json:"master"`
		Name     string `json:"name,omitempty"`
		Registry []struct {
			Label   string `json:"label"`
			Address string `json:"address"`
		} `json:"registry"`
	}
)

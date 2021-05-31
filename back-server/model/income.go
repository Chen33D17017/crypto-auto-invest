package model

type Income struct {
	CryptoName        string  `json:"crypto_name"`
	Strategy          int     `json:"strategy"`
	Amount            float64 `json:"amount"`
	Cost              float64 `json:"cost"`
	JPY               float64 `json:"JPY"`
	IncomeRate        string  `json:"income_rate"`
	Deposit           float64 `json:"deposit"`
	DepositIncomeRate string  `json:"deposit_income_rate"`
}

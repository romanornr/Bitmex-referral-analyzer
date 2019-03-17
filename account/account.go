package account

type Transaction struct {
	Time          string
	Type          string
	Amount        float64
	Fee           float64
	Address       string
	Status        string
	WalletBalance string
}

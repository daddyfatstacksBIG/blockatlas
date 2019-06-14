package blockscout

type TransactionResult struct {
	Result []Transaction `json:"result"`
}

type TokenTransactionResult struct {
	Result []TokenTransaction `json:"result"`
}

type Transaction struct {
	Value           string `json:"value"`
	TimeStamp       string `json:"timeStamp"`
	Nonce           string `json:"nonce"`
	IsError         string `json:"isError"`
	Hash            string `json:"hash"`
	GasUsed         string `json:"gasUsed"`
	GasPrice        string `json:"gasPrice"`
	Gas             string `json:"gas"`
	From            string `json:"from"`
	To              string `json:"to"`
	ContractAddress string `json:"contractAddress"`
	BlockNumber     string `json:"blockNumber"`
}

type TokenTransaction struct {
	Value           string `json:"value"`
	TokenSymbol     string `json:"tokenSymbol"`
	TokenName       string `json:"tokenName"`
	TokenDecimal    string `json:"tokenDecimal"`
	TimeStamp       string `json:"timeStamp"`
	Nonce           string `json:"nonce"`
	Hash            string `json:"hash"`
	GasUsed         string `json:"gasUsed"`
	GasPrice        string `json:"gasPrice"`
	Gas             string `json:"gas"`
	From            string `json:"from"`
	To              string `json:"to"`
	ContractAddress string `json:"contractAddress"`
	BlockNumber     string `json:"blockNumber"`
}

package blockscout

import (
	"math/big"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/trustwallet/blockatlas/models"
	"github.com/trustwallet/blockatlas/util"
)

// MakeSetup returns a function used to register an Ethereum-based platform route
func MakeSetup(coinIndex uint, platform string) func(gin.IRouter) {
	apiKey := platform + ".api"

	client := Client{
		HTTPClient: http.DefaultClient,
	}

	return func(router gin.IRouter) {
		router.Use(util.RequireConfig(apiKey))
		router.Use(func(c *gin.Context) {
			client.BaseURL = viper.GetString(apiKey)
			c.Next()
		})
		router.GET("/:address", func(c *gin.Context) {
			GetTransactions(c, coinIndex, &client)
		})
	}
}

func GetTransactions(c *gin.Context, coinIndex uint, client *Client) {
	token := c.Query("token")
	address := c.Param("address")
	var srcPage TransactionResult
	var err error

	if token != "" {
		//srcPage, err = client.GetTxsWithContract(address, token)
	} else {
		srcPage, err = client.GetTxs(address)
	}

	if apiError(c, err) {
		return
	}

	var txs []models.Tx
	for _, srcTx := range srcPage.Result {
		txs = AppendTxs(txs, &srcTx, coinIndex)
	}

	page := models.Response(txs)
	page.Sort()
	c.JSON(http.StatusOK, &page)
}

//0xa9059cbb

func extractBase(srcTx *Transaction, coinIndex uint) (base models.Tx, ok bool) {
	var status, errReason string
	if srcTx.IsError == "0" {
		status = models.StatusCompleted
	} else {
		status = models.StatusFailed
		errReason = srcTx.IsError
	}

	unix, err := strconv.ParseInt(srcTx.TimeStamp, 10, 64)
	if err != nil {
		return base, false
	}
	block, err := strconv.ParseUint(srcTx.BlockNumber, 10, 64)
	if err != nil {
		return base, false
	}
	sequence, err := strconv.ParseUint(srcTx.Nonce, 10, 64)
	if err != nil {
		return base, false
	}

	fee := calcFee(srcTx.GasPrice, srcTx.GasUsed)

	base = models.Tx{
		ID:       srcTx.Hash,
		Coin:     coinIndex,
		From:     srcTx.From,
		To:       srcTx.To,
		Fee:      models.Amount(fee),
		Date:     unix,
		Block:    block,
		Status:   status,
		Error:    errReason,
		Sequence: sequence,
	}
	return base, true
}

func AppendTxs(in []models.Tx, srcTx *Transaction, coinIndex uint) (out []models.Tx) {
	out = in
	baseTx, ok := extractBase(srcTx, coinIndex)
	if !ok {
		return
	}

	// Native ETH transaction
	//if len(srcTx.Ops) == 0 && srcTx.Input == "0x" {
	transferTx := baseTx
	transferTx.Meta = models.Transfer{
		Value: models.Amount(srcTx.Value),
	}
	out = append(out, transferTx)

	// 	// Smart Contract Call
	// 	if len(srcTx.Ops) == 0 && srcTx.Input != "0x" {
	// 		contractTx := baseTx
	// 		contractTx.Meta = models.ContractCall{
	// 			Input: srcTx.Input,
	// 			Value: srcTx.Value,
	// 		}
	// 		out = append(out, contractTx)
	// 	}

	// 	if len(srcTx.Ops) == 0 {
	// 		return
	// 	}
	// 	op := &srcTx.Ops[0]

	// 	if op.Type == models.TxTokenTransfer {
	// 		tokenTx := baseTx

	// 		tokenTx.Meta = models.TokenTransfer{
	// 			Name:     op.Contract.Name,
	// 			Symbol:   op.Contract.Symbol,
	// 			TokenID:  op.Contract.Address,
	// 			Decimals: op.Contract.Decimals,
	// 			Value:    models.Amount(op.Value),
	// 			From:     op.From,
	// 			To:       op.To,
	// 		}
	// 		out = append(out, tokenTx)
	// 	}
	return out
}

func calcFee(gasPrice string, gasUsed string) string {
	var gasPriceBig, gasUsedBig, feeBig big.Int

	gasPriceBig.SetString(gasPrice, 10)
	gasUsedBig.SetString(gasUsed, 10)

	feeBig.Mul(&gasPriceBig, &gasUsedBig)

	return feeBig.String()
}

func apiError(c *gin.Context, err error) bool {
	if err != nil {
		logrus.WithError(err).Errorf("Unhandled error")
		c.AbortWithStatus(http.StatusInternalServerError)
		return true
	}
	return false
}

package blockscout

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/trustwallet/blockatlas/models"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
}

func (c *Client) GetTxs(address string) (TransactionResult, error) {
	res, err := c.getTxs(fmt.Sprintf("%s/api?module=account&action=txlist&%s",
		c.BaseURL,
		url.Values{
			"address": {address},
		}.Encode()))
	var txs TransactionResult
	if err != nil {
		return txs, err
	}
	if res != nil {
		defer res.Body.Close()
	}
	err = json.NewDecoder(res.Body).Decode(&txs)
	return txs, err
}

func (c *Client) GetTxsWithContract(address, contract string) (TokenTransactionResult, error) {
	res, err := c.getTxs(fmt.Sprintf("%s/api?module=account&action=tokentx&%s",
		c.BaseURL,
		url.Values{
			"address":         {address},
			"contractaddress": {contract},
		}.Encode()))
	var txs TokenTransactionResult
	if err != nil {
		return txs, err
	}
	if res != nil {
		defer res.Body.Close()
	}
	err = json.NewDecoder(res.Body).Decode(&txs)

	return txs, err
}

func (c *Client) getTxs(uri string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", uri, nil)
	println(uri)
	res, err := c.HTTPClient.Do((req))
	if err != nil {
		logrus.WithError(err).Error("Ethereum/Trust Ray: Failed to get transactions")
		return nil, models.ErrSourceConn
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("http %s", res.Status)
	}
	return res, nil
}

package bitbank

import (
	"bytes"
	"crypto-auto-invest/bitbank/model"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func readUTC(timestamp int64) string {
	return time.Unix(timestamp/1000, 0).Format("2006-01-02")
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func encode(s model.Secret, content string) string {
	h := hmac.New(sha256.New, []byte(s.ApiSecret))
	h.Write([]byte(content))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func addHeader(req *http.Request, s model.Secret, content string) {
	nonce := fmt.Sprint(makeTimestamp())
	req.Header.Add("ACCESS-KEY", s.ApiKey)
	req.Header.Add("ACCESS-NONCE", nonce)
	req.Header.Add("ACCESS-SIGNATURE", encode(s, nonce+content))
	req.Header.Add("Content-Type", "application/json")
}

func apiRequest(req *http.Request, response interface{}) error {
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Fail to request: %s", err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Fail to read response body: %s", err.Error())
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	return nil
}

func getRequest(s model.Secret, query string, response interface{}) error {
	url := fmt.Sprintf("https://api.bitbank.cc%s", query)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Fail to build GET request: %s", err.Error())
	}

	addHeader(req, s, query)
	err = apiRequest(req, response)
	if err != nil {
		return err
	}
	return nil
}

func postRequest(s model.Secret, endpoint string, payload []byte, response interface{}) error {
	url := fmt.Sprintf("https://api.bitbank.cc%s", endpoint)
	payloadReader := bytes.NewReader(payload)
	req, err := http.NewRequest("POST", url, payloadReader)
	if err != nil {
		return fmt.Errorf("Fail to build POST request: %s", err)
	}

	addHeader(req, s, string(payload))
	err = apiRequest(req, response)
	if err != nil {
		return err
	}
	return nil
}

func CheckAssets(s model.Secret) ([]model.Asset, error) {
	var response model.AssetRst
	err := getRequest(s, "/v1/user/assets", &response)
	if err != nil {
		return response.Data.Assets, err
	}
	return response.Data.Assets, nil
}

func MakeTrade(s model.Secret, cryptoName string, buy_sell string, amount float64, tradeType string, postOnly bool) (model.Order, error) {
	url := fmt.Sprintf("/v1/user/spot/order")
	var response model.OrderRst

	order := model.OrderRequest{
		Pair:     fmt.Sprintf("%s_jpy", cryptoName),
		Amount:   fmt.Sprintf("%.4f", amount),
		Side:     buy_sell,
		Type:     tradeType,
		PostOnly: postOnly,
	}

	reqBody, _ := json.Marshal(order)
	err := postRequest(s, url, reqBody, &response)
	if err != nil {
		return response.Data, err
	}

	return response.Data, nil
}

func BuyWithJPY(s model.Secret, cryptoName string, JPY int64) (model.Order, error) {
	cryptmsg, err := GetPrice(cryptoName)
	if err != nil {
		fmt.Println(err.Error())
	}
	cryptPrice, _ := strconv.Atoi(cryptmsg.Buy)
	amount := float64(JPY) / float64(cryptPrice)

	return MakeTrade(s, cryptoName, "buy", amount, "market", false)
}

func SellToJPY(s model.Secret, cryptoName string, amount float64) (model.Order, error) {
	return MakeTrade(s, cryptoName, "sell", amount, "market", false)
}

func GetTradeHistory(s model.Secret, cryptoName string) ([]model.Trade, error) {
	var response model.TradeRst
	url := fmt.Sprintf("/v1/user/spot/trade_history?pair=%s_jpy", cryptoName)
	err := getRequest(s, url, &response)
	if err != nil {
		return nil, err
	}
	return response.Data.Trades, nil
}

func GetOrderInfo(s model.Secret, cryptoName, order_id string) (model.Order, error) {
	var response model.OrderRst
	url := fmt.Sprintf("/v1/user/spot/order?pair=%s_jpy&order_id=%s", cryptoName, order_id)
	err := getRequest(s, url, &response)
	if err != nil {
		return response.Data, err
	}
	return response.Data, nil
}

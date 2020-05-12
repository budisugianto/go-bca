package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/budisugianto/go-bca/auth"
	"github.com/budisugianto/go-bca/bca"
	"github.com/budisugianto/go-bca/business"
)

func main() {
	cfg := bca.Config{
		URL:          "https://devapi.klikbca.com",
		ClientID:     "dcc99ba6-3b2f-479b-9f85-86a09ccaaacf",
		ClientSecret: "5e636b16-df7f-4a53-afbe-497e6fe07edc",
		APIKey:       "b095ac9d-2d21-42a3-a70c-4781f4570704",
		APISecret:    "bedd1f8d-3bd6-4d4a-8cb4-e61db41691c9",
		CompanyCode:  "80888",
		CorporateID:  "h2hauto008", //Based on API document
		OriginHost:   "localhost",
		LogLevel:     3,
	}
	businessClient := business.NewClient(cfg)
	authClient := auth.NewClient(cfg)

	ctx := context.Background()
	authToken, err := authClient.GetToken(ctx)
	if err != nil {
		panic(err)
	}

	businessClient.AccessToken = authToken.AccessToken

	getBalanceInfo(ctx, businessClient)
	getAccountStatement(ctx, businessClient)
	//fundTransfer(ctx, businessClient)
}

func getInquiry(ctx context.Context, client business.Client) {
	response, err := client.GetBalanceInfo(ctx, []string{"8161964775"})
	if err != nil {
		panic(err)
	}
	fmt.Println(response)
}
func getBalanceInfo(ctx context.Context, client business.Client) {
	response, err := client.GetBalanceInfo(ctx, []string{"8161964775"})
	if err != nil {
		panic(err)
	}
	if len(response.AccountDetailDataFailed) > 0 {
		for i, account := range response.AccountDetailDataFailed {
			fmt.Printf("%d - Error: %s - %s", i, account.English, account.Indonesian)
		}
		return
	}
	for i, account := range response.AccountDetailDataSuccess {
		jsonStr, _ := json.Marshal(account)
		fmt.Printf("%d - Account: %s", i, jsonStr)
	}
}

func getAccountStatement(ctx context.Context, client business.Client) {
	startDate, err := time.Parse("2006-01-02", "2016-08-29")
	if err != nil {
		panic(err)
	}

	endDate, err := time.Parse("2006-01-02", "2016-09-01")
	if err != nil {
		panic(err)
	}

	response, err := client.GetAccountStatement(ctx, "0201245680", startDate, endDate)
	if err != nil {
		panic(err)
	}

	jsonStr, _ := json.Marshal(response)
	fmt.Printf("Statement: %s", jsonStr)
}

func fundTransfer(ctx context.Context, client business.Client) {
	response, err := client.FundTransfer(ctx, bca.FundTransferRequest{
		CorporateID:              "BCAAPI2016",
		SourceAccountNumber:      "0201245680",
		TransactionID:            "00000001",
		TransactionDate:          "2016-01-30",
		ReferenceID:              "12345/PO/2016",
		CurrencyCode:             "IDR",
		Amount:                   float64(100000.00),
		BeneficiaryAccountNumber: "0201245681",
		Remark1:                  "Transfer Test",
		Remark2:                  "Online Transfer",
	})

	if err != nil {
		panic(err)
	}

	jsonStr, _ := json.Marshal(response)
	fmt.Printf("Statement: %s", jsonStr)
}

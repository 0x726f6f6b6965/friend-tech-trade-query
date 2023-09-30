package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/0x726f6f6b6965/friend-tech-trade-query/api/internal/helper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TradeRecordRequest struct {
	TxHash string `json:"tx_hash,omitempty"`
}

type TradeRecordResponse struct {
	TxHash            string `json:"-" dynamodbav:"tx_hash"`
	Trader            string `json:"trader,omitempty" abi:"trader" dynamodbav:"trader"`
	Subject           string `json:"subject,omitempty" abi:"subject" dynamodbav:"subject"`
	IsBuy             bool   `json:"is_buy" abi:"isBuy" dynamodbav:"isBuy"`
	ShareAmount       string `json:"share_amount,omitempty" abi:"shareAmount" dynamodbav:"shareAmount"`
	EthAmount         string `json:"eth_amount,omitempty" abi:"ethAmount" dynamodbav:"ethAmount"`
	ProtocolEthAmount string `json:"protocol_eth_amount,omitempty" abi:"protocolEthAmount" dynamodbav:"protocolEthAmount"`
	SubjectEthAmount  string `json:"subject_eth_amount,omitempty" abi:"subjectEthAmount" dynamodbav:"subjectEthAmount"`
	Supply            string `json:"supply,omitempty" abi:"supply" dynamodbav:"supply"`
	TTL               int64  `json:"-" dynamodbav:"ttl"`
}

const (
	RPC_ENDPOINT      = "https://mainnet.base.org"
	ABI_FILE          = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"trader","type":"address"},{"indexed":false,"internalType":"address","name":"subject","type":"address"},{"indexed":false,"internalType":"bool","name":"isBuy","type":"bool"},{"indexed":false,"internalType":"uint256","name":"shareAmount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"ethAmount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"protocolEthAmount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"subjectEthAmount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"supply","type":"uint256"}],"name":"Trade","type":"event"},{"inputs":[{"internalType":"address","name":"sharesSubject","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"buyShares","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"address","name":"sharesSubject","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"getBuyPrice","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"sharesSubject","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"getBuyPriceAfterFee","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"supply","type":"uint256"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"getPrice","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"address","name":"sharesSubject","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"getSellPrice","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"sharesSubject","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"getSellPriceAfterFee","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"protocolFeeDestination","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"protocolFeePercent","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"sharesSubject","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"sellShares","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"address","name":"_feeDestination","type":"address"}],"name":"setFeeDestination","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_feePercent","type":"uint256"}],"name":"setProtocolFeePercent","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_feePercent","type":"uint256"}],"name":"setSubjectFeePercent","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"sharesBalance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"sharesSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"subjectFeePercent","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	FRIEND_TECH_TRADE = "friend-tech-trade"
)

var client, _ = ethclient.DialContext(context.Background(), RPC_ENDPOINT)

var sess, _ = session.NewSession()

var cacheClient = dynamodb.New(sess, aws.NewConfig().WithRegion("ap-northeast-1"))

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqBody := TradeRecordRequest{}
	err := json.Unmarshal([]byte(request.Body), &reqBody)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}
	if helper.Empty(reqBody.TxHash) {
		return events.APIGatewayProxyResponse{Body: "tx_hash is empty", StatusCode: 400}, nil
	}

	if !helper.IsValidTx(reqBody.TxHash) {
		return events.APIGatewayProxyResponse{Body: "invalid tx_hash value", StatusCode: 400}, nil
	}
	cache, err := cacheClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(FRIEND_TECH_TRADE),
		Key: map[string]*dynamodb.AttributeValue{
			"tx_hash": {S: aws.String(reqBody.TxHash)},
		},
	})
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	var (
		resp = &TradeRecordResponse{
			TxHash: reqBody.TxHash,
		}
		response []byte
	)

	if cache.Item != nil {
		err = dynamodbattribute.UnmarshalMap(cache.Item, &resp)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
		}
		response, err = json.Marshal(&resp)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
		}
		return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 200}, nil
	}

	receipt, err := client.TransactionReceipt(ctx, common.HexToHash(reqBody.TxHash))
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	if receipt.Status == 0 {
		return events.APIGatewayProxyResponse{Body: "the transaction's status is failed", StatusCode: 400}, nil
	}
	trade := make(map[string]interface{}, 0)
	targetABI, _ := abi.JSON(strings.NewReader(ABI_FILE))
	err = targetABI.UnpackIntoMap(trade, "Trade", receipt.Logs[0].Data)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: fmt.Errorf("unpack data error: %v", err).Error(), StatusCode: 500}, nil
	}

	err = helper.GetDataByAbi(trade, resp)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: fmt.Errorf("unpack data error: %v", err).Error(), StatusCode: 500}, nil
	}

	response, err = json.Marshal(&resp)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	resp.TTL = time.Now().Add(helper.GeneralDuration(12, 2, 12, time.Hour)).Unix()
	item, err := dynamodbattribute.MarshalMap(*resp)
	if err == nil {
		_, _ = cacheClient.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(FRIEND_TECH_TRADE),
			Item:      item,
		})
	}

	return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}

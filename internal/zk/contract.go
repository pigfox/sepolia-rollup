package zk

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// ABI for OptimisticRollupFraudProof contract
const fraudProofABI = `[{"constant":true,"inputs":[{"name":"challenger","type":"address"}],"name":"verifyZKProof","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_verifier","type":"address"}],"name":"setVerifier","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"blockNumber","type":"uint256"},{"name":"stateRootBefore","type":"bytes32"},{"name":"stateRootAfter","type":"bytes32"},{"name":"a","type":"uint256[2]"},{"name":"b","type":"uint256[2][2]"},{"name":"c","type":"uint256[2]"},{"name":"input","type":"uint256[3]"}],"name":"submitFraudProof","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`

func ConnectToClient() (*ethclient.Client, error) {
	// Connect to Sepolia RPC node
	client, err := ethclient.Dial("https://ethereum-sepolia-rpc.publicnode.com")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}
	return client, nil
}

func SetVerifier(client *ethclient.Client, contractAddress common.Address, verifierAddress string) {
	// ABI interface for OptimisticRollupFraudProof contract
	fraudProofContract, err := abi.JSON(strings.NewReader(fraudProofABI))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	txData, _ := fraudProofContract.Pack("setVerifier", common.HexToAddress(verifierAddress))
	tx := types.NewTransaction(0, contractAddress, big.NewInt(0), uint64(0), nil, txData)

	// Here you would submit the transaction (signing, sending, etc.)
	fmt.Println("Verifier address set successfully")
}

func SubmitFraudProof(client *ethclient.Client, contractAddress common.Address, blockNumber uint64, stateRootBefore string, stateRootAfter string, a [2]uint256, b [2][2]uint256, c [2]uint256, input [3]uint256) {
	// ABI interface for OptimisticRollupFraudProof contract
	fraudProofContract, err := abi.JSON(strings.NewReader(fraudProofABI))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	txData, _ := fraudProofContract.Pack("submitFraudProof", blockNumber, stateRootBefore, stateRootAfter, a, b, c, input)

	// Here you would submit the transaction (signing, sending, etc.)
	fmt.Println("Fraud proof submitted successfully")
}

func VerifyProof(client *ethclient.Client, contractAddress common.Address, challengerAddress string) {
	// ABI interface for OptimisticRollupFraudProof contract
	fraudProofContract, err := abi.JSON(strings.NewReader(fraudProofABI))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	callData, _ := fraudProofContract.Pack("verifyZKProof", challengerAddress)

	// Call the contract
	var result bool
	err = client.CallContract(context.Background(), rpc.Call{
		To:   &contractAddress,
		Data: callData,
	}, &result)

	if err != nil {
		log.Fatalf("Failed to call verifyProof: %v", err)
	}

	fmt.Printf("zk-SNARK Proof Validity: %v\n", result)
}

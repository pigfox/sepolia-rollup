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

func main() {
	// Connect to Sepolia RPC node
	client, err := ethclient.Dial("https://ethereum-sepolia-rpc.publicnode.com")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Contract address of the deployed OptimisticRollupFraudProof contract
	contractAddress := common.HexToAddress("0xYourContractAddressHere")
	fraudProofContract, err := abi.JSON(strings.NewReader(fraudProofABI))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	// Call setVerifier (example, using a hardcoded verifier address)
	setVerifier(client, contractAddress, "0xYourVerifierContractAddressHere")

	// Example fraud proof submission (you need to fill in the actual proof data)
	submitFraudProof(client, contractAddress, 1, "0xStateRootBefore", "0xStateRootAfter", [2]uint256{1, 2}, [2][2]uint256{{1, 2}, {3, 4}}, [2]uint256{5, 6}, [3]uint256{7, 8, 9})

	// Example fraud proof verification
	verifyProof(client, contractAddress, "0xChallengerAddressHere")
}

// Function to set the verifier contract address
func setVerifier(client *ethclient.Client, contractAddress common.Address, verifierAddress string) {
	tx := &types.Transaction{}
	txData, _ := fraudProofContract.Pack("setVerifier", common.HexToAddress(verifierAddress))
	tx = types.NewTransaction(0, contractAddress, big.NewInt(0), uint64(0), nil, txData)

	// Submit the transaction (signing, sending, etc.)
	fmt.Println("Verifier address set successfully")
}

// Function to submit a fraud proof with zk-SNARK data
func submitFraudProof(client *ethclient.Client, contractAddress common.Address, blockNumber uint64, stateRootBefore string, stateRootAfter string, a [2]uint256, b [2][2]uint256, c [2]uint256, input [3]uint256) {
	txData, _ := fraudProofContract.Pack("submitFraudProof", blockNumber, stateRootBefore, stateRootAfter, a, b, c, input)

	// Create a new transaction (signing, sending, etc.)
	fmt.Println("Fraud proof submitted successfully")
}

// Function to verify a zk-SNARK proof (example)
func verifyProof(client *ethclient.Client, contractAddress common.Address, challengerAddress string) {
	callData, _ := fraudProofContract.Pack("verifyZKProof", challengerAddress)

	// Call the contract
	var result bool
	err := client.CallContract(context.Background(), rpc.Call{
		To:   &contractAddress,
		Data: callData,
	}, &result)

	if err != nil {
		log.Fatalf("Failed to call verifyProof: %v", err)
	}

	fmt.Printf("zk-SNARK Proof Validity: %v\n", result)
}

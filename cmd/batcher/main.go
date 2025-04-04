package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	loadEnv()
	rpcURL := os.Getenv("RPC_URL")
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	contractAddr := os.Getenv("CONTRACT_ADDRESS")

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		log.Fatalf("Invalid private key: %v", err)
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	nonce, _ := client.PendingNonceAt(context.Background(), fromAddress)
	gasPrice, _ := client.SuggestGasPrice(context.Background())
	chainID, _ := client.NetworkID(context.Background())

	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	contractABI := `[{"inputs":[{"internalType":"bytes32","name":"_stateRoot","type":"bytes32"}],"name":"submitBatch","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	parsedABI, _ := abi.JSON(strings.NewReader(contractABI))

	contractAddress := common.HexToAddress(contractAddr)
	stateRoot := [32]byte{}
	copy(stateRoot[:], []byte("dummy_state_root"))

	input, _ := parsedABI.Pack("submitBatch", stateRoot)
	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), auth.GasLimit, gasPrice, input)

	signedTx, _ := auth.Signer(fromAddress, tx)
	client.SendTransaction(context.Background(), signedTx)

	fmt.Printf("📦 Submitted batch tx: %s\n", signedTx.Hash().Hex())
}

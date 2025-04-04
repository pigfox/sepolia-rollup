package eth

import "os"

type EthClient struct {
	RPCUrl     string
	PrivateKey string
}

func NewEthClientFromEnv() *EthClient {
	return &EthClient{
		RPCUrl:     os.Getenv("RPC_URL"),
		PrivateKey: os.Getenv("PRIVATE_KEY"),
	}
}

func (c *EthClient) SubmitStateRoot(stateRoot string) error {
	// TODO: Implement transaction logic with go-ethereum
	return nil
}

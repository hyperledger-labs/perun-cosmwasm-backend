package node

import (
	"os"

	"github.com/CosmWasm/wasmd/app"
	wtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	client "github.com/perun-network/perun-cosmwasm-backend/pkg/cosmwasm"
	"github.com/tendermint/tendermint/rpc/client/http"
)

// Client provides methods for interacting with a CosmWasm node.
type Client struct {
	wtypes.MsgClient
	wtypes.QueryClient
	tmservice.ServiceClient
}

var _ client.Client = &Client{}

// NewClient creates a new client.
func NewClient(nodeURL string, chainID string, acc types.AccAddress, kr keyring.Keyring) (*Client, error) {
	tendermintClient, err := http.New(nodeURL, "/websocket")
	if err != nil {
		return nil, err
	}

	encodingConfig := app.MakeEncodingConfig()

	clientCtx := sdkclient.Context{
		FromAddress:       acc,
		Client:            tendermintClient,
		ChainID:           chainID,
		JSONMarshaler:     encodingConfig.Marshaler,
		InterfaceRegistry: encodingConfig.InterfaceRegistry,
		Input:             os.Stdin,
		Keyring:           kr,
		Output:            os.Stdout,
		OutputFormat:      "json",
		Height:            0,
		HomeDir:           app.DefaultNodeHome,
		KeyringDir:        "",
		From:              acc.String(),
		BroadcastMode:     "block",
		FromName:          acc.String(),
		SignModeStr:       "",
		UseLedger:         false,
		Simulate:          false,
		GenerateOnly:      false,
		Offline:           false,
		SkipConfirm:       true,
		TxConfig:          encodingConfig.TxConfig,
		AccountRetriever:  authtypes.AccountRetriever{},
		NodeURI:           nodeURL,
	}

	return &Client{
		MsgClient:     wtypes.NewMsgClient(clientCtx),
		QueryClient:   wtypes.NewQueryClient(clientCtx),
		ServiceClient: tmservice.NewServiceClient(clientCtx),
	}, nil
}

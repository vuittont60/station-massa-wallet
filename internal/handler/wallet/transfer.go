package wallet

import (
	"fmt"
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/node/sendoperation/transaction"
)

type PromptRequestTransferData struct {
	NicknameFrom     string
	Amount           string
	Fee              string
	RecipientAddress string
}

func NewTransferCoin(prompterApp prompt.WalletPrompterInterface, massaClient network.NodeFetcherInterface) operations.TransferCoinHandler {
	return &transferCoin{prompterApp: prompterApp, massaClient: massaClient}
}

type transferCoin struct {
	prompterApp prompt.WalletPrompterInterface
	massaClient network.NodeFetcherInterface
}

func (t *transferCoin) Handle(params operations.TransferCoinParams) middleware.Responder {
	// params.Nickname length is already checked by go swagger
	wlt, resp := loadWallet(params.Nickname)
	if resp != nil {
		return resp
	}

	// convert amount to uint64
	amount, err := strconv.ParseUint(string(params.Body.Amount), 10, 64)
	if err != nil {
		return operations.NewTransferCoinBadRequest().WithPayload(
			&models.Error{
				Code:    errorTransferCoin,
				Message: "Error during amount conversion",
			})
	}

	// convert fee to uint64
	fee, err := strconv.ParseUint(string(params.Body.Fee), 10, 64)
	if err != nil {
		return operations.NewTransferCoinBadRequest().WithPayload(
			&models.Error{
				Code:    errorTransferCoin,
				Message: "Error during fee conversion",
			})
	}

	promptRequest := prompt.PromptRequest{
		Action:      walletapp.Transfer,
		CodeMessage: utils.MsgTransferRequest,
		Data: PromptRequestTransferData{
			NicknameFrom:     wlt.Nickname,
			Amount:           string(params.Body.Amount),
			Fee:              string(params.Body.Fee),
			RecipientAddress: *params.Body.RecipientAddress,
		},
	}

	_, err = prompt.WakeUpPrompt(t.prompterApp, promptRequest, wlt)
	if err != nil {
		return operations.NewTransferCoinUnauthorized().WithPayload(
			&models.Error{
				Code:    errorCanceledAction,
				Message: "Unable to unprotect wallet",
			})
	}

	// create the transaction and send it to the network
	operation, transferError := doTransfer(wlt, amount, fee, *params.Body.RecipientAddress, t.massaClient)
	if transferError != nil {
		errStr := fmt.Sprintf("error transferring coin: %v", transferError.Err.Error())
		t.prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: transferError.CodeErr})
		return operations.NewTransferCoinInternalServerError().WithPayload(
			&models.Error{
				Code:    errorTransferCoin,
				Message: errStr,
			})
	}

	t.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true, CodeMessage: utils.MsgTransferSuccess})
	return operations.NewTransferCoinOK().WithPayload(
		&models.OperationResponse{
			OperationID: operation.OperationID,
		})
}

func doTransfer(wlt *wallet.Wallet, amount, fee uint64, recipientAddress string, massaClient network.NodeFetcherInterface) (*sendOperation.OperationResponse, *wallet.WalletError) {
	operation, err := transaction.New(recipientAddress, amount)
	if err != nil {
		return nil, &wallet.WalletError{Err: fmt.Errorf("Error during transaction creation: %w", err), CodeErr: errorTransferCoin}
	}

	return network.SendOperation(wlt, massaClient, operation, fee)
}

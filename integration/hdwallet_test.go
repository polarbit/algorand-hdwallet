package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/polarbit/hdwallet"
)

const (
	algodAddress = "https://testnet-algorand.api.purestake.io/ps2"
	psTokenKey   = "X-API-Key"
	psToken      = "B3SU4KcVKi94Jap2VXkK83xx38bsv95K5UZm2lab"
	poolAddr     = "KI7RS3I3IHRR5T3WJ2YH2BILXF652Z2XNS45X6AXL4UVGZ76FKWKZYB4UM"
	poolMn25     = "finish when fun spatial art feed scare bomb fame hold measure hurt hill hope way warrior satisfy country inflict father option flee enlist abandon ice"
	toMn24       = "use regular hybrid trim impulse flash globe mother jealous route label road notable april local face truck obtain neglect sauce surge field gorilla hair"
	toPath       = "m/0'/0'"
)

// In each test run, the pool address sends some Algos to the derived address.
// The derived address then sends back all amount to the pool address.
// Also the pool address is closed with this latest transaction.

func TestAlgorandPayment(t *testing.T) {

	// From
	hotPriv, err := mnemonic.ToPrivateKey(poolMn25)
	if err != nil {
		fmt.Printf("error recovering private key: %s\n", err)
		return
	}

	// To
	toAddr, toXKey, err := hdwallet.DeriveAddressFromMnemonic(toMn24, toPath)

	client, err := NewClient()
	if err != nil {
		t.Error(err)
	}

	// Send payment to derived account
	amount1 := uint64(110000)
	_, err = sendPayment(client, poolAddr, toAddr, amount1, "", hotPriv)
	if err != nil {
		t.Error(err)
	}

	// Send payment back to pool account.
	amount2 := uint64(100000)
	_, err = sendPayment(client, toAddr, poolAddr, amount2, poolAddr, toXKey.CurvePrivateKey)
	if err != nil {
		t.Error(err)
	}
}

func sendPayment(client *algod.Client, from string, to string, amount uint64, closeTo string, priv []byte) (confirmedTx models.PendingTransactionInfoResponse, err error) {
	// Get the suggested transaction parameters
	params, err := client.SuggestedParams().Do(context.Background())
	if err != nil {
		fmt.Printf("error getting suggested tx params: %s\n", err)
		return
	}

	// Create payment transaction
	note := []byte("Hello World")
	genID := params.GenesisID
	genHash := params.GenesisHash
	firstValidRound := uint64(params.FirstRoundValid)
	lastValidRound := uint64(params.LastRoundValid)
	tx, err := transaction.MakePaymentTxn(from, to, uint64(params.Fee), amount, firstValidRound, lastValidRound, note, closeTo, genID, genHash)
	if err != nil {
		fmt.Printf("Error creating transaction: %s\n", err)
		return
	}

	// Sign the transaction
	txID, signedTxn, err := crypto.SignTransaction(priv, tx)
	if err != nil {
		fmt.Printf("Failed to sign transaction: %s\n", err)
		return
	}
	fmt.Printf("Signed txid: %s\n", txID)

	// Submit the transaction
	sendResponse, err := client.SendRawTransaction(signedTxn).Do(context.Background())
	if err != nil {
		fmt.Printf("failed to send transaction: %s\n", err)
		return
	}
	fmt.Printf("Submitted transaction %s\n", sendResponse)

	// Wait for confirmation
	confirmedTx, err = future.WaitForConfirmation(client, txID, 4, context.Background())
	if err != nil {
		fmt.Printf("Error waiting for confirmation on txID: %s\n", txID)
		fmt.Errorf("Error: %s", err)
		return
	}

	return
}

func NewClient() (*algod.Client, error) {
	commonClient, err := common.MakeClient(algodAddress, psTokenKey, psToken)
	if err != nil {
		fmt.Printf("Failed to create common client: %s", err)
		return nil, err
	}
	algodClient := (*algod.Client)(commonClient)

	return algodClient, nil
}

func PrettyPrintTx(tx *models.PendingTransactionInfoResponse) {
	// Display completed transaction
	txnJSON, err := json.MarshalIndent(tx.Transaction.Txn, "", "\t")
	if err != nil {
		fmt.Printf("Can not marshall txn data: %s\n", err)
	}
	fmt.Printf("Transaction information: %s\n", txnJSON)
	fmt.Printf("Decoded note: %s\n", string(tx.Transaction.Txn.Note))
	fmt.Printf("Amount sent: %d microAlgos\n", tx.Transaction.Txn.Amount)
	fmt.Printf("Fee: %d microAlgos\n", tx.Transaction.Txn.Fee)

	fmt.Printf("Decoded note: %s\n", string(tx.Transaction.Txn.Note))
}

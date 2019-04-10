package stellartx

import (
	"fmt"
	"github.com/awnumar/memguard"
	"github.com/stellar/go/build"
	"github.com/stellar/go/xdr"
	"log"
)

//DecodeTx Decodes an XDR encoded transaction string.
func DecodeTx(xdrString string) (*xdr.TransactionEnvelope, error) {
	var tx xdr.TransactionEnvelope
	err := xdr.SafeUnmarshalBase64(xdrString, &tx)
	if err != nil {
		log.Fatal("Invalid XDR encoded transaction", err)
	}

	fmt.Println("")
	fmt.Println("Decoded Transaction Details:")
	fmt.Printf("  Source Account: %s\n", tx.Tx.SourceAccount.Address())
	fmt.Printf("  Number of Operations: %d\n", len(tx.Tx.Operations))
	fmt.Printf("  Number of Signatures: %d\n", len(tx.Signatures))
	fmt.Println("")

	return &tx, nil
}

//SignTx Sign a transaction envelope given a seed. Returns a base64 encoded signed tx.
func SignTx(tx xdr.TransactionEnvelope, seed memguard.LockedBuffer) (*string, error) {

	// Sign the transaction
	txBuilder := &build.TransactionEnvelopeBuilder{E: &tx}
	txBuilder.Init()
	err := txBuilder.MutateTX(build.PublicNetwork)
	if err != nil {
		return nil, err
	}

	// TODO: This basically voids our use of memguard to protect the key in memory. Unfortunately the Stellar SDK uses
	// the golang string type which is immutable and can't be wiped from memory. But for the purposes of demonstrating
	// best practices we still use memguard to lower our threat profile.
	err = txBuilder.Mutate(build.Sign{Seed: string(seed.Buffer())})
	if err != nil {
		return nil, err
	}

	// Convert to base64
	signedTxBase64, err := xdr.MarshalBase64(txBuilder.E)
	if err != nil {
		log.Fatal(err)
	}

	return &signedTxBase64, nil
}

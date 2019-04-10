package main

import (
	"flag"
	"fmt"
	"github.com/ChainFront/stellar-xdr-signer/pkg/stellartx"
	"github.com/awnumar/memguard"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"syscall"
)

func main() {
	// Tell memguard to listen out for interrupts, and cleanup in case of one
	memguard.CatchInterrupt(func() {
		fmt.Println("Interrupt signal received. Exiting...")
	})

	// Make sure to destroy all LockedBuffers when returning
	defer memguard.DestroyAll()

	// Parse input flags
	xdrInput := flag.String("xdr", "", "Unsigned XDR-encoded transaction.")
	flag.Parse()
	if *xdrInput == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Decode the XDR
	tx, err := stellartx.DecodeTx(*xdrInput)
	if err != nil {
		log.Fatalln("Invalid XDR encoded transaction: ", err)
	}

	// Ask user to securely enter private key
	lockedSeed, err := getSeedInput()
	if err != nil {
		log.Fatalln("Unable to read input", err)
	}

	// Sign the transaction
	signedTx, err := stellartx.SignTx(*tx, *lockedSeed)
	if err != nil {
		fmt.Println("Unable to sign transaction: ", err)
		memguard.SafeExit(1)
	}

	fmt.Println()
	fmt.Printf("Signed Transaction (Base64 Encoded):\n%s", *signedTx)
}

func getSeedInput() (*memguard.LockedBuffer, error) {
	fmt.Println()
	fmt.Println("Enter the signing seed or private key: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err

	}
	defer memguard.WipeBytes(bytePassword)

	lockedSeed, err := memguard.NewImmutableFromBytes(bytePassword)
	return lockedSeed, nil
}

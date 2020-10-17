package main

import (
	"flag"
	"fmt"

	"github.com/adityak368/swissknife/crypto"
)

// A CLI to generate a pub/priv rsa key pair

func main() {

	generatersa := flag.Bool("generatersa", false, "Generate RSA Pub/Priv Key Pair")
	flag.Parse()

	if *generatersa {
		// Create the keys
		priv, pub := crypto.GenerateRsaKeyPair()

		if err := crypto.ExportRsaPrivateKeyToFile("privatekey.pem", priv); err != nil {
			panic(err)
		}

		if err := crypto.ExportRsaPublicKeyToFile("pubkey.pub", pub); err != nil {
			panic(err)
		}
		fmt.Println("Generated Pub/Priv Key Pair...")
		return
	}

	fmt.Println("No options given. Doing nothing...")

}

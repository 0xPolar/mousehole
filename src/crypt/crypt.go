package crypt

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type signer struct {
	public  string
	private string
}

type keyPair struct {
	public  ed25519.PublicKey
	private ed25519.PrivateKey
}

func genKeys() (keyPair, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return keyPair{}, err
	}

	return keyPair{pub, priv}, nil
}

func (s *signer) Save() error

func savePublicKey(key ed25519.PublicKey) error {
	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keyBytes,
	}

	file, err := os.Create("public.pem")
	if err != nil {
		return fmt.Errorf("Failed to create public key pem file: %v", err)
	}
	defer file.Close()

	if err := pem.Encode(file, block); err != nil {
		return fmt.Errorf("Failed to save private key to pem file: %v", err)
	}
	return nil
}

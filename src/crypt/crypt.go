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

func savePrivateKey(key ed25519.PrivateKey) error {
	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return fmt.Errorf("failed to marshall private key: %v", err)
	}

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	}

	file, err := os.Create("private.pem")
	if err != nil {
		return fmt.Errorf("Failed to create private key pem file: %v", err)
	}
	defer file.Close()

	if err := pem.Encode(file, block); err != nil {
		return fmt.Errorf("Failed to save privaete key to pem file: %v", err)
	}
	return nil
}

func savePublicKey(key ed25519.PublicKey) error {
	keyBytes, err := x509.MarshalPKCS8PublicKey(key)
	if err != nil {
		return fmt.Errorf("failed to marshall public key: %v", err)
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
		return fmt.Errorf("Failed to save public key to pem file: %v", err)
	}
	return nil
}

package crypt

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type DefaultSigner struct {
	keys keyPair
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

func (s *DefaultSigner) Save() error {
	if err := savePrivateKey(s.keys.private); err != nil {
		return err
	}

	if err := savePublicKey(s.keys.public); err != nil {
		return err
	}

	return nil
}

func savePrivateKey(key ed25519.PrivateKey) error {
	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return fmt.Errorf("failed to marshall private key: %v", err)
	}

	block := &pem.Block{
		Type:  "PRIVATE KEY",
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
	keyBytes, err := x509.MarshalPKIXPublicKey(key)
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

func (s *DefaultSigner) Load() error

func loadPrivateKey(filePath string) (ed25519.PrivateKey, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %v", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM blcok: %v", err)
	}

	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := parsed.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("expected Ed25519 private key, got %T", parsed)
	}

	return key, nil
}

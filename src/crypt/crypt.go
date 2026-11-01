package crypt

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

var signer = &DefaultSigner{}

type DefaultSigner struct {
	keys keyPair
}

func NewSigner() *DefaultSigner {
	return signer
}

func (s *DefaultSigner) New() error {
	keys, err := genKeys()
	if err != nil {
		return err
	}

	s.keys = keys

	return nil
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

func (s *DefaultSigner) Save(filepath string) error {
	if err := savePrivateKey(s.keys.private, filepath); err != nil {
		return err
	}

	if err := savePublicKey(s.keys.public, filepath); err != nil {
		return err
	}

	return nil
}

func savePrivateKey(key ed25519.PrivateKey, filepath string) error {
	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return fmt.Errorf("failed to marshall private key: %v", err)
	}

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	}

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("Failed to create private key pem file: %v", err)
	}
	defer file.Close()

	if err := pem.Encode(file, block); err != nil {
		return fmt.Errorf("Failed to save privaete key to pem file: %v", err)
	}
	return nil
}

func savePublicKey(key ed25519.PublicKey, filepath string) error {
	keyBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return fmt.Errorf("failed to marshall public key: %v", err)
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keyBytes,
	}

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("Failed to create public key pem file: %v", err)
	}
	defer file.Close()

	if err := pem.Encode(file, block); err != nil {
		return fmt.Errorf("Failed to save public key to pem file: %v", err)
	}
	return nil
}

func (s *DefaultSigner) Load(filepath string) error {
	var key ed25519.PrivateKey
	var err error
	if key, err = loadPrivateKey(filepath); err != nil {
		return err
	}

	s.keys = keyPair{
		private: key,
		public:  key.Public().(ed25519.PublicKey),
	}
	return nil
}

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

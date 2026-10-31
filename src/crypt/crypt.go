package crypt

import "crypto/ed25519"

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

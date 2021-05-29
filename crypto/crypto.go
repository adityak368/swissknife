package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"hash"
	"io/ioutil"

	"golang.org/x/crypto/chacha20poly1305"
)

// GenerateRsaKeyPair generates a pub/priv rsa key pair
func GenerateRsaKeyPair(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

// ExportRsaPrivateKeyAsPemStr exports the rsa priv key in pem format
func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {

	privkeyBytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkeyBytes,
		},
	)
	return string(privkeyPEM)
}

// ExportRsaPublicKeyAsPemStr exports the rsa pub key in pem format
func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) string {

	pubkeyBytes := x509.MarshalPKCS1PublicKey(pubkey)

	pubkeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkeyBytes,
		},
	)

	return string(pubkeyPEM)
}

// ExportRsaPrivateKeyToFile exports the rsa priv key to a file
func ExportRsaPrivateKeyToFile(fileName string, key *rsa.PrivateKey) error {

	// Export the keys to pem string
	priv_pem := ExportRsaPrivateKeyAsPemStr(key)

	// Write to file
	return ioutil.WriteFile(fileName, []byte(priv_pem), 0644)
}

// ExportRsaPublicKeyToFile exports the rsa pub key to a file
func ExportRsaPublicKeyToFile(fileName string, key *rsa.PublicKey) error {

	// Export the keys to pem string
	pubKey := ExportRsaPublicKeyAsPemStr(key)

	// Write to file
	return ioutil.WriteFile(fileName, []byte(pubKey), 0644)
}

// ParseRsaPrivateKeyFromPemStr parses the private key from the given pem string
func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

// ParseRsaPublicKeyFromPemStr parses the pub key from the given pem string
func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("Key type is not RSA")
}

// ParseRsaPrivateKeyFromFile parses the priv key from the given file
func ParseRsaPrivateKeyFromFile(fileName string) (*rsa.PrivateKey, error) {

	// read the whole file at once
	priv_pem, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	// Import the keys from pem string
	return ParseRsaPrivateKeyFromPemStr(string(priv_pem))
}

// GetHash returns the hash of a message
func GetHash(msg []byte) []byte {
	h := sha512.New()
	h.Write(msg)
	return h.Sum(nil)
}

// Sign returns a signature made by combining the message and the signers private key
func Sign(msg []byte, key *rsa.PrivateKey) (signature []byte, err error) {
	hs := GetHash(msg)
	return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA512, hs)
}

// Verify checks if a message is signed by a given Public Key
func Verify(msg []byte, sig []byte, pk *rsa.PublicKey) error {
	hs := GetHash(msg)
	return rsa.VerifyPKCS1v15(pk, crypto.SHA512, hs, sig)
}

// EncryptUsingSymmKey Enrypts data using ChaCha2020
// TODO: Give option to change encryption algo
func EncryptUsingSymmKey(msg, privKey []byte) ([]byte, error) {
	hs := GetHash(privKey)
	key := make([]byte, chacha20poly1305.KeySize)
	if _, err := rand.Read(hs); err != nil {
		return nil, err
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}
	var encryptedMsg []byte
	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(msg)+aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	// Encrypt the message and append the ciphertext to the nonce.
	encryptedMsg = aead.Seal(nonce, nonce, msg, nil)

	return encryptedMsg, nil
}

// DecryptUsingSymmKey Decrypts data using ChaCha2020
// TODO: Give option to change encryption algo
func DecryptUsingSymmKey(encryptedMsg, privKey []byte) ([]byte, error) {
	hs := GetHash(privKey)
	key := make([]byte, chacha20poly1305.KeySize)
	if _, err := rand.Read(hs); err != nil {
		return nil, err
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}
	if len(encryptedMsg) < aead.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	// Split nonce and ciphertext.
	nonce, ciphertext := encryptedMsg[:aead.NonceSize()], encryptedMsg[aead.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// HMAC returns the hmac of the message and key
func HMAC(message, key []byte, hashFunc func() hash.Hash) []byte {
	mac := hmac.New(hashFunc, key)
	mac.Write(message)
	return mac.Sum(nil)
}

// ValidMAC reports whether messageMAC is a valid HMAC tag for message.
func ValidMAC(message, key, messageMAC []byte, hashFunc func() hash.Hash) bool {
	expectedMAC := HMAC(message, key, hashFunc)
	return hmac.Equal(messageMAC, expectedMAC)
}

// SignEcdsa returns a signature made by combining the message and the signers private key using ecdsa
func SignEcdsa(msg []byte, key *ecdsa.PrivateKey) ([]byte, error) {
	hs := GetHash(msg)
	return ecdsa.SignASN1(rand.Reader, key, hs)
}

// VerifyEcdsa checks if a message is signed by a given public key using ecdsa
func VerifyEcdsa(msg []byte, sig []byte, pk *ecdsa.PublicKey) bool {
	hs := GetHash(msg)
	return ecdsa.VerifyASN1(pk, hs, sig)
}

// GenerateEcdsaKeyPair generates a pub/priv ecdsa key pair
func GenerateEcdsaKeyPair() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
}

// EncodeEcdsaPrivateKeyToPem encodes ecdsa private key to pem format
func EncodeEcdsaPrivateKeyToPem(privateKey *ecdsa.PrivateKey) ([]byte, []byte, error) {
	x509EncodedPriv, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}
	pemEncodedPriv := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: x509EncodedPriv})

	x509EncodedPub, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return pemEncodedPriv, pemEncodedPub, nil
}

// DecodeEcdsaPrivateKeyFromPem decodes pem private key
func DecodeEcdsaPrivateKeyFromPem(pemEncodedPrivKey []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(pemEncodedPrivKey)
	x509Encoded := block.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

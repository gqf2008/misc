package misc

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

//生成RSA密钥对
func GenRSAKey(bits int) (priv_key, pub_key []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	priv_key = pem.EncodeToMemory(block)

	derPkix, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pub_key = pem.EncodeToMemory(block)
	return priv_key, pub_key, nil
}

//RSA公钥加密
func RSAEncrypt(origData, pub_key []byte) (encrypt []byte, err error) {
	block, _ := pem.Decode(pub_key)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// RSA私钥解密
func RSADecrypt(encrypt, priv_key []byte) (data []byte, err error) {
	block, _ := pem.Decode(priv_key)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, encrypt)
}

func HMAC_SHA256(key []byte, b []byte) []byte {
	sha_256 := sha256.New
	hash := hmac.New(sha_256, key)
	hash.Write(b)
	return hash.Sum(nil)
}

func SHA512(b []byte) []byte {
	hash := sha512.New()
	hash.Write(b)
	return hash.Sum(nil)
}
func MD5(b []byte) []byte {
	hash := md5.New()
	hash.Write(b)
	return hash.Sum(nil)
}
func SHA1(b []byte) []byte {
	hash := sha1.New()
	hash.Write(b)
	return hash.Sum(nil)
}
func SHA256(b []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(b))
	return hash.Sum(nil)
}

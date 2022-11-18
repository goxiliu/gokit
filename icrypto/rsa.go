package icrypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

type Secret uint

const (
	PKCS1 Secret = 1 + iota
	PKCS8
)

func packageData(originalData []byte, packageSize int) (r [][]byte) {
	var src = make([]byte, len(originalData))
	copy(src, originalData)

	r = make([][]byte, 0)
	if len(src) <= packageSize {
		return append(r, src)
	}
	for len(src) > 0 {
		var p = src[:packageSize]
		r = append(r, p)
		src = src[packageSize:]
		if len(src) <= packageSize {
			r = append(r, src)
			break
		}
	}
	return r
}

//RSAEncrypt 公钥加密
func RSAEncrypt(plaintext, path string) (string, error) {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	pub, err := parsePKCS1PublicKey(key)
	if err != nil {
		return "", err
	}

	var data = packageData([]byte(plaintext), pub.N.BitLen()/8-11)
	var cipherData = make([]byte, 0, 0)

	for _, d := range data {
		var c, e = rsa.EncryptPKCS1v15(rand.Reader, pub, d)
		if e != nil {
			return "", e
		}
		cipherData = append(cipherData, c...)
	}

	return base64.StdEncoding.EncodeToString(cipherData), nil
}

// RSADecrypt 私钥解密
// @rsaType 一般默认PKCS1
func RSADecrypt(plaintext, path string, rsaType Secret) ([]byte, error) {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(plaintext)
	if err != nil {
		return nil, err
	}

	pri, err := parsePrivateKey(key, rsaType)
	if err != nil {
		return nil, err
	}

	var data = packageData(ciphertext, pri.PublicKey.N.BitLen()/8)
	var plainData = make([]byte, 0, 0)

	for _, d := range data {
		var p, e = rsa.DecryptPKCS1v15(rand.Reader, pri, d)
		if e != nil {
			return nil, e
		}
		plainData = append(plainData, p...)
	}
	return plainData, nil
}

func parsePrivateKey(data []byte, rsaType Secret) (key *rsa.PrivateKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, errors.New("private key error")
	}

	switch rsaType {
	case PKCS1:
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case PKCS8:
		prikey, errp := x509.ParsePKCS8PrivateKey(block.Bytes)
		if errp != nil {
			return nil, errp
		}
		key = prikey.(*rsa.PrivateKey)
	}
	return
}

func parsePKCS1PublicKey(data []byte) (key *rsa.PublicKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, errors.New("public key error")
	}

	var pubInterface interface{}
	pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key error")
	}

	return key, err
}

//Sign 私钥签名
func Sign(src, privateKey string, hash crypto.Hash, rsaType Secret, keyFromFile bool) (string, error) {
	var key []byte
	var err error
	if keyFromFile {
		key, err = ioutil.ReadFile(privateKey)
		if err != nil {
			return "", err
		}
	} else {
		key = []byte(privateKey)
	}

	pri, err := parsePrivateKey(key, rsaType)
	if err != nil {
		return "", err
	}
	val, err := signPKCS1v15WithKey([]byte(src), pri, hash)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(val), nil
}

func signPKCS1v15WithKey(src []byte, key *rsa.PrivateKey, hash crypto.Hash) ([]byte, error) {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, key, hash, hashed)
}

//Verify 公钥验签
func Verify(src, sigstr, publicKey string, hash crypto.Hash, keyFromFile bool) error {
	sig, err := base64.StdEncoding.DecodeString(sigstr)
	if err != nil {
		return err
	}

	var key []byte
	if keyFromFile {
		key, err = ioutil.ReadFile(publicKey)
		if err != nil {
			return err
		}
	} else {
		key = []byte(publicKey)
	}

	pub, err := parsePKCS1PublicKey(key)
	if err != nil {
		return err
	}
	return verifyPKCS1v15WithKey([]byte(src), sig, pub, hash)
}

func verifyPKCS1v15WithKey(src, sig []byte, key *rsa.PublicKey, hash crypto.Hash) error {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)
	return rsa.VerifyPKCS1v15(key, hash, hashed, sig)
}

func GenRsaKey(bits int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "PRIVATE KEY", //"RSA PRIVATE KEY",
		Bytes: derStream,
	}
	privFile, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	defer privFile.Close()

	err = pem.Encode(privFile, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubFile, err := os.Create("public.pem")
	if err != nil {
		return err
	}
	defer pubFile.Close()
	err = pem.Encode(pubFile, block)
	if err != nil {
		return err
	}
	return nil
}

// func MarshalPKCS8PrivateKey(key *rsa.PrivateKey) []byte {
// 	info := struct {
// 		Version             int
// 		PrivateKeyAlgorithm []asn1.ObjectIdentifier
// 		PrivateKey          []byte
// 	}{}
// 	info.Version = 0
// 	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
// 	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
// 	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)

// 	k, err := asn1.Marshal(info)
// 	if err != nil {
// 		log.Panic(err.Error())
// 	}
// 	return k
// }

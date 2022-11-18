package icrypto

import (
	"encoding/base64"
	"github.com/thinkoner/openssl"
)

func DESDecrypt(desval, desKey string) ([]byte, error) {
	crytedByte, err := base64.StdEncoding.DecodeString(desval)
	if err != nil {
		return nil, err
	}

	data, err := openssl.Des3ECBDecrypt(crytedByte, []byte(desKey), openssl.PKCS7_PADDING)
	return data, err
}

func DESEncrypt(src []byte, desKey string) (string, error) {
	body, err := openssl.Des3ECBEncrypt(src, []byte(desKey), openssl.PKCS7_PADDING)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString([]byte(body)), nil
}

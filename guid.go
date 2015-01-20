package timerWheel

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/base64"
	"io"
)

func GetMd5String (s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
	
}
//return guid str
func GetGuid () (string,error) {

	b := make([]byte,48)
	if _,err := io.ReadFull(rand.Reader,b); err != nil {
		return "", err
	}
	
	return GetMd5String(base64.URLEncoding.EncodeToString(b)),nil
}

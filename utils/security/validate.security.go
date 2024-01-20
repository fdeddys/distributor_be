package security

import (
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"regexp"
	"strings"
)

// ValidateSignature ...
func ValidateSignature(timestamp, key, signature, bodyMsg string) bool {

	regx := regexp.MustCompile("[^a-zA-Z0-9{}:.,]")
	var bodySign = strings.ToLower(regx.ReplaceAllLiteralString(bodyMsg, "")) + "&" + timestamp + "&" + key

	signatureSystem := hashSha512(key, bodySign)

	fmt.Println("\n============= Request Identity ============= ")
	fmt.Println("mkey Or sKey from DB : ", key)
	fmt.Println("Body from Requestor : \n", bodyMsg)
	fmt.Println("Body to create Signature Sistem : ", bodySign)
	fmt.Println("Signature Sistem  : ", signatureSystem)
	fmt.Println("Signature Requestor  : ", signature)
	fmt.Println("\n ")

	if signature == signatureSystem {
		return true
	}

	return false
}

func hashSha512(secret, data string) string {
	hash := hmac.New(sha512.New, []byte(secret))
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

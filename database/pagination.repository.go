package database

import (
	"crypto/aes"
	"crypto/cipher"

	// "crypto/rand"

	"distribution-system-be/constants"
	"encoding/hex"
	"fmt"

	"github.com/mergermarket/go-pkcs7"

	"math/rand"
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
)

// AsyncQueryParam ...
// type AsyncQueryParam struct {
// 	// DB gorm.DB
// 	// total int
// 	param       interface{}
// 	models      interface{}
// 	fieldLookup string
// 	// resChan chan error
// }

// AsyncQueryCount ...
func AsyncQueryCount(db *gorm.DB, total *int, param interface{}, models interface{}, fieldLookup string, resChan chan error) {
	// func AsyncQueryCount(db *gorm.DB, total *int, parameters AsyncQueryParam, resChan chan error) {
	varInterface := reflect.ValueOf(param)
	strQuery := varInterface.Field(0).Interface().(string)

	// var criteriaName = ""
	// if strings.TrimSpace(strQuery) != "" {
	// 	criteriaName = strQuery
	// }
	criteriaName := strQuery
	if criteriaName == "" {
		criteriaName = "%"
	} else {
		criteriaName = "%" + strQuery + "%"
	}

	// err := db.Model(models).Where(fieldLookup+" ~* ?", criteriaName).Count(&*total).Error
	err := db.Model(models).Where("COALESCE("+fieldLookup+", '') ILIKE ?", criteriaName).Count(&*total).Error

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuery ...
func AsyncQuery(db *gorm.DB, offset int, limit int, modelWVal interface{}, param interface{}, fieldLookup string, resChan chan error) {
	modelsDump := reflect.ValueOf(modelWVal).Interface()
	paramDump := reflect.ValueOf(param)
	strQuery := paramDump.Field(0).Interface().(string)
	// var criteriaName = ""
	// if strings.TrimSpace(strQuery) != "" {
	// 	criteriaName = strQuery //+ criteriaBrandName
	// }

	criteriaName := strQuery
	if criteriaName == "" {
		criteriaName = "%"
	} else {
		criteriaName = "%" + strQuery + "%"
	}



	var err error
	// err = db.Set("gorm:auto_preload", true).Order("name ASC").Offset(offset).Limit(limit).Find(modelsDump, fieldLookup+" ~* ?", criteriaName).Error
	err = db.Set("gorm:auto_preload", true).Order("name ASC").Offset(offset).Limit(limit).Find(modelsDump, "COALESCE("+fieldLookup+",'') ILIKE ?", criteriaName).Error

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// StringWithCharset ...
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// GeneratePassword ...
func GeneratePassword(length int) string {
	return StringWithCharset(length, charset)
}

// func EncryptTripleDES(key []byte, plaintext string) string {
// 	c, err := des.NewTripleDESCipher(key)
// 	if err != nil {
// 		// fmt.Errorf("NewTripleDESCipher(%d bytes) = %s", len(key), err)
// 		panic(err)
// 	}

// 	out := make([]byte, len(plaintext))
// 	c.Encrypt(out, []byte(plaintext))

// 	return hex.EncodeToString(out)
// }

// func DecryptTripleDES(key []byte, ct string) string {

// 	ciphertext, _ := hex.DecodeString(ct)
// 	fmt.Printf("ini chipertext %d", ciphertext)
// 	c, err := des.NewTripleDESCipher([]byte(key))
// 	if err != nil {
// 		// fmt.Errorf("NewTripleDESCipher(%d bytes) = %s", len(key), err)
// 		panic(err)
// 	}
// 	plain := make([]byte, len(ciphertext))
// 	c.Decrypt(plain, ciphertext)
// 	s := string(plain[:])
// 	fmt.Printf("3DES Decrypyed Text:  %s\n", s)
// 	return s
// }

func Encrypt(unencrypted string) string {
	key := []byte(constants.DesKey)
	plainText := []byte(unencrypted)
	plainText, err := pkcs7.Pad(plainText, aes.BlockSize)
	if err != nil {
		return ""
	}
	if len(plainText)%aes.BlockSize != 0 {
		// err := fmt.Errorf(`plainText: "%s" has the wrong block size`, plainText)
		return ""
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]

	// if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	// 	return "", err
	// }
	// fmt.Printf("%d\n", iv)
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	// return fmt.Sprintf("%x", cipherText)
	str := hex.EncodeToString(cipherText)
	// data, err := base64.StdEncoding.DecodeString(str)
	// if err != nil {
	// 	fmt.Println("error:", err)
	// 	return ""
	// }
	return str //fmt.Sprintf("%x", data)
	// return hex.EncodeToString(cipherText)
}

//Decrypt decrypts cipher text string into plain text string
func Decrypt(encrypted string) string {
	key := []byte(constants.DesKey)
	cipherText, _ := hex.DecodeString(encrypted)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(cipherText) < aes.BlockSize {
		panic("cipherText too short")
	}
	iv := cipherText[:aes.BlockSize]

	// fmt.Printf("%d", iv)
	cipherText = cipherText[aes.BlockSize:]
	if len(cipherText)%aes.BlockSize != 0 {
		panic("cipherText is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	cipherText, _ = pkcs7.Unpad(cipherText, aes.BlockSize)
	return fmt.Sprintf("%s", cipherText)
}

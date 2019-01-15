package processor

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"marisinau.com/KitaUndangAPI/helper/validation"
	"regexp"
	"time"
)

/*
RemoveDuplicatesElementInt : menghapus element yang sama pada slice
*/
func RemoveDuplicatesElementInt(elements []int) []int {
	encountered := map[int]bool{}
	result := []int{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

/*
RemoveDuplicatesElementByte : menghapus element yang sama pada slice
*/
func RemoveDuplicatesElementByte(elements []byte) []byte {
	encountered := map[byte]bool{}
	result := []byte{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

/*
EncodeBase64 : encode byte slice ke base64 byte slice
*/
func EncodeBase64(data []byte) []byte {
	base64data := []byte{}
	base64.StdEncoding.Encode(base64data, data)
	return base64data
}

/*
DecodeBase64 : decode base64 byte slice ke byte slice
*/
func DecodeBase64(base64data []byte) (data []byte) {
	base64.StdEncoding.Decode(data, base64data)
	return
}

/*
SHA1 : hash SHA1
*/
func SHA1(data string) string {
	hash := sha1.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

/*
SHA256 : hash SHA256
*/
func SHA256(data string) []byte {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hash.Sum(nil)
}

/*
SHA1WithRandomSalt : hash SHA1 dengan salt random, return salt dan hashed
*/
func SHA1WithRandomSalt(text string) (string, string) {
	var salt = fmt.Sprintf("%d", time.Now().UnixNano())
	var saltedText = fmt.Sprintf("text: '%s', salt: %s", text, salt)
	var sha = sha1.New()
	sha.Write([]byte(saltedText))
	var encrypted = sha.Sum(nil)

	return fmt.Sprintf("%x", encrypted), salt
}

/*
SHA256WithRandomSalt : hash SHA256 dengan salt random, return salt dan hashed
*/
func SHA256WithRandomSalt(text string) (string, string) {
	var salt = fmt.Sprintf("%d", time.Now().UnixNano())
	var saltedText = fmt.Sprintf("text: '%s', salt: %s", text, salt)
	var sha = sha256.New()
	sha.Write([]byte(saltedText))
	var encrypted = sha.Sum(nil)

	return fmt.Sprintf("%x", encrypted), salt
}

/*
SHA1WithDefinedSalt : hash SHA1 dengan salt manual, return hashed
*/
func SHA1WithDefinedSalt(text string, salt string) string {
	var saltedText = fmt.Sprintf("text: '%s', salt: %s", text, salt)
	var sha = sha1.New()
	sha.Write([]byte(saltedText))
	var encrypted = sha.Sum(nil)

	return fmt.Sprintf("%x", encrypted)
}

/*
SHA256WithDefinedSalt : hash SHA256 dengan salt manual, return hashed
*/
func SHA256WithDefinedSalt(text string, salt string) string {
	var saltedText = fmt.Sprintf("text: '%s', salt: %s", text, salt)
	var sha = sha256.New()
	sha.Write([]byte(saltedText))
	var encrypted = sha.Sum(nil)

	return fmt.Sprintf("%x", encrypted)
}

/*
SHA1WithDefinedSaltEncoded : hash SHA1 dengan salt manual encoded, return hashed
*/
func SHA1WithDefinedSaltEncoded(text string, salt string) string {
	var saltedText = fmt.Sprintf("text: '%s', salt: %s", text, salt)
	var sha = sha1.New()
	sha.Write([]byte(saltedText))
	var encrypted = sha.Sum(nil)

	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%x", encrypted)))
}

/*
SHA256WithDefinedSaltEncoded : hash SHA256 dengan salt manual encoded, return hashed
*/
func SHA256WithDefinedSaltEncoded(text string, salt string) string {
	var saltedText = fmt.Sprintf("text: '%s', salt: %s", text, salt)
	var sha = sha256.New()
	sha.Write([]byte(saltedText))
	var encrypted = sha.Sum(nil)

	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%x", encrypted)))
}

/*
NormalizePhoneNumber : normalisasi phone number
*/
func NormalizePhoneNumber(phone string, countryCode int) (string, error) {
	var regex, err = regexp.Compile("[^0-9]*")
	if err != nil {
		fmt.Println(err.Error())
	}

	if regex.MatchString(phone) {
		phone = regex.ReplaceAllString(phone, "")
	}

	if &phone != nil || !validation.EqualsIgnoredCase(phone, "") {
		if len(phone) > 0 {
			if phone[:1] == "0" {
				//replace
				phone = fmt.Sprintf("%d%s", countryCode, phone[1:]) //atau phone[1:len(phone)]
			}

			return phone, nil
		}
	}

	return phone, errors.New("Failed to normalize phone number")
}

/*
ExtractServerAddressPort : extract port
*/
func ExtractServerAddressPort(address string) string {
	var strPort string
	var regex, err = regexp.Compile("[\\:]([0-9]{1,})")
	if err != nil {
		fmt.Println(err.Error())
	}

	if regex.MatchString(address) {
		var getParsing = regex.FindAllStringSubmatch(address, -1) //-1 jika semua karakter
		for _, group := range getParsing {
			strPort = group[1]
		}
	}

	return strPort
}

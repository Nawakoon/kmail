package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/asn1"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"passwordless-mail-client/pkg/account"
	"passwordless-mail-client/pkg/model"
	"passwordless-mail-client/pkg/request"
	"strings"
)

func main() {
	errChan := make(chan error)
	defer close(errChan)
	// go ListenErrChan(errChan)

	credentialFlag := flag.String("user", "", "user private key")
	inboxFlag := flag.String("inbox", "", "get inbox")
	sendMailFlag := flag.String("send", "", "send mail")
	flag.Parse()

	// TestPrivateKey1 := "1baa694c49154f63b1503c7138f184c80f221670f035403ff428a65183bab147"
	// TestPrivateKey2 := "1baa694aa9154f63b1503c7138f187780f221670f035403ff428a65182bab146"
	// testAccount1, _ := account.ConnectAccount(TestPrivateKey1)
	// testAccount2, _ := account.ConnectAccount(TestPrivateKey2)
	// testAccount1Address := testAccount1.GetAddress()
	// testAccount2Address := testAccount2.GetAddress()
	// fmt.Println("Test Account 1 Address:", testAccount1Address)
	// fmt.Println("Test Account 2 Address:", testAccount2Address)
	// "cd91d7cab64774ed58db47351f885cb600df09cc44354b16560308189a7c5013f862679f86d71f264d0b18cd608e87c44ee6d18059fa7ab9c4da9a5e56b9cb57"

	if *inboxFlag != "" {
		if *credentialFlag == "" {
			fmt.Println("user credential is required")
			os.Exit(1)
			return
		}

		err := GetInboxCmd(*inboxFlag, *credentialFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}

		return
	}

	if *sendMailFlag != "" {
		if *credentialFlag == "" {
			fmt.Println("user credential is required")
			os.Exit(1)
			return
		}

		err := SendMailCmd(*sendMailFlag, *credentialFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}

		return
	}
}

func AddCmd(
	flagName string,
	description string,
	command func(flagValue string) error,
	errChan chan error,
) {
	readFlag := flag.String(flagName, "", description)
	flag.Parse()

	if *readFlag != "" {
		err := command(*readFlag)
		errChan <- err
	}
}

func GetInboxCmd(queryPath string, user string) error {
	// validate user credential should be 64 characters and hex
	if len(user) != 64 {
		return fmt.Errorf("invalid user credential: credential should be hex with 64 characters long")
	}
	for _, c := range user {
		if c < '0' || c > 'f' {
			return fmt.Errorf("invalid user credential: credential should be hex with 64 characters long")
		}
	}

	// validate query path
	if _, err := os.Stat(queryPath); os.IsNotExist(err) {
		return fmt.Errorf("query file not found, invalid path or file name")
	}

	// validate query json
	queryFile, err := os.Open(queryPath)
	if err != nil {
		return err
	}
	defer queryFile.Close()
	// validate that queryFile contains a valid json
	var query model.QueryJson
	err = json.NewDecoder(queryFile).Decode(&query)
	if err != nil || query.Page == nil || query.Limit == nil {
		return fmt.Errorf("invalid query json: please use this format\n\t{ \"page\":int, \"limit\":int }")
	}
	if *query.Page <= 0 || *query.Limit <= 0 {
		return fmt.Errorf("invalid query json: page and limit should be greater than 0")
	}

	acc, err := account.ConnectAccount(user)
	if err != nil {
		return err
	}

	message, err := request.NewGetInbox()
	if err != nil {
		return err
	}
	signedMessage, err := acc.Sign(message)
	if err != nil {
		return err
	}
	requestBody := model.RequestBody{
		Data:      string(message),
		Signature: signedMessage,
	}
	requestBodyByte, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	BaseInboxPath := "http://localhost:8080/mail/inbox"
	queryParams := fmt.Sprintf("?page=%d&limit=%d", *query.Page, *query.Limit)
	apiPath := BaseInboxPath + queryParams
	payLoad := strings.NewReader(string(requestBodyByte))
	request, err := http.NewRequest(http.MethodGet, apiPath, payLoad)
	if err != nil {
		return err
	}
	request.Header.Add("x-public-key", acc.GetAddress())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	// print response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}

func hexToBytes(hexStr string) ([]byte, error) {
	// Convert the hex string to a big integer
	bigInt, success := new(big.Int).SetString(hexStr, 16)
	if !success {
		return nil, fmt.Errorf("failed to convert hex string to big integer")
	}

	// Convert the big integer to bytes
	bytes := bigInt.Bytes()

	return bytes, nil
}

type Account struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

func ConnectAccount(privateKeyHex string) (*Account, error) {
	// Convert the private key from hex to bytes
	privateKeyBytes, err := hexToBytes(privateKeyHex)
	if err != nil {
		return nil, err
	}

	// Generate the ECDSA private key from the bytes
	privateKey := new(ecdsa.PrivateKey)
	privateKey.Curve = elliptic.P256()
	privateKey.D = new(big.Int).SetBytes(privateKeyBytes)

	// Compute the public key points
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.Curve.ScalarBaseMult(privateKeyBytes)

	return &Account{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}

func (a *Account) Sign(data []byte) ([]byte, error) {
	// Sign the data with the private key
	r, s, err := ecdsa.Sign(rand.Reader, a.PrivateKey, data)
	if err != nil {
		return nil, err
	}

	// Encode the signature
	signature, err := asn1.Marshal(struct{ R, S *big.Int }{r, s})
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func Verify(publicKey *ecdsa.PublicKey, data, signature []byte) bool {
	// Decode the signature
	var decoded struct{ R, S *big.Int }
	_, err := asn1.Unmarshal(signature, &decoded)
	if err != nil {
		return false
	}

	// Verify the signature
	return ecdsa.Verify(publicKey, data, decoded.R, decoded.S)
}

func SendMailCmd(mailPath string, user string) error {
	// validate mail path
	if _, err := os.Stat(mailPath); os.IsNotExist(err) {
		return fmt.Errorf("mail file not found, invalid path or file name")
	}

	// validate mail json
	mailFile, err := os.Open(mailPath)
	if err != nil {
		return err
	}
	defer mailFile.Close()
	var mail model.MailFileContent
	err = json.NewDecoder(mailFile).Decode(&mail)
	if err != nil ||
		mail.To == nil ||
		mail.Subject == nil ||
		mail.Body == nil {
		return fmt.Errorf("invalid kmail json: please use this format\n\t{ \"to\":string, \"subject\":string, \"body\":string }")
	}

	// validate user credential should be 64 characters and hex
	if len(user) != 64 {
		return fmt.Errorf("invalid user credential: credential should be hex with 64 characters long")
	}
	for _, c := range user {
		if c < '0' || c > 'f' {
			return fmt.Errorf("invalid user credential: credential should be hex with 64 characters long")
		}
	}

	acc, err := account.ConnectAccount(user)
	if err != nil {
		return err
	}

	message, err := request.NewSendEmail(
		*mail.To,
		*mail.Subject,
		*mail.Body,
	)
	if err != nil {
		return err
	}
	signedMessage, err := acc.Sign(message)
	if err != nil {
		return err
	}
	requestBody := model.RequestBody{
		Data:      string(message),
		Signature: signedMessage,
	}
	requestBodyByte, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	BaseSendMailPath := "http://localhost:8080/mail/send"
	payLoad := strings.NewReader(string(requestBodyByte))
	request, err := http.NewRequest(http.MethodPost, BaseSendMailPath, payLoad)
	if err != nil {
		return err
	}
	request.Header.Add("x-public-key", acc.GetAddress())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	// print response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}

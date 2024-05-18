# passwordless mail server

## run the server
```bash
cd server && go run cmd/main.gos
# or make run
```

## new mail
```bash
cd client && go run cmd/main.go -compose
# input1: receiver public key
# input2: subject
# input3: message
```

## send message
```bash
cd client && go run cmd/main.go -send mail.txt
# input1: private key to sign the mail
```
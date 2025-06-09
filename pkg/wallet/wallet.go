package wallet

import (
	// "bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"crypto/sha256"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

// 创建新钱包
func NewWallet() *Wallet {
	privKey, pubKey := generateKeyPair()
	return &Wallet{privKey, pubKey}
}

// 生成密钥对
func generateKeyPair() (*ecdsa.PrivateKey, []byte) {
	// 随机生成私钥
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	// 将公钥的x和y坐标组装成公钥
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return privateKey, publicKey
}

// 对公钥进行哈希
func HashPubKey(pubKey []byte) []byte {
	// 使用双哈希Hash160
	pubHash := sha256.Sum256(pubKey)
	RIPEMD := ripemd160.New()
	_, err := RIPEMD.Write(pubHash[:])
	if err != nil {
		panic(err)
	}
	return RIPEMD.Sum(nil)
}

// 生成钱包地址
const version = byte(0x00)
const addressChecksumLen = 4

func (w *Wallet) GetAddress() string {
	pubKeyHash := HashPubKey(w.PublicKey)
	payload := append([]byte{version}, pubKeyHash...)

	// 计算两次SHA256，并取前4字节作为校验和
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	checksum := second[0:addressChecksumLen]

	// 组装并编码
	fullPayload := append(payload, checksum...)
	address := base58.Encode(fullPayload)
	return address
}

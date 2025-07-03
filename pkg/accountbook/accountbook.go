package accountbook

import (
	"github.com/marshuni/Blockchain-AccountBook/pkg/blockchain"
	"github.com/marshuni/Blockchain-AccountBook/pkg/core/tx"
	"github.com/marshuni/Blockchain-AccountBook/pkg/core/wallet"
	"github.com/marshuni/Blockchain-AccountBook/pkg/utxo"
)

// 账本结构体，封装区块链、UTXO集等
type AccountBook struct {
	Chain   *blockchain.Blockchain
	UTXOSet *utxo.UTXOSet
}

// 初始化账本（区块链+UTXO集）
func NewAccountBook(dbPath string) *AccountBook {
	chain := blockchain.NewBlockchain(dbPath)
	utxoSet := &utxo.UTXOSet{Blockchain: chain}
	return &AccountBook{
		Chain:   chain,
		UTXOSet: utxoSet,
	}
}

// 创建新钱包
func (ab *AccountBook) NewWallet() *wallet.Wallet {
	return wallet.NewWallet()
}

// 获取钱包地址
func (ab *AccountBook) GetAddress(w *wallet.Wallet) string {
	return w.GetAddress()
}

// 查询余额
func (ab *AccountBook) GetBalance(address string) int {
	pubKeyHash := wallet.GetPubKeyHashFromAddress(address)
	utxos := ab.UTXOSet.FindUTXO(pubKeyHash)
	balance := 0
	for _, out := range utxos {
		balance += out.Value
	}
	return balance
}

// 创建交易（from向to转账amount）
func (ab *AccountBook) CreateTransaction(from, to string, amount int, w *wallet.Wallet) (*tx.Transaction, error) {
	return ab.UTXOSet.CreateTransaction(from, to, amount, w)
}

// 打包并添加区块（自动添加Coinbase奖励给minerAddress）
func (ab *AccountBook) AddBlock(txs []*tx.Transaction, minerAddress string) {
	pool := blockchain.TxPool{}
	for _, t := range txs {
		pool.AddTx(t)
	}
	ab.Chain.AddBlock(&pool, minerAddress)
}

// 查询某地址所有UTXO
func (ab *AccountBook) ListUTXO(address string) []utxo.UTXOOutput {
	pubKeyHash := wallet.GetPubKeyHashFromAddress(address)
	return ab.UTXOSet.FindUTXO(pubKeyHash)
}

// 打印区块链
func (ab *AccountBook) PrintChain() {
	ab.Chain.Print()
}

// 查询交易
func (ab *AccountBook) FindTransaction(txid []byte) *tx.Transaction {
	return ab.Chain.FindTx(txid)
}

// 创建Coinbase交易
func (ab *AccountBook) NewCoinbaseTx(to, data string) *tx.Transaction {
	return tx.NewCoinbaseTX(to, data)
}

// 验证交易签名
func (ab *AccountBook) VerifyTransaction(t *tx.Transaction) bool {
	return t.VerifyTransaction()
}

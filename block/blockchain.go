package block

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	MINING_DIFFICULTY = 3 // ディフィカルティの設定
	MINING_SENDAR = "THE BLOCKCHAIN"
	MINIG_REWARD = 1.0
)
type Block struct{
	timestamp int64
	nonce int
	previousHash [32]byte
	transactions []*Transaction
}

type Transaction struct{
	senderBlockchainAddress string
	recipientBlockchainAddress string
	value float32
}

// M, _ :=json.Marshal(b)とすると、実行時に{}のjsonが返ってくるため、MarshalJson()でオーバーライドする。
func (b *Block) Hash() [32]byte {
	m, _ := b.MarshalJson()
	return sha256.Sum256([]byte(m))

}
// 'json:'の記述でマーシャルする時にどうようにマーシャルするかを定める。
func (b *Block) MarshalJson() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp  int64  			`json:"timestamp"`
		Nonce      int				`json:"nonce"`
		PrevioushHash [32]byte 		`json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp: b.timestamp,
		Nonce :b.nonce,
		PrevioushHash: b.previousHash,
		Transactions: b.transactions,
	})
}
func (b *Block) Print() {
	fmt.Printf("timestamp        %d\n", b.timestamp)
	fmt.Printf("nonce             %d\n", b.nonce)
	fmt.Printf("previousHash     %x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}

}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block) // ポインタ型
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}

type Blockchain struct{
	transactionPool [] *Transaction
	chain []*Block
	blockchainAddress string
}

func NewBlockchain(blockchainAddress string) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}
func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

// transactionpoolをsyncしたいときに利用する。
func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions, 
			NewTransaction(t.senderBlockchainAddress,
						   t.recipientBlockchainAddress,
						   t.value),
		)
	}
	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1

	}
	return nonce
}

func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(MINING_SENDAR, bc.blockchainAddress, MINIG_REWARD)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining,, status=success")
	return true
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32  {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if blockchainAddress == t.recipientBlockchainAddress {
				totalAmount += value
			}
			if blockchainAddress == t.senderBlockchainAddress {
				totalAmount -= value
			}
		}

	}
	return totalAmount
}



func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

func (t *Transaction) Print(){
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_blockchain_address   %s\n", t.senderBlockchainAddress)
	fmt.Printf(" recipent_blockchain_address %s\n", t.recipientBlockchainAddress)
	fmt.Printf(" value                       %.1f\n", t.value)

}

func (t *Transaction) MarshalJson() ([]byte, error) {
	return json.Marshal(struct {
		Sender string `json:"sender_blockchain_address"` 
		Recipient string `json:"recipient_blockchain_address"`
		Value float32 `json:"value"`
	}{
		Sender: t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value: t.value,
	})
}

func init(){
	log.SetPrefix("Blockchain: ")
}



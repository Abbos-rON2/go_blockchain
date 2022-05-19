package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"time"
)

func main() {
	t := time.Now()
	bc := InitBlockchain()
	bc.AddBlock("First block")
	bc.AddBlock("Second block")

	for _, block := range bc.Blocks {
		printBlock(block)
	}

	fmt.Println(time.Since(t))
}

func printBlock(block *Block) {
	fmt.Printf("Index: %v \n", block.Index)
	fmt.Printf("Nonce: %v \n", block.Nonce)
	fmt.Printf("Timestamp: %v \n", block.Timestamp)
	fmt.Printf("Data: %v \n", block.Data)
	fmt.Printf("Hash: %v \n", block.Hash)
	fmt.Printf("PrevHash: %s \n", block.PrevHash)
}




type Block struct {
	Index     int
	Timestamp string
	Data      string
	Hash      string
	PrevHash  string
	Nonce     int
}
type BlockChain struct {
	Blocks []*Block
}

func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash, prevBlock.Index+1)
	bc.Blocks = append(bc.Blocks, newBlock)
}
func InitBlockchain() *BlockChain {
	bc := &BlockChain{}
	genesis := GenesisBlock()
	bc.Blocks = append(bc.Blocks, genesis)
	return bc
}
func NewBlock(data string, prevHash string, index int) *Block {
	block := Block{}
	block.Index = index
	block.Timestamp = time.Now().String()
	block.Data = data
	block.PrevHash = prevHash
	block.Nonce = 0
	pow := NewProof(&block)
	nonce, hash := pow.Run()
	block.Hash = hex.EncodeToString(hash[:])
	block.Nonce = nonce
	return &block
}
func GenesisBlock() *Block {
	return NewBlock("Genesis Block", "", 0)
}

const Difficulty = 18

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}
	return pow
}
func (pow *ProofOfWork) PrepareData(nonce int) []byte {
	data := []byte(pow.Block.PrevHash + pow.Block.Data+ pow.Block.Timestamp + string(nonce))
	return data
}
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	for nonce < math.MaxInt64 {
		data := pow.PrepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()
	return nonce, hash[:]
}
func (pow *ProofOfWork) ValidateProof() bool {
	var hashInt big.Int
	hashInt.SetBytes([]byte(pow.Block.Hash))
	return hashInt.Cmp(pow.Target) == -1
}

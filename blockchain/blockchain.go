package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevhash",omitempty`
	Height 	 int	`json:"height"`
}

type blockchain struct {
	blocks []*Block
}

var b *blockchain
var once sync.Once

func (b *Block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}
func getLastHash() string {
	totalBlocks := len(GetBlockchain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockchain().blocks[totalBlocks-1].Hash
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash(), len(GetBlockchain().blocks) +1}
	newBlock.calculateHash()
	return &newBlock
}
func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}
func GetBlockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis")
		})
	}
	return b
}

func (b *blockchain) AllBlocks() []*Block {
	return b.blocks
}

var ErrNotFound = errors.New("block not found")

func (b *blockchain) GetBlock(height int) (*Block, error){
	if height > len(b.blocks){
		return nil, ErrNotFound
	}
	return b.blocks[height-1], nil
}



// package blockchain

// import (
// 	"crypto/sha256"
// 	"fmt"
// 	"sync"
// )

// type Block struct{
// 	Data string
// 	Hash string
// 	PrevHash string
// }

// type blockchain struct{
// 	blocks []*Block
// }

// var once sync.Once //한번만 실행되는 go sync
// var b *blockchain

// func (b *Block) calculHash(){
// 	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
// 	b.Hash = fmt.Sprintf("%x", hash)
// }

// func getLastHash() string{
// 	totalBlocks := len(GetBlockchain().blocks)
// 	if totalBlocks == 0 {
// 		return ""
// 	}
// 	return GetBlockchain().blocks[totalBlocks - 1].Hash
// }

// func createBlock(data string) *Block{
// 	newBlock := Block{data, "",getLastHash()}
// 	newBlock.calculHash()
// 	return &newBlock
// }

// //첫블록
// func (b *blockchain) AddBlock(data string){
// 	b.blocks = append(b.blocks, createBlock(data))
// }


// //인스턴스 하나 싱글톤패턴
// func GetBlockchain()*blockchain{
// 	if b == nil {
// 		once.Do(func(){
// 			b = &blockchain{}
// 			b.AddBlock("Genesis")
// 		})
// 	}
// 	return b
// }

// func (b *blockchain) AllBlocks()[]*Block{
// 	return b.blocks
// }
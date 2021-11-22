package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/go_lang_coins/db"
	"github.com/go_lang_coins/utils"
)


type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevhash",omitempty`
	Height 	 int	`json:"height"`
}

var ErrNotFound = errors.New("블록을 찾을 수 없습니다")


func (b *Block) restore(data []byte){
	utils.FromBytes(b, data)
}

//FindBlock 함수는 (hash string)을 받고 *Block을 리턴해줌
func FindBlock(hash string) (*Block, error){
	blockBytes := db.Block(hash)
	if blockBytes == nil{
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}


//persist는 블록을 저장하기위해 만들어 놓은 SaveBlock함수를 호출
func (b *Block) persist(){
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

//새로운 블록을 만들면 블록을 hash하고 이 블록을 db에저장(persist 이용)
func createBlock(data string, prevHash string, height int) *Block{
	block := &Block{
		Data: data,
		Hash: "",
		PrevHash: prevHash,
		Height: height,
	}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.persist()
	return block
}
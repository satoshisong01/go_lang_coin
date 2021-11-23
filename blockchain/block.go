package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go_lang_coins/db"
	"github.com/go_lang_coins/utils"
)


type Block struct {
	Hash     string `json:"hash"`
	PrevHash string `json:"prevhash",omitempty`
	Height 	 int	`json:"height"`
	Difficulty int `json:"difficulty"`
	Nonce	int 	`json:"nonce"`
	Timestamp int	`json:"timestamp"`
	Transactions []*Tx `json:"transactions"`
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

func (b *Block) mine(){
	target := strings.Repeat("0", b.Difficulty) //Repeat 0 몇번 반복
	for{
		b.Timestamp = int(time.Now().Unix()) //unixrk int64를 리턴해준다
		hash := utils.Hash(b)
		fmt.Printf("\n\nTarget:%s\nHash:%s\nNonce:%d\n\n\n", target, hash, b.Nonce)
		//Hasprefix 찾을값 앞쪽 , HasSuffix 찾을값 뒤쪽 
		if strings.HasPrefix(hash, target){
			b.Hash = hash
			break
		} else{
			b.Nonce++
		}
	}
}

//persist는 블록을 저장하기위해 만들어 놓은 SaveBlock함수를 호출
func (b *Block) persist(){
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

//새로운 블록을 만들면 블록을 hash하고 이 블록을 db에저장(persist 이용)
func createBlock(prevHash string, height int) *Block{
	block := &Block{
		Hash: "",
		PrevHash: prevHash,
		Height: height,
		Difficulty: Blockchain().difficulty(),
		Nonce: 0,
		Transactions: []*Tx{makeCoinbaseTx("sks")},
	}
	block.mine()
	block.persist()
	return block
}

package blockchain

import (
	"fmt"
	"sync"

	"github.com/go_lang_coins/db"
	"github.com/go_lang_coins/utils"
)

const(
	defaultDifficulty int = 2
	difficultyInterval int = 5
	blockInterval int = 2
	allowedRange int = 2
)

type blockchain struct {
	NewstHash string `json:"newestHash`
	Height 	 int	`json:"height"`
	CurrentDifficulty int `json:"currentDifficulty"`
}

var b *blockchain
var once sync.Once

//byte를 다시 blockchain으로 변환하는 함수를 만든다
func (b *blockchain) restore(data []byte){
	utils.FromBytes(b, data)
}

func (b *blockchain) persist(){
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string){
	block := createBlock(data, b.NewstHash, b.Height +1)
	b.NewstHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	b.persist()
}

func (b *blockchain) Blocks() []*Block{
	var blocks []*Block //블록 포인터의 slice만든뒤
	hashCursor := b.NewstHash  //찾을 해쉬인 hashCursor만듦(초기에는 newstHash찾음)
	for {
		block, _ := FindBlock(hashCursor) //findblock함수로 NewstHash찾음
		blocks = append(blocks, block) //찾아서 블록 슬라이스에 넣고
		if block.PrevHash != ""{ //Prevhash가 빈값이 아니라면 (최초의 블록은 PrevHash가없으니 나올때까지 계속 타고들어감)
			hashCursor = block.PrevHash //찾을 해쉬를 Prevhash로바꾼다
		} else {
			break
		}
	}
	return blocks
}

func (b *blockchain) recalculateDifficulty() int{
	allBLocks := b.Blocks()
	newestBlock := allBLocks[0] //첫값을 가저옴
	lastrecalculatedBlock := allBLocks[difficultyInterval -1] //0부터 카운팅 하기때문에 -1
	actualTime := (newestBlock.Timestamp/60) - (lastrecalculatedBlock.Timestamp/60) //초단위를 분단위로 바꿈
	expectedTime := difficultyInterval * blockInterval
	if actualTime <= (expectedTime - allowedRange) { //범위를 여유를 준다 ex) 10>0 대신 11 ~ 9 >0
		return b.CurrentDifficulty + 1 //실제 예상 시간보다 적다면, 빨리 생성되니까 +1 로 늘린다
	} else if actualTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1 //실제 예상시간보다 길다면, -1로 줄인다
	} else {
		return b.CurrentDifficulty
	}
}

func (b *blockchain) difficulty() int{
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height % difficultyInterval == 0{ //5단위간격으로 확인
		return b.recalculateDifficulty()
	}else {
		return b.CurrentDifficulty
	}
}



//블록체인을 처음 만들때
func Blockchain() *blockchain {
	if b == nil { //처음에 아무것도 없을때
		once.Do(func() {
			//빈 블록체인을 만들고
			b = &blockchain{
				Height: 0,
			}
			checkpoint := db.Checkpoint()

			//DB에서 체크포인트를 찾는다
			if checkpoint == nil{
				b.AddBlock("새롭게 시작")
			} else{
				b.restore(checkpoint) //db에서 찾은 bytes를 보내준다
			}
			//체크포인트가 있다면 bytes로 부터 블록체인을 복원함
		})
	}
	fmt.Printf("뉴해쉬: %s\n", b.NewstHash)
	return b
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
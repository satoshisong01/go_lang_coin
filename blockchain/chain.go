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
	NewstHash string `json:"newestHash"`
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

func (b *blockchain) AddBlock(){
	block := createBlock(b.NewstHash, b.Height +1)
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


func (b *blockchain) txOuts() []*TxOut{ //거래 출력값들을 가져다주는 함수
	var txOuts []*TxOut
	blocks := b.Blocks() //전체블록 로드
	for _, block := range blocks {//모든 블록을 살펴봄
		for _, tx := range block.Transactions{//모든 블록안에 있는 거래내역을 살펴봄
			txOuts = append(txOuts, tx.TxOuts...)//모든거래 내역들의 출력값들을 하나의 슬라이스로 모음
		}
	}
	return txOuts
}


func (b *blockchain) TxOutsByAddress(address string) []*TxOut{ //거래출력값들을 주소에 따라 걸러내는 함수
	//거래 출력값들을 주어진 address에 따라 필터함
	var ownedTxOuts []*TxOut//address에 소속된 출력값들만 뽑아내고 변수에 지정해줌 (슬라이스)
	txOuts := b.txOuts() //출력값을 리턴했던 txOuts를 다시불러와서 사용
	for _, txOut := range txOuts { //블록 안에 모든 거래출력값에서 나온 출력값에
		if txOut.Owner == address { // 해당 출력값의 주인이 주소와 동일하다면
			ownedTxOuts = append(ownedTxOuts, txOut) // 출력값을 ownedTxOuts 슬라이스에 포함시킴
		}
	}
	return ownedTxOuts
}

//총량 을 보여주는 함수
func (b *blockchain) BalanceByAddress(address string) int { 
	txOuts := b.TxOutsByAddress(address)//주소에따라 거래 출력값들을 받아옴
	var amount int //총액 변수
	for _, txOut := range txOuts { //출력값 목록에 있는 트랜잭션 출력값마다
		amount += txOut.Amount //해당하는 출력값의 총량을 amount 변수에 더해줌
	}			//출력 값을 모두 더해 amount를 반환하는 것
	return amount //리턴
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
				b.AddBlock()
			} else{
				b.restore(checkpoint) //db에서 찾은 bytes를 보내준다
			}
			//체크포인트가 있다면 bytes로 부터 블록체인을 복원함
		})
	}
	fmt.Printf("뉴해쉬: %s\n", b.NewstHash)
	return b
}
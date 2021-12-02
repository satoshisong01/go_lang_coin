package blockchain

import (
	"sync"

	"github.com/go_lang_coins/db"
	"github.com/go_lang_coins/utils"
)

// 함수에 변경값이 없다면 method로 남아있는건 직관적이지 않고 바람직하지 않다
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


func (b *blockchain) AddBlock(){
	block := createBlock(b.NewstHash, b.Height +1, getDifficulty(b))
	b.NewstHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchain(b)
}

func persistBlockchain(b *blockchain){
	db.SaveCheckpoint(utils.ToBytes(b))
}

func Blocks(b *blockchain) []*Block{
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

func recalculateDifficulty(b *blockchain) int{
	allBLocks := Blocks(b)
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

func getDifficulty(b *blockchain) int{
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height % difficultyInterval == 0{ //5단위간격으로 확인
		return recalculateDifficulty(b)
	}else {
		return b.CurrentDifficulty
	}
}

//input에 사용되지 않은 output을 넘겨주는 함수
func UTxOutsByAddress(address string, b *blockchain) []*UTxOut { //거래출력값들을 주소에 따라 걸러내는 함수
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)//키: 트랜색션 ID[string]타입   밸류: bool

	for _, block := range Blocks(b) { //블럭을 참조
		for _, tx := range block.Transactions { //블럭안에 트랜잭션을 참조
			for _, input := range tx.TxIns {	//트랜색션 안의 트랜색션 input 추적
				if input.Owner == address {
					creatorTxs[input.TxID] = true //해당 input으로 사용하는 output을 생성한 트랜잭션을 찾음
				}
			}
			for index, output := range tx.TxOuts { //해당 output이 creatorTxs 안에 있는 트랜잭션 내에 없다는 것을 확인함
				if output.Owner == address {
					if _, ok := creatorTxs[tx.ID]; !ok { //input으로 사용하고 있는 output을 소유한 트랜잭션ID로 들어오지않으면
						uTxOut := &UTxOut{tx.ID, index, output.Amount} //새로 생성된 unspent 트랜색션 output을 확인하면서
						if !isOnMempool(uTxOut){	//이미 mempool에서 사용되고 있는지 체크함(해당 트랜색션ID를 가진 input과 index를 찾아옴)
							uTxOuts = append(uTxOuts, uTxOut) //아직할당하지 않은것이므로 uTxouts을 찾은것임
						}
					}
				}
			}
		}
	}
		return uTxOuts
}

//총량 을 보여주는 함수
func BalanceByAddress(address string, b *blockchain) int {
	txOuts := UTxOutsByAddress(address, b)//주소에따라 거래 출력값들을 받아옴
	var amount int //총액 변수
	for _, txOut := range txOuts { //출력값 목록에 있는 트랜잭션 출력값마다
		amount += txOut.Amount //해당하는 출력값의 총량을 amount 변수에 더해줌
	}			//출력 값을 모두 더해 amount를 반환하는 것
	return amount //리턴
}

//블록체인을 처음 만들때
func Blockchain() *blockchain {
	once.Do(func() {
		b = &blockchain{
			Height: 0,
		}
		checkpoint := db.Checkpoint()
		if checkpoint == nil{
			b.AddBlock()
		} else{
			b.restore(checkpoint)
		}
	})
	return b
}
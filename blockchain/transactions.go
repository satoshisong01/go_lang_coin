package blockchain

import (
	"errors"
	"time"

	"github.com/go_lang_coins/utils"
)

const (
	minerReward int = 50
)

//mempool은 아직 성립되지 않은 transaction이 들어가는 곳

type mempool struct{
	Txs []*Tx
}

//비어있는 mempool생성
var Mempool *mempool = &mempool{}

type Tx struct{
	ID string		`json:"id"`
	Timestamp int	`json:"timestamp"`
	TxIns []*TxIn	`json:"txIns"`
	TxOuts []*TxOut	`json:"txOuts"`
}

func (t *Tx) getId(){
	t.ID = utils.Hash(t)
}

type TxIn struct{
	TxID string `json:"txId"` //어떻게 이전의 트랜잭션 output을 찾을것인가 (ID를 통해서)
	Index int	`json:"index"` //이 트랜잭션의 트랜잭션 output의 어디에 해당 트랜잭션 output이 위치해 있는지 찾는다
	Owner string	`json:"owner_TxIn"`
}

type TxOut struct{
	Owner string	`json:"owner_TxOut"`
	Amount int		`json:"amount_TxOut"`
}

type UTxOut struct{
	TxID string `json:"txId"`
	Index int	`json:"index"`
	Amount int	`json:"amount"`
}

//mempool에 있는 트랜잭션에 존재하는 input들 중에 uTxOut과 같은 트랜잭션 ID와 index를 가지고 있는 항목이 있는지 확인해줌
//*mempoop에 현재 내가 추가하고자 하는 unspent 트랜잭션 output이 존재하는지 확인해주는 함수* (true라면 chain.go에서 해당함수 실행하지않음)
func isOnMempool(uTxOut *UTxOut) bool{ //추가하려는 unspent 트랜잭션 output이  mempool에 아직 없는지 확인
	exists := false
	Outer: //label 사용
		for _, tx := range Mempool.Txs { //mempool안에 들어있는 트랜잭션을 둘러봄
			for _, input := range tx. TxIns { //트랜잭션 안에있는 트랜잭션 input을 둘러봄
				if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
					exists = true
					break Outer
				}
			}
		}
	return exists
}

//채굴자를 주소로 삼는 코인베이스 거래내역을 생성해서 Tx포인터를 리턴
func makeCoinbaseTx(address string) *Tx{
	txIns := []*TxIn{ //트랜잭션 input은 ID, index, owner가 필요함
		{"", -1, "COINBASE"},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		ID: "",
		Timestamp: int(time.Now().Unix()),
		TxIns: txIns,
		TxOuts: txOuts,
	}
	tx.getId()
	return &tx
}

//transaction을 생성하는 함수
func makeTx(from, to string, amount int) (*Tx, error){
	if BalanceByAddress(from, Blockchain()) < amount {
		return nil, errors.New("돈이 충분하지 않습니다")
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0 //트랜색션 output으로 부터 받아온 총 잔고량
	uTxOuts := UTxOutsByAddress(from, Blockchain()) //unspent 만 가져온다 (아직 input과 output할장 되지않은것)
	for _, uTxOut := range uTxOuts { //각각의 트랜색션 output별로 트랜색션 input을 생성
		if total >= amount {
			break
		}			  
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount //total 에 할당되지않은 output(uTxOut)의 amount 할당
	}
	if change := total - amount; change != 0 { //잔돈(change)은 = total - amount / change 가 0이 아니라면 반환
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{ //트랜잭션 생성
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId() // 트랜잭션에 getId함수를 통해 id 생성
	return tx, nil //트랜잭션 반환, 에러는 빼고
}


//AddTx 함수는 mempool에 transaction을 추가할뿐, transaction을 생성하지는 않음
func (m *mempool) AddTx(to string, amount int) error{
	tx, err := makeTx("sks", to, amount)
	if err != nil { //에러를 반환받았을때
		return err
	}
	m.Txs = append(m.Txs, tx) //transaction을 반환받았을때
	return nil
}

func (m *mempool) TxToConfirm() []*Tx{	//모든 transaction들을 건내주고 mempool을 비워주는 함수
	coinbase := makeCoinbaseTx("sks")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil		//mempool에서 transaction을 비워줌
	return txs
}



// func makeTx(from, to string, amount int) (*Tx, error){
// 	if Blockchain().BalanceByAddress(from) < amount{ //잔금이 보내는 금액보다 적다면
// 		return nil, errors.New("돈이 충분하지 않습니다.")
// }
// var txIns []*TxIn //트랜잭션 인풋
// var txOuts []*TxOut //트랜잭션 아웃풋
// total := 0
// oldTxOuts := Blockchain().TxOutsByAddress(from)//(transaction outputs로부터 transaction inputs를 생성)
// for _, txOut := range oldTxOuts{
// 	if total > amount { //total 이 amount(잔고)보다 큰지확인
// 		break
// 	}
// 	//새로운 transaction input을 이전의 transaction outputs로 부터 생성 (트랜잭션 인풋에는 아웃풋 으로부터오는 owner가 필요하고 트랜잭션 아웃풋 Amount(금액)도 필요하다)
// 	txIn := &TxIn{txOut.Owner, txOut.Amount} //이전 트랜잭션 아웃풋으로부터 트랜잭션 인풋을 생성
// 	txIns = append(txIns, txIn)//새로만든 transaction input을 txIns 으로 새로만든다
// 	total += txIn.Amount //total에는 transaction input의 amount를 대입해준다
// }
// change := total - amount
// if change != 0 { //잔돈이 0원이 아니라면
// 	changeTxOut := &TxOut{from, change} //잔돈을 위한 transaction output 생성 (다시 from 유저에게 돌아감)
// 	txOuts = append(txOuts, changeTxOut) //txOuts를 만들고 잔돈(changeTxout)을 넣어줌
// }
// txOut := &TxOut{to, amount}	//from 유저가 to 유저에게 보내고 싶은 amount만큼의 transaction output (to유저의 주소와 보내는 amount)
// txOuts = append(txOuts, txOut) //txOuts에  txOut넣어줌
// tx := &Tx{
// 	Id: "",
// 	Timestamp: int(time.Now().Unix()),
// 	TxIns: txIns,
// 	TxOuts: txOuts,
// }
// tx.getId() //해쉬변환
// return tx, nil
// }
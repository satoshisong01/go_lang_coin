package db

import (
	"github.com/boltdb/bolt"
	"github.com/go_lang_coins/utils"
)

const (
	dbName = "blockchain.db"
	dataBucket = "data"
	blocksBucket = "blocks"

	checkpoint = "checkpoint"
)

var db *bolt.DB
//Singleton패턴을 사용해서 export되지 않는 변수를 만듬

//*버킷은 table과 비슷하지만 key/value db에서 사용됨*
//export 되는 DB함수만듦
//DB함수는 dbPointer에 접근하게함
func DB() *bolt.DB{
	//dbPointer가 존재하지 않는다면 DB를 열고 db를 dbPointer에 지정한뒤
	if db == nil{
		dbPointer, err := bolt.Open(dbName, 0600, nil) //path는 db이름, 파일이없으면 자동생성
		db = dbPointer
		utils.HandleErr(err) //에러 처리
		//버킷이 존재하지 않으면 생성시켜주는 트랜색션 생성
		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket)) //data버킷이 아니면 err리턴
			utils.HandleErr(err)
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket)) //블록버킷이 아니면 err 리턴
			// 위에 err 명시후 := 가아닌 = 만 사용
			return err
		})
		utils.HandleErr(err)
	}
	return db
}

func Close(){
	DB().Close()
}

//블록 db에 저장(키,벨류) 벨류는 byte여야만함(볼트를 쓰기때문)
func SaveBlock(hash string, data []byte){
	err := DB().Update(func(t *bolt.Tx) error {
		//만들어준 버킷지정
		bucket := t.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data) // hash : key , data : value
		return err
	})
	utils.HandleErr(err)
}

//블록체인을 저장할때는 key없이 data만 필요하다
func SaveBlockchain(data []byte){
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data) //checkpoint : key , data : value
		return err
	})
	utils.HandleErr(err)
}

//아무인자도 받지않고 byte리턴 하는 블록체인함수 생성
func Checkpoint() []byte{
	var data []byte
	//read only 뷰 를 만든다
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}

//체크포인트 처럼 해쉬로 값찾기(DB에 blockBucket에서 특정 블록을 찾을 수 있음)
func Block(hash string) []byte{
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}
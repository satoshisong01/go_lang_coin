package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

func HandleErr(err error){
	if err != nil {
		log.Panic(err)
	}
}


//바이트로 바꿔주는 함수
//*블록은 인코더 하고 aBuffer 저장*
//interface 는 타입상관없이 받아옴
func ToBytes(i interface{}) []byte{
	var aBuffer bytes.Buffer //비어있는 변수생성(type은 bytes.Buffer) Buffer에는 bytes를 넣을 수 있고 read-write가능
	encoder := gob.NewEncoder(&aBuffer) //인코더 생성
	HandleErr(encoder.Encode(i)) //블록전체 인코딩 (바이트로바꿈)
	return aBuffer.Bytes() //aBuffer bytes리턴
}

//복원해 주는 함수(FromBytes에 포인터와 복원 할 data를 보냄)
func FromBytes(i interface{}, data []byte){
	encoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(encoder.Decode(i))
}

//%v 
func Hash(i interface{}) string {
	s := fmt.Sprint("%v", i)
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}
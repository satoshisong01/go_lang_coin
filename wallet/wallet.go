package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/go_lang_coins/utils"
)

const filename string ="krono.wallet"

type wallet struct{
	privateKey *ecdsa.PrivateKey //노출 x
	Address string //외부공유
}

var w *wallet

func hasWalletFile() bool { //월렛이 있다면 해당함수 실행
	_, err := os.Stat(filename) //파일정보는 _ , err만 받아옴
	return !os.IsNotExist(err) //파일이 이미 생성되어 있을 때 발생하는 에러라면 true 반환
}

func createPriveKey() *ecdsa.PrivateKey{ //비공개키 만드는 함수 (지갑이 없을경우 생성)
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey //go언어는 어떤  variable이 type을 가지고 있는지 모르기때문에
}

func persistKey(key *ecdsa.PrivateKey){
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = os.WriteFile(filename, bytes, 0644) //읽기와 쓰기 허용 (파일있는지 체크)
	utils.HandleErr(err)
}


func restoreKey() (key *ecdsa.PrivateKey) { //지갑이 있다면 키를 복원한다
	keyAsBytes, err := os.ReadFile(filename)
	utils.HandleErr(err)
	key, err = x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)
	return
}

func encodeBigInts(a, b []byte) string {
	z := append(a, b...)
	return fmt.Sprintf("%x", z)
}


func aFromK(key *ecdsa.PrivateKey) string { //key에서부터 주소를 만들어내는 함수
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
}

func Sign(payload string, w *wallet) string{ //서명
	payloadAsB, err := hex.DecodeString(payload) //string에서 byte가져오기 (string -> byte 변환)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsB)
	utils.HandleErr(err)
	return encodeBigInts(r.Bytes(), s.Bytes())
}

func restoreBigInts(payload string) (*big.Int, *big.Int, error){
	bytes, err := hex.DecodeString(payload)
	if err != nil{
		return nil, nil, err
	}
	firstHalfBytes := bytes[:len(bytes)/2]
	sencodHalfBytes := bytes[len(bytes)/2:]
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(sencodHalfBytes)
	return &bigA, &bigB, nil
}


func Verify(signature, payload, address string) bool{ //검증 (publicKey string으로 변환한 address 대신받음)
	r, s, err := restoreBigInts(signature)
	utils.HandleErr(err)
	x, y, err := restoreBigInts(address)
	utils.HandleErr(err)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),//어떤 종류의 publicKey인지 명시
		X:     x,
		Y:     y,
	}
	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	ok := ecdsa.Verify(&publicKey, payloadBytes, r, s)
	return ok
}

func Wallet() *wallet{
	if w == nil{
		w = &wallet{}
		//사용자가 지갑을 가지고 있는지 체크
		if hasWalletFile(){
			//가지고 있다면, 그 지갑을 파일로부터 복구
			w.privateKey = restoreKey()
		} else {
			key := createPriveKey()
			persistKey(key)
			w.privateKey = key
			//지갑이 없다면, 비공개키를 생성해서 파일에 저장
		}
		w.Address = aFromK(w.privateKey)
	}
	return w
}

//1. 문자 등을 해쉬화 한다
//2. 키페어를 생성한다 (공개키, 비공개키)(지갑생성) 비공개키-서명 공개키-검증
//3. 해쉬를 서명한다 (1번 + 2번 비공개키) -> 서명
//4 검증한다 (1번해쉬 + 3번 서명 + 공개키) -> true or false 반환

//privateKey생성 (파일로 저장하려면 공개키를 문자 나 byte로 바꿀 수 있어야함)
// privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) //알고리즘 선정후 비공개키 공개키설정 rand.Reader (암호학적으로 보안된 난수생성기의 전역 공유 인스턴스)

// keyAsBytes, err := x509.MarshalECPrivateKey(privateKey) //비공개키를 받아서 byte로 변환

// fmt.Printf("%x\n\n\n\n",keyAsBytes)

// utils.HandleErr(err)

// //해싱한 메세지를 byte로 변환
// hashAsBytes, err := hex.DecodeString(hashedMessage)

// utils.HandleErr(err)
// //Sing함수를 사용해 난수생성기, 비공개키, 해싱된 byte전달
// r, s, err := ecdsa.Sign(rand.Reader, privateKey,hashAsBytes) //Sign함수는 비공개키 필요
// //서명이 r,s 로 나눠져있음 (메세지 해싱의 서명)

// singature := append(r.Bytes(), s.Bytes()...)

// fmt.Printf("%x\n", singature)

// utils.HandleErr(err)

// 복구 ----------------

// privBytes, err := hex.DecodeString(pricateKey) //비공개키를 이용해  문자열을 가져온다 singature를구함(검사) (비공개키의 인코딩 방식이 16진수 형식인지 판단)

// 	utils.HandleErr(err)

// 	restoredKey, err := x509.ParseECPrivateKey(privBytes)
// 	utils.HandleErr(err)

// 	sigBytes, err := hex.DecodeString(singature) //singature(서명) 에서 r 과 s 값을 구하기위해 바이트로 변환

// 	rBytes := sigBytes[:len(sigBytes)/2] //시작부터 길이의 반 (r의값)
// 	sBytes := sigBytes[len(sigBytes)/2:] // 반부터 끝까지 (s의값)

// 	var bigR, bigS = big.Int{}, big.Int{}

// 	bigR.SetBytes(rBytes)//byte를 big.Int로 변환 get과는 반대
// 	bigS.SetBytes(sBytes)

// 	hashBytes, err := hex.DecodeString(hashedMessage)

// 	utils.HandleErr(err)

// 	ok := ecdsa.Verify(&restoredKey.PublicKey, hashBytes, &bigR, &bigS)

// 	fmt.Println(ok)
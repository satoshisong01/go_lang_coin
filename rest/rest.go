package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go_lang_coins/blockchain"
	"github.com/go_lang_coins/utils"
	"github.com/gorilla/mux"
)

var port string

type url string //내재 되어있기 때문에 type URL string <(TextMarshaler, Stringer) 생략가능>


//원하는 url로 변경할수 있다  예) naver.com/index -< naver.com/인덱스 등
//MarshalText로 구현
func (u url) MarshalText() ([]byte, error){
	url := fmt.Sprintf("http://localhost%s%s", port, u) //url 내용변경
	return []byte(url), nil
}

type urlDescription struct{
	URL url `json:"url"` //json 표기명 변경
	Method string `json:"method"`
	Description string `json:"description"`
	Payload string `json:"payload,omitempty"`
	//json:",omitempty" 쓰면 파일이 없을때 표기안함
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type errorResponse struct{
	ErrorMessage string `json:"errorMessage"`
}

type addTxPayload struct{
	To string
	Amount int
}

//*Marshal은 Json으로 encoding한 interface(v)를 return해준다*
//*Marshal은 ~메모리형식으로 저장된 객체를, 저장/송신 할 수 있도록 변환해 준다
func documentation(rw http.ResponseWriter, r *http.Request){
	data := []urlDescription{
		{
			URL: url("/"),
			Method: "GET",
			Description: "보여줘 항목들을!",
		},
		{
			URL: url("/status"),
			Method: "GET",
			Description: "블록체인의 상태를 보여줘",
		},
		{
			URL: url("/blocks"),
			Method: "GET",
			Description: "블록을 보여줘",
		},
		{
			URL: url("/blocks"),
			Method: "POST",
			Description: "블록을 추가하자",
			Payload: "data:string",
		},
		{
			URL: url("/blocks{hash}"),
			Method: "GET",
			Description: "블록하나만 보여줘",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "해당 주소의 거래출력값들을 구함",
		},
	}
	
	// b, err := json.Marshal(data)
	// utils.HandleErr(err)
	// fmt.Fprintf(rw, "%s", b)
	utils.HandleErr(json.NewEncoder(rw).Encode(data)) //위 주석코드 3줄과 같음
}
//write -> Encoder, read -> Decode
func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method{
		case "GET":
			utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.Blockchain())))
		case "POST":
			blockchain.Blockchain().AddBlock()
			rw.WriteHeader(http.StatusCreated)
	}
}


func block(rw http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r) //mux가 request에서 변수를 추출함
	hash := vars["hash"] //vars 안에있는 hash를 가져올수있음
	fmt.Println(hash)
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound{
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	}else{
		encoder.Encode(block)
	}
	
}

func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.Blockchain())
}

//json 명시 미들웨어
func jsonContentTypeMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request){
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

//주소에 /blocks{address} 같이 넣기위에 고릴라 mux를 쓴다
func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) //변수를받고
	address := vars["address"] //map에서 뽑아서씀
	total := r.URL.Query().Get("total") //?total=ture 통해서 자산 총액을 받음
	switch total {
	case "true":
		amount := blockchain.BalanceByAddress(address, blockchain.Blockchain())
		json.NewEncoder(rw).Encode(balanceResponse{address, amount})
	default:
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(address, blockchain.Blockchain()))) //거래 출력값
	}
}

func mempool(rw http.ResponseWriter, r *http.Request){
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempool.Txs))
}

func transactions(rw http.ResponseWriter, r *http.Request){
	var payload addTxPayload //받을 사람과 금액이 있는 구조체 가져옴
	//decode된 값을 payload 변수로 넘겨준다 (JSON은 request body의 데이터를 decode해서 struct(payload)로 변환시켜준다)
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
	err := blockchain.Mempool.AddTx(payload.To, payload.Amount) //blockchain의 Mempool.AddTx()을 호출해서 보낼 금액과 대상을 넘겨줄수있다
	if err != nil {
		json.NewEncoder(rw).Encode(errorResponse{"잔액이 부족합니다 선생님"})
	}
	rw.WriteHeader(http.StatusCreated) //else인 경우에는 HTTP 상태를 created로 전송
}

func Start(aPort int){
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware) //미들웨어 사용
	router.HandleFunc("/", documentation).Methods("GET") //Methods로 (Get만 쓸지 POST도 쓸지 고를수있다)
	router.HandleFunc("/status", status)
	router.HandleFunc("/blocks", blocks).Methods("GET","POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance) // balance{address}를 통해서 거래출력값 목록을받음
	router.HandleFunc("/mempool", mempool)
	router.HandleFunc("/transactions", transactions).Methods("POST")
	fmt.Printf("http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
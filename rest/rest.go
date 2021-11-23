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

type addBlockBody struct{
	Message string
}

type errorResponse struct{
	ErrorMessage string `json:"errorMessage"`
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
	}
	
	// b, err := json.Marshal(data)
	// utils.HandleErr(err)
	// fmt.Fprintf(rw, "%s", b)
	json.NewEncoder(rw).Encode(data) //위 주석코드 3줄과 같음
}
//write -> Encoder, read -> Decode
func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method{
		case "GET":
			json.NewEncoder(rw).Encode(blockchain.Blockchain().Blocks())
		case "POST":
			var addBlockBody addBlockBody
			utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
			blockchain.Blockchain().AddBlock(addBlockBody.Message)
			rw.WriteHeader(http.StatusCreated)
	}
}

func block(rw http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	hash := vars["hash"]
	fmt.Println(hash)
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound{
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	}else{
		encoder.Encode(block)
	}
	
}
//json 명시 미들웨어
func jsonContentTypeMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request){
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.Blockchain())
}

func Start(aPort int){
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware) //미들웨어 사용
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status)
	router.HandleFunc("/blocks", blocks).Methods("GET","POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	fmt.Printf("http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
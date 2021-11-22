package main

import (
	"github.com/go_lang_coins/cli"
	"github.com/go_lang_coins/db"
)

//rw (보내고싶은 데이터)

//go explorer.Start(3000)
//rest.Start(4000)


func main(){
	defer db.Close() //defer은 함수가 종료될때 실행
	cli.Start()
	
}
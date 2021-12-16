package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/go_lang_coins/explorer"
	"github.com/go_lang_coins/rest"
)

func usage(){
	fmt.Printf("어서오세요 블록체인에\n\n")
	fmt.Printf("명령어를 입력해 주세요\n\n")
	fmt.Printf("-port: 서버의 포트\n")
	fmt.Printf("-mode: 'REST' 와 'HTML'을 고르시오\n")
	os.Exit(0)
}

func Start(){												//gohtml 로 끝나는 파일 불러오기

if len(os.Args) == 1{
	usage()
	}

	port := flag.Int("port", 4000, "서버의 포트")
	mode := flag.String("mode", "rest", "html REST API")

	flag.Parse()

	switch *mode{
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	default:
		usage()
	}
}


//FlagSet 정의
// rest :=flag.NewFlagSet("rest", flag.ExitOnError)

// portFlag := rest.Int("port", 4000, "서버의 포트")

// switch os.Args[1]{
// case "explorer":
// 	fmt.Println("Explorer 시작")
// case "rest":
// 	rest.Parse(os.Args[2:])
// default:
// 	usage()

// fmt.Println(*portFlag)
// }
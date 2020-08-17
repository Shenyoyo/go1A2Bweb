package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func gennum(length int) []int {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	b := make([]int, length, length)                     //創建length長度的陣列，儲存答案
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //放入亂數環境隨機生成
	for i := 0; i < length; i++ {
		tmp := r.Intn(len(a) - i) //亂數讀取10-i內的數字
		b[i] = a[tmp]             //將a[tmp]值放入答案
		for b[0] == 0 {
			tmp = r.Intn(len(a) - i) //第一位數不能為0
			b[0] = a[tmp]
		}
		a[tmp], a[len(a)-1-i] = a[len(a)-1-i], a[tmp] //使答案數字不重複
	}
	return b
}

func validateAnswer(answer string) (bool, string) {
	flag := true
	var isString bool
	var ErrorMessage string
	//檢查輸入是否有字串
	for _, r := range answer {
		// fmt.Printf("%c = %v\n", r, unicode.IsLetter(r))
		if unicode.IsLetter(r) {
			isString = true
			break
		} else {
			isString = false
		}
	}
	if isString {
		ErrorMessage = "輸入有誤，請輸入數字"
		flag = false
		return flag, ErrorMessage
	}
	if len(answer) != 4 {
		ErrorMessage = "輸入有誤，請輸入長度4的數字"
		flag = false
		return flag, ErrorMessage
	}
	ipulens := strings.Split(answer, "")
	if ipulens[0] == "0" {
		ErrorMessage = "第一位不能為0，請重新輸入"
		flag = false
		return flag, ErrorMessage
	}
	if len(removeDuplicateElement(ipulens)) != 4 {
		ErrorMessage = "輸入答案重複數字，請重新輸入"
		flag = false
		return flag, ErrorMessage
	}

	return flag, ErrorMessage

}

func checknum(a, b []int) (bool, string) {
	var message string
	var aa, bb int
	if len(a) != len(b) {
		message = "答案與輸入長度不一致，請聯絡管理員"
		return false, message
	}
	dict := make(map[int]int)
	for i := 0; i < len(a); i++ {
		dict[a[i]] = 0
	}
	for i := 0; i < len(a); i++ {
		if _, ok := dict[b[i]]; ok {
			bb++
		}
	}
	for i := 0; i < len(a); i++ {
		if a[i] == b[i] {
			aa++
		}
	}
	bb = bb - aa
	answerA := strconv.Itoa(aa) + "A"
	answerB := strconv.Itoa(bb) + "B"

	message = answerA + answerB
	if aa == len(a) {
		return true, message
	} else {
		return false, message
	}

}

func removeDuplicateElement(addrs []string) []string {
	result := make([]string, 0, len(addrs))
	temp := map[string]struct{}{}
	for _, item := range addrs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func main() {
	//變數宣告
	var historyAnswer []string
	var historyAB []string
	var genAnswer string
	var src []int
	var i int
	var count int
	var messageAB string
	status := false
	length := 4
	//使用gin框架
	router := gin.Default()
	//可使其他網站讀取到此API路由
	router.Use(cors.Default())
	router.GET("/game", func(c *gin.Context) {
		answer := c.Query("answer")
		//驗證輸入資料
		flag, ErrorMessages := validateAnswer(answer)
		if flag {
			//產生答案與答對初始化
			if genAnswer == "" || status {
				historyAnswer = historyAnswer[:0]
				historyAB = historyAB[:0]
				status = false
				count = 0
				genAnswer = ""
				src = gennum(length)

				for _, element := range src {
					genAnswer = genAnswer + strconv.Itoa(element)
				}
			}
			//處理輸入的答案
			req := make([]int, length, length)
			ipu, _ := strconv.Atoi(answer)
			for i = length - 1; ipu > 0; i-- {
				req[i] = ipu % 10
				ipu = ipu / 10
			}
			//檢查是否正確，並回傳提示
			status, messageAB = checknum(src, req)
			historyAB = append(historyAB, messageAB)
			historyAnswer = append(historyAnswer, answer)
			//計算次數
			count++
			//log
			fmt.Println(genAnswer)
			// fmt.Println(req)
			// fmt.Println(src)
			if status {
				c.String(http.StatusOK, "恭喜答對了，答案就是 %s，總共猜了 %v 次", answer, count)
			} else {
				for i, history := range historyAB {
					c.String(http.StatusOK, historyAnswer[i]+"  "+history+"，第"+strconv.Itoa(i)+"次的猜測<br>")
				}

			}
		} else {
			c.String(http.StatusOK, ErrorMessages)
		}

	})

	router.Run(":8081")
}

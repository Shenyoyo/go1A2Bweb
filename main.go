package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
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
	}
	return false, message

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

//AccessJsMiddleware use
func AccessJsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.Writer
		r := c.Request
		// 處理js-ajax跨域问题
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Add("Access-Control-Allow-Headers", "Access-Token")
		c.Next()
	}
}

type initial struct {
	historyAnswer []string
	historyAB     []string
	genAnswer     string
	src           []int
	i             int
	count         int
	messageAB     string
	status        bool
}

func main() {
	//變數宣告
	var init = new(initial)

	length := 4
	//使用gin框架
	router := gin.Default()
	//可使其他網站讀取到此API路由
	router.Use(AccessJsMiddleware())
	router.GET("/index", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://shenyoyo.github.io/go1A2Bweb/")
	})
	router.GET("/game", func(c *gin.Context) {
		cookie, err := c.Cookie("uuid")
		if err != nil {
			sID, _ := uuid.NewV4()
			uuid := sID.String()
			cookie = "NotSet"
			c.SetCookie("uuid", uuid, 3600, "/", "localhost", false, false)
		}
		fmt.Printf("Cookie value: %s \n", cookie)
		answer := c.Query("answer")
		//驗證輸入資料
		flag, ErrorMessages := validateAnswer(answer)
		if flag {
			//產生答案與答對初始化
			if init.genAnswer == "" || init.status {
				//初始化
				init.historyAnswer = init.historyAnswer[:0]
				init.historyAB = init.historyAB[:0]
				init.status = false
				init.count = 0
				init.genAnswer = ""
				init.src = gennum(length)

				for _, element := range init.src {
					init.genAnswer = init.genAnswer + strconv.Itoa(element)
				}
			}
			//處理輸入的答案
			req := make([]int, length, length)
			ipu, _ := strconv.Atoi(answer)
			for i := length - 1; ipu > 0; i-- {
				req[i] = ipu % 10
				ipu = ipu / 10
			}
			//檢查是否答對，並回傳AB提示
			init.status, init.messageAB = checknum(init.src, req)
			//歷史紀錄
			init.historyAB = append(init.historyAB, init.messageAB)
			init.historyAnswer = append(init.historyAnswer, answer)
			//計算次數
			init.count++
			//log
			fmt.Println(init.genAnswer)
			// fmt.Println(req)
			// fmt.Println(src)
			if init.status {
				c.String(http.StatusOK, "恭喜答對了，答案就是 %s，總共猜了 %v 次", answer, init.count)
			} else {
				for i, history := range init.historyAB {
					c.String(http.StatusOK, init.historyAnswer[i]+"  "+history+"，第"+strconv.Itoa(i)+"次的猜測<br>")
				}
			}
		} else {
			c.String(http.StatusOK, ErrorMessages)
		}

	})

	router.Run(":8085")
}

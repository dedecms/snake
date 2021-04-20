package snake

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/width"
)

type snaketext struct {
	Input string
}

// String ...
type String interface {
	Add(str ...string) String         // 在当前Text后追加字符
	Find(dst string) bool             // 判断字符串或符合正则规则的字符串是否存在
	Keep(dst string) String           // 保留符合正则规则的字符串或指定字符串
	Remove(dst ...string) String      // 删除符合正则规则的字符串或指定字符串
	Replace(src, dst string) String   // 替换字符串或符合正则规则的字符串
	Widen() String                    // 半角字符转全角字符
	Narrow() String                   // 全角字符转半角字符
	ReComment() String                // 去除注解
	Between(start, end string) String // 取A字符与B字符之间的字符
	Trim(sep string) String           // 去除首尾的特定字符
	ToLower() String                  // 英文字母全部转为小写
	ToUpper() String                  // 英文字母全部转为大写
	LcFirst() String                  // 英文首字母小写
	UcFirst() String                  // 英文首字母大写
	EnBase(base int) String           // 将Text转为2～36进制编码
	DeBase(base int) String           // 将2～36进制解码为Text
	CamelCase() String                // 将英文字符转为驼峰格式
	SnakeCase() String                // 将英文字符转为蛇形格式
	KebabCase() String                // 将英文字符转化为“烤串儿”格式
	Lines() []string                  // 将行转为数组
	MD5() string                      // 输出字符串MD5
	Split(sep string) []string        // 通过特定字符分割Text
	SplitPlace(sep []int) []string    // 根据字符串的位置进行分割
	SplitInt(sep int) []string        // 根据字数进行分割
	Get() string                      // 输出Text
}

// ---------------------------------------
// 输入 :

// Text 初始化...
func Text(str ...string) String {
	t := &snaketext{}
	if len(str) > 0 {
		t.Add(str...)
	}
	return t
}

// Add 在字符串中追加文字...
func (t *snaketext) Add(str ...string) String {
	b := bytes.NewBufferString(t.Input)
	if len(str) > 0 {
		for _, v := range str {
			b.WriteString(v)
		}
	}
	t.Input = b.String()
	return t
}

// ---------------------------------------
// 处理 :

// Replace 替换字符串或符合正则规则的字符串 ...
// snake.Text("http://www.dedecms.com").Replace("(http://).*(dedecms.com)", "${1}${2}")
// out: http://dedecms.com
// snake.Text("http://www.example.com").Replace("example", "dedecms")
// out: http://www.dedecms.com
// 如需替换$等字符，请使用\\$
// snake.Text("http://$1example.com").Replace("\\$1.*(.com)", "www.dedecms${1}")
func (t *snaketext) Replace(src, dst string) String {
	t.Input = regexp.MustCompile(src).ReplaceAllString(t.Input, dst)
	return t
}

// Find 判断字符串或符合正则规则的字符串是否存在 ...
func (t *snaketext) Find(dst string) bool {
	if d := regexp.MustCompile(dst).FindAll([]byte(t.Input), -1); len(d) > 0 {
		return true
	}
	return false
}

// Remove 根据正则规则删除字符串 ...
func (t *snaketext) Remove(dst ...string) String {
	if len(dst) > 0 {
		for _, v := range dst {
			t.Input = regexp.MustCompile(v).ReplaceAllString(t.Input, "")
		}
	}

	return t
}

// Keep 根据正则规则保留字符串 ...
func (t *snaketext) Keep(dst string) String {

	if t.Find(dst) == true {
		p := Text()
		d := regexp.MustCompile(dst).FindAll([]byte(t.Get()), -1)

		for _, v := range d {
			p.Add(string(v))
		}

		t.Input = p.Get()
	}

	return t
}

// Keep 根据正则规则保留字符串 ...
func (t *snaketext) Extract(dst string) []string {

	if t.Find(dst) == true {
		p := Text()
		d := regexp.MustCompile(dst).FindAll([]byte(t.Get()), -1)

		for _, v := range d {
			fmt.Println(v)
		}

		t.Input = p.Get()
	}

	return []string{}
}

// Narrow 全角字符转半角字符 ...
func (t *snaketext) Narrow() String {
	t.Input = width.Narrow.String(t.Input)
	return t
}

// Widen 半角字符转全角字符 ...
func (t *snaketext) Widen() String {
	t.Input = width.Narrow.String(t.Input)
	return t
}

// ReComment 去除代码注解...
func (t *snaketext) ReComment() String {
	t.Remove(
		`\/\/.*`,
		`\/\*(\s|.)*?\*\/`,
		`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)
	return t
}

// Trim 去除开始及结束出现的字符 ...
func (t *snaketext) Trim(sep string) String {
	t.Input = strings.Trim(t.Input, sep)
	return t
}

// ToLower 字符串全部小写 ...
func (t *snaketext) ToLower() String {
	t.Input = strings.ToLower(t.Input)
	return t
}

// ToUpper 字符串全部小写 ...
func (t *snaketext) ToUpper() String {
	t.Input = strings.ToUpper(t.Input)
	return t
}

// UCFirst 字符串首字母大写 ...
func (t *snaketext) UcFirst() String {
	t.Input = ucfirst(t.Input)
	return t
}

// LCFirst 字符串首字母小写 ...
func (t *snaketext) LcFirst() String {
	t.Input = lcfirst(t.Input)
	return t
}

// Between 截取区间内容 ...
func (t *snaketext) Between(start, end string) String {
	if (start == "" && end == "") || t.Input == "" {
		return t
	}
	// 处理数据，将所有文字转为小写 .
	input := strings.ToLower(t.Input)
	lowerStart := strings.ToLower(start)
	lowerEnd := strings.ToLower(end)

	var startIndex, endIndex int

	if len(start) > 0 && strings.Contains(input, lowerStart) {
		startIndex = len(start)
	}
	if len(end) > 0 && strings.Contains(input, lowerEnd) {
		endIndex = strings.Index(input, lowerEnd)
	} else if len(input) > 0 {
		endIndex = len(input)
	}
	// 输出字符A与字符B之间的字符 .
	t.Input = strings.TrimSpace(t.Input[startIndex:endIndex])
	return t
}

// EnBase Text to Base-x:  2 < base > 36 ...
// 将Text转为2～36进制编码
func (t *snaketext) EnBase(base int) String {
	var r []string
	for _, i := range []rune(t.Input) {
		r = append(r, strconv.FormatInt(int64(i), base))
	}
	t.Input = strings.Join(r, " ")
	return t
}

// DeBase Text Base-x to Text:  2 < base > 36 ...
// 将2～36进制解码为Text
func (t *snaketext) DeBase(base int) String {
	var r []rune
	for _, i := range t.Split(" ") {
		i64, err := strconv.ParseInt(i, base, 10)
		if err != nil {
			panic(err)
		}
		r = append(r, rune(i64))
	}
	t.Input = string(r)
	return t
}

// ---------------------------------------
// 分词 :

// CamelCase 驼峰分词: HelloWord ...
func (t *snaketext) CamelCase() String {
	caseWords := t.caseWords(true)
	for i, word := range caseWords {
		caseWords[i] = ucfirst(word)
	}
	t.Input = strings.Join(caseWords, "")
	return t
}

// SnakeCase 贪吃蛇分词: hello_word ...
func (t *snaketext) SnakeCase() String {
	caseWords := t.caseWords(false)
	t.Input = strings.Join(caseWords, "_")
	return t
}

// KebabCase "烤串儿"分词: hello-word ...
func (t *snaketext) KebabCase() String {
	caseWords := t.caseWords(false)
	t.Input = strings.Join(caseWords, "-")
	return t
}

// ---------------------------------------
// 输出 :

// Get 获取文本...
func (t *snaketext) Get() string {
	return t.Input
}

// Split 根据字符串进行文本分割 ...
func (t *snaketext) Split(sep string) []string {
	return strings.Split(t.Input, sep)
}

// SplitPlace 根据字符串的位置进行分割
// Text("abcdefg").SpltPlace([]int{1,3,4})
// Out: []string{"a", "bc", "d", "efg"}
func (t *snaketext) SplitPlace(sep []int) []string {
	var a []string
	b := Text()
	for k, v := range []rune(t.Input) {
		b.Add(string(v))
		for _, i := range sep {
			if i == k+1 {
				a = append(a, b.Get())
				b = Text()
			}
		}

		if len(t.Input) == k+1 {
			a = append(a, b.Get())
		}
	}
	return a
}

// SplitInt 根据设置对字符串等分
// Text("abcdefg").SpltPlace([]int{1,3,4})
// Out: []string{"a", "bc", "d", "efg"}
func (t *snaketext) SplitInt(sep int) []string {
	var a []string
	b := Text()
	i := 0
	for _, v := range []rune(t.Get()) {
		b.Add(string(v))

		i = i + len(string(v))

		bl := len(b.Get())

		if bl >= sep || i == len(t.Get()) {
			a = append(a, b.Get())
			b = Text()
		}
	}

	return a
}

// Lines 根据行进行分割字符 ...
func (t *snaketext) Lines() []string {
	return strings.Split(strings.TrimSuffix(t.Input, "\n"), "\n")
}

// MD5 获取文件的MD5
func (t *snaketext) MD5() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(t.Get())))
}

// ---------------------------------------
// 辅助函数 :

// 根据规则字符进行分词 ...
func (t *snaketext) caseWords(isCamel bool, rule ...string) []string {
	src := t.Input
	if !isCamel {
		re := regexp.MustCompile("([a-z])([A-Z])")
		src = re.ReplaceAllString(src, "$1 $2")
	}
	src = strings.Join(strings.Fields(strings.TrimSpace(src)), " ")
	rule = append(rule, ".", " ", "_", " ", "-", " ")
	replacer := strings.NewReplacer(rule...)
	src = replacer.Replace(src)
	return strings.Fields(src)
}

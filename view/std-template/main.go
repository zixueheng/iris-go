package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
)

// golang 标准 模板包 测试

// Inventory ...
type Inventory struct {
	Material string
	Count    uint
}

func demo1() {
	sweaters := Inventory{"wool", 17}
	tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, sweaters) // 将 sweaters 渲染到 模板 tmpl 并输出到 Stdout
	// 17 items are made of wool
	if err != nil {
		panic(err)
	}
}

func demo2() {
	sweaters := Inventory{"wool", 17}
	// 如果“{{”紧接着跟随“—”和“ ”的话，那么“{{”之前的文本中的空白（空格、换行符、回车符、制表符）会被移除。对应的，“ -}}”表示移除之后文本中的空白
	tmpl, err := template.New("test").Parse("{{.Count -}} items are made of {{- .Material}}\n") // - 表示去除空格
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, sweaters)
	// 17items are made ofwool
	if err != nil {
		panic(err)
	}
}

// 自定义函数
func demo3() {
	// 自定义一个夸人的模板函数
	kua := func(arg string) (string, error) {
		return arg + "真帅", nil
	}
	// 采用链式操作在Parse之前调用Funcs添加自定义的kua函数
	tmpl, err := template.New("hello").Funcs(template.FuncMap{"kua": kua}).Parse(`{{kua .Name}}`)
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return
	}

	user := struct {
		Name   string
		Gender string
		Age    int
	}{
		Name:   "小明",
		Gender: "男",
		Age:    18,
	}
	// 使用user渲染模板，并将结果写入w
	tmpl.Execute(os.Stdout, user) // 小明真帅
}

func demo4() {
	// First we create a FuncMap with which to register the function.
	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"title": strings.Title,
	}
	// A simple template definition to test our function.
	// We print the input text several ways:
	// - the original: printf("%q", .)
	// - title-cased: title(.)
	// - title-cased and then printed with %q: printf("%q", title(.))
	// - printed with %q and then title-cased: title(printf("%q", .))
	// %q 表示给值加上双引号
	const templateText = `
Input: {{printf "%q" .}}
Output 0: {{title .}}
Output 1: {{title . | printf "%q"}}
Output 2: {{printf "%q" . | title}}
`
	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New("titleTest").Funcs(funcMap).Parse(templateText)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}
	// Run the template to verify the output.
	err = tmpl.Execute(os.Stdout, "the go programming language")
	if err != nil {
		log.Fatalf("execution: %s", err)
	}
	// Input: "the go programming language"
	// Output 0: The Go Programming Language
	// Output 1: "The Go Programming Language"
	// Output 2: "The Go Programming Language"
}

func demo5() {
	// Define a template.
	const letter = `
Dear {{.Name}},
{{if .Attended}}
It was a pleasure to see you at the wedding.{{else}}
It is a shame you couldn't make it to the wedding.{{end}}
{{with .Gift}}Thank you for the lovely {{.}}.
{{end}}
Best wishes,
Josie
`
	// Prepare some data to insert into the template.
	type Recipient struct {
		Name, Gift string
		Attended   bool
	}
	var recipients = []Recipient{
		{"Aunt Mildred", "bone china tea set", true},
		{"Uncle John", "moleskin pants", false},
		{"Cousin Rodney", "", false},
	}
	// Create a new template and parse the letter into it.
	t := template.Must(template.New("letter").Parse(letter))
	// Execute the template for each recipient.
	for _, r := range recipients {
		err := t.Execute(os.Stdout, r)
		if err != nil {
			log.Println("executing template:", err)
		}
	}
}

func demo6() {
	// const str = `{{eq 1 2}} {{if eq 1 2}}yes{{else}}no{{end}}` // 比较函数
	const str = `{{$username := .}}{{with $username}}My Name: {{$username}}{{end}}` // 变量使用 My Name: Jack
	tmpl, err := template.New("test").Parse(str)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}
	// Run the template to verify the output.
	err = tmpl.Execute(os.Stdout, "Jack")
	if err != nil {
		log.Fatalf("execution: %s", err)
	}
}

// map 字典值 使用
func demo7() {
	// 使用下标引用字典值 , 遍历字典
	const str = `{{.Index1}} {{.Index2}}
{{range $index, $element := .}}{{$index}}:{{$element}} {{end}}
`
	// const str = `{{range .}}{{.}} {{end}}`
	tmpl, err := template.New("test").Parse(str)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}
	// Run the template to verify the output.
	err = tmpl.Execute(os.Stdout, map[string]string{"Index1": "Value1", "Index2": "Value2"})
	if err != nil {
		log.Fatalf("execution: %s", err)
	}
	// Value1 Value2
}

// 内嵌模板
// 当解析模板时，可以定义另一个模板，该模板会和当前解析的模板相关联。模板必须定义在当前模板的最顶层，就像go程序里的全局变量一样。
// 这种定义模板的语法是将每一个子模板的声明放在"define"和"end" action内部。
func demo8() {
	// 定义两个模板T1和T2，第三个模板T3在执行时调用这两个模板；最后该模板调用了T3
	const str = `{{define "T1"}}ONE{{end}}
{{define "T2"}}TWO{{end}}
{{define "T3"}}{{template "T1"}} {{template "T2"}}{{end}}
{{template "T3"}}`
	tmpl, err := template.New("test").Parse(str)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}
	// Run the template to verify the output.
	err = tmpl.Execute(os.Stdout, map[string]string{"Index1": "Value1", "Index2": "Value2"})
	if err != nil {
		log.Fatalf("execution: %s", err)
	}
	// ONE TWO
}

// html
func demo9() {
	t, err := template.New("foo").Parse(`ttt {{define "T"}}Hello, {{.}}!{{end}}`)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}
	err = t.ExecuteTemplate(os.Stdout, "T", "<script>alert('you have been pwned')</script>") // 从 T 模板中 解析 HTML
	if err != nil {
		log.Fatalf("execution: %s", err)
	}
	// 生成的 是安全的 html
	// Hello, &lt;script&gt;alert(&#39;you have been pwned&#39;)&lt;/script&gt;!
	// 如果使用 text/template 包 则生成 未转义的 html ： Hello, <script>alert('you have been pwned')</script>!
}

// 注意 Parse() 中的文本不会被转义，转义的仅仅是传入的文本即 Execute()的data
func demo10() {
	t, err := template.New("foo").Parse(`<script>alert('1')</script> {{.}}`)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}
	err = t.Execute(os.Stdout, "<script>alert('you have been pwned')</script>") // 从 T 模板中 解析 HTML
	if err != nil {
		log.Fatalf("execution: %s", err)
	}
	// <script>alert('1')</script> &lt;script&gt;alert(&#39;you have been pwned&#39;)&lt;/script&gt;
}

// 非字符串内容可以 自动 序列化 并嵌入JavaScript里的
func demo11() {
	t, err := template.New("foo").Parse(`<script>var obj = {{.}}</script> {{.}}`)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}
	err = t.Execute(os.Stdout, struct{ A, B string }{"foo", "bar"}) //
	if err != nil {
		log.Fatalf("execution: %s", err)
	}
	// <script>var obj = {"A":"foo","B":"bar"}</script> {foo bar}
}

func main() {
	// demo1()
	// demo2()
	// demo3()
	// demo4()
	// demo5()
	// demo6()
	// demo7()
	// demo8()
	// demo9()
	// demo10()
	demo11()
}

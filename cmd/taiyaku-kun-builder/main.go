package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Interpretation struct {
	original    string
	translation string
	comment     string
	audio       string
}

var (
	config_file       = "config.yaml"
	docs_directory    = "docs"
	raw_file          = "raw.csv"
	interpretaions    []Interpretation
	top_page_template = `<!DOCTYPE html>
<html>

<head>
    <title>対訳君(両毛弁)</title>
    <meta charset="UTF-8">
</head>

<body>
    <h1>対訳君（両毛弁）</h1>
    このページはPoCです

    <p>ここで短く能書きをたれる</p>

    <a href="words/index.html">対訳リスト</a>

    <hr>
    <a href="https://github.com/yuchiki/taiyaku-kun">このページのソースコード</a>
</body>

</html>
`

	words_page_template = `<!DOCTYPE html>
<html>

<head>
	<title>対訳君(両毛弁)</title>
	<meta charset="UTF-8">
</head>

<body>
	<a href="/taiyaku-kun/index.html">トップへ</a>
	<hr>

	このページはPoCです

	<p>ここで短く能書きを垂れる</p>

	<table>
		<tr>
			<th>原文</th>
			<th>訳文</th>
			<th>音声</th>
		</tr>
		<tr>
			<td>
				原文1
			</td>
			<td>
				<a href="1/index.html">訳文１</a>
			</td>
			<td>
				<a href="/taiyaku-kun/sounds/8.mp3">音声</a>
			</td>
		</tr>
		<tr>
			<td>
				原文2
			</td>
			<td>
				<a href=" 2/index.html">訳文2</a>
			</td>
			<td>
				未収録
			</td>
		</tr>
		<tr>
			<td>
				原文3
			</td>
			<td>
				<a href="3/index.html">訳文3</a>
			</td>
			<td>
				未収録
			</td>
		</tr>
	</table>
</body>

</html>
`

	word_page_template = `<!DOCTYPE html>
<html>

<head>
    <title>対訳君(両毛弁)</title>
    <meta charset="UTF-8">
</head>

<body>
    <a href="/taiyaku-kun/words/index.html">対訳リスト</a>
    <a href="/taiyaku-kun/words/2/index.html">次の対訳</a>
    <hr>

    このページはPoCです


    <h1>原文</h1>
    <p>原文1</p>

    <h1>対訳</h1>
    <p>対訳1</p>

    <h1>解説</h1>
    <p>解説1</p>

    <h1>音声</h1>
    <a href="/taiyaku-kun/sounds/8.mp3">音声</a>
</body>

</html>
`
)

func main() {
	// configを読む

	// raw fileを読む

	// トップページ生成
	genTopPage()

	// words 生成
	genWordsPages()
}

func genTopPage() {
	err := ioutil.WriteFile(filepath.Join(docs_directory, "index.html"), []byte(top_page_template), 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func genWordsPages() {
	genWordsListPage()

	for i := 0; i < len(interpretaions); i++ {
		genWordPage(i, i == 0, i == len(interpretaions)-1)
	}
}

func genWordsListPage() {
	err := ioutil.WriteFile(filepath.Join(docs_directory, "words", "index.html"), []byte(words_page_template), 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func genWordPage(index int, isFirst bool, isLast bool) {}

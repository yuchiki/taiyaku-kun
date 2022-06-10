package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	language string
	author   string
}

type TranslationData struct {
	original    string
	translation string
	comment     string
	audio       string
}

var (
	config_file    = "config.yaml"
	docs_directory = "docs"
	config         = Config{
		language: "両毛弁",
		author:   "ゆーちき(@yuchiki1000yen)",
	}
	raw_file         = "raw.csv"
	translationDatas = []TranslationData{
		TranslationData{
			original:    "原文1",
			translation: "訳文1",
			comment:     "解説1",
			audio:       "8",
		},
		TranslationData{
			original:    "原文2",
			translation: "訳文2",
			comment:     "解説2",
			audio:       "",
		},
		TranslationData{
			original:    "原文3",
			translation: "訳文3",
			comment:     "解説3",
			audio:       "",
		},
	}
	top_page_template = `<!DOCTYPE html>
<html>

<head>
    <title>対訳君(%s)</title>
    <meta charset="UTF-8">
</head>

<body>
    <h1>対訳君（%s）</h1>
	<p>編集者: %s</p> <br>
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
	<title>対訳君(%s)</title>
	<meta charset="UTF-8">
</head>

<body>
	<a href="../index.html">トップへ</a>
	<hr>

	このページはPoCです

	<p>ここで短く能書きを垂れる</p>
%s

</body>

</html>
`

	word_page_template = `<!DOCTYPE html>
<html>

<head>
    <title>対訳君(%s)</title>
    <meta charset="UTF-8">
</head>

<body>
%s
    <hr>

    このページはPoCです


    <h1>原文</h1>
    <p>%s</p>

    <h1>対訳</h1>
    <p>%s</p>

    <h1>解説</h1>
    <p>%s</p>

    <h1>音声</h1>
    %s
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
	genWordsPages(translationDatas)
}

func genTopPage() {
	html := fmt.Sprintf(top_page_template, config.language, config.language, config.author)

	err := ioutil.WriteFile(filepath.Join(docs_directory, "index.html"), []byte(html), 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func genWordsPages(TranslationDatas []TranslationData) {
	genWordsListPage()

	for i := 0; i < len(translationDatas); i++ {
		genWordPage(translationDatas[i], i, i == 0, i == len(translationDatas)-1)
	}
}

func genWordsListPage() {

	table_template := `<table>
	<tr>
		<th>原文</th>
		<th>訳文</th>
		<th>音声</th>
	</tr>
%s
</table>`

	entry_template := `<tr>
	<td>
		%s
	</td>
	<td>
		<a href="%d/index.html">%s</a>
	</td>
	<td>
		%s
	</td>
</tr>`

	entries := []string{}

	for i := 0; i < len(translationDatas); i++ {
		translationData := translationDatas[i]

		var audio_link string
		if translationData.audio == "" {
			audio_link = "未収録"
		} else {
			audio_link = fmt.Sprintf(`<a href="../sounds/%s.mp3">音声</a>`, translationData.audio)
		}

		entries = append(
			entries,
			fmt.Sprintf(
				entry_template,
				translationData.original,
				i,
				translationData.translation,
				audio_link))
	}

	table := fmt.Sprintf(table_template, strings.Join(entries, "\n"))

	html := fmt.Sprintf(words_page_template, config.language, table)

	err := ioutil.WriteFile(filepath.Join(docs_directory, "words", "index.html"), []byte(html), 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func genWordPage(translationData TranslationData, index int, isFirst bool, isLast bool) {
	var link_to_before string
	var link_to_next string
	link_to_upper := `    <a href="../index.html">対訳リスト</a>`

	if isFirst {
		link_to_before = ""
	} else {
		link_to_before = fmt.Sprintf(`    <a href="../%d/index.html">前の対訳</a>`, index-1)
	}

	if isLast {
		link_to_next = ""
	} else {
		link_to_next = fmt.Sprintf(`    <a href="../%d/index.html">次の対訳</a>`, index+1)
	}

	header := strings.Join([]string{link_to_before, link_to_upper, link_to_next}, "\n")

	var audio_link string
	if translationData.audio == "" {
		audio_link = ""
	} else {
		audio_link = fmt.Sprintf(`<a href="../../sounds/%s.mp3">音声</a>`, translationData.audio)
	}

	page_html := fmt.Sprintf(
		word_page_template,
		config.language,
		header,
		translationData.original,
		translationData.translation,
		translationData.comment,
		audio_link)

	err := ioutil.WriteFile((filepath.Join(docs_directory, "words", strconv.Itoa(index), "index.html")), []byte(page_html), 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

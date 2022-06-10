package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Language string
	Author   string
}

type TranslationData struct {
	original    string
	translation string
	comment     string
	audio       string
}

const TranslationsDirectory = "translations"

var (
	config_file       = "config.yaml"
	docs_directory    = "docs"
	config            Config
	raw_file          = "raw.csv"
	top_page_template = `<!DOCTYPE html>
<html>

<head>
    <title>対訳君(%s)</title>
    <meta charset="UTF-8">
</heFad>

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
	readConfig(config_file)

	translationDatas := readTranslationDatas(raw_file)

	cleanUpDocsExceptForSounds()

	genTopPage(translationDatas)

	genWordsPages(translationDatas)
}

func cleanUpDocsExceptForSounds() {
	err := os.RemoveAll(path.Join(docs_directory, "index.html"))
	if err != nil {
		log.Fatal(err)
	}

	err = os.RemoveAll(path.Join(docs_directory, TranslationsDirectory))
	if err != nil {
		log.Fatal(err)
	}
}

func readConfig(filepath string) {

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal()
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

}

func readTranslationDatas(filepath string) []TranslationData {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	records = records[1:]

	translationDatas := []TranslationData{}

	for _, record := range records {
		if len(record) == 1 {
			continue
		}

		record = append(record, "", "") // 後半の要素は省略されうるので、空白フィールドを継ぎ足す
		translationDatas = append(translationDatas, TranslationData{
			original:    strings.Trim(record[0], " "),
			translation: strings.Trim(record[1], " "),
			comment:     strings.Trim(record[2], " "),
			audio:       strings.Trim(record[3], " "),
		})
	}

	return translationDatas
}

func genTopPage(translationDatas []TranslationData) {
	html := fmt.Sprintf(top_page_template, config.Language, config.Language, config.Author)

	err := ioutil.WriteFile(filepath.Join(docs_directory, "index.html"), []byte(html), 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func genWordsPages(translationDatas []TranslationData) {
	err := os.Mkdir(path.Join(docs_directory, TranslationsDirectory), 0777)
	if err != nil {
		log.Fatal(err)
	}

	genWordsListPage(translationDatas)

	for i := 0; i < len(translationDatas); i++ {
		genWordPage(translationDatas[i], i, i == 0, i == len(translationDatas)-1)
	}
}

func genWordsListPage(translationDatas []TranslationData) {

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

	html := fmt.Sprintf(words_page_template, config.Language, table)

	err := ioutil.WriteFile(filepath.Join(docs_directory, TranslationsDirectory, "index.html"), []byte(html), 0666)
	if err != nil {
		log.Fatal(err)
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
		config.Language,
		header,
		translationData.original,
		translationData.translation,
		translationData.comment,
		audio_link)

	err := os.Mkdir(path.Join(docs_directory, TranslationsDirectory, strconv.Itoa(index)), 0777)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile((filepath.Join(docs_directory, TranslationsDirectory, strconv.Itoa(index), "index.html")), []byte(page_html), 0666)
	if err != nil {
		log.Fatal(err)
	}
}

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
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Language          string
	Author            string
	Github            string
	Introduction      string
	WordsIntroduction string
}

type TranslationData struct {
	original    string
	translation string
	comment     string
	audio       string
}

const (
	TranslationsDirectory = "translations"
	RecorderDirectory     = "recorder"
)

var (
	build_timestamp   string
	config_file       = "config.yaml"
	docs_directory    = "docs"
	config            Config
	raw_file          = "raw.csv"
	top_page_template = `<!DOCTYPE html>
<html>

<head>
    <title>対訳君(%s)</title>
    <meta charset="UTF-8">
</head>

<body>
    <h1>対訳君（%s）</h1>
	<p>編集者: %s</p> <br>

	%s

    <a href="%s/index.html">対訳リスト</a>

    <hr>
    <a href="%s">このページのソースコード</a>
	<br>
	最終更新日時: %s
	<br>
	<a href="recorder/index.html"  >管理者録音用ページ</a>
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
%s
	<hr>
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
	build_timestamp = time.Now().Format(time.RFC822)

	readConfig(config_file)

	translationDatas := readTranslationDatas(raw_file)

	cleanUpDocsExceptForSounds()

	genTopPage(translationDatas)

	genWordsPages(translationDatas)

	genRecorderPages(translationDatas)
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
	r.FieldsPerRecord = -1
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	records = records[1:]

	translationDatas := []TranslationData{}

	for _, record := range records {
		record = append(record, "", "", "")
		if strings.Trim(record[1], " ") == "" {
			continue
		}

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
	html := fmt.Sprintf(
		top_page_template,
		config.Language,
		config.Language,
		config.Author,
		config.Introduction,
		TranslationsDirectory,
		config.Github,
		build_timestamp,
	)

	err := ioutil.WriteFile(filepath.Join(docs_directory, "index.html"), []byte(html), 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func genRecorderPages(translationDatas []TranslationData) {
	htmlTemplate1 := `
	<!--
	https://github.com/sozysozbot/recording-pekzep/blob/master/index.html
	を丸々コピーしたもの
	-->

	<!DOCTYPE html>
	<html lang="en-us">

	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width">
		<title>対訳君</title>

		<link href="normalize.css" rel="stylesheet" type="text/css">
		<link href="app.css" rel="stylesheet" type="text/css">

	</head>

	<body>
		<div class="wrapper">

			<header>
				<h1>録音ページ</h1>
			</header>

			<div style="font-size: 200%">
			<label for="which_sentence">録音する文を選択：</label>

				<select name="which_sentence" id="which_sentence">
`

	htmlTemplate2 :=
		`
				</select>
			</div>
			<section class="main-controls">
				<canvas class="visualizer"></canvas>
				<button class="record">Record</button>
				<button class="stop">Stop</button>

			</section>

			<section class="sound-clips">

				<!-- <article class="clip">
			<audio controls></audio>
			<a href="">Download clip</a>
		</article> -->

			</section>

		</div>

		<!-- Below is your custom application script -->

		<script src="app.js"></script>

	</body>

	</html>

`

	var options []string
	for _, translationData := range translationDatas {
		if translationData.translation != "" {
			options = append(options, fmt.Sprintf("<option>%s</option>", translationData.translation))
		}
	}

	html := fmt.Sprintf("%s\n%s\n%s", htmlTemplate1, strings.Join(options, "\n"), htmlTemplate2)

	err := ioutil.WriteFile(filepath.Join(docs_directory, RecorderDirectory, "index.html"), []byte(html), 0666)
	if err != nil {
		log.Fatal(err)
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
			soundsTemplate := `
<audio style="vertical-align: middle;" controls="">
	<source src="%s" type="audio/mpeg">
	Your browser does not support the audio element.
</audio>`
			audio_link = fmt.Sprintf(soundsTemplate, translationData.audio)
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

	html := fmt.Sprintf(words_page_template, config.Language, config.WordsIntroduction, table, build_timestamp)

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
		soundsTemplate := `
		<audio style="vertical-align: middle;" controls="">
			<source src="%s" type="audio/mpeg">
			Your browser does not support the audio element.
		</audio>`
		audio_link = fmt.Sprintf(soundsTemplate, translationData.audio)
	}

	page_html := fmt.Sprintf(
		word_page_template,
		config.Language,
		header,
		translationData.original,
		translationData.translation,
		translationData.comment,
		audio_link,
	)

	err := os.Mkdir(path.Join(docs_directory, TranslationsDirectory, strconv.Itoa(index)), 0777)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile((filepath.Join(docs_directory, TranslationsDirectory, strconv.Itoa(index), "index.html")), []byte(page_html), 0666)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	config_file    = "config.yaml"
	docs_directory = "docs"
	raw_file       = "raw.csv"
)

var top_page_template = `<!DOCTYPE html>
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

func main() {
	// configを読む

	// raw fileを読む

	// トップページ生成

	err := ioutil.WriteFile(filepath.Join(docs_directory, "index.html"), []byte(top_page_template), 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// words 生成
}

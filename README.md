# 対訳君

対訳を公開するページを自動生成します



# 自分の対訳作成プロジェクトを作成したい方へ

この リポジトリをforkして自分の手元へ持ってくることで、あなた自身が記述したい対訳集を生成することができます。
以下にその手順を記します。

## 初回設定

1. あなたの手元のパソコンで、gitコマンドとmakeコマンドが使えることようにします。
1. github アカウントを作成します。　TODO:　参考リンクを貼る
1. このリポジトリを自分の名前空間のもとにforkします。 TODO:参考リンクを貼る
1. forkしたリポジトリをローカルにcloneしてきます。
1. ローカルのリポジトリで、`make clean-install` します。　TODO:　未実装


##





## 人間がいじるべき元データ

- [raw.csv](./raw.csv)
- [config.yaml](./config.yaml)

以上を更新して、 `make build` を行い、git push を行うとdocs 以下が更新されます。

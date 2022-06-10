.PHONY: taiyaku-kun-builder build clean-install

taiyaku-kun-builder:
	go build -o taiyaku-kun-builder cmd/taiyaku-kun-builder/main.go

build: taiyaku-kun-builder
	./taiyaku-kun-builder

clean-install:
	@echo あなたの手元の今までの対訳データをすべて消去して、新たに対訳作りを1から始められる準備をします。
	@echo 消去した後でも、git commit/push したことのある対訳データは、コミット履歴から復元することができます。
	@read -p "本当に clean install しますか? (yes/no) :" str; if [ "$$str" != "yes" ]; then echo " clean install を中止します"; false; fi
	@echo

	cp initial_settings/config.yaml ./config.yaml
	cp initial_settings/raw.csv ./raw.csv
	rm -rf docs/sounds/*
	cp initial_settings/taiyaku1.mp3 docs/sounds/taiyaku1.mp3

	@echo
	@echo 初期化が完了しました。
	@echo 対訳君を使ってくれてありがとう! 一緒に対訳づくりに励みましょう！

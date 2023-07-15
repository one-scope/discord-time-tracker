# discord-time-tracker

## 概要
https://github.com/orgs/one-scope/discussions/4

## セットアップ

config.yml 

discord_bot_token以外は変更する必要なし
```
discordbot:
  discord_bot_token: {DiscordBotトークン} #現在はAdministratorの権限を持ったdiscordbotのトークン
  data_dir: /var/lib #データの保存先ディレクトリ. docker-composeでvolumeしている
  execution_timing: "* * * * *" # cron 形式のファイルに出力する頻度

log:
  file_path: /var/log/log.txt #ログの出力先ディレクトリ. docker-compose.ymlでvolumeしている
```

起動
```
docker compose up -d 
```


## 仕様
discussionで公開されていた要件から以下のような仕様にしました。

discordイベントが発生するたびにメモリにデータを保存する。

config.ymlで設定されたcron形式のexecution_timingの頻度でメモリのデータをファイルに出力する。

メンバーを記録するユーザーファイル(1個)、disocrdでのイベントを記録するステータスファイル(n個)の2種類のファイルを作成

ステータスファイルは日付ごとに作成し、1時間間隔でデータを集計し記録する。

以下のようなイメージ
```
20230623_status.json
    20230623-00:00
        user1
            id
            ...
        user2
    20230623-01:00
        user1
        ...
```

ユーザーファイル
```
ユーザーID
ユーザー名
ロール 
ユーザーが現在ギルドにいるか
```
ステータスファイル
```
時間帯 #1時間ごとに分類して集計
ユーザーID
オンライン/オフライン #1時間当たりのオンライン時間(秒)
ボイスチャンネル # 1時間当たりに滞在していた全てのチャンネルIDとそれぞれの滞在時間(秒)
マイクON/マイクOFF/スピーカーOFF # 
```


### ざっくりとした依存関係

大きい

main：appの初期化、Bot開始

app：ログの初期化、Botの初期化

discord > discord.go：リスナー

discord > datamanager.go：リスナーのデータをメモリに保存、ファイル出力用の形にメモリのデータを整形

discord > model.go：Botやファイル出力する際の構造体、Botのリスナー登録用の関数を定義

db：メモリのデータをファイルに出力するための関数

config：config.ymlファイルのデータから受け取るための構造体

小さい
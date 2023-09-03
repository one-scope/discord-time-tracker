# discord-time-tracker

## 概要
discordサーバーのメンバーの状態を保存し集計します。

## セットアップ

設定ファイル`.env`

example.envを参考に作成します。


起動
```
docker compose up -d 
```

## 使用方法

任意のTextチャンネルにコマンドを送信しするとBotが集計し同チャンネルに結果を送信します。

### コマンド
- 全てのユーザー情報
```
tracker,users
```
- 集計範囲を分単位で分割した全てのユーザーのステータス

    tracker,statuses,YearMonthDay,YearMonthDay,minute
```
tracker,statuses,20230101,20230131,1440
```
- 集計範囲を分単位で分割した指定したユーザーのステータス

    tracker,statuses,YearMonthDay,YearMonthDay,minute,userIDs...

    userIDはコンマ区切り
```
tracker,status,20230101,20230131,1440,141241241,345345345
```
- 集計範囲を分単位で分割した指定したロールを所持するユーザーのステータス

    tracker,statuses,YearMonthDay,YearMonthDay,minute,roleID

```
tracker,statusesbyrole,20230101,20230131,1440,1117731122005164063
```
- 全てのロール情報
```
tracker,roles
```
- ヘルプを表示
```
tracker,help
```
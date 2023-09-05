# discord-time-tracker

## 概要
discordサーバーのメンバーの状態を保存し集計します。

## セットアップ

### 設定ファイル`.env`

example.envを参考に作成します。

### Botの権限(作成中)

tSession.Identify.Intents = discordgo.IntentGuildMembers | discordgo.IntentGuildPresences | discordgo.IntentGuildVoiceStates | discordgo.IntentGuilds | discordgo.IntentGuildMessages

テキストチャンネルにメッセージを送信
GUILD_CREATE
// ボットが起動したとき。
aBot.onEvent(aBot.event)
// 誰かがサーバーに参加したとき。
aBot.onGuildMemberAdd(aBot.guildMemberAdd)
//誰かのロールが変わったとき。
aBot.onGuildMemberUpdate(aBot.guildMemberUpdate)
// 誰かがサーバーから退出したとき。
aBot.onGuildMemberRemove(aBot.guildMemberRemove)
//誰かのオンラインになったとき
aBot.onPresenceUpdate(aBot.presenceUpdate)
// 誰かの音声通話が更新されたとき。// 接続、切断もこれ。切断時は ChannelID が空文字。
aBot.onVoiceStateUpdate(aBot.voiceStateUpdate)
// 誰かがメッセージを送信したとき。
aBot.onMessageCreate(aBot.messageCreate)
//ロールが作成されたとき。
aBot.onGuildRoleCreate(aBot.roleCreate)
//ロールが更新されたとき。
aBot.onGuildRoleUpdate(aBot.roleUpdate)
//ロールが削除されたとき。
aBot.onGuildRoleDelete(aBot.roleDelete)


### 起動
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
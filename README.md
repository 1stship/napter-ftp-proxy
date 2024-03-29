# napter-ftp-proxy

napter-ftp-proxyはSORACOM Napterを利用したリモートデバイスへのFTP接続ツールです。

ソラコム社のオンデマンドリモートアクセスサービスSORACOM Napterを用いることにより、リモートデバイスへの安全なFTP接続を実現します。

SORACOM Napterの説明はこちら

https://soracom.jp/services/napter/

## 取得方法
go getコマンドで取得できます。
```sh
go get -u github.com/1stship/napter-ftp-proxy
```

また、リリースページから各環境のバイナリ実行ファイルがダウンロードできます。

## 使用方法
以下のコマンドを実行します。IMSIには接続したいデバイスで使用されているSORACOM Air SIMのIMSIを指定します。
```sh
napter-ftp-proxy --target IMSI
```

SORACOMのアカウントとパスワードを入力すると、ローカルホストにてFTPの待ち受け状態になります。

FTPクライアントツールを使用してローカルホストの21番ポートにFTP接続すると、透過的にリモートデバイスにFTP接続されます。

## コマンドラインオプション

|オプション  |説明  |
|---|---|
|--help  |ヘルプを表示します。  |
|--version  |バージョンを表示します。  |
|--target IMSI |接続先のデバイスのIMSIを指定します。  |
|--listen IP_ADDRESS |待ち受けするIPアドレスを指定します。自機以外からのアクセスを受け付ける際に使用します。  |
|--local PORT |待ち受けするポートを指定します。権限の関係などで21番で待ち受けできない場合に使用します。  |
|--remote PORT |接続先のポートを指定します。デバイス側のポートが21番でない時に使用します。  |

## 環境変数

|環境変数  |説明  |
|---|---|
|SORACOM_EMAIL  |ソラコムアカウントのメールアドレスを指定します。メールアドレス入力を省略できます。  |
|SORACOM_PASSWORD  |ソラコムアカウントのパスワードを指定します。パスワード入力を省略できます。  |

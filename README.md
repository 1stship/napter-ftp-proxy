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

## 使用方法
以下のコマンドを実行します。${IMSI}には接続したいデバイスで使用されているSORACOM Air SIMのIMSIを指定します。
```sh
napter-ftp-proxy --target ${IMSI}
```

SORACOMのアカウントとパスワードを入力すると、ローカルホストにてFTPの待ち受け状態になります。
アカウントとパスワードはそれぞれ環境変数SORACOM_EMAIL、SORACOM_PASSWORDでも指定できます。

PC側で入力したコマンドがデバイス側で実行され、コマンドの実行結果を表示します。

## 複数デバイス対応

デフォルト設定では、エンドポイント名：inventory-terminalのデバイスを生成し、そのデバイスに対しアクセスします。

```sh
inventory-terminal --mode daemon --endpoint <任意のエンドポイント名>
```

```sh
inventory-terminal --endpoint <任意のエンドポイント名>
```

とすることで、複数デバイスに対応できます。

## ネットワーク環境について

- デバイス側 : SORACOM Airネットワーク
- PC側 : 外向きのポートが制限されていないネットワーク(TURN非対応のため)

## TODO

- アクセスID、アクセスキー認証およびSAMユーザー認証の対応
- Goのパッケージ管理

---
id: 0000000021
tags: [Go,Memo,Module]
title: ForkされたレポジトリをGoパッケージとしてimportする方法
description: |-
  Github上に公開されている、ForkされたGoのパッケージのレポジトリをGoパッケージとしてimportする方法。
date: 2021-03-21
---

Github上に公開されている、forkされたGoのパッケージのレポジトリをGoパッケージとしてimportする方法。

`replace`ディレクティヴを使用する。

方法を例を用いて説明する。

次のような`go.mod`があるとする。

```
module github.com/suzuito/fuga-go

go 1.15

require (
	firebase.google.com/go/v4 v4.2.0
)
```

Goのパッケージ`firebase.google.com/go/v4`をforkしたレポジトリが`github.com/maku693/firebase-admin-go`にある。
今回は`github.com/maku693/firebase-admin-go`の`auth-emulator-token-verification`ブランチをGoパッケージとしてimportし、使用したい。

まず、forkしたレポジトリ`github.com/maku693/firebase-admin-go`をcloneし、ローカルへダウンロードする。

```
git clone https://github.com/maku693/firebase-admin-go.git
```

cloneしてダウンロードされたレポジトリのパスを確認する。

```bash
% pwd firebase-admin-go
/Users/taro/firebase-admin-go
```

`replace`ディレクティヴを`go.mod`の中に追加する。

```
module github.com/suzuito/fuga-go

go 1.15

require (
	firebase.google.com/go/v4 v4.2.0
)

replace firebase.google.com/go/v4 => /Users/taro/firebase-admin-go
```

---
id: 0000000018
tags: [Angular,monorepo]
title: Angularでmonorepo
description: |-
  現在は、残念ながら「実存に対する疑問」を「不安神経症」であると言わざるを得ない状況だ。
date: 2021-03-03
---

# Angular でモノレポ

Angularにてモノレポするための方法を書きます。
モノレポとは、ソースコードを1つのレポジトリで管理する手法だ。

Angularはフロントエンドのソースコードを綺麗に構造化できるフレームワークだ。
俗人的な書き方を防止するようなあらゆる仕組みが用意されている。
ソースコードの構造についても、Angularについてある程度の知識がある人であれば、誰が書いても同じような構造となる。
ビューは`component`という単位で、複雑なロジックは`service`という単位で、`component`と`service`を組み合わせて実装された小さなアプリケーションを`module`という独立したライブラリの形にすることもできる。
いかにも`Google`が作りそうなフレームワークという感じだ。

それ故に、僕はしばしば、再利用可能なソースコードの一部をmoduleとして抽出し、別のAngularアプリケーションでも使用できるようにと、余計なことをしたくなってしまう。開発者のさがというやつだ。

これをやるとき、再利用可能なソースコードだけが含まれるモジュールを、独立したレポジトリとして、管理したくなる。

しかし、ライブラリを独立したレポジトリで管理するようなワークフローは、効率が悪い。
ライブラリに変更を加えた場合、ライブラリのソースコードをリモートレポジトリにコミットし、コミットされたソースコードをAngularアプリケーション側でフェッチしなければならないからだ。
もちろん、それを回避する方法はあるのだろうが、回避するためにはひと工夫が必要だ。
その工夫のために頭を使うのが面倒臭い。

このような理由から、1つのレポジトリで全てのソースコードを管理したくなる。

## スコープ

この記事では、Angularで書かれたフロントエンドのソースコードをモノレポで管理するまで。
サーバーサイドとフロントエンドのソースコードをモノレポで管理するということはやらない。

また、Angular CLIは次のバージョンのものを用いた。

```bash
% ng version

     _                      _                 ____ _     ___
    / \   _ __   __ _ _   _| | __ _ _ __     / ___| |   |_ _|
   / △ \ | '_ \ / _` | | | | |/ _` | '__|   | |   | |    | |
  / ___ \| | | | (_| | |_| | | (_| | |      | |___| |___ | |
 /_/   \_\_| |_|\__, |\__,_|_|\__,_|_|       \____|_____|___|
                |___/
    

Angular CLI: 11.2.2
Node: 12.12.0
OS: darwin x64

Angular: 
... 
Ivy Workspace: 

Package                      Version
------------------------------------------------------
@angular-devkit/architect    0.1100.1
@angular-devkit/core         11.0.1
@angular-devkit/schematics   11.0.1
@schematics/angular          11.0.1
@schematics/update           0.1100.1
```

## 仮想プロジェクト

本記事では、仮想的なプロジェクトを作りながら、そのサンプルプロジェクトをソースコードをモノレポで管理する方法を説明する。

サンプルプロジェクトは、以下のようなものを想定する。

- site1 `アプリケーション1`
- site2 `アプリケーション2`
- lib1 `site1,site2の両方にて使用される共通ソースコード`

本記事の目的は、`site1`,`site2`,`lib1`のソースコードを1つのレポジトリで管理する方法を示すことだ。

## 作り方

### レポジトリを作る

まずはレポジトリを作る。
今後、このレポジトリで全てのソースコードを管理していく。
レポジトリの名前は`monorepo-angular`だ。

以下のコマンドを実行し、レポジトリを作る。

```bash
ng new monorepo-angular --create-application=false
```

コマンド実行後、`monorepo-angular`ディレクトリ配下にファイルが作成される。

```bash
% tree -I node_modules ./monorepo-angular
./monorepo-angular
├── README.md
├── angular.json
├── package-lock.json
├── package.json
├── tsconfig.json
└── tslint.json

0 directories, 6 files
```

生成された各ファイルの役割は後ほど説明する。

ご覧の通り、このレポジトリはアプリケーション部分を構成するソースコードは含まれていない。
この後、ライブラリ`lib1`や、アプリケーション`site1`、`site2`を作っていく。
これらのソースコードは全て`projects`ディレクトリ配下に置かれる。

### ライブラリ`lib1`を作る

以下のコマンドを実行する。

```bash
cd monorepo-angular
ng generate library lib1
```

コマンド実行後、`monorepo-angular`ディレクトリ配下にファイルが作成される。

```bash
% tree -I node_modules . 
.
├── README.md
├── angular.json
├── package-lock.json
├── package.json
├── projects
│   └── lib1
│       ├── README.md
│       ├── karma.conf.js
│       ├── ng-package.json
│       ├── package.json
│       ├── src
│       │   ├── lib
│       │   │   ├── lib1.component.spec.ts
│       │   │   ├── lib1.component.ts
│       │   │   ├── lib1.module.ts
│       │   │   ├── lib1.service.spec.ts
│       │   │   └── lib1.service.ts
│       │   ├── public-api.ts
│       │   └── test.ts
│       ├── tsconfig.lib.json
│       ├── tsconfig.lib.prod.json
│       ├── tsconfig.spec.json
│       └── tslint.json
├── tsconfig.json
└── tslint.json

4 directories, 21 files
```

### アプリケーション`site1`と`site2`を作る

以下のコマンドを実行する。

```bash
ng generate application site1
ng generate application site2
```

コマンド実行後、`monorepo-angular`ディレクトリ配下にファイルが作成される。

```bash
% tree -I node_modules .
.
├── README.md
├── angular.json
├── package-lock.json
├── package.json
├── projects
│   ├── lib1
...省略...
│   ├── site1
...省略...
│   └── site2
...省略...
├── tsconfig.json
└── tslint.json

18 directories, 65 files
```

### serve方法（ローカルでの開発）

では、ローカルで開発するフローを説明する。

#### ライブラリ側のソースコード

ライブラリの`Lib1Service`に`hello`関数を作る。

```typescript
import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class Lib1Service {

  constructor() { }

  hello(): string {
    return 'Hello world!';
  }
}
```

#### アプリケーション側のソースコード

この関数を`site1`,`site2`の中で呼んでみる。
`site1`の`app.module.ts`の中でライブラリの`Lib1Service`を`providers`へ追加する。

```typescript
import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { Lib1Service } from 'lib1';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule
  ],
  providers: [
    Lib1Service,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
```

`site1`の`application/app.component.ts`

```typescript
import { Component } from '@angular/core';
import { Lib1Service } from 'lib1';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'site1';
  constructor(public lib1Service: Lib1Service) {}
}
```

`site1`の`application/app.component.html`

```html
<html>
  <head></head>
  <body>
    <div>{{title}} using {{lib1Service.hello()}}</div>
  </body>
</html>
```

#### ライブラリのビルドと`--watch`

以下のコマンドで、ライブラリをアプリケーションで利用する方法。

```bash
# lib1のビルド
ng build lib1 --watch
```

ビルド後のコードは`dist/lib1`ディレクトリに出力される。

#### アプリケーションのserve

以下のコマンドを用いれば、アプリケーションをserveできる。

site1をポート4201、site2をポート4202でserveする。

```bash
# site1の起動
ng serve site1 --port=4201
# site2の起動
ng serve site2 --port=4202
```

### アプリケーションのbuildとdeploy

#### build

以下のコマンドを用いれば、アプリケーションをbuildできる。

```bash
ng build lib1 --prod
```

ビルド後のコードは`dist/lib1`に出力される。

```bash
ng build site1 --prod
ng build site2 --prod
```

ビルド後のコードは`dist/site1`,`dist/site2`に出力される。

#### deploy

あとは、生成されたコード`dist/site1`,`dist/site2`をデプロイするだけ。

ライブラリ`dist/lib1`のコードはデプロイしなくても良い。
必要であればライブラリを、npmパッケージとして、npmレジストリにpublishすることができる。
本記事では、ライブラリの生成方法については省略する。
詳しくは[こちら](https://angular.io/guide/using-libraries)。

この記事で作成した`site1`をGitHub pagesにデプロイする方法は以下。

```bash
ng build site1 --prod --base-href=/monorepo-angular --deploy-url=/monorepo-angular/
git add docs
git push origin main
```

- https://suzuito.github.io/monorepo-angular/

## 参考

- https://angular.io/guide/file-structure#multiple-projects
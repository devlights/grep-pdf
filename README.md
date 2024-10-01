# これは何？

指定ディレクトリの下に存在する PDF ファイル (*.pdf) をGREPコマンドのように検索します。

対象のPDFファイルが存在していて、どこに書いてあるのかが分からない時などに良かったらご利用ください。

## インストール

```sh
go install github.com/devlights/grep-pdf/cmd/grep-pdf@latest
```

## 使い方

```sh
$ ./grep-pdf.exe -h
Usage of ./grep-pdf.exe:
  -debug
        debug mode
  -dir string
        directory (default ".")
  -json
        output as JSON
  -only-hit
        show ONLY HIT (default true)
  -text string
        search text
  -verbose
        verbose mode
```

ヒットした文書のパスが知りたい場合は以下のようにします。

```sh
$ ./grep-pdf.exe -dir ~/path/to/documents -text "データベースサイズ"
test.pdf: HIT
```

ヒットした箇所も見たい場合は ```-verbose``` オプションを付与するとみることが出来ます。

```sh
$ ./grep-pdf.exe -dir ~/path/to/documents -text "データベースサイズ" -verbose  
```

結果をjsonで出力したい場合は ```-json``` オプションを付与します。

```sh
$ ./grep-pdf.exe -dir ~/path/to/documents -text "データベースサイズ" -verbose  -json
```

## ビルド方法

[Task](https://taskfile.dev/#/) を使っています。詳細は [Taskfile.yml](./Taskfile.yml) を参照ください。

```sh
$ task build
```

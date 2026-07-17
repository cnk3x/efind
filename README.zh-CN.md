# efind

[English](README.md)

一个用 Go 编写的命令行小工具，用于在 `PATH` 中查找可执行文件 —— 可以看作是 `which` / `where` 的增强版。

`efind` 能够定位命令的绝对路径，可选地解析符号链接，并原生支持由 [Scoop](https://scoop.sh/) 和 [mise](https://mise.jdx.dev/) 管理的工具。

## 功能特性

- 按名称查找一个或多个可执行文件（类似 `which`）。
- 通过 `--eval` 选项解析符号链接 / shim 指向的真实文件。
- 无缝支持通过 **Scoop**（`scoop which`）和 **mise**（`mise which`）安装的可执行文件。
- 多种输出格式：对齐的可读文本表格、纯路径，或结构化的 **JSON**。
- 支持在单次调用中查询多个名称。

## 安装

需要 Go 1.26 及以上版本。

```bash
go install github.com/cnk3x/efind@latest
```

或从源码构建：

```bash
git clone https://github.com/cnk3x/efind
cd efind
go build -o efind .
```

## 使用说明

```
efind [参数] <名称> [<名称> ...]
```

### 参数

| 参数 | 简写 | 说明 |
|------|------|------|
| `--eval` | `-e` | 解析符号链接 / shim 指向，输出真实路径。 |
| `--json` | `-j` | 以 JSON 格式输出结果。 |
| `--noe` |      | 抑制未找到可执行文件时的错误信息。 |
| `--verbose` | `-v` | 即使只查询单个名称也强制使用详细（对齐）输出。 |

如果不提供任何名称，将打印使用说明。

### 示例

查找单个命令：

```bash
$ efind go
C:\Users\wen\sdk\go\bin\go.exe
```

解析 shim / 符号链接背后的真实路径：

```bash
$ efind -e node
node => C:\Users\wen\bin\mise\bin\node => C:\Users\wen\.local\share\mise\installs\node\21.0.0\bin\node.exe
```

一次性查询多个命令（多个名称时会自动启用对齐的详细输出）：

```bash
$ efind git python rustc
git    => C:\Program Files\Git\cmd\git.exe
python => C:\Users\wen\AppData\Local\Programs\Python\python.exe
rustc  => executable file not found
```

忽略未找到的错误：

```bash
$ efind --noe rustc
```

获取机器可读的输出：

```bash
$ efind -j git
[
  {
    "name": "git",
    "full": "C:\\Program Files\\Git\\cmd\\git.exe"
  }
]
```

## 工作原理

1. `efind` 使用 `exec.LookPath` 在 `PATH` 中查找可执行文件。
2. 当指定 `--eval` 时，按以下顺序尝试解析真实位置：
   - 通过 `scoop which <名称>` 解析 **Scoop** shim
   - 通过 `mise which <名称>` 解析 **mise** shim
   - 否则回退到 `filepath.EvalSymlinks`

## 开源协议

[MIT](LICENSE)

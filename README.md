
# git-cmd

一个简洁的 Git 辅助命令行工具，将常用的 Git 操作封装为简单的子命令，提升日常开发效率。

## 安装

确保 `$GOPATH/bin` 或 `$HOME/go/bin` 已加入 `$PATH` 环境变量，然后执行：

---
```bash
go install github.com/aid297/git-cmd@latest
```
---
## 快速开始

---
```bash
# 查看帮助
git-cmd

# 查看具体命令帮助
git-cmd help push
```
---
## 命令一览

| 命令       | 说明               |
| ---------- | ------------------ |
| `push`     | 一键提交并推送代码 |
| `rebase`   | rebase             |
| `merge`    | 合并分支并推送     |
| `tag-list` | 查看标签列表       |
| `tag-last` | 查看最新标签       |
| `tag-push` | 创建并推送标签     |
| `help`     | 显示帮助信息       |

所有命令均支持 `-v` 参数，开启后会显示每一步执行的具体命令、输出和结果。

## 命令详解

### rebase
一键完成 **(**`git add --all` → `git commit`  **|** `git stash`**)** → `git pull` → `git fetch origin` → `git rebase origin/目标分支` → `git push --force-with-lease` → **[**`git stash pop`**]** 。

---
```bash
git-cmd rebase --branch <目标分支> [stash] [--force | --force-with-lease] [-v]
```
---

| 参数                        | 说明             | 必填 |
| --------------------------- | ---------------- | ---- |
| --branch                    | 目标分支名称     | 是   |
| --stash                     | 使用 stash       | 否   |
| --force\|--force-with-lease | 合并后推送参数   | 否   |
| -v                          | 显示详细执行过程 | 否   |

**示例：**

---
```bash
git-cmd rebase --branch main
git-cmd rebase --branch main --stash
git-cmd rebase --branch main --force
git-cmd rebase --branch main --force-with-lease -v
```
---

### push

一键完成 `git add --all` → `git commit` → `git pull` → `git push`。

---
```bash
git-cmd push --comment <提交信息> [--branch <分支名>] [-v]
```
---
| 参数        | 说明                         | 必填 |
| ----------- | ---------------------------- | ---- |
| `--comment` | 提交注释                     | 是   |
| `--branch`  | 推送到的远程分支（默认当前分支） | 否   |
| `-v`        | 显示详细执行过程             | 否   |

**示例：**

---
```bash
git-cmd push --comment "修复登录bug"
git-cmd push --comment "新增功能" --branch dev
git-cmd push --comment "优化性能" -v
```
---
### merge

将源分支合并到目标分支并推送，完成后自动切回源分支。

---
```bash
git-cmd merge --branch "<源分支> -> <目标分支>" [-v]
```
---
| 参数       | 说明                             | 必填 |
| ---------- | -------------------------------- | ---- |
| `--branch` | 合并方向，格式: `src -> dst`     | 是   |
| `-v`       | 显示详细执行过程                 | 否   |

**示例：**

---
```bash
git-cmd merge --branch "feature -> main"
git-cmd merge --branch "dev -> release" -v
```
---
### tag-list

查看最近的标签列表，按版本号降序排列。

---
```bash
git-cmd tag-list [--count <数量>] [-v]
```
---
| 参数      | 说明             | 必填 | 默认值 |
| --------- | ---------------- | ---- | ------ |
| `--count` | 显示标签数量     | 否   | 1      |
| `-v`      | 显示详细执行过程 | 否   | -      |

**示例：**

---
```bash
git-cmd tag-list
git-cmd tag-list --count 3
git-cmd tag-list --count 5 -v
```
---
### tag-last

查看仓库最新的标签。

---
```bash
git-cmd tag-last [-v]
```
---
| 参数 | 说明             | 必填 |
| ---- | ---------------- | ---- |
| `-v` | 显示详细执行过程 | 否   |

**示例：**

---
```bash
git-cmd tag-last
git-cmd tag-last -v
```
---
### tag-push

创建标签并推送到远程仓库。

---
```bash
git-cmd tag-push --tag <标签名> [-v]
```
---
| 参数   | 说明             | 必填 |
| ------ | ---------------- | ---- |
| `--tag`| 标签名           | 是   |
| `-v`   | 显示详细执行过程 | 否   |

**示例：**

---
```bash
git-cmd tag-push --tag v1.0.0
git-cmd tag-push --tag v2.3.1 -v
```
---
## verbose 模式 (`-v`)

所有写操作命令均支持 `-v` 参数，开启后会逐步打印：

---
```
$ git-cmd tag-list --count 2 -v
[执行] git tag --sort=-v:refname
[输出] v0.0.5
v0.0.4
v0.0.3
[结果] 成功
最近 2 个标签:
  v0.0.5
  v0.0.4
```
---
## 许可证

MIT

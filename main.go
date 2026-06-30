package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func printUsage() {
	fmt.Println(`git-cmd - Git 辅助命令行工具

用法:
  git-cmd <command> [options]

可用命令:
  push        一键提交并推送代码
  rebase      将当前分支 rebase 到目标分支
  merge       合并分支并推送
  tag-last    查看最新标签
  tag-push    创建并推送标签
  help        显示帮助信息

全局选项:
  -v          显示详细执行过程

使用 git-cmd help <command> 查看具体命令的用法`)
}

func printPushHelp() {
	fmt.Println(`push - 一键提交并推送代码

用法:
  git-cmd push --comment <提交信息> [--branch <分支名>] [-v]

选项:
  --comment   提交注释（必填）
  --branch    推送到的远程分支名（可选，默认推送当前分支）
  --tag       创建并推送标签（可选）
  -v          显示详细执行过程

示例:
  git-cmd push --comment "修复bug"
  git-cmd push --comment "新功能" --branch dev
  git-cmd push --comment "优化" --tag v1.0.0
  git-cmd push --comment "优化" -v`)
}

func printMergeHelp() {
	fmt.Println(`merge - 合并分支并推送

用法:
  git-cmd merge --branch [<源分支>|不写默认当前分支合并] -> <目标分支> [-v]

选项:
  --branch   分支合并方向，格式: src -> dst（必填）
  -v         显示详细执行过程

示例:
  git-cmd merge --branch "feature -> main"
  git-cmd merge --branch "feature -> main" -v`)
}

func printTagPushHelp() {
	fmt.Println(`tag-push - 创建并推送标签

用法:
  git-cmd tag-push --tag <标签名> [-v]

选项:
  --tag   标签名（必填）
  -v      显示详细执行过程

示例:
  git-cmd tag-push --tag v1.0.0
  git-cmd tag-push --tag v1.0.0 -v`)
}

func printTagListHelp() {
	fmt.Println(`tag-list - 查看标签列表

用法:
  git-cmd tag-list [--count <数量>] [-v]

选项:
  --count   显示标签数量（默认 1）
  -v        显示详细执行过程

示例:
  git-cmd tag-list
  git-cmd tag-list --count 3
  git-cmd tag-list --count 5 -v`)
}

func printTagLastHelp() {
	fmt.Println(`tag-last - 查看最新标签

用法:
  git-cmd tag-last [-v]

选项:
  -v   显示详细执行过程

示例:
  git-cmd tag-last
  git-cmd tag-last -v`)
}

func printRebaseHelp() {
	fmt.Println(`rebase - 将当前分支 rebase 到目标分支

用法:
  git-cmd rebase --branch <目标分支> [--comment <提交注释>] [--stash] [--force-with-lease | --force] [-v]

选项:
  --branch   要 rebase 的目标分支（必填）
  --comment  合并时提交注释
  --stash    是否 stash 变更
  --force-with-lease   强制推送 lease
  --force              强制推送
  -v         显示详细执行过程

流程:
  git fetch → git rebase origin/<目标分支> → git push --force-with-lease

冲突处理:
  遇到冲突时，手动解决后执行:
    git add .
    git rebase --continue
    git push --force-with-lease

示例:
  git-cmd rebase --branch main
  git-cmd rebase --branch main -v`)
}

// verbose 控制是否显示详细执行过程
var verbose bool

// runCmd 封装命令执行，verbose 模式下打印执行的命令和输出
func runCmd(name string, args ...string) ([]byte, error) {
	cmdStr := name + " " + strings.Join(args, " ")
	if verbose {
		fmt.Printf("[执行] %s\n", cmdStr)
	}

	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()

	if verbose {
		output := strings.TrimSpace(string(out))
		if output != "" {
			fmt.Printf("[输出] %s\n", output)
		}
		if err != nil {
			fmt.Printf("[错误] %v\n", err)
		} else {
			fmt.Println("[结果] 成功")
		}
	}

	return out, err
}

// gitPull 执行 git pull 并处理冲突
// pullArgs: git pull 的参数（如 "origin", "main"）
// conflictHint: 遇到冲突时的解决提示
func gitPull(pullArgs []string, conflictHint string) {
	args := append([]string{"pull"}, pullArgs...)
	out, err := runCmd("git", args...)
	if err != nil {
		if strings.Contains(string(out), "CONFLICT") {
			fmt.Println("git pull 遇到冲突，请手动解决冲突后执行:")
			fmt.Println(conflictHint)
			os.Exit(1)
		}
		log.Fatalf("执行 git pull 失败：%v", err)
	}
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "push":
		pushCmd := flag.NewFlagSet("push", flag.ExitOnError)
		pushBranch := pushCmd.String("branch", "", "分支")
		pushComment := pushCmd.String("comment", "", "合并注释")
		pushTag := pushCmd.String("tag", "", "标签")
		pushVerbose := pushCmd.Bool("v", false, "显示详细执行过程")
		pushCmd.Parse(os.Args[2:])
		verbose = *pushVerbose

		if *pushComment == "" {
			log.Fatalln("提交注释不能为空")
		}

		if _, err := runCmd("git", "add", "--all"); err != nil {
			log.Fatalf("执行 git add 失败：%v", err)
		}

		out, err := runCmd("git", "commit", "-m", *pushComment)
		if err != nil {
			// 没有变更时 git commit 返回非零退出码，但不应视为错误
			if strings.Contains(string(out), "nothing to commit") {
				fmt.Println("没有变更需要提交，跳过提交和推送")
				return
			}
			log.Fatalf("执行 git commit 失败：%v", err)
		}

		if *pushBranch != "" {
			gitPull(nil, "  git add . && git commit && git push")
			if _, err := runCmd("git", "push", "-u", "origin", *pushBranch); err != nil {
				log.Fatalf("执行 git push 失败：%v", err)
			}
		} else {
			gitPull([]string{"origin", *pushBranch}, "  git add . && git commit && git push")
			if _, err := runCmd("git", "push"); err != nil {
				log.Fatalf("执行 git push 失败：%v", err)
			}
		}

		if *pushTag != "" {
			if _, err := runCmd("git", "tag", *pushTag); err != nil {
				log.Fatalf("执行 git tag 失败：%v", err)
			}
			if _, err := runCmd("git", "push", "origin", *pushTag); err != nil {
				log.Fatalf("执行 git push 失败：%v", err)
			}
		}

	case "rebase":
		var (
			err error
			out []byte
		)

		rebaseCmd := flag.NewFlagSet("rebase", flag.ExitOnError)
		rebaseBranch := rebaseCmd.String("branch", "", "要 rebase 的目标分支")
		rebaseVerbose := rebaseCmd.Bool("v", false, "显示详细执行过程")
		rebaseForceWithLease := rebaseCmd.Bool("force-with-lease", false, "强制推送 lease")
		rebaseForce := rebaseCmd.Bool("force", false, "强制推送")
		rebaseStash := rebaseCmd.Bool("stash", false, "是否 stash 变更")
		rebaseComment := rebaseCmd.String("comment", "", "合并时提交注释")

		needStashPop := true

		rebaseCmd.Parse(os.Args[2:])
		verbose = *rebaseVerbose

		if *rebaseBranch == "" {
			log.Fatalln("目标分支不能为空")
		}

		if *rebaseComment != "" {
			if _, err := runCmd("git", "add", "--all"); err != nil {
				log.Fatalf("执行 git add 失败：%v", err)
			}

			if _, err := runCmd("git", "commit", "-m", *rebaseComment); err != nil {
				log.Fatalf("执行 git commit 失败：%v", err)
			}

			if _, err := runCmd("git", "push"); err != nil {
				log.Fatalf("执行 git push 失败：%v", err)
			}

			// 先 fetch 获取远程最新状态
			if _, err := runCmd("git", "fetch", "origin"); err != nil {
				log.Fatalf("执行 git fetch 失败：%v", err)
			}

			// 执行 rebase
			if out, err = runCmd("git", "rebase", fmt.Sprintf("origin/%s", *rebaseBranch)); err != nil {
				if strings.Contains(string(out), "CONFLICT") {
					fmt.Println("git rebase 遇到冲突，请手动解决冲突后执行:")
					fmt.Println("  git add .")
					fmt.Println("  git rebase --continue")
					fmt.Println("  git push --force-with-lease")
					os.Exit(1)
				}
				log.Fatalf("执行 git rebase 失败：%v", err)
			}

			if *rebaseForceWithLease {
				if _, err := runCmd("git", "push", "--force-with-lease"); err != nil {
					log.Fatalf("执行 git push 失败：%v", err)
				}
			} else if *rebaseForce {
				if _, err := runCmd("git", "push", "--force"); err != nil {
					log.Fatalf("执行 git push 失败：%v", err)
				}
			}

			return
		}

		if *rebaseStash {
			out, err = runCmd("git", "stash")
			if err != nil {
				log.Fatalf("git stash 失败：%v", err)
			}

			if strings.Contains(string(out), "No local changes to save") {
				needStashPop = false
			}

			// 先 fetch 获取远程最新状态
			if _, err := runCmd("git", "fetch", "origin"); err != nil {
				log.Fatalf("执行 git fetch 失败：%v", err)
			}

			// 执行 rebase
			if out, err = runCmd("git", "rebase", fmt.Sprintf("origin/%s", *rebaseBranch)); err != nil {
				if strings.Contains(string(out), "CONFLICT") {
					fmt.Println("git rebase 遇到冲突，请手动解决冲突后执行:")
					fmt.Println("  git add .")
					fmt.Println("  git rebase --continue")
					fmt.Println("  git push --force-with-lease")
					os.Exit(1)
				}
				log.Fatalf("执行 git rebase 失败：%v", err)
			}

			if *rebaseForceWithLease {
				if _, err := runCmd("git", "push", "--force-with-lease"); err != nil {
					log.Fatalf("执行 git push 失败：%v", err)
				}
			} else if *rebaseForce {
				if _, err := runCmd("git", "push", "--force"); err != nil {
					log.Fatalf("执行 git push 失败：%v", err)
				}
			}

			if needStashPop {
				if _, err := runCmd("git", "stash", "pop"); err != nil {
					log.Fatalf("git stash pop 失败：%v", err)
				}
			}
		}

		fmt.Println("rebase 完成并已推送")

	case "merge":
		mergeCmd := flag.NewFlagSet("merge", flag.ExitOnError)
		mergeBranch := mergeCmd.String("branch", "", "分支，格式: src -> dst")
		mergeVerbose := mergeCmd.Bool("v", false, "显示详细执行过程")
		mergeCmd.Parse(os.Args[2:])
		verbose = *mergeVerbose

		if *mergeBranch == "" {
			log.Fatalln("分支不能为空")
		}

		var branchSrc, branchDst string

		if !strings.Contains(*mergeBranch, " -> ") {
			log.Fatalln("分支格式错误，请使用 src -> dst")
		}
		branchParts := strings.Split(*mergeBranch, " -> ")
		branchSrc = branchParts[0]
		branchDst = branchParts[1]

		if _, err := runCmd("git", "checkout", branchDst); err != nil {
			log.Fatalf("执行 git checkout %s 失败：%v", branchDst, err)
		}

		if _, err := runCmd("git", "fetch", "origin"); err != nil {
			log.Fatalf("执行 git fetch 失败：%v", err)
		}

		gitPull([]string{"origin", branchDst}, fmt.Sprintf(" git add --all && git commit && git push origin %s\n  git checkout %s", branchDst, branchSrc))

		if _, err := runCmd("git", "fetch", "origin"); err != nil {
			log.Fatalf("执行 git fetch 失败：%v", err)
		}

		out, err := runCmd("git", "merge", fmt.Sprintf("origin/%s", branchSrc))
		if err != nil {
			if strings.Contains(string(out), "CONFLICT") {
				fmt.Printf("git merge 遇到冲突，请手动解决冲突后执行:\n")
				fmt.Printf("  git add . && git commit && git push origin %s\n", branchDst)
				fmt.Printf("  git checkout %s\n", branchSrc)
				os.Exit(1)
			}
			log.Fatalf("执行 git merge 失败：%v", err)
		}

		if _, err := runCmd("git", "push", "origin", branchDst); err != nil {
			log.Fatalf("执行 git push origin %s 失败：%v", branchDst, err)
		}

		if _, err := runCmd("git", "checkout", branchSrc); err != nil {
			log.Fatalf("执行 git checkout %s 失败：%v", branchSrc, err)
		}

	case "tag-list":
		tagListCmd := flag.NewFlagSet("tag-list", flag.ExitOnError)
		tagListCount := tagListCmd.Int("count", 1, "显示标签数量")
		tagListVerbose := tagListCmd.Bool("v", false, "显示详细执行过程")
		tagListCmd.Parse(os.Args[2:])
		verbose = *tagListVerbose

		out, err := runCmd("git", "tag", "--sort=-v:refname")
		if err != nil {
			log.Fatalf("执行 git tag 失败：%v", err)
		}
		tags := strings.Split(strings.TrimSpace(string(out)), "\n")
		// 取前 count 条（最新版本）
		if len(tags) > *tagListCount {
			tags = tags[:*tagListCount]
		}
		fmt.Printf("最近 %d 个标签:\n", *tagListCount)
		for _, t := range tags {
			fmt.Println(" ", t)
		}

	case "tag-last":
		tagLastCmd := flag.NewFlagSet("tag-last", flag.ExitOnError)
		tagLastVerbose := tagLastCmd.Bool("v", false, "显示详细执行过程")
		tagLastCmd.Parse(os.Args[2:])
		verbose = *tagLastVerbose

		out, err := runCmd("git", "describe", "--tags")
		if err != nil {
			log.Fatalf("执行 git describe 失败：%v", err)
		}
		fmt.Println("最新标签:", strings.TrimSpace(string(out)))

	case "tag-push":
		tagCmd := flag.NewFlagSet("tag-push", flag.ExitOnError)
		tagName := tagCmd.String("tag", "", "标签")
		tagVerbose := tagCmd.Bool("v", false, "显示详细执行过程")
		tagCmd.Parse(os.Args[2:])
		verbose = *tagVerbose

		if *tagName == "" {
			log.Fatalln("标签不能为空")
		}

		if _, err := runCmd("git", "tag", *tagName); err != nil {
			log.Fatalf("执行 git tag %s 失败：%v", *tagName, err)
		}

		if _, err := runCmd("git", "push", "origin", *tagName); err != nil {
			log.Fatalf("执行 git push 失败：%v", err)
		}

	case "help":
		if len(os.Args) >= 3 {
			switch os.Args[2] {
			case "push":
				printPushHelp()
			case "rebase":
				printRebaseHelp()
			case "merge":
				printMergeHelp()
			case "tag-push":
				printTagPushHelp()
			case "tag-list":
				printTagListHelp()
			case "tag-last":
				printTagLastHelp()
			default:
				fmt.Printf("未知命令: %s\n", os.Args[2])
				printUsage()
			}
		} else {
			printUsage()
		}

	default:
		fmt.Printf("未知命令: %s\n", os.Args[1])
		printUsage()
	}
}

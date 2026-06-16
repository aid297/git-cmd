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
  merge       合并分支并推送
  tag-last    查看最新标签
  tag-push    创建并推送标签
  help        显示帮助信息

使用 git-cmd help <command> 查看具体命令的用法`)
}

func printPushHelp() {
	fmt.Println(`push - 一键提交并推送代码

用法:
  git-cmd push --comment <提交信息> [--branch <分支名>]

选项:
  --comment   提交注释（必填）
  --branch    推送到的远程分支名（可选，默认推送当前分支）

示例:
  git-cmd push --comment "修复bug"
  git-cmd push --comment "新功能" --branch dev`)
}

func printMergeHelp() {
	fmt.Println(`merge - 合并分支并推送

用法:
  git-cmd merge --branch <源分支> -> <目标分支>

选项:
  --branch   分支合并方向，格式: src -> dst（必填）

示例:
  git-cmd merge --branch "feature -> main"`)
}

func printTagPushHelp() {
	fmt.Println(`tag-push - 创建并推送标签

用法:
  git-cmd tag-push --tag <标签名>

选项:
  --tag   标签名（必填）

示例:
  git-cmd tag-push --tag v1.0.0`)
}

func printTagLastHelp() {
	fmt.Println(`tag-last - 查看最新标签

用法:
  git-cmd tag-last

示例:
  git-cmd tag-last`)
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
		pushCmd.Parse(os.Args[2:])

		err := exec.Command("git", "add", "--all").Run()
		if err != nil {
			log.Fatalf("执行 git add 失败：%v", err)
		}

		if *pushComment == "" {
			log.Fatalln("提交注释不能为空")
		}
		err = exec.Command("git", "commit", "-m", *pushComment).Run()
		if err != nil {
			log.Fatalf("执行 git commit 失败：%v", err)
		}

		if *pushBranch == "" {
			err = exec.Command("git", "pull").Run()
			if err != nil {
				log.Fatalf("执行 git pull 失败：%v", err)
			}

			err = exec.Command("git", "push").Run()
			if err != nil {
				log.Fatalf("执行 git push 失败：%v", err)
			}
		} else {
			err = exec.Command("git", "pull").Run()
			if err != nil {
				log.Fatalf("执行 git pull 失败：%v", err)
			}

			err = exec.Command("git", "push", "-u", "origin", *pushBranch).Run()
			if err != nil {
				log.Fatalf("执行 git push 失败：%v", err)
			}
		}

	case "merge":
		mergeCmd := flag.NewFlagSet("merge", flag.ExitOnError)
		mergeBranch := mergeCmd.String("branch", "", "分支，格式: src -> dst")
		mergeCmd.Parse(os.Args[2:])

		if *mergeBranch == "" {
			log.Fatalln("分支不能为空")
		}

		branchParts := strings.Split(*mergeBranch, " -> ")
		if len(branchParts) != 2 {
			log.Fatalln("分支格式错误，需要为 src -> dst")
		}
		branchSrc := branchParts[0]
		branchDst := branchParts[1]

		err := exec.Command("git", "checkout", branchDst).Run()
		if err != nil {
			log.Fatalf("执行 git checkout %s 失败：%v", branchDst, err)
		}

		err = exec.Command("git", "fetch", "origin").Run()
		if err != nil {
			log.Fatalf("执行 git fetch 失败：%v", err)
		}

		err = exec.Command("git", "pull", "origin", branchDst).Run()
		if err != nil {
			log.Fatalf("执行 git pull 失败：%v", err)
		}

		err = exec.Command("git", "merge", fmt.Sprintf("origin/%s", branchSrc)).Run()
		if err != nil {
			log.Fatalf("执行 git merge 失败：%v", err)
		}

		err = exec.Command("git", "push", "origin", branchDst).Run()
		if err != nil {
			log.Fatalf("执行 git push origin %s 失败：%v", branchDst, err)
		}

		err = exec.Command("git", "checkout", branchSrc).Run()
		if err != nil {
			log.Fatalf("执行 git checkout %s 失败：%v", branchSrc, err)
		}

	case "tag-last":
		cmd := exec.Command("git", "describe", "--tags")
		out, err := cmd.Output()
		if err != nil {
			log.Fatalf("执行 git describe 失败：%s -> %v", out, err)
		}
		fmt.Println("最新标签:", strings.TrimSpace(string(out)))

	case "tag-push":
		tagCmd := flag.NewFlagSet("tag-push", flag.ExitOnError)
		tagName := tagCmd.String("tag", "", "标签")
		tagCmd.Parse(os.Args[2:])

		if *tagName == "" {
			log.Fatalln("标签不能为空")
		}

		err := exec.Command("git", "tag", *tagName).Run()
		if err != nil {
			log.Fatalf("执行 git tag %s 失败：%v", *tagName, err)
		}

		err = exec.Command("git", "push", "origin", *tagName).Run()
		if err != nil {
			log.Fatalf("执行 git push 失败：%v", err)
		}

	case "help":
		if len(os.Args) >= 3 {
			switch os.Args[2] {
			case "push":
				printPushHelp()
			case "merge":
				printMergeHelp()
			case "tag-push":
				printTagPushHelp()
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

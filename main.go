package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("参数不足，至少需要一个子命令")
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
		fmt.Println("最新标签:", string(out))

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
	}
}

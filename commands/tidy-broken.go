package commands

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"github.com/zyrouge/pho/core"
	"github.com/zyrouge/pho/utils"
)

var TidyBrokenCommand = cli.Command{
	Name:    "tidy-broken",
	Aliases: []string{},
	Usage:   "Remove incomplete files",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "assume-yes",
			Aliases: []string{"y"},
			Usage:   "Automatically answer yes for questions",
		},
	},
	Action: func(ctx *cli.Context) error {
		utils.LogDebug("reading transactions")
		transactions, err := core.GetTransactions()
		if err != nil {
			return err
		}

		reader := bufio.NewReader(os.Stdin)
		assumeYes := ctx.Bool("assume-yes")
		utils.LogDebug(fmt.Sprintf("argument assume-yes: %v", assumeYes))

		utils.LogLn()
		utils.LogInfo("List of affected directories and files:")
		involvedIds := []string{}
		involvedDirs := []string{}
		involvedFiles := []string{}
		for k, v := range transactions.PendingInstallations {
			involvedIds = append(involvedIds, k)
			involvedDirs = append(involvedDirs, v.InvolvedDirs...)
			involvedFiles = append(involvedFiles, v.InvolvedFiles...)
			for _, x := range v.InvolvedDirs {
				utils.LogInfo(
					fmt.Sprintf("%s %s", color.HiBlackString("D"), color.RedString(x)),
				)
			}
			for _, x := range v.InvolvedFiles {
				utils.LogInfo(
					fmt.Sprintf("%s %s", color.HiBlackString("F"), color.RedString(x)),
				)
			}
		}

		if len(involvedDirs)+len(involvedFiles) == 0 {
			utils.LogInfo(color.HiBlackString("no directories or files are affected"))
			utils.LogLn()
			utils.LogInfo(
				fmt.Sprintf("%s Everything is working perfectly!", utils.LogTickPrefix),
			)
			return nil
		}

		if !assumeYes {
			utils.LogLn()
			proceed, err := utils.PromptYesNoInput(reader, "Do you want to proceed?")
			if err != nil {
				return err
			}
			if !proceed {
				utils.LogWarning("aborted...")
				return nil
			}
		}

		removedDirsCount := 0
		removedFilesCount := 0
		utils.LogLn()
		for _, x := range involvedDirs {
			utils.LogDebug(fmt.Sprintf("removing %s", x))
			if err := os.RemoveAll(x); err != nil {
				utils.LogError(err)
				continue
			}
			removedDirsCount++
		}
		for _, x := range involvedFiles {
			utils.LogDebug(fmt.Sprintf("removing %s", x))
			if err := os.Remove(x); err != nil {
				utils.LogError(err)
				continue
			}
			removedFilesCount++
		}
		core.UpdateTransactions(func(transactions *core.Transactions) error {
			for _, x := range involvedIds {
				delete(transactions.PendingInstallations, x)
			}
			return nil
		})

		utils.LogLn()
		utils.LogInfo(
			fmt.Sprintf(
				"%s Removed %d directories and %d files successfully!",
				utils.LogTickPrefix,
				removedDirsCount,
				removedFilesCount,
			),
		)

		return nil
	},
}

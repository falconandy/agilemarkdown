package commands

import (
	"fmt"
	"github.com/mreider/agilemarkdown/backlog"
	"github.com/mreider/agilemarkdown/git"
	"github.com/mreider/agilemarkdown/utils"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var chartColorCodeRe = regexp.MustCompile(`.\[\d+m`)

var SyncCommand = cli.Command{
	Name:      "sync",
	Usage:     "Sync state",
	ArgsUsage: " ",
	Action: func(c *cli.Context) error {
		action := &SyncAction{}
		return action.Execute()
	},
}

type SyncAction struct {
}

func (a *SyncAction) Execute() error {
	rootDir, _ := filepath.Abs(".")
	if err := checkIsBacklogDirectory(); err == nil {
		rootDir = filepath.Dir(rootDir)
	} else if err := checkIsRootDirectory(); err != nil {
		return err
	}

	err := a.updateOverviews(rootDir)
	if err != nil {
		return err
	}

	err = a.updateHome(rootDir)
	if err != nil {
		return err
	}

	err = a.updateSidebar(rootDir)
	if err != nil {
		return err
	}

	return a.syncToGit()
}

func (a *SyncAction) updateOverviews(rootDir string) error {
	backlogDirs, err := a.backlogDirs(rootDir)
	if err != nil {
		return err
	}
	for _, backlogDir := range backlogDirs {
		overview, err := backlog.LoadBacklogOverview(filepath.Join(backlogDir, backlog.OverviewFileName))
		if err != nil {
			return err
		}
		bck, err := backlog.LoadBacklog(backlogDir)
		if err != nil {
			return err
		}

		items := bck.Items()
		overview.Update(items)
	}
	return nil
}

func (a *SyncAction) updateHome(rootDir string) error {
	var lines []string
	backlogDirs, err := a.backlogDirs(rootDir)
	if err != nil {
		return err
	}
	for _, backlogDir := range backlogDirs {
		overview, err := backlog.LoadBacklogOverview(filepath.Join(backlogDir, backlog.OverviewFileName))
		if err != nil {
			return err
		}
		lines = append(lines, fmt.Sprintf("### [%s](%s)", overview.Title(), strings.TrimSuffix(backlog.OverviewFileName, ".md")))
		bck, err := backlog.LoadBacklog(backlogDir)
		if err != nil {
			return err
		}

		progressAction := &ProgressAction{}
		chart, err := progressAction.Execute(backlogDir, 12)
		if err != nil {
			return err
		}

		chart = chartColorCodeRe.ReplaceAllString(chart, "")
		lines = append(lines, utils.WrapLinesToMarkdownCodeBlock(strings.Split(chart, "\n"))...)

		flying := backlog.StatusByCode("f")
		items := bck.ItemsByStatus(flying.Code)
		itemsLines := backlog.BacklogView{}.WriteBacklogItems(items, fmt.Sprintf("Status: %s", flying.Name))
		itemsLines = utils.WrapLinesToMarkdownCodeBlock(itemsLines)
		lines = append(lines, itemsLines...)
	}
	err = ioutil.WriteFile(filepath.Join(rootDir, "Home.md"), []byte(strings.Join(lines, "  \n")), 0644)
	return err
}

func (a *SyncAction) updateSidebar(rootDir string) error {
	var lines []string
	backlogDirs, err := a.backlogDirs(rootDir)
	if err != nil {
		return err
	}
	for _, backlogDir := range backlogDirs {
		overview, err := backlog.LoadBacklogOverview(filepath.Join(backlogDir, backlog.OverviewFileName))
		if err != nil {
			return err
		}
		lines = append(lines, fmt.Sprintf("[%s](%s)", overview.Title(), strings.TrimSuffix(backlog.OverviewFileName, ".md")))
	}
	err = ioutil.WriteFile(filepath.Join(rootDir, "_Sidebar.md"), []byte(strings.Join(lines, "  \n")), 0644)
	return err
}

func (a *SyncAction) syncToGit() error {
	err := git.AddAll()
	if err != nil {
		return err
	}
	git.Commit("sync") // TODO commit message
	err = git.Fetch()
	if err != nil {
		return fmt.Errorf("can't fetch: %v", err)
	}
	output, err := git.Merge()
	if err != nil {
		status, _ := git.Status()
		if !strings.Contains(status, "Your branch is based on 'origin/master', but the upstream is gone.") {
			fmt.Println(output)
			git.AbortMerge()
			return fmt.Errorf("can't merge: %v", err)
		}
	}
	err = git.Push()
	if err != nil {
		return fmt.Errorf("can't push: %v", err)
	}
	return nil
}

func (a *SyncAction) backlogDirs(rootDir string) ([]string, error) {
	infos, err := ioutil.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(infos))
	for _, info := range infos {
		if !info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			continue
		}
		result = append(result, filepath.Join(rootDir, info.Name()))
	}
	sort.Strings(result)
	return result, nil
}

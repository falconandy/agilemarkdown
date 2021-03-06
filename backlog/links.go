package backlog

import (
	"fmt"
	"github.com/mreider/agilemarkdown/utils"
	"path/filepath"
	"strings"
)

func MakeItemLink(item *BacklogItem, baseDir string) string {
	itemPath := item.markdown.contentPath
	if itemPath == "" {
		itemPath = item.Name()
	}
	return utils.MakeMarkdownLink(item.Title(), itemPath, baseDir)
}

func MakeIdeaLink(idea *BacklogIdea, baseDir string) string {
	ideaPath := idea.markdown.contentPath
	if ideaPath == "" {
		ideaPath = idea.Name()
	}
	return utils.MakeMarkdownLink(idea.Title(), ideaPath, baseDir)
}

func MakeOverviewLink(overview *BacklogOverview, baseDir string) string {
	overviewPath := overview.markdown.contentPath
	return utils.MakeMarkdownLink(overview.Title(), overviewPath, baseDir)
}

func MakeArchiveLink(archive *BacklogOverview, title string, baseDir string) string {
	archivePath := archive.markdown.contentPath
	return utils.MakeMarkdownLink(title, archivePath, baseDir)
}

func MakeIndexLink(rootDir, baseDir string) string {
	return utils.MakeMarkdownLink("home", filepath.Join(rootDir, IndexFileName), baseDir)
}

func MakeIdeasLink(rootDir, baseDir string) string {
	return utils.MakeMarkdownLink("idea list", filepath.Join(rootDir, IdeasFileName), baseDir)
}

func MakeTagsLink(rootDir, baseDir string) string {
	return utils.MakeMarkdownLink("tag list", filepath.Join(rootDir, TagsFileName), baseDir)
}

func MakeTagLink(tag, tagsDir, baseDir string) string {
	return utils.MakeMarkdownLink(tag, filepath.Join(tagsDir, fmt.Sprintf("%s.md", utils.GetValidFileName(strings.ToLower(tag)))), baseDir)
}

func MakeTagLinks(tags []string, tagsDir, baseDir string) string {
	links := make([]string, 0, len(tags))
	for _, tag := range tags {
		links = append(links, MakeTagLink(tag, tagsDir, baseDir))
	}

	return strings.Join(links, " ")
}

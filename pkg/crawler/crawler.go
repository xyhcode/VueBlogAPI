package crawler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

// Friend 表示一个友链
type Friend struct {
	Name   string
	Link   string
	Avatar string
	Descr  string
}

// Post 表示一篇文章
type Post struct {
	Title        string
	Link         string
	Created      string
	Updated      string
	Rules        []string
	Author       string
	Avatar       string
	FriendLink   string
	UsedCSSRules map[string]string
}

// CSSRule 定义CSS选择规则
type CSSRule struct {
	Selector string `yaml:"selector"`
	Attr     string `yaml:"attr"`
}

// CSSRules 定义完整的规则集合
type CSSRules struct {
	PostPageRules map[string]map[string][]CSSRule `yaml:"post_page_rules"`
}

// Crawler 爬虫结构体
type Crawler struct {
	Client        *http.Client
	CSSRules      CSSRules
	MaxPostsNum   int
	MaxConcurrent int
}

// NewCrawler 创建新的爬虫实例
func NewCrawler() (*Crawler, error) {
	// 使用硬编码的CSS规则
	cssRules := GetDefaultCSSRules()

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	return &Crawler{
		Client:        client,
		CSSRules:      cssRules,
		MaxPostsNum:   5,
		MaxConcurrent: 3,
	}, nil
}

// CrawlAllFriends 爬取所有友链站点
func (c *Crawler) CrawlAllFriends(friends []Friend) ([]Post, int, error) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, c.MaxConcurrent)
	postsChan := make(chan []Post, len(friends))
	errChan := make(chan error, len(friends))

	for _, friend := range friends {
		wg.Add(1)
		go func(friend Friend) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			posts, err := c.CrawlPosts(friend)
			if err != nil {
				errChan <- fmt.Errorf("爬取 %s 失败: %w", friend.Link, err)
				return
			}

			postsChan <- posts
		}(friend)
	}

	// 等待所有爬取完成
	go func() {
		wg.Wait()
		close(postsChan)
		close(errChan)
	}()

	// 收集结果
	allPosts := []Post{}
	for posts := range postsChan {
		allPosts = append(allPosts, posts...)
	}

	// 收集错误
	errors := []error{}
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return allPosts, len(errors), fmt.Errorf("部分爬取失败: %v", errors)
	}

	return allPosts, 0, nil
}

// CrawlPosts 爬取单个友链的文章
func (c *Crawler) CrawlPosts(friend Friend) ([]Post, error) {
	return c.crawlPostPages(friend.Link, friend)
}

// crawlPostPages 爬取单个站点的文章页面
func (c *Crawler) crawlPostPages(baseURL string, friend Friend) ([]Post, error) {
	// 尝试Feed解析
	posts, err := c.crawlWithFeed(baseURL)
	if err == nil && len(posts) > 0 {
		// 添加作者信息
		for i := range posts {
			posts[i].Author = friend.Name
			posts[i].Avatar = friend.Avatar
			// 标准化 friend_link，移除 trailing slash
			friendLink := friend.Link
			if len(friendLink) > 0 && friendLink[len(friendLink)-1] == '/' {
				friendLink = friendLink[:len(friendLink)-1]
			}
			posts[i].FriendLink = friendLink
		}
		// 应用最大文章数限制
		if c.MaxPostsNum > 0 && len(posts) > c.MaxPostsNum {
			posts = posts[:c.MaxPostsNum]
		}
		return posts, nil
	}

	// Feed解析失败，尝试CSS规则解析
	posts, err = c.crawlWithCSSRules(baseURL)
	if err != nil {
		return nil, err
	}

	// 添加作者信息
	for i := range posts {
		posts[i].Author = friend.Name
		posts[i].Avatar = friend.Avatar
		// 标准化 friend_link，移除 trailing slash
		friendLink := friend.Link
		if len(friendLink) > 0 && friendLink[len(friendLink)-1] == '/' {
			friendLink = friendLink[:len(friendLink)-1]
		}
		posts[i].FriendLink = friendLink
	}

	// 应用最大文章数限制
	if c.MaxPostsNum > 0 && len(posts) > c.MaxPostsNum {
		posts = posts[:c.MaxPostsNum]
	}

	return posts, nil
}

// crawlWithFeed 使用Feed解析文章
func (c *Crawler) crawlWithFeed(baseURL string) ([]Post, error) {
	feedSuffixes := []string{
		"atom.xml",
		"feed/atom",
		"rss.xml",
		"rss2.xml",
		"feed",
		"index.xml",
	}

	for _, suffix := range feedSuffixes {
		feedURL := fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), suffix)

		resp, err := c.Client.Get(feedURL)
		if err != nil {
			continue
		}

		if resp.StatusCode != 200 {
			resp.Body.Close()
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		// 使用专业的RSS/Atom解析库
		feedParser := gofeed.NewParser()
		feed, err := feedParser.ParseString(string(body))
		if err != nil {
			continue
		}

		// 转换为内部Post结构
		posts := []Post{}
		for _, item := range feed.Items {
			link := item.Link
			if link == "" && item.GUID != "" {
				link = item.GUID
			}

			created := ""
			if item.Published != "" {
				created = item.Published
			} else if item.Updated != "" {
				created = item.Updated
			}

			updated := created
			if item.Updated != "" && item.Updated != created {
				updated = item.Updated
			}

			if item.Title != "" && link != "" {
				posts = append(posts, Post{
					Title:        item.Title,
					Link:         link,
					Created:      formatDate(created),
					Updated:      formatDate(updated),
					Rules:        []string{"feed"},
					UsedCSSRules: make(map[string]string),
				})
			}
		}

		if len(posts) > 0 {
			return posts, nil
		}
	}

	return nil, fmt.Errorf("所有feed后缀都无法解析")
}

// crawlWithCSSRules 使用CSS规则解析文章
func (c *Crawler) crawlWithCSSRules(baseURL string) ([]Post, error) {
	resp, err := c.Client.Get(baseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	baseURLParsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	usedRules := []string{}

	// 遍历所有主题规则
	for themeName, themeRules := range c.CSSRules.PostPageRules {
		usedRules = append(usedRules, themeName)

		results := make(map[string][]string)
		usedFieldRules := make(map[string]string) // 记录每个字段使用的规则
		allFieldsFound := true

		// 按顺序查找必须字段：title, link, created, updated
		requiredFields := []string{"title", "link", "created", "updated"}
		for _, field := range requiredFields {
			fieldRules, exists := themeRules[field]
			if !exists {
				allFieldsFound = false
				break
			}

			fieldValues := []string{}
			ruleApplied := false
			appliedRuleSelector := ""

			// 尝试每个规则直到找到内容
			for _, rule := range fieldRules {
				selection := doc.Find(rule.Selector)
				if selection.Length() == 0 {
					continue
				}

				selection.Each(func(i int, s *goquery.Selection) {
					var value string
					switch rule.Attr {
					case "text":
						value = s.Text()
					case "time":
						// 尝试获取datetime属性
						if val, exists := s.Attr("datetime"); exists {
							value = val
						} else {
							// 否则获取text内容
							value = s.Text()
						}
					default:
						if val, exists := s.Attr(rule.Attr); exists {
							value = val
						} else {
							return
						}
					}

					// 清理空白字符
					value = strings.TrimSpace(value)
					if value != "" {
						fieldValues = append(fieldValues, value)
					}
				})

				if len(fieldValues) > 0 {
					ruleApplied = true
					appliedRuleSelector = rule.Selector // 记录应用的规则选择器
					break                               // 找到内容就停止尝试其他规则
				}
			}

			if !ruleApplied {
				allFieldsFound = false
				break
			}

			results[field] = fieldValues
			usedFieldRules[field] = appliedRuleSelector // 记录该字段使用的规则
		}

		// 检查所有字段长度是否一致
		if allFieldsFound {
			length := len(results["title"])
			if len(results["link"]) != length || len(results["created"]) != length {
				allFieldsFound = false
			}
		}

		if allFieldsFound {
			// 创建文章对象
			posts := []Post{}
			maxLen := len(results["title"])

			for i := 0; i < maxLen; i++ {
				title := results["title"][i]
				linkStr := results["link"][i]
				created := formatDate(results["created"][i])

				updated := created
				if i < len(results["updated"]) {
					updated = formatDate(results["updated"][i])
				}

				// 处理相对链接
				parsedLink, err := url.Parse(linkStr)
				if err == nil && !parsedLink.IsAbs() {
					linkStr = baseURLParsed.ResolveReference(parsedLink).String()
				}

				posts = append(posts, Post{
					Title:        title,
					Link:         linkStr,
					Created:      created,
					Updated:      updated,
					Rules:        []string{themeName},
					UsedCSSRules: usedFieldRules,
				})
			}

			return posts, nil
		}
	}

	return []Post{}, nil
}

// formatDate 格式化日期字符串
func formatDate(dateStr string) string {
	if dateStr == "" {
		return time.Now().Format("2006-01-02 15:04:05")
	}

	// 移除可能的多余空格
	dateStr = strings.TrimSpace(dateStr)

	// 尝试解析常见的日期格式并转换为统一格式
	formats := []string{
		"Mon, 02 Jan 2006 15:04:05 MST",
		"Mon, 02 Jan 2006 15:04:05 GMT",
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 MST",
		"Mon, 2 Jan 2006 15:04:05 GMT",
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05+00:00",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"Jan 2, 2006",
		"January 2, 2006",
		"02 Jan 2006",
		"2 Jan 2006",
		"2006/01/02 15:04:05",
		"2006/01/02",
		"02/01/2006",
		"2/01/2006",
		"02-01-2006",
		"2-01-2006",
		"Monday, 02-Jan-06 15:04:05 MST",
		"Monday, 2-Jan-06 15:04:05 MST",
		"Mon Jan _2 15:04:05 2006",
		"Mon Jan _2 15:04:05 MST 2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Format("2006-01-02 15:04:05")
		}
	}

	// 尝试处理带毫秒的格式
	millisecondFormats := []string{
		"2006-01-02T15:04:05.999Z",
		"2006-01-02T15:04:05.999+00:00",
		"2006-01-02 15:04:05.999",
	}

	for _, format := range millisecondFormats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Format("2006-01-02 15:04:05")
		}
	}

	return dateStr
}

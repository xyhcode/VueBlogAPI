package crawler

// GetDefaultCSSRules 返回默认的CSS规则
func GetDefaultCSSRules() CSSRules {
	return CSSRules{
		PostPageRules: map[string]map[string][]CSSRule{
			"anzhiyu": {
				"title": {
					{Selector: "#recent-posts>.recent-post-item .recent-post-info .article-title", Attr: "text"},
				},
				"link": {
					{Selector: "#recent-posts>.recent-post-item .recent-post-info .article-title", Attr: "href"},
				},
				"created": {
					{Selector: "#recent-posts > div.recent-post-item.lastestpost-item > div.recent-post-info > div.article-meta-wrap > span.post-meta-date > time", Attr: "time"},
				},
				"updated": {
					{Selector: "#recent-posts .recent-post-info .post-meta-date time:nth-of-type(2)", Attr: "time"},
				},
			},
			"butterfly": {
				"title": {
					{Selector: "#recent-posts .recent-post-info a:first-child", Attr: "text"},
				},
				"link": {
					{Selector: "#recent-posts .recent-post-info a:first-child", Attr: "href"},
				},
				"created": {
					{Selector: "#recent-posts .recent-post-info .post-meta-date time:first-of-type", Attr: "text"},
				},
				"updated": {
					{Selector: "#recent-posts .recent-post-info .post-meta-date time:nth-of-type(2)", Attr: "text"},
				},
			},
			"fluid": {
				"title": {
					{Selector: "#board .index-header a", Attr: "text"},
				},
				"link": {
					{Selector: "#board .index-header a", Attr: "href"},
				},
				"created": {
					{Selector: "#board .post-meta time", Attr: "text"},
				},
				"updated": {
					{Selector: "#board .post-meta time", Attr: "text"},
				},
			},
			"matery": {
				"title": {
					{Selector: "#articles .card .card-title", Attr: "text"},
				},
				"link": {
					{Selector: "#articles .card a:first-child", Attr: "href"},
				},
				"created": {
					{Selector: "#articles .card span.publish-date", Attr: "text"},
				},
				"updated": {
					{Selector: "#articles .card span.publish-date", Attr: "text"},
				},
			},
			"sakura": {
				"title": {
					{Selector: "#main a.post-title h3", Attr: "text"},
				},
				"link": {
					{Selector: "#main a.post-title", Attr: "href"},
				},
				"created": {
					{Selector: "#main .post-date", Attr: "text"},
				},
				"updated": {
					{Selector: "#main .post-date", Attr: "text"},
				},
			},
			"volantis": {
				"title": {
					{Selector: ".post-list .article-title a", Attr: "text"},
				},
				"link": {
					{Selector: ".post-list .article-title a", Attr: "href"},
				},
				"created": {
					{Selector: ".post-list .meta-v3 time", Attr: "text"},
				},
				"updated": {
					{Selector: ".post-list .meta-v3 time", Attr: "text"},
				},
			},
			"nexmoe": {
				"title": {
					{Selector: "section.nexmoe-posts .nexmoe-post h1", Attr: "text"},
				},
				"link": {
					{Selector: "section.nexmoe-posts .nexmoe-post>a", Attr: "href"},
				},
				"created": {
					{Selector: "section.nexmoe-posts .nexmoe-post-meta a:first-child", Attr: "text"},
				},
				"updated": {
					{Selector: "section.nexmoe-posts .nexmoe-post-meta a:first-child", Attr: "text"},
				},
			},
			"stun": {
				"title": {
					{Selector: "article .post-title__link", Attr: "text"},
				},
				"link": {
					{Selector: "article .post-title__link", Attr: "href"},
				},
				"created": {
					{Selector: "article .post-meta .post-meta-item--createtime .post-meta-item__value", Attr: "text"},
				},
				"updated": {
					{Selector: "article .post-meta .post-meta-item--updatetime .post-meta-item__value", Attr: "text"},
				},
			},
			"stellar": {
				"title": {
					{Selector: ".post-list .post-card:not(.photo) .post-title", Attr: "text"},
				},
				"link": {
					{Selector: ".post-list .post-card:not(.photo)", Attr: "href"},
				},
				"created": {
					{Selector: ".post-list .post-card:not(.photo) #post-meta time", Attr: "text"},
				},
				"updated": {
					{Selector: ".post-list .post-card:not(.photo) #post-meta time", Attr: "text"},
				},
			},
			"next": {
				"title": {
					{Selector: "article h2 a:first-child", Attr: "text"},
					{Selector: "article .post-title", Attr: "text"},
					{Selector: "article .post-title-link", Attr: "text"},
				},
				"link": {
					{Selector: "article h2 a:first-child", Attr: "href"},
					{Selector: "article .post-title", Attr: "href"},
					{Selector: "article .post-title-link", Attr: "href"},
				},
				"created": {
					{Selector: "article time[itemprop*='dateCreated']", Attr: "text"},
				},
				"updated": {
					{Selector: "article time[itemprop='dateModified']", Attr: "text"},
				},
			},
			"default": {
				"title": {
					{Selector: ".post-title", Attr: "text"},
					{Selector: "h1", Attr: "text"},
				},
				"link": {
					{Selector: ".post-title a", Attr: "href"},
					{Selector: "h1 a", Attr: "href"},
				},
				"created": {
					{Selector: ".post-date", Attr: "text"},
					{Selector: ".date", Attr: "text"},
				},
				"updated": {
					{Selector: ".post-updated", Attr: "text"},
					{Selector: ".updated", Attr: "text"},
				},
			},
		},
	}
}

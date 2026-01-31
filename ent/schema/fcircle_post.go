package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// FCirclePost holds the schema definition for the FCirclePost entity.
type FCirclePost struct {
	ent.Schema
}

// Annotations of the FCirclePost.
func (FCirclePost) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("朋友圈文章表"),
	}
}

// Fields of the FCirclePost.
func (FCirclePost) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").Comment("文章标题").NotEmpty(),
		field.String("link").Comment("文章链接").NotEmpty().Unique(), // 为 link 字段添加唯一索引，确保文章链接唯一
		field.Time("created").Comment("创建时间").Optional(),
		field.Time("updated").Comment("更新时间").Optional(),
		field.String("author").Comment("作者名称").NotEmpty(),
		field.String("avatar").Comment("作者头像链接").Optional(),
		field.String("friend_link").Comment("友链用户网站链接").NotEmpty(),
		field.Time("crawled_at").Comment("爬取时间").Default(time.Now),
		field.String("rules").Comment("使用的规则类型").Optional(),
	}
}

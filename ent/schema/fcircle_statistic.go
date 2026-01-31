package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

// FCircleStatistic holds the schema definition for the FCircleStatistic entity.
type FCircleStatistic struct {
	ent.Schema
}

// Annotations of the FCircleStatistic.
func (FCircleStatistic) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("朋友圈统计信息表"),
	}
}

// Fields of the FCircleStatistic.
func (FCircleStatistic) Fields() []ent.Field {
	return []ent.Field{
		field.Int("friends_num").Comment("已收录的友链数量").Default(0),
		field.Int("active_num").Comment("正常运行的友链数量").Default(0),
		field.Int("error_num").Comment("异常/失效的友链数量").Default(0),
		field.Int("article_num").Comment("总文章数").Default(0),
		field.Time("last_updated_time").Comment("数据最后更新时间").Default(time.Now),
	}
}

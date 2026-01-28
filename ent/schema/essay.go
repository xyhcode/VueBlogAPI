/*
 * @Description: Essay schema for Ent ORM
 * @Author: Qwenjie
 * @Date: 2026-01-27
 * @LastEditTime: 2026-01-27
 * @LastEditors: Qwenjie
 */
package schema

import (
	"time"

	"github.com/anzhiyu-c/anheyu-app/ent/schema/mixin"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Essay holds the schema definition for the Essay entity.
type Essay struct {
	ent.Schema
}

// Annotations of the Essay.
func (Essay) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("随笔记录表"),
	}
}

// Mixin of the Essay.
func (Essay) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.SoftDeleteMixin{},
	}
}

// Fields of the Essay.
func (Essay) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id").
			Comment("主键ID"),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("创建时间"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("更新时间"),
		field.Text("content").
			Comment("随笔内容"),
		field.Time("date").
			Default(time.Now).
			Comment("随笔日期"),
		field.Text("images").
			Optional().
			Comment("随笔图片列表(JSON格式字符串)"),
		field.String("link").
			Optional().
			Comment("随笔链接"),
	}
}

// Edges of the Essay.
func (Essay) Edges() []ent.Edge {
	return nil
}

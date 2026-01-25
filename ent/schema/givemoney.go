/*
 * @Description: GiveMoney schema for Ent ORM
 * @Author: Qwenjie
 * @Date: 2026-01-24
 * @LastEditTime: 2026-01-24
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

// GiveMoney holds the schema definition for the GiveMoney entity.
type GiveMoney struct {
	ent.Schema
}

// Annotations of the GiveMoney.
func (GiveMoney) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("打赏记录表"),
	}
}

// Mixin of the GiveMoney.
func (GiveMoney) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.SoftDeleteMixin{},
	}
}

// Fields of the GiveMoney.
func (GiveMoney) Fields() []ent.Field {
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
		field.String("nickname").
			MaxLen(50).
			NotEmpty().
			Comment("昵称"),
		field.Int("figure").
			Comment("金额"),
	}
}

// Edges of the GiveMoney.
func (GiveMoney) Edges() []ent.Edge {
	// GiveMoney 模型通常是独立的，没有与其他模型的关联，所以这里返回 nil
	return nil
}
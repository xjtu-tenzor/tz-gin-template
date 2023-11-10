package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseModel struct {
	ID        int64          `gorm:"type:BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey;comment:主键"`
	CreatedAt time.Time      `gorm:"type:DATETIME(3) NOT NULL;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"type:DATETIME(3) NOT NULL;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"type:DATETIME(3) NULL;index;comment:删除时间"`
}

type Fields json.RawMessage

func (n Fields) GormDataType() string {
	return "TEXT"
}

func (n Fields) GormValue(_ context.Context, _ *gorm.DB) clause.Expr {
	if len(n) == 0 {
		return clause.Expr{SQL: "?", Vars: []any{"null"}}
	}
	return clause.Expr{SQL: "?", Vars: []any{string(n)}}
}

func (n *Fields) Scan(value any) error {
	*n = []byte(fmt.Sprintf("%s", value))
	return nil
}

func (n Fields) MarshalJSON() ([]byte, error) {
	if len(n) == 0 {
		return []byte("null"), nil
	}
	return n, nil
}

func (n *Fields) UnmarshalJSON(data []byte) error {
	if n == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*n = append((*n)[0:0], data...)
	return nil
}

package entities

import "time"

type URLIDGenerator struct {
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement"`
	CurrentID uint      `gorm:"column:current_id;not null"`
	RangeEnd  uint      `gorm:"column:range_end;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime:nano;not null"`
}

func (URLIDGenerator) TableName() string {
	return "url_id"
}

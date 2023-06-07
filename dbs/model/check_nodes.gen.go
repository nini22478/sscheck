// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameCheckNode = "check_nodes"

// CheckNode mapped from table <check_nodes>
type CheckNode struct {
	ID           int64      `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true" json:"id"`
	Host         *string    `gorm:"column:host;type:varchar(50)" json:"host"`
	ListPath     *string    `gorm:"column:list_path;type:text" json:"list_path"`
	LimitWait    *int64     `gorm:"column:limit_wait;type:int;default:20" json:"limit_wait"`
	NodeType     *string    `gorm:"column:node_type;type:varchar(50)" json:"node_type"`
	ReqEncode    *string    `gorm:"column:req_encode;type:varchar(255)" json:"req_encode"`
	ReqEncodeKey *string    `gorm:"column:req_encode_key;type:varchar(255)" json:"req_encode_key"`
	CreatedAt    *time.Time `gorm:"column:created_at;type:datetime;unsigned:autoUpdateTime" json:"created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at;type:datetime;unsigned:autoUpdateTime" json:"updated_at"`
	IsShow       *int64     `gorm:"column:is_show;type:tinyint" json:"is_show"`
}

// TableName CheckNode's table name
func (*CheckNode) TableName() string {
	return TableNameCheckNode
}

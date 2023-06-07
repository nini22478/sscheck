// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

var (
	Q            = new(Query)
	CheckHistory *checkHistory
	CheckNode    *checkNode
	DocheckNode  *docheckNode
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
	CheckHistory = &Q.CheckHistory
	CheckNode = &Q.CheckNode
	DocheckNode = &Q.DocheckNode
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:           db,
		CheckHistory: newCheckHistory(db, opts...),
		CheckNode:    newCheckNode(db, opts...),
		DocheckNode:  newDocheckNode(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	CheckHistory checkHistory
	CheckNode    checkNode
	DocheckNode  docheckNode
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:           db,
		CheckHistory: q.CheckHistory.clone(db),
		CheckNode:    q.CheckNode.clone(db),
		DocheckNode:  q.DocheckNode.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:           db,
		CheckHistory: q.CheckHistory.replaceDB(db),
		CheckNode:    q.CheckNode.replaceDB(db),
		DocheckNode:  q.DocheckNode.replaceDB(db),
	}
}

type queryCtx struct {
	CheckHistory ICheckHistoryDo
	CheckNode    ICheckNodeDo
	DocheckNode  IDocheckNodeDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		CheckHistory: q.CheckHistory.WithContext(ctx),
		CheckNode:    q.CheckNode.WithContext(ctx),
		DocheckNode:  q.DocheckNode.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
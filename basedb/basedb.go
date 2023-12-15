package basedb

import (
	"context"
	"go/ast"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"reflect"
)

type Table interface {
	schema.Tabler
	IdValue() int
}
type IDB interface {
	DB(ctx context.Context, isWrite bool) *gorm.DB
}

type txContextKey struct{}

type IBaseStore[T any] interface {
	IDB
	Get(ctx context.Context, id int) (*T, error)
	Save(ctx context.Context, t *T) error
	UpdateByUnique(ctx context.Context, t *T, uniques []string) error
	Delete(ctx context.Context, id int) (int, error)
	UpdateWhere(ctx context.Context, t *T, m map[string]interface{}) (int, error)
	Update(ctx context.Context, t *T) (int, error)
	DeleteWhere(ctx context.Context, m map[string]interface{}) (int, error)
	Insert(ctx context.Context, t ...*T) error
	List(ctx context.Context, m map[string]interface{}, order string) ([]*T, error)
	ListQuery(ctx context.Context, where string, args []interface{}, order string) ([]*T, error)
	First(ctx context.Context, m map[string]interface{}) (*T, error)
	FirstQuery(ctx context.Context, where string, args []interface{}, order string) (*T, error)
	ListPage(ctx context.Context, where string, pageNum, pageSize int, args []interface{}, order string) ([]*T, int, error)
	Transaction(ctx context.Context, f func(txCtx context.Context) error) error
}

func CreateStore[T any](db IDB) *BaseStore[T] {
	b := &BaseStore[T]{
		IDB:        db,
		TargetType: new(T),
	}
	modelType := reflect.TypeOf(new(T)).Elem()
	for i := 0; i < modelType.NumField(); i++ {
		if fieldStruct := modelType.Field(i); ast.IsExported(fieldStruct.Name) {
			tagSetting := schema.ParseTagSetting(fieldStruct.Tag.Get("gorm"), ";")
			if _, ok := tagSetting["DBUNIQUEINDEX"]; ok {
				b.UniqueList = append(b.UniqueList, tagSetting["COLUMN"])
			}
		}
	}

	return b
}

type ReadWriteDB struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func (d *ReadWriteDB) DB(ctx context.Context, isWrite bool) *gorm.DB {
	if tx, ok := ctx.Value(txContextKey{}).(*gorm.DB); ok {
		return tx
	}
	if isWrite {
		return d.writeDB.WithContext(ctx)
	}
	return d.readDB.WithContext(ctx)
}

func NewReadWriteDB(readDB, writeDB *gorm.DB) IDB {
	db := &ReadWriteDB{
		readDB:  readDB,
		writeDB: writeDB,
	}
	return db
}

type BaseStore[T any] struct {
	IDB
	UniqueList []string
	TargetType *T
}

func (b *BaseStore[T]) Get(ctx context.Context, id int) (*T, error) {
	value := new(T)
	err := b.DB(ctx, false).First(value, id).Error
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (b *BaseStore[T]) Save(ctx context.Context, t *T) error {

	var v interface{} = t
	if table, ok := v.(Table); ok {

		if table.IdValue() != 0 {
			return b.DB(ctx, true).Save(t).Error
		}
		//没查到主键ID的数据 看看有没有唯一索引 有唯一索引 用唯一索引更新所有字段
		if len(b.UniqueList) > 0 {
			return b.UpdateByUnique(ctx, t, b.UniqueList)
		}
	}
	return b.Insert(ctx, t)
}

func (b *BaseStore[T]) UpdateByUnique(ctx context.Context, t *T, uniques []string) error {
	columns := make([]clause.Column, 0, len(uniques))
	for _, unique := range uniques {
		columns = append(columns, clause.Column{
			Name: unique,
		})
	}
	return b.DB(ctx, true).Clauses(clause.OnConflict{
		Columns:   columns,
		UpdateAll: true,
	}).Create(t).Error
}

func (b *BaseStore[T]) Delete(ctx context.Context, id int) (int, error) {

	result := b.DB(ctx, true).Delete(b.TargetType, id)

	return int(result.RowsAffected), result.Error
}

func (b *BaseStore[T]) UpdateWhere(ctx context.Context, t *T, m map[string]interface{}) (int, error) {

	result := b.DB(ctx, true).Model(t).Updates(m)

	return int(result.RowsAffected), result.Error
}

func (b *BaseStore[T]) Update(ctx context.Context, t *T) (int, error) {

	result := b.DB(ctx, true).Updates(t)

	return int(result.RowsAffected), result.Error
}
func (b *BaseStore[T]) DeleteWhere(ctx context.Context, m map[string]interface{}) (int, error) {
	if len(m) == 0 {
		return 0, gorm.ErrMissingWhereClause
	}
	result := b.DB(ctx, true).Where(m).Delete(b.TargetType)

	return int(result.RowsAffected), result.Error
}

func (b *BaseStore[T]) Insert(ctx context.Context, t ...*T) error {

	return b.DB(ctx, true).Create(t).Error
}

func (b *BaseStore[T]) List(ctx context.Context, m map[string]interface{}, order string) ([]*T, error) {
	list := make([]*T, 0)

	db := b.DB(ctx, false).Where(m)
	if order != "" {
		db = db.Order(order)
	}

	err := db.Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}
func (b *BaseStore[T]) ListQuery(ctx context.Context, where string, args []interface{}, order string) ([]*T, error) {
	list := make([]*T, 0)
	db := b.DB(ctx, false)
	db = db.Where(where, args...)
	if order != "" {
		db = db.Order(order)
	}
	err := db.Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (b *BaseStore[T]) First(ctx context.Context, m map[string]interface{}) (*T, error) {
	value := new(T)
	db := b.DB(ctx, false)

	err := db.Where(m).First(value).Error
	if err != nil {
		return nil, err
	}

	return value, nil
}
func (b *BaseStore[T]) FirstQuery(ctx context.Context, where string, args []interface{}, order string) (*T, error) {
	value := new(T)
	db := b.DB(ctx, false)
	if order != "" {
		db = db.Order(order)
	}
	err := db.Where(where, args...).Take(value).Error
	if err != nil {
		return nil, err
	}

	return value, nil
}
func (b *BaseStore[T]) ListPage(ctx context.Context, where string, pageNum, pageSize int, args []interface{}, order string) ([]*T, int, error) {
	list := make([]*T, 0, pageSize)
	db := b.DB(ctx, false)
	if where != "" && len(args) > 0 {
		db = db.Where(where, args...)
	}
	if order != "" {
		db = db.Order(order)
	}
	count := int64(0)
	err := db.Model(list).Count(&count).Limit(pageSize).Offset(pageIndex(pageNum, pageSize)).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, int(count), nil
}

// Transaction 执行事务
func (b *BaseStore[T]) Transaction(ctx context.Context, f func(context.Context) error) error {
	return b.DB(ctx, true).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txContextKey{}, tx)
		return f(txCtx)
	})
}

// PageIndex 用于gorm  MYSQL 分页查询
// 示例如下
// db.Limit(pageSize).Offset(entry.PageIndex(pageNum,pageSize))
func pageIndex(pageNum, pageSize int) int {
	if pageNum == 0 {
		pageNum = 1
	}
	return (pageNum - 1) * pageSize
}

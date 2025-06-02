package orm

import (
	"github.com/rwpp/go-orm/orm/internal/errs"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

const (
	tagColumn = "column"
)

type Registry interface {
	Get(val any) (*Model, error)
	Register(val any, opts ...ModelOption) (*Model, error)
}
type Model struct {
	tableName string
	fields    map[string]*Field
}

type ModelOption func(m *Model) error
type Field struct {
	colName string
}

// registry元数据注册中心
type registry struct {
	//lock   sync.RWMutex
	models sync.Map
}

func newRegistry() *registry {
	return &registry{}
}

func (r *registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}
	m, err := r.Register(val)
	if err != nil {
		return nil, err
	}
	r.models.Store(typ, m)
	return m.(*Model), nil
}

//	func (r *registry) get1(val any) (*Model, error) {
//		typ := reflect.TypeOf(val)
//		r.lock.RLock()
//		m, ok := r.models[typ]
//		r.lock.RUnlock()
//		if ok {
//			return m, nil
//		}
//		r.lock.Lock()
//		defer r.lock.Unlock()
//		m, ok = r.models[typ]
//		if ok {
//			return m, nil
//		}
//		m, err := r.Register(val)
//		if err != nil {
//			return nil, err
//		}
//		r.models[typ] = m
//		return m, nil
//	}

func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag, ok := tag.Lookup("orm")
	if !ok {
		return map[string]string{}, nil
	}
	pairs := strings.Split(ormTag, ",")
	res := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		segs := strings.Split(pair, "=")
		if len(segs) != 2 {
			return nil, errs.NewErrInvalidTagContent(pair)
		}
		key := segs[0]
		val := segs[1]
		res[key] = val
	}
	return res, nil
}

func ModelWithTableName(tableName string) ModelOption {
	return func(m *Model) error {
		m.tableName = tableName
		return nil
	}
}
func ModelWithColumnName(field string, colName string) ModelOption {
	return func(m *Model) error {
		fd, ok := m.fields[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		fd.colName = colName
		return nil
	}
}
func (r *registry) Register(entity any, opts ...ModelOption) (*Model, error) {
	typ := reflect.TypeOf(entity)
	for typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrorPointerOnly
	}
	elemType := typ.Elem()
	numField := elemType.NumField()
	fieldMap := make(map[string]*Field, numField)
	for i := 0; i < numField; i++ {
		fd := elemType.Field(i)
		pair, err := r.parseTag(fd.Tag)
		if err != nil {
			return nil, err
		}
		colName := pair[tagColumn]
		if colName == "" {
			colName = underscoreName(fd.Name)
		}
		fieldMap[fd.Name] = &Field{
			colName: colName,
		}
	}
	var tableName string
	if tbl, ok := entity.(TableName); ok {
		tableName = tbl.TableName()
	}
	if tableName == "" {
		tableName = underscoreName(elemType.Name())
	}
	res := &Model{
		tableName: tableName,
		fields:    fieldMap,
	}
	for _, opt := range opts {
		opt(res)
	}
	r.models.Store(typ, res)
	return res, nil
}

func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}
	}
	return string(buf)
}

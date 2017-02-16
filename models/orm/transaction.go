package orm

import "github.com/astaxie/beego/orm"

const (
	INSERT int = iota
	UPDATE
	DELETE
)

/**
*  事物封装类型
 */
type Transaction struct {
	models []modeloper
}

type modeloper struct {
	model   interface{}
	operate int
}

func (t *Transaction) Add(model interface{}, operate int) {

	if t.models == nil {
		t.models = make([]modeloper, 0)
	}
	t.models = append(t.models, modeloper{model, operate})
}

func (t *Transaction) Clear() {
	t.models = nil
}

func (t *Transaction) Commit() error {
	o := orm.NewOrm()
	o.Begin()
	for _, md := range t.models {
		var err error
		switch md.operate {
		case 0:
			_, err = o.Insert(md.model)
		case 1:
			_, err = o.Update(md.model)
		case 2:
			_, err = o.Delete(md.model)
		}
		if err != nil {
			o.Rollback()
			return err
		}
	}

	return o.Commit()
}

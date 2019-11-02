package models

import (
	"time"

	"BeeCustom/enums"

	"github.com/astaxie/beego/orm"
)

type BaseModel struct {
	Id        int64
	CreatedAt time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt time.Time `orm:"column(updated_at);type(timestamp);null"`
}

// JsonResult 用于返回ajax请求的基类
type JsonResult struct {
	Status enums.JsonResultCode `json:"status"`
	Msg    string               `json:"msg"`
	Obj    interface{}          `json:"obj"`
}

// BaseQueryParam 用于查询的类
type BaseQueryParam struct {
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
}

type Page struct {
	PageNo     int
	PageSize   int
	TotalPage  int
	TotalCount int
	FirstPage  bool
	LastPage   bool
	List       interface{}
}

func PageUtil(count int, pageNo int, pageSize int, list interface{}) Page {
	tp := count / pageSize
	if count%pageSize > 0 {
		tp = count/pageSize + 1
	}
	return Page{PageNo: pageNo, PageSize: pageSize, TotalPage: tp, TotalCount: count, FirstPage: pageNo == 1, LastPage: pageNo == tp, List: list}
}

//默认列表数据
func BaseListQuery(query orm.QuerySeter, sort, order string, limit, offset int64) orm.QuerySeter {
	//默认排序
	sortorder := "Id"
	if len(sort) > 0 {
		sortorder = sort
	}

	if order == "desc" {
		sortorder = "-" + sortorder
	}

	query = query.OrderBy(sortorder)

	if limit != -1 {
		query = query.Limit(limit, (offset-1)*limit)
	}

	return query
}

//删除
func BaseDelete(m interface{}) (num int64, err error) {
	o := orm.NewOrm()
	if num, err := o.Delete(m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}

//删除
func BaseDeleteAll(clearanceType int8) (num int64, err error) {
	o := orm.NewOrm()
	if num, err := o.QueryTable(ClearanceTBName()).Filter("type", clearanceType).Delete(); err != nil {
		return num, err
	} else {
		return num, nil
	}
}

//批量插入
func BaseInsertMulti(bulk int, m interface{}) (num int64, err error) {
	o := orm.NewOrm()
	if num, err := o.InsertMulti(bulk, m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}

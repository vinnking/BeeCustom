package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

// TableName 设置Company表名
func (u *Company) TableName() string {
	return CompanyTBName()
}

// Company 实体类
type Company struct {
	BaseModel

	Number              string       `orm:"column(number);size(10)" description:"海关编号"`
	Name                string       `orm:"column(name);size(200)" description:"全称"`
	Short               string       `orm:"column(short);size(255);null" description:"简称"`
	Registration        string       `orm:"column(registration);size(10);null" description:"商检注册号"`
	Address             string       `orm:"column(address);size(200);null" description:"地址"`
	DeclareType         int8         `orm:"column(declare_type);null" description:"申报方式:代理，自理"`
	RegistrationCode    string       `orm:"column(registration_code);size(100);null" description:"产地证注册号"`
	Phone               string       `orm:"column(phone);size(20);null" description:"电话"`
	Fax                 string       `orm:"column(fax);size(20);null" description:"传真"`
	CreditCode          string       `orm:"column(credit_code);size(18)" description:"信用代码"`
	BusinessName        string       `orm:"column(business_name);size(100);null" description:"经营单位名称"`
	BusinessCode        string       `orm:"column(business_code);size(10);null" description:"经营单位代码"`
	Bank                string       `orm:"column(bank);size(50);null" description:"银行账号"`
	CustomCertification int8         `orm:"column(custom_certification);null" description:"海关认证：高级认证，一般认证，一般信用，失信企业"`
	CompanyType         int8         `orm:"column(company_type);null" description:"企业类别：往来客户，保税仓，供应商，代理报关公司，代理报检公司，物流公司"`
	CompanyKind         int8         `orm:"column(company_kind);null" description:"企业性质：国有，合作，合资， 独资， 集体，私营"`
	ControlRating       int8         `orm:"column(control_rating);null" description:"风控评级"`
	Remark              string       `orm:"column(remark);size(1000);null" description:"备注"`
	SubPhone            string       `orm:"column(sub_phone);null" description:"订阅手机 多个用，隔开 [is_open , data ]"`
	SubEmail            string       `orm:"column(sub_email);null" description:"订阅邮箱 多个用，隔开 [is_open , data ]"`
	SealDatas           string       `orm:"column(seal_datas);null" description:"公章 array[seal_file,seal_type] 多个用，隔开"`
	SubContentCheck     int8         `orm:"column(sub_content_check);null" description:"订阅内容 审核通过"`
	SubContentSubmit    int8         `orm:"column(sub_content_submit);null" description:"订阅内容 已提交海关处理"`
	SubContentReject    int8         `orm:"column(sub_content_reject);null" description:"订阅内容 驳回信息"`
	SubContentPass      int8         `orm:"column(sub_content_pass);null" description:"订阅内容 机关放行"`
	UserId              *BackendUser `orm:"column(user_id);rel(fk)"`
	StatementDate       int8         `orm:"column(statement_date)" description:"生成账单日期"`
	IsTrade             int8         `orm:"column(is_trade)" description:"境内收发货单位 是否开启"`
	IsOwner             int8         `orm:"column(is_owner)" description:"生产销售单位 是否开启"`
	Business            string       `orm:"column(business);size(255);null" description:"营业执照 file_path"`
	BusinessAuditStatus int8         `orm:"column(business_audit_status);null" description:"营业执照审核状态"`
	BusinessAuditAt     time.Time    `orm:"column(business_audit_at);type(datetime);null" description:"营业执照审核时间"`
	Tax                 int8         `orm:"column(tax)" description:"税率"`
}

// CompanyQueryParam 用于查询的类
type CompanyQueryParam struct {
	BaseQueryParam

	NameLike string //模糊查询

}

func NewCompany(id int64) Company {
	return Company{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//查询参数
func NewCompanyQueryParam() CompanyQueryParam {
	return CompanyQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// CompanyPageList 获取分页数据
func CompanyPageList(params *CompanyQueryParam) ([]*Company, int64) {
	query := orm.NewOrm().QueryTable(CompanyTBName())
	datas := make([]*Company, 0)

	if len(params.NameLike) > 0 {
		cond := orm.NewCondition()
		cond1 := cond.And("number__istartswith", params.NameLike).
			Or("name__istartswith", params.NameLike).
			Or("credit_code__istartswith", params.NameLike).
			Or("business_code__istartswith", params.NameLike)
		query = query.SetCond(cond1)
	}

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

// CompanyOne 根据id获取单条
func CompanyOne(id int64) (*Company, error) {
	m := NewCompany(0)
	o := orm.NewOrm()
	if err := o.QueryTable(CompanyTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("获取失败")
	}

	return &m, nil
}

//Save 添加、编辑页面 保存
func CompanySave(m *Company) (*Company, error) {
	o := orm.NewOrm()
	if m.Id == 0 {
		if _, err := o.Insert(m); err != nil {
			return nil, err
		}
	} else {
		if _, err := o.Update(m); err != nil {
			return nil, err
		}
	}

	return m, nil
}

//删除
func CompanyDelete(id int64) (num int64, err error) {
	m := NewCompany(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}
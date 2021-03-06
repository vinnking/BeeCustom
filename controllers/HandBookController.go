package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/transforms"
	"BeeCustom/utils"
	"BeeCustom/xlsx"
	gtf "github.com/snowlyg/gotransformer"
)

type HandBookController struct {
	BaseController
}

func (c *HandBookController) Prepare() {
	//  先执行
	c.BaseController.Prepare()
	//  如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	perms := []string{
		"Index",
		"Create",
		"Edit",
		"Delete",
	}
	c.checkAuthor(perms)

	//  如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//  权限控制里会进行登录验证，因此这里不用再作登录验证
	//  c.checkLogin()

}

func (c *HandBookController) Index() {

	params := models.NewCompanyQueryParam()
	limit, err := c.GetInt64("limit", 10)
	offset, err := c.GetInt64("offset", 1)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "关联关系获取失败", nil)
	}

	searchWord := c.GetString("searchWord", "")
	params.SearchWord = searchWord
	params.Limit = limit
	params.Offset = offset

	companies, count := models.CompanyPageList(&params)

	cs, err := models.CompaniesGetRelations(companies, "HandBooks")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "关联关系获取失败", nil)
	}
	// 页面模板设置
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "handbook/index_footerjs.html"
	c.Data["m"] = cs
	c.Data["count"] = count
	c.Data["searchWord"] = searchWord

	// 页面里按钮权限控制
	c.getActionData("", "Delete", "Import")

	c.GetXSRFToken()
}

// handbookgoods 列表数据
func (c *HandBookController) GoodDataGrid() {
	// 直接获取参数 GoodDataGrid()
	params := models.NewHandBookGoodQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	// 获取数据列表和总数
	data, total := models.HandBookGoodPageList(&params)
	c.ResponseList(c.TransformHandBookGoodsList(data), total)
	c.ServeJSON()
}

//  根据 handbookid 获取 handbookgoods
func (c *HandBookController) GetHandBookGoodByHandBookId() {
	params := models.NewHandBookGoodQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	data, _ := models.GetHandBookGoodById(&params)

	c.Data["json"] = c.TransformHandBookGood(data)
	c.ServeJSON()
}

// HandBook 列表数据
func (c *HandBookController) DataGrid() {
	params := models.NewHandBookQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	data, total := models.HandBookPageList(&params)
	c.ResponseList(data, total)
	c.ServeJSON()
}

//  Ullage 列表数据
func (c *HandBookController) UllageDataGrid() {
	params := models.NewHandBookUllageQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	data, total := models.HandBookUllagePageList(&params)
	c.ResponseList(c.TransformHandBookUllageList(data), total)
	c.ServeJSON()
}

//  Edit 添加 编辑 页面
func (c *HandBookController) Show() {
	Id, _ := c.GetInt64(":id", 0)
	m, err := models.HandBookOne(Id, "Company")
	if m != nil && Id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}

	var html, showFooterjs string
	handBookTypeS, err := c.getHandBookTypes()
	chandBookType, err := enums.TransformCnToInt(handBookTypeS, "手册")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("TransformCnToInt:%v", err), nil)
	}

	ahandBookType, err := enums.TransformCnToInt(handBookTypeS, "账册")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("TransformCnToInt:%v", err), nil)
	}

	if m.Type == chandBookType {
		html = "handbook/manual/show.html"
		showFooterjs = "handbook/manual/show_footerjs.html"
	} else if m.Type == ahandBookType {
		html = "handbook/account/show.html"
		showFooterjs = "handbook/account/show_footerjs.html"
	}

	c.setTpl(html)
	c.Data["m"] = c.TransformHandBook(m)
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = showFooterjs
	c.GetXSRFToken()
}

func (c *HandBookController) getHandBookTypes() (map[string]string, error) {
	return models.GetSettingRValueByKey("handBookType", false)
}

// 删除
func (c *HandBookController) Delete() {
	id, _ := c.GetInt64(":id")
	_, os, err := models.GetOrders(id)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "查询关联货物出错", err)
	}
	if os != 0 {
		c.jsonResult(enums.JRCodeFailed, "不能删除有关联货物的手账册", err)
	}

	_, as, err := models.GetAnnotations(id)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "查询关联清单出错", err)
	}
	if as != 0 {
		c.jsonResult(enums.JRCodeFailed, "不能删除有关联清单的手账册", err)
	}

	if num, err := models.HandBookDelete(id); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}
}

// 导入
func (c *HandBookController) Import() {
	importType, _ := c.GetInt8(":type")

	fileType := "handBook/" + strconv.FormatInt(c.curUser.Id, 10) + "/"
	fileNamePath, err := c.BaseUpload(fileType)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("BaseUpload:%v", err))
		c.jsonResult(enums.JRCodeFailed, "上传失败", nil)
	}

	fileNamePaths := strings.Split(fileNamePath, ".")
	fileExt := fileNamePaths[len(fileNamePaths)-1]
	if fileExt != "xlsx" {
		c.jsonResult(enums.JRCodeFailed, "文件格式错误，只能导入 xlsx 文件", nil)
	}

	var sheetName, sheet1Title, sheet2Name, sheet2Title, sheet3Name, sheet3Title, sheet4Name, sheet4Title string
	handBookTypeS, err := c.getHandBookTypes()
	importManualType, err := enums.TransformCnToInt(handBookTypeS, "手册")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", err)
	}
	importAccountType, err := enums.TransformCnToInt(handBookTypeS, "账册")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", err)
	}

	if importType == importManualType {
		sheetName = "handbookManualExcelSheetName"
		sheet1Title = "handbookManualExcelSheet1Title"
		sheet2Title = "handbookManualExcelSheet2Title"
		sheet3Title = "handbookManualExcelSheet3Title"
		sheet4Title = "handbookManualExcelSheet4Title"
	} else if importType == importAccountType {
		sheetName = "handbookAccountExcelSheetName"
		sheet1Title = "handbookAccountExcelSheet1Title"
		sheet2Title = "handbookAccountExcelSheet2Title"
		sheet3Title = "handbookAccountExcelSheet3Title"
		sheet4Title = "handbookAccountExcelSheet4Title"
	}

	handBookName, err := models.GetSettingRValueByKey(sheetName, false)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", err)
	}
	sheet2Name = handBookName["2"]
	sheet3Name = handBookName["3"]
	sheet4Name = handBookName["4"]

	handBook1Title, err := models.GetSettingRValueByKey(sheet1Title, false)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", err)
	}

	hIP := models.HandBookImport{
		BaseImport: xlsx.BaseImport{
			ExcelTitle:   handBook1Title,
			ExcelName:    handBookName["1"],
			FileNamePath: fileNamePath,
		},
		HandBook: models.NewHandBook(0),
	}

	hIP.HandBook.Type = importType

	c.ImportHandBookXlsxByCell(&hIP)

	m, err := models.HandBookSave(&hIP.HandBook)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", err)
	}
	hIP.HandBook = *m

	hBGIP := models.HandBookGoodImportParam{
		sheet2Name,
		sheet2Title,
		"料件",
	}

	c.InsertHandBookGoods(&hIP, &hBGIP)

	hBGIP = models.HandBookGoodImportParam{
		sheet3Name,
		sheet3Title,
		"成品",
	}

	c.InsertHandBookGoods(&hIP, &hBGIP)

	hBGIP = models.HandBookGoodImportParam{
		sheet4Name,
		sheet4Title,
		"",
	}

	c.InsertHandBookGoods(&hIP, &hBGIP)

	c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("导入成功"), m.Id)

}

// 导入账册表体
func (c *HandBookController) InsertHandBookGoods(hIP *models.HandBookImport, hBGIP *models.HandBookGoodImportParam) {
	handBookSheetTitle, err := models.GetSettingRValueByKey(hBGIP.ExcelTitleString, false)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", err)
	}

	hIP.ExcelName = hBGIP.ExcelName
	hIP.ExcelTitle = handBookSheetTitle

	c.ImportHandBookXlsxByRow(hIP, hBGIP.HandBookTypeString)

}

// 获取 xlsx 文件内容
func (c *HandBookController) InsertHandBookGood(hIP *models.HandBookImport, handBookGoods []*models.HandBookGood) error {

	num, err := models.InsertHandBookGoodMulti(handBookGoods)
	if err != nil {
		return err
	}

	if num == 0 {
		return errors.New("InsertHandBookGoodMulti error")
	}

	return nil
}

// 获取 xlsx 文件内容
func (c *HandBookController) InsertHandBookUllage(hIP *models.HandBookImport, handBookUllages []*models.HandBookUllage) error {
	num, err := models.InsertHandBookUllageMulti(handBookUllages)
	if err != nil {
		return err
	}

	if num == 0 {
		return errors.New("InsertHandBookGoodMulti:导入失败")
	}

	return nil
}

// 导入基础参数 xlsx 文件内容
func (c *HandBookController) ImportHandBookXlsxByCell(hIP *models.HandBookImport) {

	f, err := xlsx.GetExcel(hIP.FileNamePath)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("GetExcel:%v", err), nil)
	}
	x := gtf.NewXlxsTransform(&hIP.HandBook, hIP.ExcelTitle, nil, hIP.ExcelName, "", f)
	err = x.XlxsCellTransformer()
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("XlxsTransformer:%v", err), nil)
	}

	hB, err := models.GetHandBookByContractNumber(hIP.HandBook.ContractNumber)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", err)
	}

	if hB != nil && hB.Id != 0 {
		var errMsg string
		if hB.Type == 1 {
			errMsg = "手册已存在"
		} else if hB.Type == 2 {
			errMsg = "账册已存在"
		}
		c.jsonResult(enums.JRCodeFailed, errMsg, nil)
	}

	CompanyManageCode := hIP.HandBook.CompanyManageCode //  经营单位代码
	company, err := models.CompanyByManageCode(CompanyManageCode)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", err)
	}

	hIP.HandBook.Company = company

}

// 导入基础参数 xlsx 文件内容
func (c *HandBookController) ImportHandBookXlsxByRow(hIP *models.HandBookImport, handBookTypeString string) {
	rows, err := xlsx.GetExcelRows(hIP.FileNamePath, hIP.ExcelName)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", err)
	}

	if len(handBookTypeString) > 0 { // 表体
		hbs := make([]*models.HandBookGood, 0)
		handBookGoodTypeS, err := models.GetSettingRValueByKey("handBookGoodType", false)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("getHandBookGoodTypes:%v", err), nil)
		}
		handBookGoodType, err := enums.TransformCnToInt(handBookGoodTypeS, handBookTypeString)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("TransformCnToInt:%v", err), nil)
		}

		for roI, row := range rows {
			if roI > 1 { // 忽略标题和表头 2 行
				hb := models.NewHandBookGood(0)
				hb.Type = handBookGoodType
				x := gtf.NewXlxsTransform(&hb, hIP.ExcelTitle, row, "", "", nil)
				err := x.XlxsTransformer()
				if err != nil {
					c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("XlxsTransformer:%v", err), nil)
				}
				hb.HandBook = &hIP.HandBook
				hbs = append(hbs, &hb)
			}
		}

		err = c.InsertHandBookGood(hIP, hbs)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "导入失败", err)
		}

	} else { // 单损

		var hbus []*models.HandBookUllage
		for roI, row := range rows {
			if roI > 1 { // 忽略标题行
				hbu := models.NewHandBookUllage(0)
				x := gtf.NewXlxsTransform(&hbu, hIP.ExcelTitle, row, "", "", nil)
				err := x.XlxsTransformer()
				if err != nil {
					c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("XlxsTransformer:%v", err), nil)
				}
				hbg, err := models.GetHandBookGoodBySerial(hbu.FinishProNo)
				if err != nil {
					c.jsonResult(enums.JRCodeFailed, "导入失败", err)
				}

				hbu.HandBookGood = hbg
				hbus = append(hbus, &hbu)
			}
		}

		err = c.InsertHandBookUllage(hIP, hbus)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "导入失败", err)
		}
	}
}

//  TransformHandBookGoodsList 格式化列表数据
func (c *HandBookController) TransformHandBookGoodsList(ms []*models.HandBookGood) []*transforms.HandBookGood {
	var uts []*transforms.HandBookGood
	for _, v := range ms {
		ut := transforms.HandBookGood{}
		g := gtf.NewTransform(&ut, v, enums.BaseDateTimeFormat)
		_ = g.Transformer()

		uts = append(uts, &ut)
	}
	return uts
}

//  TransformHandBookGoodUllageList 格式化列表数据
func (c *HandBookController) TransformHandBookUllageList(ms []*models.HandBookUllage) []*transforms.HandBookUllage {
	var uts []*transforms.HandBookUllage
	for _, v := range ms {
		ut := transforms.HandBookUllage{}
		g := gtf.NewTransform(&ut, v, enums.BaseDateTimeFormat)
		_ = g.Transformer()

		uts = append(uts, &ut)
	}
	return uts
}

// TransformHandBookGood 格式化列表数据
func (c *HandBookController) TransformHandBookGood(v *models.HandBookGood) *transforms.HandBookGood {
	ut := &transforms.HandBookGood{}
	g := gtf.NewTransform(ut, v, enums.BaseDateTimeFormat)
	_ = g.Transformer()

	return ut
}

// TransformHandBookGood 格式化列表数据
func (c *HandBookController) TransformHandBook(v *models.HandBook) *transforms.HandBook {
	ut := &transforms.HandBook{}
	g := gtf.NewTransform(ut, v, enums.BaseDateTimeFormat)
	_ = g.Transformer()

	return ut
}

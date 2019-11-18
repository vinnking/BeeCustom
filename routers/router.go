package routers

import (
	"BeeCustom/controllers"
	"strings"

	"github.com/astaxie/beego"
)

func init() {

	//beego.Router("/rabc/test", &controllers.RABCController{}, "*:Test")
	////新建角色
	//beego.Router("/rabc/post", &controllers.RABCController{}, "post:Post")
	//beego.Router("/rabc/update", &controllers.RABCController{}, "put:Update")
	//beego.Router("/rabc/delete", &controllers.RABCController{}, "post:Delete")
	//beego.Router("/rabc/get/?:id:string", &controllers.RABCController{}, "get:Get")
	////添加用户角色
	//beego.Router("/rabc/userrole", &controllers.RABCController{}, "post:UserRole")
	////添加角色对项目目录文件操作权限
	////beego.Router("/rabc/permission", &controllers.RABCController{}, "post:RolePermission")
	////查询角色对项目目录文件操作的权限
	//beego.Router("/rabc/getpermission", &controllers.RABCController{}, "get:GetRolePermission")

	flags := [2]string{"I", "E"}
	/*清单管理*/
	for _, flag := range flags {
		lFkag := strings.ToLower(flag)
		//清单
		beego.Router("/annotation/index/"+lFkag, &controllers.AnnotationController{}, "Get:"+flag+"Index")
		//代客下单
		beego.Router("/annotation/create/"+lFkag, &controllers.AnnotationController{}, "Get:"+flag+"Create")
		//列表
		beego.Router("/annotation/datagrid/"+lFkag, &controllers.AnnotationController{}, "Post:"+flag+"DataGrid")
		//数量统计
		beego.Router("/annotation/statuscount/"+lFkag, &controllers.AnnotationController{}, "Post:"+flag+"StatusCount")
		//保存
		beego.Router("/annotation/store/"+lFkag, &controllers.AnnotationController{}, "Post:"+flag+"Store")
		//开始审单
		beego.Router("/annotation/edit/"+lFkag+"/?:id", &controllers.AnnotationController{}, "Get:"+flag+"Edit")
		//开始制单
		beego.Router("/annotation/make/"+lFkag+"/?:id", &controllers.AnnotationController{}, "Get:"+flag+"Make")
		//取消订单
		beego.Router("/annotation/cancel/"+lFkag+"/?:id", &controllers.AnnotationController{}, "Get:"+flag+"Cancel")
		//审核通过
		beego.Router("/annotation/audit/"+lFkag+"/?:id", &controllers.AnnotationController{}, "Get:"+flag+"Audit")
		//更新
		beego.Router("/annotation/update/"+lFkag+"/?:id", &controllers.AnnotationController{}, "Patch:"+flag+"Update")
		//派单
		beego.Router("/annotation/distribute/"+lFkag+"/?:id", &controllers.AnnotationController{}, "Post:"+flag+"Distribute")
		//删除
		beego.Router("/annotation/delete/"+lFkag+"/?:id", &controllers.AnnotationController{}, "Delete:"+flag+"Delete")
	}

	//手账册
	beego.Router("/handbook/index", &controllers.HandBookController{}, "*:Index")
	beego.Router("/handbook/show/?:id", &controllers.HandBookController{}, "Get:Show")
	beego.Router("/handbook/delete/?:id", &controllers.HandBookController{}, "Delete:Delete")
	beego.Router("/handbook/datagrid", &controllers.HandBookController{}, "Post:DataGrid")
	beego.Router("/handbook/gooddatagrid", &controllers.HandBookController{}, "Post:GoodDataGrid")
	beego.Router("/handbook/ullagedatagrid", &controllers.HandBookController{}, "Post:UllageDataGrid")
	beego.Router("/handbook/import/?:type", &controllers.HandBookController{}, "Post:Import")

	//客户关联公司管理
	beego.Router("/company_seal/create/?:cid", &controllers.CompanySealController{}, "Get:Create")
	beego.Router("/company_seal/store", &controllers.CompanySealController{}, "Post:Store")
	beego.Router("/company_seal/edit/?:id", &controllers.CompanySealController{}, "Get:Edit")
	beego.Router("/company_seal/update/?:id", &controllers.CompanySealController{}, "Patch:Update")
	beego.Router("/company_seal/delete/?:id", &controllers.CompanySealController{}, "Delete:Delete")

	//客户关联公司管理
	beego.Router("/company_foreign/create/?:cid", &controllers.CompanyForeignController{}, "Get:Create")
	beego.Router("/company_foreign/store", &controllers.CompanyForeignController{}, "Post:Store")
	beego.Router("/company_foreign/datagrid", &controllers.CompanyForeignController{}, "Post:DataGrid")
	beego.Router("/company_foreign/edit/?:id", &controllers.CompanyForeignController{}, "Get:Edit")
	beego.Router("/company_foreign/update/?:id", &controllers.CompanyForeignController{}, "Patch:Update")
	beego.Router("/company_foreign/delete/?:id", &controllers.CompanyForeignController{}, "Delete:Delete")

	//客户联系人管理
	beego.Router("/company_contact/create/?:cid", &controllers.CompanyContactController{}, "Get:Create")
	beego.Router("/company_contact/store", &controllers.CompanyContactController{}, "Post:Store")
	beego.Router("/company_contact/datagrid", &controllers.CompanyContactController{}, "Post:DataGrid")
	beego.Router("/company_contact/edit/?:id", &controllers.CompanyContactController{}, "Get:Edit")
	beego.Router("/company_contact/update/?:id", &controllers.CompanyContactController{}, "Patch:Update")
	beego.Router("/company_contact/delete/?:id", &controllers.CompanyContactController{}, "Delete:Delete")

	//客户管理
	beego.Router("/company/index", &controllers.CompanyController{}, "*:Index")
	beego.Router("/company/create/", &controllers.CompanyController{}, "Get:Create")
	beego.Router("/company/store", &controllers.CompanyController{}, "Post:Store")
	beego.Router("/company/datagrid", &controllers.CompanyController{}, "Post:DataGrid")
	beego.Router("/company/edit/?:id", &controllers.CompanyController{}, "Get:Edit")
	beego.Router("/company/update/?:id", &controllers.CompanyController{}, "Patch:Update")
	beego.Router("/company/delete/?:id", &controllers.CompanyController{}, "Delete:Delete")
	beego.Router("/company/import", &controllers.CompanyController{}, "Post:Import")

	//商品编码管理
	beego.Router("/hs_code/index", &controllers.HsCodeController{}, "*:Index")
	beego.Router("/hs_code/datagrid", &controllers.HsCodeController{}, "Post:DataGrid")
	beego.Router("/hs_code/import", &controllers.HsCodeController{}, "Post:Import")

	//商检编码管理
	beego.Router("/ciq/index", &controllers.CiqController{}, "*:Index")
	beego.Router("/ciq/datagrid", &controllers.CiqController{}, "Post:DataGrid")
	beego.Router("/ciq/import", &controllers.CiqController{}, "Post:Import")

	//基础参数
	beego.Router("/clearance/index", &controllers.ClearanceController{}, "*:Index")
	beego.Router("/clearance/get_clearance_update_time", &controllers.ClearanceController{}, "Get:GetClearanceUpdateTime")
	beego.Router("/clearance/last_update_time/?:type", &controllers.ClearanceController{}, "Get:GetClearanceUpdateTimeByType")
	beego.Router("/clearance/create/?:type", &controllers.ClearanceController{}, "Get:Create")
	beego.Router("/clearance/store", &controllers.ClearanceController{}, "Post:Store")
	beego.Router("/clearance/datagrid", &controllers.ClearanceController{}, "Post:DataGrid")
	beego.Router("/clearance/edit/?:id", &controllers.ClearanceController{}, "Get:Edit")
	beego.Router("/clearance/update/?:id", &controllers.ClearanceController{}, "Patch:Update")
	beego.Router("/clearance/delete/?:id", &controllers.ClearanceController{}, "Delete:Delete")
	beego.Router("/clearance/import/?:type", &controllers.ClearanceController{}, "Post:Import")

	//文件上传
	beego.Router("/file/upload", &controllers.FileController{}, "Post:Upload")
	beego.Router("/orderfile/upload/?:id", &controllers.FileController{}, "Post:OrderDataUpload")

	//后台用户路由
	beego.Router("/backenduser/index", &controllers.BackendUserController{}, "*:Index")
	beego.Router("/backenduser/create", &controllers.BackendUserController{}, "Get:Create")
	beego.Router("/backenduser/store", &controllers.BackendUserController{}, "Post:Store")
	beego.Router("/backenduser/datagrid", &controllers.BackendUserController{}, "Post:DataGrid")
	beego.Router("/backenduser/edit/?:id", &controllers.BackendUserController{}, "Get:Edit")
	beego.Router("/backenduser/freeze/?:id", &controllers.BackendUserController{}, "Get:Freeze")
	beego.Router("/backenduser/update/?:id", &controllers.BackendUserController{}, "Patch:Update")
	beego.Router("/backenduser/delete/?:id", &controllers.BackendUserController{}, "Delete:Delete")
	beego.Router("/backenduser/profile", &controllers.BackendUserController{}, "Get:Profile")

	//用户角色路由
	beego.Router("/role/index", &controllers.RoleController{}, "*:Index")
	beego.Router("/role/create", &controllers.RoleController{}, "Get:Create")
	beego.Router("/role/perm_lists/?:id", &controllers.RoleController{}, "Get:PermLists")
	beego.Router("/role/store", &controllers.RoleController{}, "Post:Store")
	beego.Router("/role/datagrid", &controllers.RoleController{}, "Post:DataGrid")
	beego.Router("/role/edit/?:id", &controllers.RoleController{}, "Get:Edit")
	beego.Router("/role/update/?:id", &controllers.RoleController{}, "Patch:Update")
	beego.Router("/role/delete/?:id", &controllers.RoleController{}, "Delete:Delete")
	beego.Router("/role/datalist", &controllers.RoleController{}, "Post:DataList")

	//资源路由
	beego.Router("/resource/index", &controllers.ResourceController{}, "*:Index")
	beego.Router("/resource/create", &controllers.ResourceController{}, "GET:Create")
	beego.Router("/resource/store", &controllers.ResourceController{}, "POST:Store")
	beego.Router("/resource/treegrid", &controllers.ResourceController{}, "POST:TreeGrid")
	beego.Router("/resource/edit/?:id", &controllers.ResourceController{}, "GET:Edit")
	beego.Router("/resource/update/?:id", &controllers.ResourceController{}, "PATCH:Update")
	beego.Router("/resource/delete/?:id", &controllers.ResourceController{}, "Delete:Delete")

	//通用选择面板
	//beego.Router("/resource/select", &controllers.ResourceController{}, "Get:Select")
	//用户有权管理的菜单列表（包括区域）
	//beego.Router("/resource/usermenutree", &controllers.ResourceController{}, "POST:UserMenuTree")
	//beego.Router("/resource/checkurlfor", &controllers.ResourceController{}, "POST:CheckUrlFor")

	beego.Router("/home/control", &controllers.HomeController{}, "*:Control")
	beego.Router("/home/login", &controllers.HomeController{}, "*:Login")
	beego.Router("/home/dologin", &controllers.HomeController{}, "Post:DoLogin")
	beego.Router("/home/logout", &controllers.HomeController{}, "*:Logout")
	beego.Router("/home/datareset", &controllers.HomeController{}, "Get:DataReset")

	beego.Router("/home/404", &controllers.HomeController{}, "*:Page404")
	beego.Router("/home/error/?:error", &controllers.HomeController{}, "*:Error")

	beego.Router("/", &controllers.HomeController{}, "*:Index")

}

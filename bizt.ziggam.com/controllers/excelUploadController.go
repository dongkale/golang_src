package controllers

// ExcelUploadController ...
type ExcelUploadController struct {
	BaseController
}

// Get ...
func (c *ExcelUploadController) Get() {

	c.Data["TMenuId"] = "M00"
	c.Data["SMenuId"] = "M00"

	c.TplName = "ldk_test/excel_upload.html"
}

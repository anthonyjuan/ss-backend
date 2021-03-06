package controllers

import (
	"encoding/json"
	"ss-backend/models"

	"github.com/astaxie/beego"
)

type (
	// PemesananController ...
	PemesananController struct {
		beego.Controller
	}
)

// Get all data product
func (c *PemesananController) Get() {
	var resp RespData
	var order models.Pemesanan
	var reqDt = models.RequestGet{
		FromDate: c.Ctx.Input.Query("fromDate"),
		ToDate:   c.Ctx.Input.Query("toDate"),
		Query:    c.Ctx.Input.Query("query"),
	}

	beego.Debug(reqDt)
	res, errGet := order.GetAll(reqDt)
	if errGet != nil {
		resp.Error = errGet
	} else {
		resp.Body = res
	}
	err := c.Ctx.Output.JSON(resp, false, false)
	if err != nil {
		panic("ERROR OUTPUT JSON LEVEL MIDDLEWARE")
	}
	// c.TplName = "index.tpl"
}

// Post add new order
func (c *PemesananController) Post() {
	var resp RespData
	var order models.Pemesanan

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &order)

	if err != nil {
		beego.Warning("unmarshall req body failed")
	}

	errAdd := order.AddPemesanan()

	if errAdd != nil {
		resp.Error = errAdd
	} else {
		resp.Body = order

	}
	err = c.Ctx.Output.JSON(resp, false, false)
	if err != nil {
		beego.Warning("failed giving output", err)
	}
	// c.TplName = "index.tpl"
}

// Put to update existing order
func (c *PemesananController) Put() {
	var resp RespData
	var order models.Pemesanan
	var req models.RequestUpdate

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)

	if err != nil {
		beego.Warning("unmarshall req body failed")
	}

	errAdd := order.UpdatePesanan(req)

	if errAdd != nil {
		resp.Error = errAdd
	} else {
		resp.Body = order

	}
	err = c.Ctx.Output.JSON(resp, false, false)
	if err != nil {
		beego.Warning("failed giving output", err)
	}
	// c.TplName = "index.tpl"
}

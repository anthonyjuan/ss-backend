package models

import (
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type (
	// Pemesanan is struct for barang masuk
	Pemesanan struct {
		ID             int64     `json:"id" orm:"column(id);auto"`
		SKU            string    `json:"sku" orm:"column(sku)"`
		NamaItem       string    `json:"nama_item" orm:"column(nama_item)"`
		NoKwitansi     string    `json:"no_kwitansi" orm:"column(no_kwitansi)"`
		JumlahPesanan  int64     `json:"jumlah_pesanan" orm:"column(jumlah_pesanan)"`
		JumlahDiterima int64     `json:"jumlah_diterima" orm:"column(jumlah_diterima)"`
		Harga          float64   `json:"harga" orm:"column(harga)"`
		Catatan        string    `json:"catatan" orm:"column(catatan)"`
		Waktu          time.Time `json:"waktu" orm:"column(waktu);auto_now_add;type(datetime)"`
		Total          float64   `json:"total" orm:"column(total)"`
		Status         string    `json:"status" orm:"column(status)"`
	}

	// RequestGet ...
	RequestGet struct {
		FromDate string
		ToDate   string
		Query    string
	}
)

// TableName return the table name
func (p *Pemesanan) TableName() string {
	return "pemesanan"
}

// GetAll record barang masuk ...
func (p *Pemesanan) GetAll(query RequestGet) ([]Pemesanan, error) {
	var brMasuk []Pemesanan
	o := orm.NewOrm()
	like := "%" + query.Query + "%"
	qb := []string{
		"SELECT *",
		"FROM",
		p.TableName(),
		"WHERE (waktu >= ? AND waktu <= ?)",
		"AND nama_item LIKE ?",
	}
	sql := strings.Join(qb, " ")

	count, err := o.Raw(sql, query.FromDate, query.ToDate,
		like).QueryRows(&brMasuk)
	if err != nil {
		beego.Warning("Failed get all data product", err)
		return []Pemesanan{}, err
	}
	beego.Debug("jumlah data = ", count)
	return brMasuk, nil
}

// AddPemesanan untuk memasukkan record barang masuk ...
func (p *Pemesanan) AddPemesanan() error {
	o := orm.NewOrm()
	// Insert
	if p.JumlahPesanan == p.JumlahDiterima {
		p.Status = "sukses"
	} else {
		p.Status = "pending"
	}

	p.Total = float64(p.JumlahPesanan) * p.Harga

	id, err := o.Insert(p)
	if err != nil {
		beego.Debug("error insert", err)
		return err
	}
	p.ID = id

	// Update Product
	var req = RequestUpdate{
		Jumlah:  p.JumlahDiterima,
		ID:      id,
		Catatan: p.Catatan,
	}
	errUpdateProd := updateProduct(req, p.SKU, "beli")
	if errUpdateProd != nil {
		beego.Warning("ERror update product", errUpdateProd)
		return errUpdateProd
	}
	// Update Product

	return nil
}

// UpdatePesanan to update pemesanana status etc
func (p *Pemesanan) UpdatePesanan(req RequestUpdate) error {
	// var resUpdate ResponseUpdatePemesanan
	o := orm.NewOrm()
	qb, errQB := orm.NewQueryBuilder("mysql")
	if errQB != nil {
		beego.Warning("Query builder failed")
		beego.Warning(errQB)
		return errQB
	}
	errGetOrder := getOneOrder(p, req.ID)
	if errGetOrder != nil {
		beego.Warning("Query get order failed")
		beego.Warning(errGetOrder)
		return errGetOrder
	}
	totalTerkini := p.JumlahDiterima + req.Jumlah
	var stat string
	if totalTerkini == p.JumlahPesanan {
		stat = "sukses"
	} else {
		stat = "pending"
	}
	newCatatan := p.Catatan + ";" + req.Catatan
	qb.Update(p.TableName()).Set(
		"jumlah_diterima = jumlah_diterima + ?",
		"catatan = ?",
		"status = ?",
	).Where(
		"id = ? ",
	)

	p.Status = stat
	p.JumlahDiterima = totalTerkini
	sqlForOrder := qb.String()

	beego.Debug(sqlForOrder)
	_, errSQLOrder := o.Raw(sqlForOrder, req.Jumlah, newCatatan,
		stat, req.ID).Exec()
	if errSQLOrder != nil {
		beego.Debug("error while updating product")
		beego.Debug(errSQLOrder)
		return errSQLOrder
	}

	// Update Product after update pemesanan selesai
	errUpdateProd := updateProduct(req, p.SKU, "beli")
	if errUpdateProd != nil {
		beego.Warning("ERror update product", errUpdateProd)
		return errUpdateProd
	}
	// Update Product

	return nil
}

// updateProduct that is used after update pemesanan done
func updateProduct(req RequestUpdate, sku string, tipe string) error {
	o := orm.NewOrm()
	var prod Product
	qb, errQB := orm.NewQueryBuilder("mysql")
	if errQB != nil {
		beego.Warning("Query builder failed")
		beego.Warning(errQB)
		return errQB
	}

	if tipe == "jual" {
		qb.Update(prod.TableName()).Set(
			"jumlah = jumlah - ?",
		).Where("sku = ?")
	} else if tipe == "beli" {
		qb.Update(prod.TableName()).Set(
			"jumlah = jumlah + ?",
		).Where("sku = ?")
	}

	sqlProd := qb.String()

	res, errSQL := o.Raw(sqlProd, req.Jumlah, sku).Exec()
	if errSQL != nil {
		beego.Debug("error while updating product")
		beego.Debug(errSQL)
		return errSQL
	}
	_, errRow := res.RowsAffected()
	if errRow != nil {
		beego.Debug("error get rows affected")
		beego.Debug(errRow)
		return errRow
	}

	return nil
}

func getOneOrder(p *Pemesanan, id int64) error {
	o := orm.NewOrm()

	qb := []string{
		"SELECT *",
		"FROM",
		p.TableName(),
		"WHERE id = ?",
	}
	sqlForGetOrder := strings.Join(qb, " ")
	errSQLGetOrder := o.Raw(sqlForGetOrder, id).QueryRow(p)

	if errSQLGetOrder != nil {
		beego.Debug("error while get Order")
		beego.Debug(errSQLGetOrder)
		return errSQLGetOrder
	}

	return nil

}

package models

import (
	"database/sql"
	_ "errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/config"

	"database/sql/driver"

	"io/ioutil"
	"reflect"

	"github.com/astaxie/beego/orm"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
)

type CMN_FILEINFO_TB struct {
	Filename       string    `orm:"pk;column(filename)"`
	Filesize       int64     `orm:"column(filesize)"`
	Fileext        string    `orm:"column(fileext)"`
	Filepath       string    `orm:"column(filepath);null"`
	Filerights     string    `orm:"column(filerights);null"`
	Expired        time.Time `orm:"column(expired);null"`
	Downloadstatus string    `orm:"column(downloadstatus);default(0)"`
	Createuser     string    `orm:"column(createuser);null"`
	Createtime     time.Time `orm:"column(createtime);null"`
	Updateuser     string    `orm:"column(updateuser);null"`
	Updatetime     time.Time `orm:"column(updatetime);null"`
}
type FILELIST struct {
	Uid      int
	Name     string
	Type     string
	Size     int64
	Status   string
	Response string
	Url      string
}
type OPTIONS struct {
	Value   string `json:"value"`
	Label   string `json:"label"`
	Checked bool   `json:"checked"`
}

func (u *CMN_FILEINFO_TB) TableName() string {
	return "cmn_fileinfo_tb"
}
func AddCMN_FILEINFO_TB(u *CMN_FILEINFO_TB) error {
	o := orm.NewOrm()
	err := o.Begin()
	_, err = o.Insert(u)

	if err != nil {
		//fmt.Println(err)
		err2 := o.Rollback()
		if err2 != nil {
			err = err2
		}
	} else {
		err = o.Commit()
	}
	return err
}
func DeleteCMN_FILEINFO_TB(u *CMN_FILEINFO_TB) error {
	o := orm.NewOrm()
	err := o.Begin()
	_, err = o.Delete(u)

	if err != nil {
		//fmt.Println(err)
		err2 := o.Rollback()
		if err2 != nil {
			err = err2
		}
	} else {
		err = o.Commit()
	}
	return err
}
func GetAllCMN_FILEINFO_TB() (admins []CMN_FILEINFO_TB, err error) {
	admins = make([]CMN_FILEINFO_TB, 0)
	o := orm.NewOrm()

	sql := "select filerights,expired,downloadstatus from cmn_fileinfo_tb "

	_, err = o.Raw(sql).QueryRows(&admins)

	return admins, err
}

func GetCMN_FILEINFO_TB(u *CMN_FILEINFO_TB) (admins []CMN_FILEINFO_TB, err error) {
	admins = make([]CMN_FILEINFO_TB, 0)
	o := orm.NewOrm()
	sql := "select * from cmn_fileinfo_tb where 1=1 "

	if u.Filerights != "" {
		sql = sql + " and filerights='" + u.Filerights + "'"
	}

	//if u.Expired != "" {
	//sql = sql + " and expired='" + u.Expired + "'"
	//}

	if u.Downloadstatus != "" {
		sql = sql + " and downloadstatus='" + u.Downloadstatus + "'"
	}

	_, err = o.Raw(sql).QueryRows(&admins)

	return admins, err
}

type CMN_TEMPLATE_TB struct {
	Templateid   string `orm:"pk;column(templateid)"`
	Templatename string `orm:"column(templatename)"`
	Formate      string `orm:"column(formate)"`
}

func (u *CMN_TEMPLATE_TB) TableName() string {
	return "cmn_template_tb"
}

type CMN_EXPORTTEMPLATE_TB struct {
	Exporttemplateid   string `orm:"pk;column(exporttemplateid)"`
	Exporttemplatename string `orm:"column(exporttemplatename)"`
	Templateid         string `orm:"column(templateid)"`
	Exporttitle        string `orm:"column(exporttitle)"`
	Exporttype         string `orm:"column(exporttype)"`
	Exportsql          string `orm:"column(exportsql)"`
	Exportfilepath     string `orm:"column(exportfilepath)"`
	Exportfilename     string `orm:"column(exportfilename)"`
	Accessmethod       string `orm:"column(accessmethod)"`
	Emailtitle         string `orm:"column(emailtitle)"`
}

func (u *CMN_EXPORTTEMPLATE_TB) TableName() string {
	return "cmn_exporttemplate_tb"
}

type CMN_TEMPLATEITEM_TB struct {
	Templateid   string `orm:"pk;column(templateid)"`
	Colid        string `orm:"pk;column(colid)"`
	Colname      string `orm:"column(colname)"`
	Coltype      string `orm:"column(coltype)"`
	Required     string `orm:"column(required)"`
	Length       string `orm:"column(length)"`
	Accuracy     string `orm:"column(accuracy)"`
	Defaultvalue string `orm:"column(defaultvalue)"`
	Pretype      string `orm:"column(pretype)"`
	Sep          string `orm:"column(sep)"`
}

func (u *CMN_TEMPLATEITEM_TB) TableName() string {
	return "cmn_templateitem_tb"
}
func AddCMN_TEMPLATE_TB(template1 *CMN_TEMPLATE_TB, templateitems []CMN_TEMPLATEITEM_TB) error {
	rows := 0
	db, _ := orm.GetDB("default")
	tr, _ := db.Begin()
	querysql := "select *  from cmn_template_tb where templateid=? "
	querysql = ConvertSQL(querysql, Getdbtype())
	result, err := tr.Query(querysql, template1.Templateid)
	if result.Next() {
		rows = 1
		result.Close()
	}

	fmt.Println(rows)
	if rows > 0 {
		deletesql := "delete  from cmn_template_tb where templateid=? "
		deletesql = ConvertSQL(deletesql, Getdbtype())
		_, err = tr.Exec(deletesql, template1.Templateid)

		if err != nil {
			fmt.Println("delete cmn_template_tb fail:==>")
			fmt.Println(err)
			err = tr.Rollback()
		}
	}

	sql := "insert into cmn_template_tb values(?,?,?) "
	sql = ConvertSQL(sql, Getdbtype())
	preparestatment, _ := tr.Prepare(sql)
	_, err = preparestatment.Exec(template1.Templateid, template1.Templatename, template1.Formate)
	if err != nil {
		fmt.Println("insert into cmn_template_tb values(?,?,?) ")
		fmt.Println(err)
		err2 := tr.Rollback()
		if err2 != nil {
			err = err2
		}
	}

	err = AddMultiCMN_TEMPLATEITEM_TB(tr, rows, template1.Templateid, templateitems)
	if err == nil {
		err = tr.Commit()
	}
	//defer db.Close()
	return err
}
func AddMultiCMN_TEMPLATEITEM_TB(tr *sql.Tx, rows int, templateid string, templateitems []CMN_TEMPLATEITEM_TB) error {
	var deletesql string = ""
	var insertsql string = ""

	//fuck有一个陷阱，不同的数据库写法不同，坑爹。
	//MySQL               PostgreSQL            Oracle
	//WHERE col = ?       WHERE col = $1        WHERE col = :col
	//VALUES(?, ?, ?)     VALUES($1, $2, $3)    VALUES(:val1, :val2, :val3)
	iniconf, err := config.NewConfig("ini", "conf/myconf.ini")
	if err != nil {
		fmt.Println(err)
	}
	dbtype := iniconf.String("dbtype")
	switch dbtype {
	case "mysql":
		deletesql = "delete  from cmn_templateitem_tb where templateid=? "
		insertsql = "insert into cmn_templateitem_tb(templateid,colid,colname,coltype,required,length,accuracy,defaultvalue,pretype,sep) values(?,?,?,?,?,?,?,?,?,?)"
		break
	case "postgres":
		//deletesql = "delete  from cmn_flowuser_tb where flowid=$1"
		insertsql = "insert into cmn_flowuser_tb(flowid,varyid,varyname,varyvalue,pretype,sep) values($1,$2,$3,$4,$5,$6)"
		break
	case "sqlite3":
		//deletesql = "delete  from cmn_flowuser_tb where flowid=?"
		insertsql = "insert into cmn_flowuser_tb(flowid,varyid,varyname,varyvalue,pretype,sep) values(?,?,?,?,?,?)"
		break
	case "oracle":
		//deletesql = "delete  from cmn_flowuser_tb where flowid=:val1"
		insertsql = "insert into cmn_flowuser_tb(flowid,varyid,varyname,varyvalue,pretype,sep) values(:val1,:val2,:val3,:val4,:val5,:val6)"
		break
	}

	if rows > 0 {
		_, err = tr.Exec(deletesql, templateid)

		if err != nil {
			fmt.Println("delete cmn_templateitem_tb fail:==>")
			fmt.Println(err)
			err = tr.Rollback()
		}
	}

	for _, templateitem := range templateitems {

		_, err = tr.Exec(insertsql, templateitem.Templateid, templateitem.Colid, templateitem.Colname, templateitem.Coltype, templateitem.Required, templateitem.Length, templateitem.Accuracy, templateitem.Defaultvalue, templateitem.Pretype, templateitem.Sep)
		if err != nil {

			err = tr.Rollback()
		}
	}

	return err
}
func GetCMN_TEMPLATEITEM_TB(templateid string) (templateitems []CMN_TEMPLATEITEM_TB, err error) {
	var sql string
	templateitems = make([]CMN_TEMPLATEITEM_TB, 0)
	o := orm.NewOrm()
	if templateid != "" {
		sql = "select * from cmn_templateitem_tb where templateid=? order by colid"

	} else {
		return nil, nil
	}
	sql = ConvertSQL(sql, Getdbtype())
	_, err = o.Raw(sql, templateid).QueryRows(&templateitems)

	return templateitems, err
}
func GetCMN_TEMPLATE_TB() (templates []CMN_TEMPLATE_TB, err error) {
	var sql string
	templates = make([]CMN_TEMPLATE_TB, 0)
	o := orm.NewOrm()

	sql = "select * from cmn_template_tb "

	_, err = o.Raw(sql).QueryRows(&templates)

	return templates, err
}
func GetCMN_TEMPLATE_TBbyid(templateid string) (template CMN_TEMPLATE_TB, err error) {
	var sql string

	o := orm.NewOrm()

	sql = "select * from cmn_template_tb where templateid=?"
	sql = ConvertSQL(sql, Getdbtype())
	err = o.Raw(sql, templateid).QueryRow(&template)

	return template, err
}
func DeleteCMN_TEMPLATE_TB(templateid string) error {

	db, _ := orm.GetDB("default")
	tr, _ := db.Begin()

	deletesql := "delete  from cmn_template_tb where templateid=? "
	deletesql = ConvertSQL(deletesql, Getdbtype())
	_, err := tr.Exec(deletesql, templateid)

	if err != nil {
		fmt.Println("delete cmn_template_tb fail:==>")
		fmt.Println(err)
		err = tr.Rollback()
	}
	deletesql = "delete  from cmn_templateitem_tb where templateid=? "
	deletesql = ConvertSQL(deletesql, Getdbtype())
	_, err = tr.Exec(deletesql, templateid)

	if err != nil {
		fmt.Println("delete cmn_templateitem_tb fail:==>")
		fmt.Println(err)
		err = tr.Rollback()
	}
	if err == nil {
		err = tr.Commit()
	}
	//defer db.Close()
	return err
}
func DeleteCMN_EXPORTTEMPLATE_TB(u *CMN_EXPORTTEMPLATE_TB) error {

	db, _ := orm.GetDB("default")
	tr, _ := db.Begin()

	deletesql := "delete  from cmn_exporttemplate_tb where exporttemplateid=? "
	deletesql = ConvertSQL(deletesql, Getdbtype())
	_, err := tr.Exec(deletesql, u.Exporttemplateid)

	if err != nil {
		fmt.Println("delete cmn_exporttemplate_tb fail:==>")
		fmt.Println(err)
		err = tr.Rollback()
	}

	if err == nil {
		err = tr.Commit()
	}
	//defer db.Close()
	return err
}
func AddCMN_EXPORTTEMPLATE_TB(u *CMN_EXPORTTEMPLATE_TB) error {
	o := orm.NewOrm()
	err := o.Begin()
	_, err = o.Delete(u)

	if err != nil {
		//fmt.Println(err)
		err2 := o.Rollback()
		if err2 != nil {
			err = err2
		}
	}

	_, err = o.Insert(u)

	if err != nil {
		//fmt.Println(err)
		err2 := o.Rollback()
		if err2 != nil {
			err = err2
		}
	} else {
		err = o.Commit()
	}
	return err
}
func DeleteCMN_EXPORTTEMPLATE_TB2(u *CMN_EXPORTTEMPLATE_TB) error {
	o := orm.NewOrm()
	err := o.Begin()
	_, err = o.Delete(u)

	if err != nil {
		//fmt.Println(err)
		err2 := o.Rollback()
		if err2 != nil {
			err = err2
		}
	} else {
		err = o.Commit()
	}
	return err
}
func GetCMN_EXPORTTEMPLATE_TB(Exporttemplateid string) (templates []CMN_EXPORTTEMPLATE_TB, err error) {
	var sql string
	templates = make([]CMN_EXPORTTEMPLATE_TB, 0)
	o := orm.NewOrm()
	sql = "select * from cmn_exporttemplate_tb "
	if Exporttemplateid != "" {
		sql = "select * from cmn_exporttemplate_tb where exporttemplateid='" + Exporttemplateid + "'"
	}

	_, err = o.Raw(sql).QueryRows(&templates)

	return templates, err
}

type CMN_IMPORTTEMPLATE_TB struct {
	Importtemplateid   string `orm:"pk;column(importtemplateid)"`
	Importtemplatename string `orm:"column(importtemplatename)"`
	Templateid         string `orm:"column(templateid)"`
	Importtable        string `orm:"column(importtable)"`
	Importtype         string `orm:"column(importtype)"`
	Importsql          string `orm:"column(importsql)"`
}

func (u *CMN_IMPORTTEMPLATE_TB) TableName() string {
	return "cmn_importtemplate_tb"
}

func GetCMN_IMPORTTEMPLATE_TB(Importtemplateid string) (templates []CMN_IMPORTTEMPLATE_TB, err error) {
	var sql string
	templates = make([]CMN_IMPORTTEMPLATE_TB, 0)
	o := orm.NewOrm()
	sql = "select * from cmn_importtemplate_tb "
	if Importtemplateid != "" {
		sql = "select * from cmn_importtemplate_tb where importtemplateid='" + Importtemplateid + "'"
	}

	_, err = o.Raw(sql).QueryRows(&templates)

	return templates, err
}
func DeleteCMN_IMPORTTEMPLATE_TB(u *CMN_IMPORTTEMPLATE_TB) error {

	db, _ := orm.GetDB("default")
	tr, _ := db.Begin()

	deletesql := "delete  from cmn_importtemplate_tb where importtemplateid=? "
	deletesql = ConvertSQL(deletesql, Getdbtype())
	_, err := tr.Exec(deletesql, u.Importtemplateid)

	if err != nil {
		fmt.Println("delete cmn_importtemplate_tb fail:==>")
		fmt.Println(err)
		err = tr.Rollback()
	}

	if err == nil {
		err = tr.Commit()
	}
	//defer db.Close()
	return err
}
func AddCMN_IMPORTTEMPLATE_TB(u *CMN_IMPORTTEMPLATE_TB) error {
	o := orm.NewOrm()
	err := o.Begin()
	_, err = o.Delete(u)

	if err != nil {
		//fmt.Println(err)
		err2 := o.Rollback()
		if err2 != nil {
			err = err2
		}
	}

	_, err = o.Insert(u)

	if err != nil {
		//fmt.Println(err)
		err2 := o.Rollback()
		if err2 != nil {
			err = err2
		}
	} else {
		err = o.Commit()
	}
	return err
}
func DeleteCMN_IMPORTTEMPLATE_TB2(u *CMN_IMPORTTEMPLATE_TB) error {
	o := orm.NewOrm()
	err := o.Begin()
	_, err = o.Delete(u)

	if err != nil {
		//fmt.Println(err)
		err2 := o.Rollback()
		if err2 != nil {
			err = err2
		}
	} else {
		err = o.Commit()
	}
	return err
}
func Getmetadata(tablename string) {
	//var cols interface{}

	db, _ := orm.GetDB("default")
	st := db.Stats()
	fmt.Println(st.OpenConnections)
	sql := "select * from " + tablename + " limit 0,1"
	fmt.Println(sql)
	rows, _ := db.Query(sql)

	cols, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
	}
	for _, col := range cols {
		fmt.Println(col)
	}

	var drv []driver.Value
	if rows.Next() {
		//drv = rows.GetLastcols()
		//fmt.Println(drv)
		for _, sv := range drv {
			//fmt.Println(sv)

			switch sv.(type) {
			case string:
				fmt.Println("string")
			case []byte:
				fmt.Println("string")
			case time.Time:
				fmt.Println("time")
			case bool:
				fmt.Println("bool")
			case float64:
				fmt.Println("float64")
			case int64:
				fmt.Println("int64")
			default:
				fmt.Println("string")

			}
			//			rv := reflect.ValueOf(sv)
			//			switch rv.Kind() {
			//			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			//				fmt.Println("Int")
			//			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			//				fmt.Println("Uint")
			//			case reflect.Float64:
			//				fmt.Println("Float64")
			//			case reflect.Float32:
			//				fmt.Println("Float32")
			//			case reflect.Bool:
			//				fmt.Println("Bool")
			//			}
		}
	}

}

func asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	rv := reflect.ValueOf(src)
	fmt.Println(rv.Kind())
	fmt.Println(rv.Type())
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	return fmt.Sprintf("%v", src)
}

type TABLEINF struct {
	Tables string `json:"text"`
}
type COLINF struct {
	Id           string `json:"id"`
	Name         string `json:"text"`
	Type         string `json:"coltype"`
	Length       string `json:"length"`
	Isnull       string `json:"isnull"`
	Isprimary    string `json:"isprimary"`
	Isautoinc    string `json:"isautoinc"`
	Defaultvalue string `json:"defaultvalue"`
	Comment      string `json:"comment"`
}

func GetTABLEINF() (tableinfs []TABLEINF, err error) {
	var sql string
	tableinfs = make([]TABLEINF, 0)
	o := orm.NewOrm()
	sql = "show tables "

	_, err = o.Raw(sql).QueryRows(&tableinfs)

	db, _ := orm.GetDB("default")
	rows, err := db.Query(sql)

	cols, err := rows.Columns()
	fmt.Println(cols)
	tableinfss := make([]string, 0)
	err = rows.Scan(&tableinfss)
	fmt.Println(tableinfss)
	//values := rows.GetLastcols()
	//fmt.Println(values)
	return tableinfs, err
}
func Convert2time(input string) time.Time {
	var cvtvalue time.Time
	var err error
	//月份 1,01,Jan,January
	//日　 2,02,_2
	//时　 3,03,15,PM,pm,AM,am
	//分　 4,04
	//秒　 5,05
	//年　 06,2006
	//周几 Mon,Monday
	//时区时差表示 -07,-0700,Z0700,Z07:00,-07:00,MST
	//时区字母缩写 MST
	if input == "" {
		input = "9999-01-01"
	}
	cvtvalue, err = time.Parse("2006-01-02", input)
	if err != nil {
		fmt.Println(err)
		cvtvalue = time.Now()
		cvtvalue.Format("2006-01-02")
	}

	return cvtvalue
}
func GetYYYYMMDDstring() string {

	//月份 1,01,Jan,January
	//日　 2,02,_2
	//时　 3,03,15,PM,pm,AM,am
	//分　 4,04
	//秒　 5,05
	//年　 06,2006
	//周几 Mon,Monday
	//时区时差表示 -07,-0700,Z0700,Z07:00,-07:00,MST
	//时区字母缩写 MST

	return time.Now().Format("2006-01-02")

}
func GetHHmmssstring() string {

	//月份 1,01,Jan,January
	//日　 2,02,_2
	//时　 3,03,15,PM,pm,AM,am
	//分　 4,04
	//秒　 5,05
	//年　 06,2006
	//周几 Mon,Monday
	//时区时差表示 -07,-0700,Z0700,Z07:00,-07:00,MST
	//时区字母缩写 MST

	return time.Now().Format("03:04:05")

}
func GetYMDtime() time.Time {
	var cvtvalue time.Time

	//月份 1,01,Jan,January
	//日　 2,02,_2
	//时　 3,03,15,PM,pm,AM,am
	//分　 4,04
	//秒　 5,05
	//年　 06,2006
	//周几 Mon,Monday
	//时区时差表示 -07,-0700,Z0700,Z07:00,-07:00,MST
	//时区字母缩写 MST

	cvtvalue = time.Now()
	cvtvalue.Format("2006-01-02")

	return cvtvalue
}
func GetYYYY() string {

	//月份 1,01,Jan,January
	//日　 2,02,_2
	//时　 3,03,15,PM,pm,AM,am
	//分　 4,04
	//秒　 5,05
	//年　 06,2006
	//周几 Mon,Monday
	//时区时差表示 -07,-0700,Z0700,Z07:00,-07:00,MST
	//时区字母缩写 MST

	cvtvalue := time.Now()
	return cvtvalue.Format("2006")

}
func GetMM() string {

	//月份 1,01,Jan,January
	//日　 2,02,_2
	//时　 3,03,15,PM,pm,AM,am
	//分　 4,04
	//秒　 5,05
	//年　 06,2006
	//周几 Mon,Monday
	//时区时差表示 -07,-0700,Z0700,Z07:00,-07:00,MST
	//时区字母缩写 MST

	cvtvalue := time.Now()
	return cvtvalue.Format("01")

}
func GetDD() string {

	//月份 1,01,Jan,January
	//日　 2,02,_2
	//时　 3,03,15,PM,pm,AM,am
	//分　 4,04
	//秒　 5,05
	//年　 06,2006
	//周几 Mon,Monday
	//时区时差表示 -07,-0700,Z0700,Z07:00,-07:00,MST
	//时区字母缩写 MST

	cvtvalue := time.Now()
	return cvtvalue.Format("02")

}
func GetYYYYMM() string {

	//月份 1,01,Jan,January
	//日　 2,02,_2
	//时　 3,03,15,PM,pm,AM,am
	//分　 4,04
	//秒　 5,05
	//年　 06,2006
	//周几 Mon,Monday
	//时区时差表示 -07,-0700,Z0700,Z07:00,-07:00,MST
	//时区字母缩写 MST

	cvtvalue := time.Now()
	return cvtvalue.Format("200601")

}

func GetYYYYMMDD() string {

	//月份 1,01,Jan,January
	//日　 2,02,_2
	//时　 3,03,15,PM,pm,AM,am
	//分　 4,04
	//秒　 5,05
	//年　 06,2006
	//周几 Mon,Monday
	//时区时差表示 -07,-0700,Z0700,Z07:00,-07:00,MST
	//时区字母缩写 MST

	cvtvalue := time.Now()
	return cvtvalue.Format("2006-01-02")

}

//返回20181118175858格式的数据，即2018年11月18日17点58分58秒
func GetYYYYMMDDHHMMSS() string {

	//月份 1,01,Jan,January
	//日　 2,02,_2
	//时　 3,03,15,PM,pm,AM,am
	//分　 4,04
	//秒　 5,05
	//年　 06,2006
	//周几 Mon,Monday
	//时区时差表示 -07,-0700,Z0700,Z07:00,-07:00,MST
	//时区字母缩写 MST

	cvtvalue := time.Now()
	return cvtvalue.Format("20060102030405")

}
func Convert2YYYYMMDD(times time.Time) string {

	//月份 1,01,Jan,January
	//日　 2,02,_2
	//时　 3,03,15,PM,pm,AM,am
	//分　 4,04
	//秒　 5,05
	//年　 06,2006
	//周几 Mon,Monday
	//时区时差表示 -07,-0700,Z0700,Z07:00,-07:00,MST
	//时区字母缩写 MST

	return times.Format("2006-01-02")

}
func Convert2YYYYMM(times time.Time) string {

	//月份 1,01,Jan,January
	//日　 2,02,_2
	//时　 3,03,15,PM,pm,AM,am
	//分　 4,04
	//秒　 5,05
	//年　 06,2006
	//周几 Mon,Monday
	//时区时差表示 -07,-0700,Z0700,Z07:00,-07:00,MST
	//时区字母缩写 MST

	return times.Format("2006-01")

}
func Convert2int64(input string) int64 {
	var cvtvalue int64
	var err error
	if input == "" {
		input = "-9999"
	}
	cvtvalue, err = strconv.ParseInt(input, 10, 64)
	if err != nil {
		cvtvalue = -9999
	}

	return cvtvalue
}
func Convert2float64(input string) float64 {
	var cvtvalue float64
	var err error

	cvtvalue, err = strconv.ParseFloat(input, 64)
	if err != nil {
		cvtvalue = 0.0
	}

	return cvtvalue
}
func Convert2bool(input string) bool {
	var cvtvalue bool
	var err error

	cvtvalue, err = strconv.ParseBool(input)
	if err != nil {
		cvtvalue = false
	}

	return cvtvalue
}
func Outputconvertleft(datatype string) string {
	switch datatype {
	case "time.Time":
		return "models.Convert2time("
	case "int64":
		return "models.Convert2int64("
	case "float64":
		return "models.Convert2float64("
	case "bool":
		return "models.Convert2bool("
	default:
		return ""
	}
}
func Outputconvertright(datatype string) string {
	switch datatype {
	case "string":
		return ""
	default:
		return ")"
	}
}
func Getexportfileinfomap(sql string) (exportfileinfmap []orm.Params, err error) {
	var expfmp []orm.Params
	expfmp = make([]orm.Params, 0)
	o := orm.NewOrm()
	_, err = o.Raw(sql).Values(&expfmp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return expfmp, nil
}

//case "mysql":
//deletesql = "delete  from cmn_flowaction_tb where flowid=? and taskid=?"
//insertsql = "insert into cmn_flowaction_tb(flowid,taskid,actionid,nexttaskid,backtotaskid,taskstatus,dispatcher) values(?,?,?,?,?,?,?)"

//case "postgres":
//deletesql = "delete  from cmn_flowaction_tb where flowid=$1 and taskid=$2"
//insertsql = "insert into cmn_flowaction_tb(flowid,taskid,actionid,nexttaskid,backtotaskid,taskstatus,dispatcher) values($1,$2,$3,$4,$5,$6,$7)"

//case "sqlite3":
//deletesql = "delete  from cmn_flowaction_tb where flowid=? and taskid=?"
//insertsql = "insert into cmn_flowaction_tb(flowid,taskid,actionid,nexttaskid,backtotaskid,taskstatus,dispatcher) values(?,?,?,?,?,?,?)"

//case "oracle":
//deletesql = "delete  from cmn_flowaction_tb where flowid=:val1 and taskid=:val2"
//insertsql = "insert into cmn_flowaction_tb(flowid,taskid,actionid,nexttaskid,backtotaskid,taskstatus,dispatcher) values(:val1,:val3,:val3,:val4,:val5,:val6,:val7)"

func ConvertSQL(sql string, databasetype string) string {
	//mysql:DATE_FORMAT(calltime,'%Y-%m-%d')
	//sqlite3:strftime('%Y-%m-%d',calltime)
	var cvtsql string
	symbol := "?"
	if databasetype == "oracle" {

		symbol = ":val"

		idx := strings.Index(sql, "?")
		if idx == -1 {
			return sql
		}
		r := strings.Split(sql, "?")
		length := len(r)
		if length == 1 {
			cvtsql = cvtsql + r[0]
			cvtsql = cvtsql + symbol + strconv.Itoa(1)
			return cvtsql
		}
		lastidx := strings.LastIndex(sql, "?")

		for i := 0; i < length; i++ {
			if i > 0 {
				cvtsql = cvtsql + symbol + strconv.Itoa(i)
			}
			cvtsql = cvtsql + r[i]
			if lastidx == len(sql)-1 && i == length-1 {
				cvtsql = cvtsql + symbol + strconv.Itoa(i)
			}
		}
	} else {
		if databasetype == "sqlite3" {
			reg := regexp.MustCompile(`DATE_FORMAT\(calltime,'\%Y-\%m-\%d'\)`)
			cvtsql = reg.ReplaceAllString(sql, "strftime('%Y-%m-%d',calltime)")
			reg = regexp.MustCompile(`DATE_FORMAT\(attdate,'\%Y-\%m-\%d'\)`)
			cvtsql = reg.ReplaceAllString(cvtsql, "strftime('%Y-%m-%d',attdate)")
			reg = regexp.MustCompile(`DATE_FORMAT\(attdate,'\%Y-\%m'\)`)
			cvtsql = reg.ReplaceAllString(cvtsql, "strftime('%Y-%m',attdate)")
			reg = regexp.MustCompile(`DATE_FORMAT\(calltime,'\%Y-\%m'\)`)
			cvtsql = reg.ReplaceAllString(cvtsql, "strftime('%Y-%m',calltime)")
			reg = regexp.MustCompile(`DATE_FORMAT\(flowstarttime,'\%Y-\%m-\%d'\)`)
			cvtsql = reg.ReplaceAllString(cvtsql, "strftime('%Y-%m-%d',flowstarttime)")
			reg = regexp.MustCompile(`DATE_FORMAT\(flowfinishtime,'\%Y-\%m-\%d'\)`)
			cvtsql = reg.ReplaceAllString(cvtsql, "strftime('%Y-%m-%d',flowfinishtime)")
		} else {
			cvtsql = sql
		}
	}

	return cvtsql
}
func Getdbtype() string {
	iniconf, err := config.NewConfig("ini", "conf/myconf.ini")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	dbtype := iniconf.String("dbtype")
	return dbtype
}
func SQLBRACKET2SPACE(sql string, databasetype string) string {
	var cvtsql string

	if databasetype == "postgres" || databasetype == "mysql" || databasetype == "sqlserver" {
		return sql
	}
	reg := regexp.MustCompile(`\(|\)`)
	cvtsql = reg.ReplaceAllString(sql, " ")
	return cvtsql
}

//文件转换为字符串
func Readfile2string(filePath string, charset string) (s string, err1 error) {

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	switch charset {
	case "GBK":
		decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(b)
		s = string(decodeBytes)
	case "TGBK":
		decodeBytes, _ := traditionalchinese.Big5.NewDecoder().Bytes(b)
		s = string(decodeBytes)
	default:
		s = string(b)

	}

	return s, nil
}

package example

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jsix/gof/db/orm"
	"github.com/jsix/gof/shell"
	"github.com/jsix/gof/web/form"
	"log"
	"testing"
)

var (
	connString = "root:@tcp(172.16.69.128:3306)/txmall?charset=utf8"
	genDir     = "gen/"
)

// 生成数据库所有的代码文件
func TestGenAll(t *testing.T) {
	// 初始化工具
	d := &orm.MySqlDialect{}
	tool := orm.NewTool(getDb(), d)
	// 设置变量
	modelPkg := "github.com/jsix/gof/example/gen/model"
	modelPkgName := "model"
	tool.Var(orm.V_ModelPkgName, modelPkgName)
	tool.Var(orm.V_ModelPkg, modelPkg)
	tool.Var(orm.V_ModelPkgIRepo, modelPkg)

	// 获取所有表格
	tables, err := tool.Tables("")
	if err == nil {
		fe := &form.Engine{}
		for _, tb := range tables {
			entityPath := genDir + modelPkgName + "/" + tb.Name + ".go"
			iRepPath := genDir + "repo/i_" + tb.Name + "_repo.go"
			repPath := genDir + "repo/" + tb.Name + "_repo.go"
			dslPath := genDir + "form/" + tb.Name + ".form"
			htmPath := genDir + "html/" + tb.Name + ".html"
			//生成实体
			str := tool.TableToGoStruct(tb)
			tool.SaveFile(str, entityPath)
			//生成仓储结构
			str = tool.TableToGoRepo(tb, true, modelPkgName+".")
			tool.SaveFile(str, repPath)
			//生成仓储接口
			str = tool.TableToGoIRepo(tb, true, modelPkgName+".")
			tool.SaveFile(str, iRepPath)
			//生成表单DSL
			f := fe.TableToForm(tb)
			err = fe.SaveDSL(f, dslPath)
			//生成表单
			if err == nil {
				_, err = fe.SaveHtmlForm(f, form.TDefaultFormHtml, htmPath)
			}
		}
		//格式化代码
		shell.Run("gofmt -w gen/")
	}
	if err != nil {
		t.Error(err)
	} else {
		t.Log("生成成功")
	}
}

func getDb() *sql.DB {
	db, err := sql.Open("mysql", connString)
	if err == nil {
		err = db.Ping()
	}

	if err != nil {
		defer db.Close()
		//如果异常，则显示并退出
		log.Fatalln("[ DBC][ MySQL] " + err.Error())
		return nil
	}
	return db
}
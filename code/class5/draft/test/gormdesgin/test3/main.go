package test3

import (
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getTenantID(ctx context.Context) (uint, error) {
	return 0, nil
}

func main() {
	db, _ := gorm.Open(mysql.Open("user:password@tcp(127.0.0.1:3306/hello"))
	// 根据 TenantID 过滤
	var setTenantScope = func(db *gorm.DB) {
		if tenantID, err := getTenantID(db.Statement.Context); err != nil {
			db.Where("tenant_id = ?", tenantID)
		} else {
			db.AddError(err)
		}
	}

	db.Callback().Query().Before("gorm:query").Register("set_tenant_scope", setTenantScope)
	db.Callback().Delete().Before("gorm:delete").Register("set_tenant_scope", setTenantScope)
	db.Callback().Update().Before("gorm:update").Register("set_tenant_scope", setTenantScope)

	// 设置 TenantID
	var setTenantID = func(db *gorm.DB) {
		tenantID, _ := getTenantID(db.Statement.Context)
		db.Statement.SetColumn("tenant_id", tenantID)
		// ...
	}

	db.Callback().Update().Before("gorm:create").Register("set_tenant_id", setTenantID)
}

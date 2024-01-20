package database

import (
	"distribution-system-be/constants"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"fmt"
	"strconv"
	"time"
)

// GetAllActiveMenu ...
func GetAllActiveMenu() ([]dbmodels.Menu, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var menus []dbmodels.Menu
	var err error

	err = db.Find(&menus, "status = ?", 1).Error

	fmt.Println("Menus => ", menus)

	if err != nil {
		return menus, err
	}
	return menus, nil
}

// GetMenuByRole ...
func GetMenuByRole(roleID int) ([]dto.RoleMenuDto, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var menus []dto.RoleMenuDto
	var err error

	// err = db.Find(&menus, "status = ?", 1).Error

	err = db.Raw(`
		select b.description as menuDescription, a.status, a.menu_id as menuId, b.parent_id as parentid
		from m_role_menu a
		left join m_menus b on a.menu_id = b.id
		where a.role_id = ?
		and b.status = 1
		order by b.ordering
	`, roleID).Scan(&menus).Error

	fmt.Println("Menus => ", menus)

	if err != nil {
		return menus, err
	}
	return menus, nil
}

// SaveMenuByRole ...
func SaveMenuByRole(roleID int, menuIds []int) models.Response {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var res models.Response

	if len(menuIds) == 0 {
		res.ErrCode = "05"
		res.ErrDesc = "Error, Menu Id empty"
		return res
	}

	var roleIDStr string
	roleIDStr = strconv.Itoa(roleID)
	query := "DELETE FROM m_role_menu " +
		" WHERE role_id = " + roleIDStr + ";"
	dbres := db.Exec(query)
	if dbres.Error != nil {
		res.ErrCode = "05"
		res.ErrDesc = "Error, Role Id not found"
		return res
	}

	errCount := false
	for _, id := range menuIds {
		saveData := dbmodels.RoleMenu{}
		saveData.RoleID = roleID
		saveData.MenuID = id
		saveData.Status = 1
		saveData.LastUpdateBy = dto.CurrUser
		saveData.LastUpdate = time.Now()
		dbres := db.Save(&saveData)
		if dbres.Error != nil {
			errCount = true
			break
		}
	}

	if errCount {
		var roleIDStr string
		roleIDStr = strconv.Itoa(roleID)
		query := "DELETE FROM m_role_menu " +
			" WHERE role_id = " + roleIDStr + ";"
		dbres := db.Exec(query)
		if dbres.Error != nil {
			res.ErrCode = "05"
			res.ErrDesc = "Error, Role Id not found"
			return res
		}
		res.ErrCode = "05"
		res.ErrDesc = "Error save menu"
		return res
	}

	res.ErrCode = "00"
	res.ErrDesc = "Save Success"
	return res
}

// CreateRoleMenuByRoleId ...
func CreateRoleMenuByRoleId(roleID int64) error {
	db := GetDbCon()
	db.Debug().LogMode(true)

	sqlStat := fmt.Sprintf(" insert into m_role_menu "+
		" (role_id, menu_id, status, last_update_by ,last_update) "+
		" select %v,m_menus.id, 1, '%v', current_timestamp from m_menus "+
		" ", roleID, dto.CurrUser)

	fmt.Println(roleID, "  =>Sql ====> ", sqlStat)
	err := db.Exec(sqlStat).Error

	return err
}

// CreateRoleMenuByRoleId ...
func UpdateMenuByRole(req dto.RequestUpdateRole) dto.NoContentResponse {
	db := GetDbCon()
	db.Debug().LogMode(true)
	res := dto.NoContentResponse{}

	var roleMenu dbmodels.RoleMenu

	err := db.Where("role_id = ? and menu_id = ?  ", req.RoleID, req.MenuID).First(&roleMenu).Error

	// err := db.Where("id = ?  ", 103).First(&roleMenu).Error

	if err != nil {
		res.ErrCode = constants.ERR_CODE_51
		res.ErrDesc = constants.ERR_CODE_51_MSG
		return res
	}

	fmt.Println("update data ", roleMenu)
	// roleMenu.Status = req.Status
	// roleMenu.LastUpdateBy = dto.CurrUser
	// roleMenu.LastUpdate = time.Now()

	// r := db.Save(&roleMenu)
	r := db.Model(&roleMenu).Where("role_id = ? and menu_id = ?  ", req.RoleID, req.MenuID).Updates(map[string]interface{}{"status": req.Status, "last_update_by": dto.CurrUser, "last_update": time.Now()})

	if r.Error != nil {
		res.ErrCode = constants.ERR_CODE_80
		res.ErrDesc = constants.ERR_CODE_80_MSG
		return res
	}

	// db.Model(&user).Updates(User{Name: "hello", Age: 18, Active: false})

	// r := db.Model(&dbmodels.RoleMenu{}).Where("role_id = ? and menu_id = ?  ", req.RoleID, req.MenuID).Updates(dbmodels.RoleMenu{
	// 	Status:       req.Status,
	// 	LastUpdateBy: dto.CurrUser,
	// 	LastUpdate:   time.Now(),
	// })

	// if r.Error != nil {
	// 	fmt.Println("Error save db", r.Error)
	// 	res.ErrCode = constants.ERR_CODE_80
	// 	res.ErrDesc = constants.ERR_CODE_80_MSG
	// 	return res
	// }

	res.ErrCode = constants.ERR_CODE_00
	res.ErrDesc = constants.ERR_CODE_00_MSG
	return res
}

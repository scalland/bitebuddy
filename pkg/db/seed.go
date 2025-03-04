package db

import (
	"fmt"
	"github.com/scalland/bitebuddy/pkg/utils"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var seedDBQueries = []string{
	// User types master table.
	`INSERT INTO user_types (user_type_id,user_type_name)
	VALUES(1, "__superadmin__"),(2,"__admin__"),(3,"__restaurant_owner__"),(4,"__customer__")
	ON DUPLICATE KEY UPDATE user_type_id=user_type_id;`,

	// Users table.
	`INSERT INTO users (user_id,email,mobile_number,user_type_id,is_active,created_at,last_login,last_accessed_from)
	VALUES
	   (1,"tech@scalland.com","+919899670183",1,true,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,'0.0.0.0'),
	   (2,"admin@example.com","+919999341745",1,true,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,'0.0.0.0'),
	   (3,"user@example.com","+919844629772",1,true,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,'0.0.0.0')
	ON DUPLICATE KEY UPDATE user_id=user_id;`,
}

func SeedDB(u *utils.Utils) error {
	db, dbErr := u.ConnectDB()
	if dbErr != nil {
		return fmt.Errorf("error connecting to database: %s", dbErr.Error())
	}

	defer db.Close()

	for _, query := range seedDBQueries {
		log.Printf("Executing DB seed queries: %s", query)
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("DB seed query error: %s", err.Error())
		}
	}
	return nil
}

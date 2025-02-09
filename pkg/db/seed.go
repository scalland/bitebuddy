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

	//// Restaurants table.
	//`CREATE TABLE IF NOT EXISTS restaurants (
	//	restaurant_id INT AUTO_INCREMENT PRIMARY KEY,
	//	name VARCHAR(255) NOT NULL,
	//	address TEXT NOT NULL,
	//	latitude DECIMAL(10,8) NOT NULL,
	//	longitude DECIMAL(11,8) NOT NULL,
	//	overall_rating DECIMAL(3,2) DEFAULT 0,
	//	price_for_two DECIMAL(10,2),
	//	image_url VARCHAR(255),
	//	discount_available BOOLEAN DEFAULT FALSE,
	//	alcohol_available BOOLEAN DEFAULT FALSE,
	//	portion_size_large BOOLEAN DEFAULT FALSE
	//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
	//
	//// Display types master table.
	//`CREATE TABLE IF NOT EXISTS display_types (
	//	display_type_id INT AUTO_INCREMENT PRIMARY KEY,
	//	display_type_name VARCHAR(50) NOT NULL UNIQUE
	//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
	//
	//// Metric types master table.
	//`CREATE TABLE IF NOT EXISTS metric_types (
	//	metric_type_id INT AUTO_INCREMENT PRIMARY KEY,
	//	metric_type_name VARCHAR(50) NOT NULL UNIQUE
	//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
	//
	//// Metrics table.
	//`CREATE TABLE IF NOT EXISTS metrics (
	//	metric_id INT AUTO_INCREMENT PRIMARY KEY,
	//	metric_name VARCHAR(255) NOT NULL,
	//	parent_metric_id INT DEFAULT NULL,
	//	is_sub_metric BOOLEAN DEFAULT FALSE,
	//	display_type_id INT NOT NULL,
	//	metric_type_id INT NOT NULL,
	//	FOREIGN KEY (parent_metric_id) REFERENCES metrics(metric_id) ON DELETE CASCADE,
	//	FOREIGN KEY (display_type_id) REFERENCES display_types(display_type_id) ON DELETE CASCADE,
	//	FOREIGN KEY (metric_type_id) REFERENCES metric_types(metric_type_id) ON DELETE CASCADE
	//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
	//
	//// Restaurant Metrics table.
	//`CREATE TABLE IF NOT EXISTS restaurant_metrics (
	//	restaurant_metric_id INT AUTO_INCREMENT PRIMARY KEY,
	//	restaurant_id INT NOT NULL,
	//	metric_id INT NOT NULL,
	//	average_score DECIMAL(3,2) DEFAULT 0,
	//	FOREIGN KEY (restaurant_id) REFERENCES restaurants(restaurant_id) ON DELETE CASCADE,
	//	FOREIGN KEY (metric_id) REFERENCES metrics(metric_id) ON DELETE CASCADE
	//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
	//
	//// Reviews table.
	//`CREATE TABLE IF NOT EXISTS reviews (
	//	review_id INT AUTO_INCREMENT PRIMARY KEY,
	//	restaurant_id INT NOT NULL,
	//	user_id INT NOT NULL,
	//	overall_score DECIMAL(3,2) NOT NULL,
	//	review_text TEXT,
	//	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	//	FOREIGN KEY (restaurant_id) REFERENCES restaurants(restaurant_id) ON DELETE CASCADE,
	//	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
	//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
	//
	//// Metric Reviews table.
	//`CREATE TABLE IF NOT EXISTS metric_reviews (
	//	metric_review_id INT AUTO_INCREMENT PRIMARY KEY,
	//	review_id INT NOT NULL,
	//	metric_id INT NOT NULL,
	//	score DECIMAL(3,2) NOT NULL,
	//	FOREIGN KEY (review_id) REFERENCES reviews(review_id) ON DELETE CASCADE,
	//	FOREIGN KEY (metric_id) REFERENCES metrics(metric_id) ON DELETE CASCADE
	//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
	//
	//// Filter types master table.
	//`CREATE TABLE IF NOT EXISTS filter_types (
	//	filter_type_id INT AUTO_INCREMENT PRIMARY KEY,
	//	filter_type_name VARCHAR(50) NOT NULL UNIQUE
	//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
	//
	//// Filters table.
	//`CREATE TABLE IF NOT EXISTS filters (
	//	filter_id INT AUTO_INCREMENT PRIMARY KEY,
	//	filter_type_id INT NOT NULL,
	//	filter_value VARCHAR(255),
	//	FOREIGN KEY (filter_type_id) REFERENCES filter_types(filter_type_id) ON DELETE CASCADE
	//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
	//
	//// OTP Requests table – note the generated column.
	//`CREATE TABLE IF NOT EXISTS otp_requests (
	//	otp_request_id INT AUTO_INCREMENT PRIMARY KEY,
	//	user_id INT NOT NULL,
	//	otp_code VARCHAR(10) NOT NULL,
	//	requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	//	delivery_method VARCHAR(10),
	//	valid_till INT UNSIGNED GENERATED ALWAYS AS (UNIX_TIMESTAMP(requested_at)+33) VIRTUAL,
	//	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
	//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
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

package db

import (
	"database/sql"
	"fmt"
	"github.com/scalland/bitebuddy/pkg/utils"
	"github.com/spf13/viper"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var migrationQueries = []string{
	// User types master table.
	`CREATE TABLE IF NOT EXISTS user_types (
		user_type_id INT AUTO_INCREMENT PRIMARY KEY,
		user_type_name VARCHAR(50) NOT NULL UNIQUE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,

	// Users table.
	`CREATE TABLE IF NOT EXISTS users (
		user_id INT AUTO_INCREMENT PRIMARY KEY,
		email VARCHAR(255) UNIQUE,
		mobile_number VARCHAR(20) UNIQUE,
		user_type_id INT NOT NULL,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		last_login TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		last_accessed_from VARCHAR(39) DEFAULT '',
		FOREIGN KEY (user_type_id) REFERENCES user_types(user_type_id) ON DELETE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,

	// Restaurants table.
	`CREATE TABLE IF NOT EXISTS restaurants (
		restaurant_id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		address TEXT NOT NULL,
		latitude DECIMAL(10,8) NOT NULL,
		longitude DECIMAL(11,8) NOT NULL,
		overall_rating DECIMAL(3,2) DEFAULT 0,
		price_for_two DECIMAL(10,2),
		image_url VARCHAR(255),
		discount_available BOOLEAN DEFAULT FALSE,
		alcohol_available BOOLEAN DEFAULT FALSE,
		portion_size_large BOOLEAN DEFAULT FALSE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,

	// Display types master table.
	`CREATE TABLE IF NOT EXISTS display_types (
		display_type_id INT AUTO_INCREMENT PRIMARY KEY,
		display_type_name VARCHAR(50) NOT NULL UNIQUE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,

	// Metric types master table.
	`CREATE TABLE IF NOT EXISTS metric_types (
		metric_type_id INT AUTO_INCREMENT PRIMARY KEY,
		metric_type_name VARCHAR(50) NOT NULL UNIQUE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,

	// Metrics table.
	`CREATE TABLE IF NOT EXISTS metrics (
		metric_id INT AUTO_INCREMENT PRIMARY KEY,
		metric_name VARCHAR(255) NOT NULL,
		parent_metric_id INT DEFAULT NULL,
		is_sub_metric BOOLEAN DEFAULT FALSE,
		display_type_id INT NOT NULL,
		metric_type_id INT NOT NULL,
		FOREIGN KEY (parent_metric_id) REFERENCES metrics(metric_id) ON DELETE CASCADE,
		FOREIGN KEY (display_type_id) REFERENCES display_types(display_type_id) ON DELETE CASCADE,
		FOREIGN KEY (metric_type_id) REFERENCES metric_types(metric_type_id) ON DELETE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,

	// Restaurant Metrics table.
	`CREATE TABLE IF NOT EXISTS restaurant_metrics (
		restaurant_metric_id INT AUTO_INCREMENT PRIMARY KEY,
		restaurant_id INT NOT NULL,
		metric_id INT NOT NULL,
		average_score DECIMAL(3,2) DEFAULT 0,
		FOREIGN KEY (restaurant_id) REFERENCES restaurants(restaurant_id) ON DELETE CASCADE,
		FOREIGN KEY (metric_id) REFERENCES metrics(metric_id) ON DELETE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,

	// Reviews table.
	`CREATE TABLE IF NOT EXISTS reviews (
		review_id INT AUTO_INCREMENT PRIMARY KEY,
		restaurant_id INT NOT NULL,
		user_id INT NOT NULL,
		overall_score DECIMAL(3,2) NOT NULL,
		review_text TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (restaurant_id) REFERENCES restaurants(restaurant_id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,

	// Metric Reviews table.
	`CREATE TABLE IF NOT EXISTS metric_reviews (
		metric_review_id INT AUTO_INCREMENT PRIMARY KEY,
		review_id INT NOT NULL,
		metric_id INT NOT NULL,
		score DECIMAL(3,2) NOT NULL,
		FOREIGN KEY (review_id) REFERENCES reviews(review_id) ON DELETE CASCADE,
		FOREIGN KEY (metric_id) REFERENCES metrics(metric_id) ON DELETE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,

	// Filter types master table.
	`CREATE TABLE IF NOT EXISTS filter_types (
		filter_type_id INT AUTO_INCREMENT PRIMARY KEY,
		filter_type_name VARCHAR(50) NOT NULL UNIQUE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,

	// Filters table.
	`CREATE TABLE IF NOT EXISTS filters (
		filter_id INT AUTO_INCREMENT PRIMARY KEY,
		filter_type_id INT NOT NULL,
		filter_value VARCHAR(255),
		FOREIGN KEY (filter_type_id) REFERENCES filter_types(filter_type_id) ON DELETE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,

	// OTP Requests table – note the generated column.
	`CREATE TABLE IF NOT EXISTS otp_requests (
		otp_request_id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		otp_code VARCHAR(24) NOT NULL,
		requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		delivery_method VARCHAR(10),
		valid_till INT UNSIGNED GENERATED ALWAYS AS (UNIX_TIMESTAMP(requested_at)+33) VIRTUAL,
		session_id VARCHAR(1024) NOT NULL,
		UNIQUE KEY(user_id,otp_code,session_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
}

func MigrateDB(u *utils.Utils, createDB bool) error {
	// Migration queries that need to be generated with dynamic values using fmt.Sprintf()
	var sprintFFD = []string{
		// Session Store for Gorilla Sessions in MySQL
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s\n(id INT NOT NULL AUTO_INCREMENT,\nsession_data LONGBLOB,\ncreated_on TIMESTAMP DEFAULT NOW(),\nmodified_on TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE CURRENT_TIMESTAMP,\nexpires_on TIMESTAMP DEFAULT NOW(),\nPRIMARY KEY(id))ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin", viper.GetString(
			"session_db_table")),
	}
	migrationQueries = append(migrationQueries, sprintFFD...)
	if createDB {
		log.Printf("User requested to create the database as well. Trying that now...")
		db, dbErr := u.ConnectSansDB()
		if dbErr != nil {
			return fmt.Errorf("error connecting to database: %s", dbErr.Error())
		}
		log.Printf("successfully connected to database server without a database")
		createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", viper.GetString("db_database"))
		log.Printf("Creating Database:\n%s\n", createDBQuery)
		_, err := db.Exec(createDBQuery)
		if err != nil {
			return fmt.Errorf("migration error: %s", err.Error())
		}
		dbErr = db.Close()
		if dbErr != nil {
			log.Printf("error closing DB connection after creating database: %s", dbErr.Error())
			log.Printf("ignoring last error and proceeding further...")
		}
	}

	db, dbErr := u.ConnectDB()
	if dbErr != nil {
		return fmt.Errorf("error connecting to database: %s", dbErr.Error())
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("error closing connection to database...")
		}
	}(db)

	for _, query := range migrationQueries {
		log.Printf("Executing migration query: %s", query)
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("migration error: %s", err)
		}
	}
	return nil
}

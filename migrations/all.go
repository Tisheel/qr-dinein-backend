package migrations

import "gofr.dev/pkg/gofr/migration"

func All() map[int64]migration.Migrate {
	return map[int64]migration.Migrate{
		1: createRestaurantsTable(),
		2: createCategoriesTable(),
		3: createProductsTable(),
		4: createStaffTable(),
		5: createOrdersTable(),
		6: createSettingsTable(),
		7: addPrepTimeAndEstimatedReadyAt(),
		8: createOrderRatingsTable(),
	}
}

func addPrepTimeAndEstimatedReadyAt() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(`ALTER TABLE products ADD COLUMN prep_time INT DEFAULT 15`)
			if err != nil {
				return err
			}

			_, err = d.SQL.Exec(`ALTER TABLE orders ADD COLUMN estimated_ready_at TIMESTAMP NULL`)
			return err
		},
	}
}

func createRestaurantsTable() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(`CREATE TABLE IF NOT EXISTS restaurants (
				id INT AUTO_INCREMENT PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				slug VARCHAR(255) NOT NULL UNIQUE,
				address TEXT,
				phone VARCHAR(20) DEFAULT '',
				logo VARCHAR(500) DEFAULT '',
				currency VARCHAR(10) DEFAULT 'INR',
				tax_rate DECIMAL(5,2) DEFAULT 0.00,
				active BOOLEAN DEFAULT TRUE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			)`)
			return err
		},
	}
}

func createCategoriesTable() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(`CREATE TABLE IF NOT EXISTS categories (
				id INT AUTO_INCREMENT PRIMARY KEY,
				restaurant_id INT NOT NULL,
				name VARCHAR(255) NOT NULL,
				` + "`order`" + ` INT DEFAULT 0,
				image VARCHAR(500) DEFAULT '',
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE,
				INDEX idx_categories_restaurant (restaurant_id)
			)`)
			return err
		},
	}
}

func createProductsTable() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(`CREATE TABLE IF NOT EXISTS products (
				id INT AUTO_INCREMENT PRIMARY KEY,
				restaurant_id INT NOT NULL,
				category_id INT NOT NULL,
				name VARCHAR(255) NOT NULL,
				description TEXT,
				price DECIMAL(10,2) NOT NULL,
				image VARCHAR(500) DEFAULT '',
				veg BOOLEAN DEFAULT TRUE,
				available BOOLEAN DEFAULT TRUE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE,
				FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
				INDEX idx_products_restaurant (restaurant_id),
				INDEX idx_products_category (restaurant_id, category_id)
			)`)
			return err
		},
	}
}

func createStaffTable() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(`CREATE TABLE IF NOT EXISTS staff (
				id INT AUTO_INCREMENT PRIMARY KEY,
				restaurant_id INT NOT NULL,
				username VARCHAR(255) NOT NULL,
    			pin VARCHAR(10) NOT NULL,
				role VARCHAR(50) NOT NULL DEFAULT 'chef',
				active BOOLEAN DEFAULT TRUE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE,
				INDEX idx_staff_restaurant (restaurant_id)
			)`)
			return err
		},
	}
}

func createOrdersTable() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(`CREATE TABLE IF NOT EXISTS orders (
				id INT AUTO_INCREMENT PRIMARY KEY,
				restaurant_id INT NOT NULL,
				table_number VARCHAR(50) DEFAULT NULL,
				customer_mobile VARCHAR(20) DEFAULT '',
				customer_name VARCHAR(255) DEFAULT '',
				items JSON NOT NULL,
				status VARCHAR(50) DEFAULT 'pending',
				special_instructions TEXT,
				total DECIMAL(10,2) NOT NULL DEFAULT 0.00,
				assigned_chef_id INT DEFAULT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE,
				FOREIGN KEY (assigned_chef_id) REFERENCES staff(id) ON DELETE SET NULL,
				INDEX idx_orders_restaurant (restaurant_id),
				INDEX idx_orders_status (restaurant_id, status)
			)`)
			return err
		},
	}
}

func createOrderRatingsTable() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(`CREATE TABLE IF NOT EXISTS order_ratings (
				id INT AUTO_INCREMENT PRIMARY KEY,
				order_id INT NOT NULL UNIQUE,
				restaurant_id INT NOT NULL,
				rating INT NOT NULL,
				comment TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
				FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE,
				INDEX idx_order_ratings_restaurant (restaurant_id)
			)`)
			return err
		},
	}
}

func createSettingsTable() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(`CREATE TABLE IF NOT EXISTS settings (
				id INT AUTO_INCREMENT PRIMARY KEY,
				restaurant_id INT NOT NULL,
				` + "`key`" + ` VARCHAR(255) NOT NULL,
				value TEXT,
				FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE,
				UNIQUE KEY unique_restaurant_key (restaurant_id, ` + "`key`" + `)
			)`)
			return err
		},
	}
}

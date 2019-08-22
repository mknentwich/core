package database

//Defines all create Statements that are necessary for the data model
const (
	createUsers      = "CREATE TABLE IF NOT EXISTS users (email TEXT PRIMARY KEY, admin INTEGER NOT NULL, name TEXT NOT NULL, password TEXT NOT NULL, lastChange INTEGER NOT NULL, lastLogin INTEGER NOT NULL);"
	createCategories = "CREATE TABLE IF NOT EXISTS categories (id INTEGER PRIMARY KEY, name TEXT NOT NULL, parentId INTEGER, FOREIGN KEY(parentId) REFERENCES categories(id));"
	createScores     = "CREATE TABLE IF NOT EXISTS scores (id INTEGER PRIMARY KEY, title TEXT NOT NULL, price REAL NOT NULL, difficulty INTEGER NOT NULL, categoryId INTEGER, FOREIGN KEY(categoryId) REFERENCES categories(id));"
	createAddresses  = "CREATE TABLE IF NOT EXISTS addresses (id INTEGER PRIMARY KEY, city TEXT NOT NULL, postCode TEXT NOT NULL, state TEXT NOT NULL, street TEXT NOT NULL, streetNumber TEXT NOT NULL);"
	createOrders     = "CREATE TABLE IF NOT EXISTS orders (id INTEGER PRIMARY KEY, billingAddressId INTEGER, company TEXT, date INTEGER NOT NULL, deliveryAddressId INTEGER NOT NULL, email TEXT NOT NULL, firstName TEXT NOT NULL, lastName TEXT NOT NULL, payed INTEGER NOT NULL, referenceCount TEXT NOT NULL, salutation TEXT, scoreAmount INTEGER NOT NULL, scoreId INTEGER NOT NULL, telephone TEXT, FOREIGN KEY(billingAddressId) REFERENCES addresses(id), FOREIGN KEY(deliveryAddressId) REFERENCES addresses(id), FOREIGN KEY(scoreId) REFERENCES scores(id));"
)

//contains all create statements for better handling
var createStatements = [...]string{createUsers, createCategories, createScores, createAddresses, createOrders}

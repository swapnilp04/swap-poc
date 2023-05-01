package models

func init() {
	migrateUser()
	migrateSession()
	migrateStudent()
	migrateBatch()
	migrateBatchStandard()
	migrateHostel()
	migrateHostelRoom()
	migrateHostelStudent()
	migrateStandard()
	migrateBatch()
	migrateHostelStudent()
	migrateBatchStandardStudent()
	migrateTransactionCategory()
	migrateTransaction()
}

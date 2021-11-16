package main

func (d *Database) RecentUser() string {
	var totalAttacks string

	d.db.QueryRow("SELECT username FROM users ORDER BY id DESC LIMIT 1").Scan(
		&totalAttacks,
	)

	return totalAttacks
}
func (d *Database) Datelol() int {
	var totalAttacks int

	d.db.QueryRow("SELECT last_paid FROM users ORDER BY id DESC LIMIT 1").Scan(
		&totalAttacks,
	)

	return totalAttacks
}

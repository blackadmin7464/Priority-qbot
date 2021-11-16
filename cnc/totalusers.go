package main

func (d *Database) getTotalUsers() int {
	var totalAttacks int

	d.db.QueryRow("SELECT COUNT(*) FROM `users`").Scan(
		&totalAttacks,
	)

	return totalAttacks
}

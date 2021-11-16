package main

func (d *Database) LastUserID() int {
	var totalAttacks int

	d.db.QueryRow("SELECT user_id FROM history ORDER BY id DESC LIMIT 1").Scan(
		&totalAttacks,
	)

	return totalAttacks
}

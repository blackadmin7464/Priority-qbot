package main

import (
	"time"
)

func (d *Database) getTotalAttacks() int {
	var totalAttacks int

	d.db.QueryRow("SELECT COUNT(*) FROM `history`").Scan(
		&totalAttacks,
	)

	return totalAttacks
}

func (d *Database) getTotalAttacksRunning() int {
	var totalAttacks int

	d.db.QueryRow("SELECT COUNT(*) FROM `history` WHERE `time_sent` + `duration` > ?", time.Now().Unix()).Scan(
		&totalAttacks,
	)

	return totalAttacks
}

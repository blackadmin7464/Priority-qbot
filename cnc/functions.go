package main

import (
	"crypto/rand"
	"encoding/base32"
	"log"
	"math/big"
	"strings"
)

func fillSpace(input string, length int, filler string) string {

var output = input
if len(input) > length {
output = output[:length]
}

for {
if len(output) >= length {
break
}

output = output + filler
}

return output
}

func splitNewline(input string, length int) string {

if len(input) < length-1 {
return input
}

var output string
for {
if len(input) > length {

output += input[:length] + "\r\n"
input = input[length:]

} else {
return output + input
}
}

return string(output)
}

func censorString(input string, censor string) string {

cut := float32(len(input)) * 0.65

section := input[:int(cut)]

for {
if len(section) >= len(input) {
break
}

section += censor
}

return section
}

func splitColor(input string) string {

if len(input) < 2 {
return input
}

startIndex := len(input) / 2

start := input[:startIndex]
end := input[startIndex:]

return "\033[00;95m" + start + "\033[00;91m" + end + "\033[00;97m"
}

func GenTOTPSecret() string {

data, err := cryptoRand(32)
if err != nil {
log.Println(err)
return strings.ReplaceAll(base32.StdEncoding.EncodeToString([]byte("0A9SF870A9SDUF09SDF234")), "=", "D")
}

return strings.ReplaceAll(base32.StdEncoding.EncodeToString([]byte(data)), "=", "D")
}

func cryptoRand(length int) (string, error) {
result := ""
for {
if len(result) >= length {
return result, nil
}

num, err := rand.Int(rand.Reader, big.NewInt(int64(127)))
if err != nil {
return "", err
}

n := num.Int64()

if n > 32 && n < 127 {
result += string(n)
}
}
}

func formatBool(input bool) string {

if input == false {
return "\x1b[31mFalse\x1b[0m"
}

return "\x1b[32mTrue\x1b[0m"
}

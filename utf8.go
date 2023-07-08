package main

import (
	"bytes"
	"fmt"
)

func encode(utf32Buf []rune) []byte {
	var utf8Buf []byte

	for _, r := range utf32Buf {
		switch {
		case r <= 0x7F:
			utf8Buf = append(utf8Buf, byte(r))
		case r <= 0x07FF:
			utf8Buf = append(utf8Buf, byte(0xC0|((r>>6)&0x1F)))
			utf8Buf = append(utf8Buf, byte(0x80|(r&0x3F)))
		case r <= 0xFFFF:
			utf8Buf = append(utf8Buf, byte(0xE0|((r>>12)&0x0F)))
			utf8Buf = append(utf8Buf, byte(0x80|((r>>6)&0x3F)))
			utf8Buf = append(utf8Buf, byte(0x80|(r&0x3F)))
		default:
			utf8Buf = append(utf8Buf, byte(0xF0|((r>>18)&0x07)))
			utf8Buf = append(utf8Buf, byte(0x80|((r>>12)&0x3F)))
			utf8Buf = append(utf8Buf, byte(0x80|((r>>6)&0x3F)))
			utf8Buf = append(utf8Buf, byte(0x80|(r&0x3F)))
		}
	}
	return utf8Buf
}

func decode(a []byte) []rune {
	var utf32Buf []rune

	for i := 0; i < len(a); i++ {
		switch {
		case a[i] < 0x80:
			utf32Buf = append(utf32Buf, rune(a[i]))
		case a[i]&0x000000E0 == 0xC0:
			utf32Buf = append(utf32Buf, rune(a[i]-0xC0)*0x40+rune(a[i+1]-0x80))
			i++
		case a[i]&0x000000F0 == 0xE0:
			utf32Buf = append(utf32Buf, rune(a[i]-0xE0)*0x1000+rune(a[i+1]-0x80)*0x40+rune(a[i+2]-0x80))
			i += 2
		default:
			utf32Buf = append(utf32Buf, rune(a[i]-0xF0)*0x40000+rune(a[i+1]-0x80)*0x1000+rune(a[i+2]-0x80)*0x40+rune(a[i+3]-0x80))
			i += 3
		}
	}
	return utf32Buf
}

func main() {
	words := [][]rune{
		[]rune("Hello, World!"),
		[]rune("Привет, мир!"),
		[]rune("こんにちは"),
		[]rune("你好"),
		[]rune("안녕하세요"),
	}

	for _, utf32 := range words {
		fmt.Printf("Original: %v\n", utf32)

		utf8Buf := encode(utf32)
		fmt.Printf("Encoded: %v\n", utf8Buf)

		utf32New := decode(utf8Buf)
		fmt.Printf("Decoded: %v\n", utf32New)

		success := bytes.Equal([]byte(string(utf32)), utf8Buf) && (fmt.Sprint(utf32) == fmt.Sprint(utf32New))

		if success {
			fmt.Println("Passed!")
		} else {
			fmt.Println("Failed!")
		}

		fmt.Println()
	}
}

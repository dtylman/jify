package main

func transliterate(hebrew string) string {
	word := ""
	for _, c := range hebrew {
		heb := c
		eng := string(c)
		if heb == 1488 { // א
			eng = "a"
		} else if heb == 1489 { // ב
			eng = "b"
		} else if heb == 1490 { // ג
			eng = "g"
		} else if heb == 1491 { // ד
			eng = "d"
		} else if heb == 1492 { // ה
			eng = "h"
		} else if heb == 1493 { // ו
			eng = "o"
		} else if heb == 1494 { // ז
			eng = "z"
		} else if heb == 1495 { // ח
			eng = "7"
		} else if heb == 1496 { // ט
			eng = "t"
		} else if heb == 1497 { // י
			eng = "y"
		} else if heb == 1499 { // כ
			eng = "kh"
		} else if heb == 1500 { // ל
			eng = "l"
		} else if heb == 1502 { // מ
			eng = "m"
		} else if heb == 1504 { // נ
			eng = "n"
		} else if heb == 1505 { // ס
			eng = "s"
		} else if heb == 1506 { // ע
			eng = "3"
		} else if heb == 1508 { // פ
			eng = "p"
		} else if heb == 1510 { // צ
			eng = "ts"
		} else if heb == 1511 { // ק
			eng = "q"
		} else if heb == 1512 { // ר
			eng = "r"
		} else if heb == 1513 { // ש
			eng = "sh"
		} else if heb == 1514 { // ת
			eng = "t"
		} else if heb == 1509 { // ץ
			eng = "ts"
		} else if heb == 1507 { // ף
			eng = "f"
		} else if heb == 1498 { // ך
			eng = "kh"
		} else if heb == 1501 { //mem
			eng = "m"
		}
		word += string(eng)
	}
	return word
}

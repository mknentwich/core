package utils

import "strings"

var safePendants = map[rune]string{
	'ä':  "ae",
	'ö':  "oe",
	'ü':  "ue",
	'ß':  "ss",
	' ':  "-",
	'/':  "-",
	'\\': "-",
	'\n': "-",
	'\r': "-",
	':':  "-",
	'.':  "",
	'?':  "",
	'!':  "",
	'&':  "und",
	',':  "",
	';':  "",
	'%':  "p",
	'#':  "",
	'$':  "dollar",
	'(':  "",
	')':  "",
	'[':  "",
	']':  "",
	'{':  "",
	'}':  "",
	'á':  "a",
	'č':  "c",
	'ď':  "d",
	'é':  "e",
	'ě':  "e",
	'í':  "i",
	'ň':  "n",
	'ó':  "o",
	'ř':  "r",
	'š':  "s",
	'ť':  "t",
	'ú':  "u",
	'ů':  "u",
	'ý':  "y",
	'ž':  "z"}

func SanitizePath(path string) string {
	path = strings.ToLower(path)
	runes := []rune(path)
	for i := len(runes); i >= 0; i-- {
		if val, ok := safePendants[runes[i]]; ok {
			path = path[:i] + val + path[i+1:]
		}
	}
	return path
}

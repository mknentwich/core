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
	pathRunes := []rune(path)
	for i := len(pathRunes) - 1; i >= 0; i-- {
		if val, ok := safePendants[pathRunes[i]]; ok {
			pathRunes = append(pathRunes[:i], append([]rune(val), pathRunes[i+1:]...)...)
		}
	}
	return string(pathRunes)
}

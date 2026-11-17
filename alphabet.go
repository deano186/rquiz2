package main

// The Russian (Cyrillic) alphabet: the -A table and the -a quizzes.

import (
	"bufio"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// letter is one letter of the alphabet. sounds lists the accepted quiz
// answers; the first is the one shown. lookalike, when set, warns about a
// letter that looks like something it isn't.
type letter struct {
	upper, lower string
	sounds       []string
	like         string
	lookalike    string
}

var alphabet = []letter{
	{"А", "а", []string{"a", "ah"}, `"a" in "father"`, ""},
	{"Б", "б", []string{"b"}, `"b" in "bed"`, ""},
	{"В", "в", []string{"v"}, `"v" in "van"`, "looks like English B, but sounds like V"},
	{"Г", "г", []string{"g"}, `"g" in "go"`, ""},
	{"Д", "д", []string{"d"}, `"d" in "dog"`, ""},
	{"Е", "е", []string{"ye", "yeh", "e"}, `"ye" in "yes"`, "looks like English E, but sounds like YE"},
	{"Ё", "ё", []string{"yo", "yoh"}, `"yo" in "yonder"; always stressed`, ""},
	{"Ж", "ж", []string{"zh"}, `"s" in "measure"`, ""},
	{"З", "з", []string{"z"}, `"z" in "zoo"`, "looks like the number 3, but sounds like Z"},
	{"И", "и", []string{"ee", "i"}, `"ee" in "meet"`, "looks like a backwards N, but sounds like EE"},
	{"Й", "й", []string{"y", "i", "j"}, `"y" in "boy"`, ""},
	{"К", "к", []string{"k"}, `"k" in "kite"`, ""},
	{"Л", "л", []string{"l"}, `"l" in "lamp"`, ""},
	{"М", "м", []string{"m"}, `"m" in "map"`, ""},
	{"Н", "н", []string{"n"}, `"n" in "net"`, "looks like English H, but sounds like N"},
	{"О", "о", []string{"o", "oh"}, `"o" in "more"; unstressed, "a" in "about"`, ""},
	{"П", "п", []string{"p"}, `"p" in "pen"`, ""},
	{"Р", "р", []string{"r"}, `rolled "r", as in Spanish "pero"`, "looks like English P, but sounds like R"},
	{"С", "с", []string{"s"}, `"s" in "sun"`, "looks like English C, but always sounds like S"},
	{"Т", "т", []string{"t"}, `"t" in "top"`, ""},
	{"У", "у", []string{"oo", "u"}, `"oo" in "food"`, "looks like English Y, but sounds like OO"},
	{"Ф", "ф", []string{"f"}, `"f" in "fun"`, ""},
	{"Х", "х", []string{"kh", "h"}, `"ch" in Scottish "loch"`, "looks like English X, but sounds like KH"},
	{"Ц", "ц", []string{"ts"}, `"ts" in "cats"`, ""},
	{"Ч", "ч", []string{"ch"}, `"ch" in "chair"`, ""},
	{"Ш", "ш", []string{"sh"}, `"sh" in "shop", but harder`, ""},
	{"Щ", "щ", []string{"shch", "sch", "shsh"}, `long, soft "sh", as in "fresh sheets"`, ""},
	{"Ъ", "ъ", []string{"hard sign", "hard"}, `hard sign: no sound; keeps letters apart`, ""},
	{"Ы", "ы", []string{"y", "i", "ih"}, `"i" in "bit", tongue pulled back`, ""},
	{"Ь", "ь", []string{"soft sign", "soft"}, `soft sign: no sound; softens the letter before`, "looks like a small b, but has no sound"},
	{"Э", "э", []string{"e", "eh"}, `"e" in "met"`, ""},
	{"Ю", "ю", []string{"yu", "you"}, `"u" in "use"`, ""},
	{"Я", "я", []string{"ya", "yah"}, `"ya" in "yard"`, "looks like a backwards R, but sounds like YA"},
}

// plain removes the stress marks from a word, leaving one rune per letter.
func plain(s string) string {
	return strings.ReplaceAll(s, "́", "") // U+0301 is the stress mark
}

// lowerCyrillic reports whether w is made only of small Russian letters,
// which leaves out names (Москва), abbreviations (США) and hyphenated words.
func lowerCyrillic(w string) bool {
	for _, r := range w {
		if !(r >= 'а' && r <= 'я' || r == 'ё') {
			return false
		}
	}
	return true
}

// example returns a common word that shows the letter: the most common
// word of three or more letters that starts with it, or failing that, the
// most common one that contains it.
func example(l letter) entry {
	var contains *entry
	for i, e := range words {
		w := plain(e.russian)
		if !lowerCyrillic(w) || len([]rune(w)) < 3 {
			continue
		}
		if strings.HasPrefix(w, l.lower) {
			return e
		}
		if contains == nil && strings.Contains(w, l.lower) {
			contains = &words[i]
		}
	}
	return *contains
}

// sound is the letter's sound for display; the hard and soft signs have
// none.
func (l letter) sound() string {
	if strings.Contains(l.sounds[0], " ") {
		return "—"
	}
	return l.sounds[0]
}

// shortMeaning is the first n English meanings of a word.
func shortMeaning(e entry, n int) string {
	m := strings.Split(meanings(e), ", ")
	if len(m) > n {
		m = m[:n]
	}
	return strings.Join(m, ", ")
}

func showAlphabet() {
	defer optionsBanner()
	width := 0
	for _, l := range alphabet {
		width = max(width, len(l.like))
	}
	fmt.Println("The Russian (Cyrillic) alphabet")
	fmt.Println()
	fmt.Printf("%-6s %-5s %-*s  %s\n", "Letter", "Sound", width, "As in", "Example")
	for _, l := range alphabet {
		e := example(l)
		fmt.Printf("%-6s %-5s %-*s  %s (%s)\n",
			l.upper+" "+l.lower, l.sound(), width, l.like, e.russian, shortMeaning(e, 1))
	}
	fmt.Println()
	fmt.Println("Watch out for letters that look like something else:")
	for _, l := range alphabet {
		if l.lookalike != "" {
			fmt.Printf("  %s %s %s\n", l.upper, l.lower, l.lookalike)
		}
	}
	fmt.Println()
	fmt.Println("Stress changes some sounds: unstressed О sounds like \"a\", and")
	fmt.Println("unstressed Е and Я sound closer to \"yi\". The word quiz marks the")
	fmt.Println("stressed vowel with an accent (´).")
}

// letterQuiz shows letters in random order and asks for their sounds.
func letterQuiz(reader *bufio.Reader, count int) {
	if count <= 0 || count > len(alphabet) {
		count = len(alphabet)
	}
	fmt.Println("=== Russian Alphabet Quiz ===")
	fmt.Println("Type the sound each letter makes, in English letters (e.g. zh).")
	fmt.Println("For Ъ and Ь, type 'hard sign' or 'soft sign'.")
	fmt.Println("Type 'q' to quit, or press Enter to skip a letter.")
	fmt.Println()

	score, asked := 0, 0
	for i, idx := range rand.Perm(len(alphabet))[:count] {
		l := alphabet[idx]
		fmt.Printf("[%d/%d]  %s %s\n> ", i+1, count, l.upper, l.lower)
		line, err := reader.ReadString('\n')
		answer := strings.Trim(normalize(line), `"'`)
		if err != nil && answer == "" {
			fmt.Println()
			break
		}
		if answer == "q" {
			break
		}
		asked++
		switch {
		case answer == "":
			fmt.Println("Skipped.")
		case hasSound(l, answer):
			score++
			fmt.Printf("Correct! Score: %d\n", score)
		default:
			fmt.Println("Wrong.")
			if l.lookalike != "" {
				fmt.Printf("%s %s %s.\n", l.upper, l.lower, l.lookalike)
			}
		}
		e := example(l)
		if l.sound() == "—" {
			fmt.Printf("%s %s is the %s.\n", l.upper, l.lower, l.like)
		} else {
			fmt.Printf("%s %s = %s (%s).\n", l.upper, l.lower, l.sound(), l.like)
		}
		fmt.Printf("Example: %s [%s] means %s.\n\n", e.russian, e.pronunciation, shortMeaning(e, 2))
	}
	finalScore(score, asked)
}

func hasSound(l letter, answer string) bool {
	for _, s := range l.sounds {
		if answer == s {
			return true
		}
	}
	return false
}

// orderWord picks a word for the alphabet-order quiz: a common word of up
// to 9 letters with the letter in it exactly once, so there is one right
// position. It picks at random from the 50 most common such words.
func orderWord(l letter) (entry, int) {
	type pick struct {
		e   entry
		pos int
	}
	var picks []pick
	for _, e := range words {
		w := []rune(plain(e.russian))
		if len(w) > 9 || !lowerCyrillic(string(w)) || strings.Count(string(w), l.lower) != 1 {
			continue
		}
		for p, r := range w {
			if string(r) == l.lower {
				picks = append(picks, pick{e, p + 1})
			}
		}
		if len(picks) == 50 {
			break
		}
	}
	p := picks[rand.Intn(len(picks))]
	return p.e, p.pos
}

func ordinal(n int) string {
	switch {
	case n%100 >= 11 && n%100 <= 13:
		return fmt.Sprintf("%dth", n)
	case n%10 == 1:
		return fmt.Sprintf("%dst", n)
	case n%10 == 2:
		return fmt.Sprintf("%dnd", n)
	case n%10 == 3:
		return fmt.Sprintf("%drd", n)
	}
	return fmt.Sprintf("%dth", n)
}

// orderQuiz goes through the alphabet in order. For letter N it shows a
// word containing the Nth letter and asks where in the word it is.
func orderQuiz(reader *bufio.Reader, count int) {
	if count <= 0 || count > len(alphabet) {
		count = len(alphabet)
	}
	fmt.Println("=== Russian Alphabet Order Quiz ===")
	fmt.Println("For each number, think of that letter of the alphabet, then type")
	fmt.Println("its position in the word shown. Type 'q' to quit, or Enter to skip.")
	fmt.Println()

	score, asked := 0, 0
	for i := 0; i < count; i++ {
		l := alphabet[i]
		e, pos := orderWord(l)
		w := plain(e.russian)
		fmt.Printf("[%d/%d] Find the %s letter of the alphabet in this word:\n", i+1, count, ordinal(i+1))
		fmt.Printf("  %s\n  %s\n> ", w, "123456789"[:len([]rune(w))])
		line, err := reader.ReadString('\n')
		answer := normalize(line)
		if err != nil && answer == "" {
			fmt.Println()
			break
		}
		if answer == "q" {
			break
		}
		asked++
		switch n, _ := strconv.Atoi(answer); {
		case answer == "":
			fmt.Println("Skipped.")
		case n == pos:
			score++
			fmt.Printf("Correct! Score: %d\n", score)
		default:
			fmt.Println("Wrong.")
		}
		fmt.Printf("The %s letter is %s %s (%s), at position %d.\n",
			ordinal(i+1), l.upper, l.lower, l.sounds[0], pos)
		fmt.Printf("%s [%s] means %s.\n", e.russian, e.pronunciation, shortMeaning(e, 2))
		var sofar []string
		for _, a := range alphabet[:i+1] {
			sofar = append(sofar, a.upper)
		}
		fmt.Printf("Alphabet so far: %s\n\n", strings.Join(sofar, " "))
	}
	finalScore(score, asked)
}

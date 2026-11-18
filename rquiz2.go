// rquiz2 - a Russian vocabulary quiz
// Version: see the VERSION constant below (rquiz2 -v shows it).
//
// Mostly written by AI (model ID claude-opus-5-5)
// Directed by and based on prior work by: deano (deano186@gmail.com)
//
// On November 18, 2026, Claude Opus 5.5 (claude-opus-5-5), in Claude Code,
// added the -v option and moved the version into the VERSION constant.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// VERSION is the program version, shown by rquiz2 -v: Ver.YYYYMMDDx,
// the date of the update and a letter that starts at a for each day.
// Update it here, and only here, with every change. AI names the model
// that made the latest change.
const (
	VERSION = "Ver.20261118d"
	AI      = "Claude Opus 5.5 (model ID claude-opus-5-5), in Claude Code"
)

func normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	return strings.Join(strings.Fields(s), " ")
}

// clean is normalize for typed meanings: it also turns anything but
// letters, digits, apostrophes and hyphens into spaces, so stray
// punctuation like "home." or "(house)" doesn't make an answer wrong.
func clean(s string) string {
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '\'' || r == '-' {
			return r
		}
		return ' '
	}, s)
	return normalize(s)
}

// maxParts is the most meanings an answer may have, as in rquiz; anything
// after the fourth comma is ignored. It is also the most meanings shown
// for a word, so the answers in words.go list the most important first.
const maxParts = 4

// shown returns the English meanings of a word, as displayed and as
// expected in a complete answer. It leaves out variants that only exist so
// more answers are accepted: "eat" when "to eat" is listed, digits like
// "1", and "youre" spellings of "you're". Every answer is still accepted.
func shown(e entry) []string {
	listed := map[string]bool{}
	for _, a := range e.answers {
		listed[a] = true
	}
	var out []string
	seen := map[string]bool{}
	for _, a := range e.answers {
		k := key(a)
		if seen[k] || listed["to "+a] {
			continue
		}
		if _, err := strconv.Atoi(a); err == nil {
			continue
		}
		seen[k] = true
		out = append(out, a)
		if len(out) == maxParts {
			break
		}
	}
	return out
}

func meanings(e entry) string {
	return strings.Join(shown(e), ", ")
}

// key is the form of a meaning used to tell whether two are the same:
// "to eat" and "eat" are one meaning, and so are "you're" and "youre".
func key(a string) string {
	return strings.ReplaceAll(strings.TrimPrefix(a, "to "), "'", "")
}

// accepted reports whether one typed meaning is right for the word.
// "to" before a verb is optional either way.
func accepted(part string, e entry) bool {
	for _, a := range e.answers {
		if part == a || "to "+part == a || strings.TrimPrefix(part, "to ") == a {
			return true
		}
	}
	return false
}

// parts splits an answer at its commas into up to maxParts meanings.
// Extra info written as --text, as rquiz showed it, is ignored.
func parts(answer string) []string {
	var out []string
	for i, p := range strings.Split(answer, ",") {
		if i == maxParts {
			break
		}
		p = clean(p)
		if p != "" && !strings.HasPrefix(p, "-") {
			out = append(out, p)
		}
	}
	return out
}

type result int

const (
	noAnswer result = iota
	wrong
	partial
	complete
)

// grade marks an answer. It is complete when it names every shown meaning
// (in any order) and nothing wrong; partial when at least one meaning is
// right; wrong when none is. It also returns the meanings that were wrong.
func grade(answer string, e entry) (result, []string) {
	ps := parts(answer)
	if len(ps) == 0 {
		return noAnswer, nil
	}
	need := map[string]bool{}
	for _, m := range shown(e) {
		need[key(m)] = true
	}
	covered := map[string]bool{}
	var bad []string
	for _, p := range ps {
		if !accepted(p, e) {
			bad = append(bad, p)
			continue
		}
		if need[key(p)] {
			covered[key(p)] = true
		}
	}
	switch {
	case len(bad) == len(ps):
		return wrong, bad
	case len(bad) == 0 && len(covered) == len(need):
		return complete, nil
	}
	return partial, bad
}

// lookalikes returns the other quiz words that a wrong meaning belongs to,
// e.g. "cat" when the question was собака returns ко́шка. The word list is
// in frequency order, so the first matches found are the most common
// words; at most maxHints of them are returned.
const maxHints = 2

func lookalikes(bad []string, current entry) []int {
	var found []int
	for i, e := range words {
		if e.russian == current.russian {
			continue
		}
		for _, p := range bad {
			if p != "to" && accepted(p, e) {
				found = append(found, i)
				break
			}
		}
		if len(found) == maxHints {
			break
		}
	}
	return found
}

// describe is how a word is shown in answers and lists:
// "[12] э́то: this, it, this is, that  [EH-tuh]".
func describe(i int) string {
	e := words[i]
	return fmt.Sprintf("[%d] %s: %s  [%s]", i+1, e.russian, meanings(e), e.pronunciation)
}

// cyrillicOrder sorts the word list into Russian alphabetical order for
// -C, ignoring stress marks and putting ё right after е.
func cyrillicOrder() {
	rank := map[rune]int{}
	for i, l := range alphabet {
		rank[[]rune(l.lower)[0]] = i + 1
	}
	sortKey := func(s string) []int {
		var k []int
		for _, r := range plain(strings.ToLower(s)) {
			k = append(k, rank[r])
		}
		return k
	}
	sort.SliceStable(words, func(a, b int) bool {
		ka, kb := sortKey(words[a].russian), sortKey(words[b].russian)
		for i := 0; i < len(ka) && i < len(kb); i++ {
			if ka[i] != kb[i] {
				return ka[i] < kb[i]
			}
		}
		return len(ka) < len(kb)
	})
}

// find prints the words that have the given English word in their
// meanings, for -f. It searches every accepted answer, so when the match
// is one that isn't shown (past the first 4), it adds it as "(also: ...)".
func find(what string) {
	what = clean(what)
	n := 0
	for i, e := range words {
		inShown := map[string]bool{}
		for _, m := range shown(e) {
			inShown[key(m)] = true
		}
		found := false
		var also []string
		for _, a := range e.answers {
			if strings.Contains(" "+a+" ", " "+what+" ") {
				found = true
				if !inShown[key(a)] {
					also = append(also, a)
					inShown[key(a)] = true
				}
			}
		}
		if !found {
			continue
		}
		line := "Found " + describe(i)
		if len(also) > 0 {
			line += "  (also: " + strings.Join(also, ", ") + ")"
		}
		fmt.Println(line)
		n++
	}
	if n == 0 {
		fmt.Printf("No words have %q in their meanings.\n", what)
	}
}

func usage() {
	fmt.Printf(`rquiz2 - shows Russian words and asks for 1 to 4 English meanings,
	separated by commas. The Russian words are numbered in usage order, 
	most common first.

  rquiz2              Easy Mode: the same as rquiz2 -s 1 -c 10 -p -m
                      (the 10 most frequently used words in order)
  rquiz2 -r 1         Hard Mode: random words from the whole list, until
                      all are used or you type q
  rquiz2 -s 100       words in order, starting at word 100, until q or
                      the end of the list
  rquiz2 -r 1 -e 40   random words from 1-40, until all are used or q
  rquiz2 -r 50 -e 250 -c 25
                      25 random words from 50-250, or until q
  rquiz2 -C           use Russian alphabetical (Cyrillic) order instead
                      of frequency order
  rquiz2 -p           show how many parts the expected answer has
  rquiz2 -m           Easy Match: a partial answer counts as complete
  rquiz2 -f high      find Russian words with "high" in English meanings
  rquiz2 -A           show the Russian alphabet, with sounds and examples
  rquiz2 -a           alphabet quiz: type the sound of each letter
  rquiz2 -a -o        alphabet order quiz: find the Nth letter in a word
  rquiz2 -v           show the program version and AI version
  rquiz2 -h           show this help

Answers can have up to 4 meanings, separated by commas, in any order.
Anything after a fourth comma is ignored. "to" before a verb is
optional. Type q to quit.

  Example question:   ка́ждый
  Example answer:     each, every

After each answer you'll see the meanings and how to say the word:
CAPITALS mark the stressed syllable, and ' marks a soft consonant.

Scoring: a complete answer +1, a partial answer +1/2, no answer -1, a
wrong answer -1 1/2. A wrong, missing or partial answer also gets a
short pause.

Suggestion: cover a small range of words repeatedly, then widen it:
  rquiz2 -s 1 -e 20, then (after learning those) rquiz2 -s 1 -e 30
  rquiz2 -r 1 -e 20, then rquiz2 -r 1 -e 30, and so on

There are currently [ %d ] words in the quiz.
`, len(words))
}

// splitFlags lets single-letter switches be run together as in rquiz,
// so -pm means -p -m.
func splitFlags(args []string) []string {
	var out []string
	for _, a := range args {
		if len(a) > 2 && a[0] == '-' && strings.Trim(a[1:], "CAaopmhv") == "" {
			for _, c := range a[1:] {
				out = append(out, "-"+string(c))
			}
			continue
		}
		out = append(out, a)
	}
	return out
}

func main() {
	seqStart := flag.Int("s", 0, "")
	randStart := flag.Int("r", 0, "")
	end := flag.Int("e", 0, "")
	count := flag.Int("c", 0, "")
	cyrillic := flag.Bool("C", false, "")
	showParts := flag.Bool("p", false, "")
	easyMatch := flag.Bool("m", false, "")
	search := flag.String("f", "", "")
	showAlpha := flag.Bool("A", false, "")
	alphaQuiz := flag.Bool("a", false, "")
	orderMode := flag.Bool("o", false, "")
	showVersion := flag.Bool("v", false, "")
	flag.CommandLine.SetOutput(os.Stdout)
	flag.Usage = usage
	flag.CommandLine.Parse(splitFlags(os.Args[1:]))

	if *showVersion {
		fmt.Printf("rquiz2 %s\nWritten by %s\nDirected by and based on prior work by: deano (deano186@gmail.com).\n", VERSION, AI)
		return
	}

	set := map[string]bool{}
	flag.Visit(func(f *flag.Flag) { set[f.Name] = true })
	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "unexpected argument: %s (see rquiz2 -h)\n", flag.Arg(0))
		os.Exit(2)
	}
	if set["s"] && set["r"] {
		fmt.Fprintln(os.Stderr, "use -s (in order) or -r (random), not both")
		os.Exit(2)
	}
	if (*showAlpha || *alphaQuiz) && (set["s"] || set["r"] || set["e"] || set["C"]) {
		fmt.Fprintln(os.Stderr, "-s, -r, -e and -C are for the word quiz; they can't be used with -A or -a")
		os.Exit(2)
	}
	if *orderMode && !*alphaQuiz {
		fmt.Fprintln(os.Stderr, "-o goes with -a: use rquiz2 -a -o")
		os.Exit(2)
	}

	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)
	fmt.Println()

	switch {
	case *showAlpha:
		showAlphabet()
		return
	case *alphaQuiz && *orderMode:
		orderQuiz(reader, *count)
		return
	case *alphaQuiz:
		letterQuiz(reader, *count)
		return
	}

	if *cyrillic {
		cyrillicOrder()
	}
	if set["f"] {
		find(*search)
		return
	}

	// With no options (other than -C), use Easy Mode, as rquiz did: the
	// first 10 words in order, with the parts hint and Easy Match.
	easyMode := !(set["s"] || set["r"] || set["e"] || set["c"] || set["p"] || set["m"])
	start, sequential := 1, true
	n := *count
	switch {
	case easyMode:
		n, *showParts, *easyMatch = 10, true, true
	case set["s"]:
		start = *seqStart
	case set["r"]:
		start, sequential = *randStart, false
	}
	start = max(start, 1)
	last := len(words)
	if set["e"] && *end < last {
		last = *end
	}
	if start > last {
		fmt.Fprintf(os.Stderr, "start word %d is after end word %d (there are %d words)\n", start, last, len(words))
		os.Exit(2)
	}
	size := last - start + 1
	if n <= 0 || n > size {
		n = size
	}
	order := make([]int, n)
	for i := range order {
		order[i] = start - 1 + i
	}
	if !sequential {
		for i, p := range rand.Perm(size)[:n] {
			order[i] = start - 1 + p
		}
	}

	quiz(reader, order, *showParts, *easyMatch, easyMode)
}

// quiz asks the words at the given list positions and scores the answers.
// The score is kept in half points.
func quiz(reader *bufio.Reader, order []int, showParts, easyMatch, easyMode bool) {
	halves, asked := 0, 0
	type miss struct {
		word   int
		answer string
	}
	var missed []miss

	for i, idx := range order {
		e := words[idx]
		fmt.Printf("%d/%d [%d] - %s\n", i+1, len(order), idx+1, e.russian)
		if showParts {
			if k := len(shown(e)); k == 1 {
				fmt.Println("Expected answer is in 1 part.")
			} else {
				fmt.Printf("Expected answer is in %d parts.\n", k)
			}
		}
		line, err := reader.ReadString('\n')
		answer := strings.TrimSpace(line)
		if err != nil && answer == "" || strings.EqualFold(answer, "q") {
			break
		}
		asked++
		fmt.Printf(" - - - ANSWER [%d]: - - %s:  %s  [%s]\n", i+1, e.russian, meanings(e), e.pronunciation)

		r, bad := grade(answer, e)
		if r == partial && easyMatch {
			fmt.Println("                 [ CORRECT +1 ] (Easy Match is ON)")
			halves += 2
		} else {
			switch r {
			case complete:
				fmt.Println("                 [ CORRECT +1 ]")
				halves += 2
			case partial:
				time.Sleep(time.Second)
				fmt.Println("      [ PARTIALLY CORRECT +.5 ]")
				halves++
			case wrong:
				time.Sleep(2 * time.Second)
				fmt.Println("       [ *** WRONG -1.5 *** ]")
				halves -= 3
			case noAnswer:
				time.Sleep(time.Second)
				fmt.Println("            [ NOT ANSWERED -1 ]")
				halves -= 2
			}
		}
		if r == wrong {
			for _, o := range lookalikes(bad, e) {
				fmt.Println("Maybe you were thinking of - " + describe(o))
			}
		}
		if r != complete {
			missed = append(missed, miss{idx, answer})
		}
		fmt.Println()
	}

	if len(missed) > 0 {
		sort.SliceStable(missed, func(a, b int) bool { return missed[a].word < missed[b].word })
		fmt.Println("Missed (or incomplete) ---------------------:")
		for _, m := range missed {
			fmt.Println()
			fmt.Println(describe(m.word))
			if m.answer == "" {
				m.answer = "(none)"
			}
			fmt.Println("  Your Answer: " + m.answer)
		}
		fmt.Println()
	}

	if easyMatch {
		fmt.Println("  (Easy Match Mode is ON)")
	}
	fmt.Printf("   Questions: %d\n", asked)
	if asked == 0 {
		fmt.Println("  No questions answered.")
	} else {
		fmt.Printf("       SCORE: %s\n\n", strconv.FormatFloat(float64(halves)/2, 'f', -1, 64))
		pct := float64(halves) * 50 / float64(asked)
		fmt.Printf("       GRADE: %.0f%%\n", pct)
		if pct >= highScore {
			fmt.Printf("\n  - %s -\n", randomPraise())
		}
	}
	fmt.Println()
	if easyMode {
		fmt.Println("Using default (EASY) mode, try ' rquiz2 -h ' to see other options.")
		fmt.Println()
	}
}

// A final grade of highScore percent or better gets one of these, picked
// at random.
const highScore = 95

var praise = []string{
	"Хоро́шая рабо́та! (Good work!)",
	"Молоде́ц! (Well done!)",
	"Отли́чно! (Excellent!)",
	"Здо́рово! (Great!)",
	"Great job! Your Russian is coming along.",
	"Impressive! Keep it up.",
	"Superb work!",
}

func randomPraise() string {
	return praise[rand.Intn(len(praise))]
}

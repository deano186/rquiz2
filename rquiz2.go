// rquiz2 - a Russian vocabulary quiz.
// Ver.20261115a
//
// Written by Claude, Anthropic's AI assistant. The first version came from
// the free version of Claude on claude.ai. On November 13, 2026, Claude
// Opus 5.5 (model ID claude-opus-5-5), working in Claude Code, expanded it
// with pronunciation guides, stress marks, lookalike-word hints, the 2000
// most common Russian words, and the -s/-r/-e/-c options. The same day,
// Claude Opus 5.5 added the alphabet table (-A) and alphabet quizzes (-a),
// and the Windows and Android builds.
//
// With help from Dean Hathaway, who came up with the ideas, tested it,
// and wrote rquiz, the Bash quiz this one grew out of. (Dean helped. A
// lot. Mostly by telling Claude what to do.)
//
// Originally named clauderquiz as a prototype; renamed rquiz2 as the
// successor to rquiz.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	return strings.Join(strings.Fields(s), " ")
}

func correct(input string, e entry) bool {
	in := normalize(input)
	if in == "" {
		return false
	}
	for _, a := range e.answers {
		if in == a {
			return true
		}
	}
	return false
}

// meanings returns the English translations for display, leaving out
// variants that only exist so more answers are accepted: "eat" when
// "to eat" is listed, digits like "1", and "youre" spellings of "you're".
func meanings(e entry) string {
	listed := map[string]bool{}
	for _, a := range e.answers {
		listed[a] = true
	}
	var shown []string
	seen := map[string]bool{}
	for _, a := range e.answers {
		key := strings.ReplaceAll(a, "'", "")
		if seen[key] || listed["to "+a] {
			continue
		}
		if _, err := strconv.Atoi(a); err == nil {
			continue
		}
		seen[key] = true
		shown = append(shown, a)
	}
	return strings.Join(shown, ", ")
}

// lookalikes returns the other quiz words that the English answer
// belongs to, e.g. "dog" when the question was кошка returns собака.
// The answer "eat" also matches "to eat", and "to eat" matches "eat".
// The word list is in frequency order, so the first matches found are the
// most common words; at most maxHints of them are returned.
// A blank answer gets no hints: there is no guess to compare.
const maxHints = 3

func lookalikes(input string, current entry) []entry {
	in := normalize(input)
	if in == "" {
		// Without this, the form "to " below would match words that
		// mean "to" (like к), as if the user had typed "to".
		return nil
	}
	forms := []string{in, "to " + in, strings.TrimPrefix(in, "to ")}
	var found []entry
	for _, e := range words {
		if e.russian == current.russian {
			continue
		}
		for _, f := range forms {
			if correct(f, e) {
				found = append(found, e)
				break
			}
		}
		if len(found) == maxHints {
			break
		}
	}
	return found
}

func usage() {
	out := flag.CommandLine.Output()
	fmt.Fprintf(out, `rquiz2 - Russian vocabulary quiz (%d words)

Shows Russian words and asks for their English meaning. The words are
numbered 1-%d, most commonly used first. Words 1-2000 follow the New
Frequency Dictionary of Russian (Lyashevskaya & Sharoff, 2009).

Usage: rquiz2 [options]

Options:
  -s N   sequential: ask words in order, starting at word N
  -r N   random: ask words in random order, from word N onwards
  -e N   end: stop at word N (use with -s or -r)
  -c N   count: ask N questions (default: every word in the range)
  -A     show the Russian alphabet, with sounds and example words
  -a     alphabet quiz: see each letter, type its sound (use -c for fewer)
  -a -o  alphabet order quiz: find the Nth letter of the alphabet in a word
  -h     show this help

With no options, it asks how many questions you want and picks them at
random from the whole list.

Examples:
  rquiz2 -s 1 -c 10     the 10 most common words, in order
  rquiz2 -r 1 -e 40     words 1-40 in random order
  rquiz2 -s 100         from word 100 to the end, in order
  rquiz2 -r 50 -e 250 -c 25
                        25 random words from words 50-250
  rquiz2 -a -c 10       10 random letters of the alphabet

During the quiz, type the English meaning, 'skip' to skip a word, or
'quit' to stop. Typing "to" before a verb is optional.

Suggestion: learn a small range first (e.g. -r 1 -e 20), then widen it
(-r 1 -e 30, and so on).
`, len(words), len(words))
}

// pickOrder returns the list positions (0-based) of the words to ask, in
// the order to ask them. start and end are 1-based word numbers.
func pickOrder(start, end, count int, sequential bool) []int {
	size := end - start + 1
	if count > size {
		count = size
	}
	order := make([]int, count)
	if sequential {
		for i := range order {
			order[i] = start - 1 + i
		}
		return order
	}
	for i, p := range rand.Perm(size)[:count] {
		order[i] = start - 1 + p
	}
	return order
}

func main() {
	seqStart := flag.Int("s", 0, "")
	randStart := flag.Int("r", 0, "")
	end := flag.Int("e", len(words), "")
	count := flag.Int("c", 0, "")
	showAlpha := flag.Bool("A", false, "")
	alphaQuiz := flag.Bool("a", false, "")
	orderMode := flag.Bool("o", false, "")
	flag.Usage = usage
	flag.Parse()

	set := map[string]bool{}
	flag.Visit(func(f *flag.Flag) { set[f.Name] = true })
	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "unexpected argument: %s\n\n", flag.Arg(0))
		usage()
		os.Exit(2)
	}
	if set["s"] && set["r"] {
		fmt.Fprintln(os.Stderr, "use -s (sequential) or -r (random), not both")
		os.Exit(2)
	}

	if (*showAlpha || *alphaQuiz) && (set["s"] || set["r"] || set["e"]) {
		fmt.Fprintln(os.Stderr, "-s, -r and -e are for the word quiz; they can't be used with -A or -a")
		os.Exit(2)
	}
	if *orderMode && !*alphaQuiz {
		fmt.Fprintln(os.Stderr, "-o goes with -a: use rquiz2 -a -o")
		os.Exit(2)
	}

	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)

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

	sequential := set["s"]
	start := 1
	if set["s"] {
		start = *seqStart
	} else if set["r"] {
		start = *randStart
	}
	if start < 1 {
		start = 1
	}
	if *end > len(words) {
		*end = len(words)
	}
	if start > *end {
		fmt.Fprintf(os.Stderr, "start word %d is after end word %d (there are %d words)\n", start, *end, len(words))
		os.Exit(2)
	}

	// fmt.Println("=== Russian Vocabulary Quiz ===")
	// fmt.Println("Type the English translation of each Russian word.")
	// fmt.Println("Type 'quit' to stop, or 'skip' to skip a word.")
	// fmt.Println("The accent mark (´) shows which vowel is stressed.")
	// fmt.Println("After each answer you'll see how to say the word:")
	// fmt.Println("  CAPITALS = stressed syllable, ' = soft consonant.")
	// fmt.Println()

	n := *count
	if len(set) == 0 {
		// No options: ask how many questions, as before.
		fmt.Printf("How many questions? (1-%d, default 10): ", len(words))
		line, _ := reader.ReadString('\n')
		n = 10
		if v, err := strconv.Atoi(strings.TrimSpace(line)); err == nil && v > 0 {
			n = v
		}
		fmt.Println()
	} else if n <= 0 {
		n = *end - start + 1
	}

	order := pickOrder(start, *end, n, sequential)
	n = len(order)
	if len(set) > 0 && n > 0 {
		if sequential {
			// In order, the quiz stops after n words, which may be
			// before the end of the range.
			fmt.Printf("%d questions: words %d-%d, in order.\n\n", n, start, start+n-1)
		} else {
			fmt.Printf("%d questions from words %d-%d, random.\n\n", n, start, *end)
		}
	}
	score, asked := 0, 0

	for i, idx := range order {
		e := words[idx]
		fmt.Printf("[%d/%d] #%d  %s\n> ", i+1, n, idx+1, e.russian)
		line, err := reader.ReadString('\n')
		answer := normalize(line)
		if err != nil && answer == "" {
			fmt.Println()
			break
		}
		if answer == "quit" || answer == "exit" {
			break
		}
		asked++
		switch {
		case answer == "skip":
			fmt.Println("Skipped.")
		case correct(answer, e):
			score++
			fmt.Printf("Correct! Score: %d\n", score)
		default:
			fmt.Println("Wrong.")
			for _, o := range lookalikes(answer, e) {
				fmt.Printf("You might be thinking of %s [%s], which means %s.\n",
					o.russian, o.pronunciation, meanings(o))
			}
		}
		fmt.Printf("Meaning: %s\n", meanings(e))
		fmt.Printf("Pronunciation: %s  [%s]\n\n", e.russian, e.pronunciation)
	}
	finalScore(score, asked)
}

func finalScore(score, asked int) {
	fmt.Println("=== Quiz over ===")
	if asked > 0 {
		fmt.Printf("Final score: %d/%d (%.0f%%)\n", score, asked, float64(score)/float64(asked)*100)
	} else {
		fmt.Println("No questions answered.")
	}
}

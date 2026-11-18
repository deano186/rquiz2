# rquiz2

A command-line quiz for learning Russian vocabulary, rebuilt in Go from
rquiz, the original Bash quiz. It shows a Russian word, you type its English
meanings, and it scores your answer and shows how to pronounce the word.

- 2006 words, numbered from most to least common. Words 1–2000 follow the
  *New Frequency Dictionary of Russian Vocabulary* (Lyashevskaya & Sharoff, 2009).
- Stress marks on every word of more than one syllable (ко́шка).
- A pronunciation guide after each answer: CAPITALS mark the stressed syllable,
  and `'` marks a soft consonant (`KOHSH-kuh`, `byt'`).
- Answers of up to 4 meanings, with partial credit.
- If you give the meaning of a different word, it tells you which word you
  were thinking of.
- At the end, a list of the words you missed, with your answers.

## Example

```
2/10 [2] - в
Expected answer is in 4 parts.
on, in
 - - - ANSWER [2]: - - в:  in, into, at, to  [v]
      [ PARTIALLY CORRECT +.5 ]
Maybe you were thinking of - [4] на: on, onto, at, to  [nah]
```

## Usage

Run `rquiz2` with no options for Easy Mode: the 10 most common words, in
order, with the parts hint and Easy Match on (the same as
`rquiz2 -s 1 -c 10 -p -m`).

| Option | Meaning |
|--------|---------|
| `-s N` | ask words in order, starting at word N |
| `-r N` | ask words in random order, from word N onwards |
| `-e N` | stop at word N |
| `-c N` | ask N questions (default: every word in the range) |
| `-C`   | use Russian alphabetical (Cyrillic) order instead of frequency order |
| `-p`   | show how many parts the expected answer has |
| `-m`   | Easy Match: a partial answer counts as complete |
| `-f W` | find words with W in their English meanings |
| `-A`   | show the Russian alphabet, with sounds and example words |
| `-a`   | alphabet quiz: see each letter, type its sound |
| `-a -o`| alphabet order quiz: find the Nth letter of the alphabet in a word |
| `-v`   | show the program version and AI version |
| `-h`   | show help |

Switches can be run together, so `-pm` is the same as `-p -m`.

```
rquiz2 -r 1                # Hard Mode: random words from the whole list
rquiz2 -s 100              # words in order from word 100
rquiz2 -r 1 -e 40          # words 1-40 in random order
rquiz2 -r 50 -e 250 -c 25  # 25 random words from words 50-250
rquiz2 -f high             # which words mean "high"?
```

### Answers and scoring

Give up to 4 meanings, separated by commas, in any order. Anything after a
fourth comma is ignored, and "to" before a verb is optional. Type `q` to quit.

| Answer | Points |
|--------|--------|
| complete: every meaning shown, nothing wrong | +1 |
| partial: at least one right meaning | +½ (+1 with Easy Match) |
| no answer | −1 |
| wrong | −1½ |

A wrong, missing or partial answer also gets a short pause. The final grade
is your score as a percentage of the questions asked. A grade of 95% or more
gets a word of praise.

### Learning the alphabet

`rquiz2 -A` shows all 33 letters with their sounds, an example word for each,
and a list of letters that look like English ones but sound different
(Р is "r", Н is "n", В is "v", and so on).

`rquiz2 -a` shows the letters in random order and asks for each one's sound,
such as `zh` for Ж. Use `-c 10` for a shorter quiz.

`rquiz2 -a -o` goes through the alphabet in order. For each number N it shows
a word, and you type where in the word the Nth letter of the alphabet is.
The words are picked at random from the word list, so they change every time.

A good way to learn: cover a small range of words repeatedly
(`-s 1 -e 20`), then widen it (`-s 1 -e 30`, and so on).

## Building

You need [Go](https://go.dev/dl/) 1.26 or later.

```
make            # build all three versions
make build      # build rquiz2 for this computer
make windows    # build rquiz2.exe for 64-bit Windows
make android    # build rquiz2-android for ARM64 Android (run it in Termux)
make clean      # remove built programs
```

On Windows, the Russian text displays correctly in Windows Terminal. In the
older `cmd.exe` window, run `chcp 65001` first.

## Credits

Written by Claude (Anthropic's AI assistant) with deano (deano186@gmail.com), who came up
with the ideas, tested it, and wrote rquiz, the Bash quiz this grew out of.

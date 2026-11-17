# rquiz2

A command-line quiz for learning Russian vocabulary. It shows a Russian word,
you type the English meaning, and it tells you whether you were right and how
to pronounce the word.

- 2005 words, numbered from most to least common. Words 1–2000 follow the
  *New Frequency Dictionary of Russian Vocabulary* (Lyashevskaya & Sharoff, 2009).
- Stress marks on every word of more than one syllable (ко́шка).
- A pronunciation guide after each answer: CAPITALS mark the stressed syllable,
  and `'` marks a soft consonant (`KOHSH-kuh`, `byt'`).
- If you give the meaning of a different word, it tells you which word you
  were thinking of.

## Example

```
[2/3] #2  в
> cat
Wrong.
You might be thinking of ко́шка [KOHSH-kuh], which means cat.
Meaning: in, into, at, to
Pronunciation: в  [v]
```

## Usage

Run `rquiz2` with no options and it asks the 10 most common words, in
order, like the original rquiz (the same as `rquiz2 -s 1 -c 10`). Each run
ends with a reminder that `rquiz2 -h` lists the options. Options let you
choose the words:

| Option | Meaning |
|--------|---------|
| `-s N` | ask words in order, starting at word N |
| `-r N` | ask words in random order, from word N onwards |
| `-e N` | stop at word N (use with `-s` or `-r`) |
| `-c N` | ask N questions (default: every word in the range) |
| `-A`   | show the Russian alphabet, with sounds and example words |
| `-a`   | alphabet quiz: see each letter, type its sound |
| `-a -o`| alphabet order quiz: find the Nth letter of the alphabet in a word |
| `-h`   | show help |

```
rquiz2 -s 1 -c 10        # the 10 most common words, in order
rquiz2 -r 1 -e 40        # words 1-40 in random order
rquiz2 -r 50 -e 250 -c 25  # 25 random words from words 50-250
```

During the quiz, press Enter to skip a word or type `q` to quit. Typing "to"
before a verb is optional.

### Learning the alphabet

`rquiz2 -A` shows all 33 letters with their sounds, an example word for each,
and a list of letters that look like English ones but sound different
(Р is "r", Н is "n", В is "v", and so on).

`rquiz2 -a` shows the letters in random order and asks for each one's sound,
such as `zh` for Ж. Use `-c 10` for a shorter quiz.

`rquiz2 -a -o` goes through the alphabet in order. For each number N it shows
a word, and you type where in the word the Nth letter of the alphabet is.
The words are picked at random from the word list, so they change every time.

A good way to learn: start with a small range (`-r 1 -e 20`), then widen it
(`-r 1 -e 30`, and so on).

## Building

You need [Go](https://go.dev/dl/) 1.26 or later.

```
make            # build rquiz2 for this computer
make windows    # build rquiz2.exe for 64-bit Windows
make android    # build rquiz2-android for ARM64 Android (run it in Termux)
make clean      # remove built programs
```

On Windows, the Russian text displays correctly in Windows Terminal. In the
older `cmd.exe` window, run `chcp 65001` first.

## Credits

Written by Claude (Anthropic's AI assistant) with Dean Hathaway, who came up
with the ideas, tested it, and wrote rquiz, the Bash quiz this grew out of.

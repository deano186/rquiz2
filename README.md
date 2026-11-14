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

Run `rquiz2` with no options and it asks how many questions you want, then
picks them at random from the whole list. Options let you choose the words:

| Option | Meaning |
|--------|---------|
| `-s N` | ask words in order, starting at word N |
| `-r N` | ask words in random order, from word N onwards |
| `-e N` | stop at word N (use with `-s` or `-r`) |
| `-c N` | ask N questions (default: every word in the range) |
| `-h`   | show help |

```
rquiz2 -s 1 -c 10        # the 10 most common words, in order
rquiz2 -r 1 -e 40        # words 1-40 in random order
rquiz2 -r 50 -e 250 -c 25  # 25 random words from words 50-250
```

During the quiz, type `skip` to skip a word or `quit` to stop. Typing "to"
before a verb is optional.

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

# Exercise: implement WordCount

Complete the function that counts word frequencies in a text.

## APIs used in this exercise

| API | Role |
| --- | --- |
| `strings.Fields(s)` | Split into words around white space |
| `strings.ToLower(s)` | Ignore case differences |
| `m[w]++` | Increase the count in a map |

## What to do

Implement `WordCount(text)` in `wordcount.go`.

1. Lower-case the text and split it around white space.
2. Count each word in a map and return it. For an empty text return an empty (not nil) map.

## How it is graded

A hidden test compares the resulting map with the expected one for several inputs. On failure it tells you which input differed and how.

## Hint

- Create the map with `make`. You cannot write to a `nil` map.

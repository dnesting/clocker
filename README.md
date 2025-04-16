# clocker

Clocker is a simple command-line tool for marking and annotating timestamps.
When you press Enter on a blank line, it marks a new timestamp.
If you type text, and press Enter, it will annotate the previously-marked timestamp and won't create a new timestamp.
If your text starts with a number, it will annotate the numbered timestamp, preserving the annotation that was there.

# output format

0: 2025-01-01T13:14:15.123456-0400 +0: start
1: 2025-01-01T13:15:15.123456-0400 +01:00.000000: annotation
2: 2025-01-01T13:17:15.123456-0400 +02:00.000000: annotation; with additional annotation

# usage

```
clocker [-T] [-D] [FILE]

  -T, --no-time
      don't output a timestamp

  -D, --no-delta
      don't output deltas
```

If a file is provided, output will be written to this file instead of STDOUT.

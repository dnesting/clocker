package clocker

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

type Options struct {
	// NoTimestamp disables the timestamp in the output.
	NoTimestamp bool
	// NoDelta disables the delta in the output.
	NoDelta bool
	// Out is the output writer.
	// If nil, os.Stdout is used.
	Out io.Writer
}

type Timestamp struct {
	T          time.Time
	Annotation []string
}

func formatTimestamp(idx int, ts Timestamp, options Options, start time.Time) string {
	// NN: YYYY-MM-DD HH:MM:SS.ffffff +[HH:[MM:]]SS.ffffff[: ANNOTATION]
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d", idx))
	if !options.NoTimestamp {
		sb.WriteString(": ")
		sb.WriteString(ts.T.Format("2006-01-02 15:04:05.000000"))
	}
	if !options.NoDelta {
		if options.NoTimestamp {
			sb.WriteString(": ")
		}
		sb.WriteString(" ")
		delta := ts.T.Sub(start)
		sb.WriteString(fmt.Sprintf("+%s", delta))
	}
	if len(ts.Annotation) > 0 {
		sb.WriteString(": ")
		sb.WriteString(strings.Join(ts.Annotation, "; "))
	}
	return sb.String()
}

type Clocker struct {
	Options    Options
	Timestamps []Timestamp
}

func (c *Clocker) ToString(idx int) string {
	return formatTimestamp(idx, c.Timestamps[idx], c.Options, c.Timestamps[0].T)
}

func (c *Clocker) Add(ts Timestamp) int {
	c.Timestamps = append(c.Timestamps, ts)
	return len(c.Timestamps) - 1
}

func (c *Clocker) Annotate(index int, text string) (int, error) {
	if index < 0 {
		index = len(c.Timestamps) - 1
	}
	if index < 0 || index >= len(c.Timestamps) {
		return index, fmt.Errorf("index out of range: %d", index)
	}
	c.Timestamps[index].Annotation = append(c.Timestamps[index].Annotation, text)
	return index, nil
}

func (c *Clocker) AnnotateLast(text string) (int, error) {
	return c.Annotate(-1, text)
}

func (c *Clocker) WriteTo(w io.Writer) (written int64, err error) {
	for i := range c.Timestamps {
		line := c.ToString(i)
		n, err := w.Write([]byte(line + "\n"))
		if err != nil {
			return written, err
		}
		written += int64(n)
	}
	return written, nil
}

func readLine() (int, string, error) {
	// Accept any of these input lines:
	// "NNN annotation text here"
	// "annotation text here"
	// ""

	var line string

	// Read a line from stdin
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return 0, "", err
		}
		return 0, "", fmt.Errorf("read line: %w", err)
	}

	line = strings.TrimSpace(line)

	// "NNN annotation text here"
	var idx int
	var text string
	if n, err := fmt.Sscanf(line, "%d %s", &idx, &text); err == nil && n == 2 {
		return idx, strings.TrimSpace(text), nil
	}
	return -1, line, nil
}

// Run waits for input from the user.
// When the user presses Enter on a blank line, a new timestamp will be created.
// If the user types text before pressing Enter, instead of creating a new timestamp, the last created timestamp will be annotated with the text.
// If the user types a number followed by a space and then text, the timestamp referenced by the number will be annotated with the text.
// If the user interrupts (^C) or closes input (^D), the timestamps will be written to the output writer and the function will terminate.
func Run(options Options) {
	var cl Clocker
	out := options.Out
	if out == nil {
		out = os.Stdout
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		close(c)
		cl.WriteTo(out)
	}()

	go func() {
		for range c {
			signal.Stop(c)
			cl.WriteTo(out)
			os.Exit(0)
		}
	}()

	idx := cl.Add(Timestamp{T: time.Now(), Annotation: []string{"started"}})
	fmt.Println(cl.ToString(idx))

	for {
		idx, text, err := readLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		if text != "" {
			idx, err := cl.Annotate(idx, text)
			if err != nil {
				fmt.Printf("%d: %s\n", idx, err)
			} else {
				fmt.Println(cl.ToString(idx))
			}
		} else {
			idx := cl.Add(Timestamp{T: time.Now()})
			fmt.Println(cl.ToString(idx))
		}
	}
}

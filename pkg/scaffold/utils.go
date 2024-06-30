package scaffold

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

func prompt(ctx context.Context, stdin io.Reader, text string) (string, error) {
	fmt.Fprint(os.Stderr, text)

	var err error
	var answer string
	done := make(chan struct{})
	go func() {
		defer close(done)

		rdr := bufio.NewReader(stdin)
		var answerBytes []byte
		answerBytes, err = rdr.ReadBytes('\n')
		if err == nil {
			answer = string(answerBytes[:len(answerBytes)-1])
		}
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-done:
		return answer, err
	}
}

func promptYesNo(ctx context.Context, stdin io.Reader, text string) (bool, error) {
	answer, err := prompt(ctx, stdin, text+" [y/N] ")
	if err != nil {
		return false, err
	}
	return strings.EqualFold(answer, "y"), nil
}

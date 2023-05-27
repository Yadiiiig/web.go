package library

import (
	"bufio"
	"fmt"
	"os"
)

func (s *Settings) GenLibrary(files []File) error {
	var output string

	for _, f := range files {
		for _, r := range f.Internal.Requests {
			tmp := GenRequest(
				r.Name,
				"POST",
				fmt.Sprintf(
					"%s/%s/%s",
					s.Endpoint,
					f.Name,
					r.Name,
				),
				r.Params,
			)

			output = fmt.Sprintf("%s\n%s", output, tmp)
		}
	}

	output = fmt.Sprintf(base, s.Endpoint, output)
	output = fmt.Sprintf("%s%s", intro, output)

	f, err := os.OpenFile("wg.js", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(output)
	if err != nil {
		return err
	}

	w.Flush()

	return nil
}

func GenFile(url, body string) error {
	output := fmt.Sprintf(base, url, body)
	fmt.Println(output)

	return nil
}

func GenRequest(name, method, url string, args []string) string {
	body := ""

	for _, arg := range args {
		body += fmt.Sprintf("{'%v': args['%v']}, ", arg, arg)
	}

	body = fmt.Sprintf("[%s]", body[:len(body)-2])
	body = GenVariable("let", "opts", fmt.Sprintf(request, method, body))
	body = fmt.Sprintf("%s%s", body, fmt.Sprintf(fetch, url))

	body = fmt.Sprintf(function, name, "", body)

	return body
}

func GenVariable(tp, name, args string) string {
	var tmp string

	if args == "" {
		tmp = ";"
	} else {
		tmp = fmt.Sprintf(" = %s;", args)
	}

	return fmt.Sprintf("\n\t%s %s%s", tp, name, tmp)
}

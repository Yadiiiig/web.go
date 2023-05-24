package library

import (
	"fmt"
)

func GenFile(url, body string) error {
	output := fmt.Sprintf(base, url, body)
	fmt.Println(output)

	return nil
}

func GenRequest(name, method, url, rdr string, args []string) string {
	body := ""

	for _, arg := range args {
		body += fmt.Sprintf("{'%v': args['%v']}, ", arg, arg)
	}

	body = fmt.Sprintf("{%s}", body[:len(body)-2])
	body = GenVariable("let", "opts", fmt.Sprintf(request, method, body, rdr))
	body = fmt.Sprintf("%s%s", body, fmt.Sprintf(fetch, url))

	body = fmt.Sprintf(function, name, args, body)

	return body
}

func GenVariable(tp, name, args string) string {
	var tmp string

	if args == "" {
		tmp = ";"
	} else {
		tmp = fmt.Sprintf(" = %s;", args)
	}

	return fmt.Sprintf(`%s %s%s`, tp, name, tmp)
}

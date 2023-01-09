package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	help := flag.Bool("help", false, "shows help")
	uri := flag.String("url", "", "input source, URL or file or stdin(default)")
	timeout := flag.Int("timeout", 0, "request timeout in seconds")
	query := flag.String("query", "", "parsing query")
	values := make(Values, 0)
	flag.Var(&values, "values", "values form of elements(text(default)/html/attr:some-html-attr), default text")
	delim := flag.String("delim", ",", "delimiter of each values(default: ',')")
	flag.Parse()
	if *help {
		line := flag.CommandLine
		line.Usage()
		os.Exit(10)
	}
	app := App{
		uri:     Uri(*uri),
		timeout: time.Duration(*timeout),
		query:   Query(*query),
		values:  values,
		delim:   *delim,
	}
	errors := app.Validate()
	if len(errors) != 0 {
		fmt.Fprintf(os.Stderr, "Invalid parameter\n")
		fmt.Fprintf(os.Stderr, "  errors:\n")
		for i, err := range errors {
			fmt.Fprintf(os.Stderr, "    %2d %s\n", i+1, err.Error())
		}
		os.Exit(1)
	}
	texts, err := app.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Runtime errors\n")
		fmt.Fprintln(os.Stderr, err.Error())
	}
	for _, text := range texts {
		fmt.Println(text)
	}
}

type Uri string

func (uri Uri) IsHttp() bool {
	return strings.Index(string(uri), "http") == 0
}

func (uri Uri) Load(timeout time.Duration) (io.ReadCloser, error) {
	if len(uri) == 0 {
		return os.Stdin, nil
	}
	if !uri.IsHttp() {
		return uri.OpenFile()
	} else {
		return uri.OpenHttp(timeout)
	}
}

func (uri Uri) OpenFile() (io.ReadCloser, error) {
	fileName := string(uri)
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("uri#OpenFile :%w", err)
	}
	return file, nil
}

func (uri Uri) OpenHttp(timeout time.Duration) (io.ReadCloser, error) {
	var client http.Client
	client.Timeout = timeout
	response, err := client.Get(string(uri))
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

type Query string
type Value string

func (v Value) Get(elem *goquery.Selection) (string, error) {
	value := string(v)
	switch value {
	case "text":
		return elem.Text(), nil
	case "html":

		return elem.Html()
	default:
		attr, _ := elem.Attr(value)
		return attr, nil
	}
}

type Values []Value

func (vs *Values) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	s := len(*vs)
	for i, v := range *vs {
		sb.WriteString(string(v))
		if i+1 <= s {
			sb.WriteString(",")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func (vs *Values) Set(v string) error {
	*vs = append(*vs, Value(v))
	return nil
}

func (vs *Values) Get(elem *goquery.Selection, delim string) (string, error) {
	var sb strings.Builder
	es := make([]error, 0)
	s := len(*vs)
	for index, value := range *vs {
		v, err := value.Get(elem)
		if strings.ToLower(string(value)) == "html" {
			return v, err
		}
		if err != nil {
			es = append(es, fmt.Errorf("at %d, value: %s, %w", index, value, err))
		}
		sb.WriteString(v)
		if index+1 < s {
			sb.WriteString(delim)
		}
	}
	result := sb.String()
	if 0 < len(es) {
		var eb strings.Builder
		t := len(es)
		for i, err := range es {
			eb.WriteString(fmt.Sprintf("\t\tat %d:%s", i, err))
			if i+1 < t {
				eb.WriteString("\n")
			}
		}
		return result, fmt.Errorf("error to get values: \n%s", eb)
	}
	return result, nil
}

type App struct {
	uri     Uri
	timeout time.Duration
	query   Query
	values  Values
	delim   string
}

func (app App) Validate() []error {
	errors := make([]error, 0)
	if app.timeout < 0 {
		errors = append(errors, fmt.Errorf("timeout must be positive values"))
	}
	if len(app.values) == 0 {
		errors = append(errors, fmt.Errorf("values should be one of these: text(default), html, attr:some-html-attribute"))
	}
	return errors
}

func (app App) Run() ([]string, error) {
	input, err := app.uri.Load(app.timeout)
	if err != nil {
		return []string{}, err
	}
	defer func() { _ = input.Close() }()
	document, err := goquery.NewDocumentFromReader(input)
	if err != nil {
		return []string{}, err
	}
	results := make([]string, 0)
	es := make([]error, 0)
	selection := document.Find(string(app.query))
	selection.Each(func(index int, elem *goquery.Selection) {
		value, err := app.values.Get(elem, app.delim)
		if err != nil {
			es = append(es, fmt.Errorf("at(%d): %w", index+1, err))
		}
		if value != "" {
			results = append(results, value)
		}
	})
	if len(es) > 0 {
		var esb strings.Builder
		for _, err := range es {
			esb.WriteString(fmt.Sprintf("\t%s", err.Error()))
		}
		return results, fmt.Errorf("error while getting values\n%s", esb.String())
	}
	return results, nil
}

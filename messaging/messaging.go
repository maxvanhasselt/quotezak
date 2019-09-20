package messaging

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"

	"git.code-cloppers.com/max/quotezak/models"
)

// Messenger handles the commands that get fed into the bot
type Messenger struct {
	db *sql.DB
}

type Message struct {
	prefix  string
	command string
	params  []string
}

// NewMessenger returns a pointer to a new Messenger object
func NewMessenger(db *sql.DB) *Messenger {
	return &Messenger{
		db: db,
	}
}

// GenerateMessage takes input and returns a message if neccesary
func (m *Messenger) GenerateMessage(msg string) *string {
	ms := m.parseMessage(msg)
	return m.routeMessage(ms)
}

func isCommand(c string) bool {
	for _, a := range []string{"quote", "addquote"} {
		if a == c {
			return true
		}
	}
	return false
}

func (m *Messenger) routeMessage(msg *Message) *string {
	var target string
	if len(msg.params) >= 2 {
		switch msg.params[0] {
		case "quotezak":
			re := regexp.MustCompile("^[^!]*")
			target = re.FindStringSubmatch(msg.prefix)[0]
		default:
			target = msg.params[0]
			//fmt.Println(msg.params[0])
			//return nil
		}
		command := strings.Split(msg.params[1], " ")
		if isCommand(command[0]) {
			str := fmt.Sprintf("PRIVMSG %s :%s\n", target, m.handleCommand(command))
			return &str
		}
	}
	return nil
}

func (m *Messenger) handleCommand(params []string) string {
	switch params[0] {
	case "quote":
		return m.handleGetQuote(params)
	case "addquote":
		return m.handleAddQuote(params[1:])
	}
	return "niks nier"
}

func (m *Messenger) handleGetQuote(params []string) string {
	var q *models.Quote
	var err error
	if len(params) == 1 {
		q, err = models.GetRandomQuote(m.db)
		if err != nil {
			log.Println(err)
			return ""
		}
		return q.String()
	}
	switch params[1][0] {
	case '#':
		re := regexp.MustCompile("#")
		params[1] = re.ReplaceAllString(params[1], "")
		q, err = models.GetRandomCategory(m.db, params[1])
		if err != nil {
			log.Println(err)
			return ""
		}
		return q.String()
	case '@':
		re := regexp.MustCompile("@")
		params[1] = re.ReplaceAllString(strings.Join(params[1:], " "), "")
		q, err = models.GetRandomOwner(m.db, params[1])
		if err != nil {
			log.Println(err)
			return ""
		}
		return q.String()
	default:
		q, err = models.GetQuoteByName(m.db, params[1])
		if err != nil {
			log.Println(err)
			return ""
		}
		return q.String()
	}
}

func (m *Messenger) handleAddQuote(params []string) string {

	re := regexp.MustCompile("(^\"[^\"]*\") (\"[^\"]*\") ([0-9]*) (\"[^\"]*\") (#.*)")
	results := re.FindStringSubmatch(strings.Join(params, " "))

	if len(results) == 0 {
		return "Format not recognized, please use '\"Quote\" \"Who said it\" <year> \"quote name\" #category'"
	}
	re = regexp.MustCompile("[^a-zA-z0-9 ?!.,]*")
	for i := range results {
		results[i] = re.ReplaceAllString(results[i], "")
	}

	quote := models.NewQuote(results[4], results[1], results[2], results[3], results[5])
	err := quote.Save(m.db)
	if err != nil {
		fmt.Println(err)
		return fmt.Sprintf("Error saving quote: %s", err)
	}

	return "Quote saved!"
}

func (m *Messenger) parseMessage(message string) *Message {
	message = strings.Trim(message, "\r\n")

	fmt.Printf("[message] %s\n", message)

	tokens := []string{}
	if strings.Index(message, ":") == 0 {
		message = message[1:]
	} else {
		tokens = append(tokens, "")
	}

	re := regexp.MustCompile(`(?::(.*$)|(\S+))(.*)$`)
	for {
		m := re.FindStringSubmatch(message)
		if m == nil || len(m) != 4 {
			break
		}
		tokens = append(tokens, m[1]+m[2])
		message = m[3]
	}
	fmt.Println("'" + strings.Join(tokens, "', '") + "'")

	return &Message{tokens[0], tokens[1], tokens[2:]}
}

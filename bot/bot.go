package bot

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"regexp"
	"sync"
	"time"

	"git.code-cloppers.com/max/quotezak/messaging"
)

type Bot struct {
	Cfg *Config
	msr *messaging.Messenger

	conn     *net.Conn
	incoming chan string
	outgoing chan string
	wg       *sync.WaitGroup
	reader   *bufio.Reader
	writer   *bufio.Writer
}

func New(Cfg *Config, msr *messaging.Messenger) *Bot {
	return &Bot{
		Cfg: Cfg,
		msr: msr,
	}
}

func (b *Bot) Start() error {
	var wg sync.WaitGroup
	b.wg = &wg

	wg.Add(1)

	b.incoming = make(chan string)
	b.outgoing = make(chan string)

	fmt.Printf("tcp %s %s", b.Cfg.Server, b.Cfg.Port)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", b.Cfg.Server, b.Cfg.Port))
	if err != nil {
		return err
	}
	b.conn = &conn
	b.reader = bufio.NewReader(*b.conn)
	b.writer = bufio.NewWriter(*b.conn)

	go b.ReadSocket()
	go b.HandleRecieve()
	go b.HandleSend()
	b.Login()
	b.JoinChannels()
	wg.Wait()

	return nil
}

func (b *Bot) ReadSocket() {
	defer b.wg.Done()
	for {
		if line, err := b.reader.ReadString('\n'); err != nil {
			fmt.Println(err)
			err := b.restartConnection()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			b.incoming <- line
		}
	}
}

func (b *Bot) restartConnection() error {
	for try := 1; try <= 10; try++ {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", b.Cfg.Server, b.Cfg.Port))
		if err != nil {
			fmt.Printf("Error %s when reconnecting, retrying %d times", err, (10 - try))
			time.Sleep(time.Duration(2*try) * time.Second)
		} else {
			b.conn = &conn
			b.reader = bufio.NewReader(*b.conn)
			b.writer = bufio.NewWriter(*b.conn)
		}
	}
	return fmt.Errorf("Could not reconnect to remote service, shutting down")
}

func (b *Bot) HandleSend() {
	for {
		message := <-b.outgoing
		b.writer.WriteString(message + "\n")
		b.writer.Flush()
		fmt.Printf("--> %s", message)
	}
}

func (b *Bot) HandleRecieve() {
	message := <-b.incoming
	for {
		fmt.Printf("<-- %s", message)
		b.HandleCommand(message)
		message = <-b.incoming
	}
}

func (b *Bot) HandleCommand(msg string) {
	exp := regexp.MustCompile(`^PING :(.*)`)
	if matches := exp.FindStringSubmatch(msg); len(matches) > 0 {
		b.outgoing <- fmt.Sprintf("PONG :%s", matches[1])
	}
	m := b.msr.GenerateMessage(msg)
	if m != nil {
		b.outgoing <- *m
	}
}

func (b *Bot) Login() {
	b.outgoing <- fmt.Sprintf("NICK %s\n", b.Cfg.Nick)
	b.outgoing <- fmt.Sprintf("USER %s 0 * :%s\n", b.Cfg.Identity, b.Cfg.Realname)
}

func (b *Bot) JoinChannels() {
	for _, channel := range b.Cfg.Channels {

		b.outgoing <- fmt.Sprintf("JOIN #%s", channel)
	}
}

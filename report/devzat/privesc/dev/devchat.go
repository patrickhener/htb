package main

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/gliderlabs/ssh"
	terminal "golang.org/x/term"
)

var (
	port = 8443

	devbot = "" // initialized in main

	mainRoom = &room{"#main", make([]*user, 0, 10), sync.Mutex{}}
	rooms    = map[string]*room{mainRoom.name: mainRoom}

	allUsers      = make(map[string]string, 400) //map format is u.id => u.name
	allUsersMutex = sync.Mutex{}

	backlog      = make([]backlogMessage, 0, scrollback)
	backlogMutex = sync.Mutex{}

	logfile, _ = os.OpenFile("log.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	l          = log.New(io.MultiWriter(logfile, os.Stdout), "", log.Ldate|log.Ltime|log.Lshortfile)

	idsInMinToTimes = make(map[string]int, 10)
	idsInMinMutex   = sync.Mutex{}

	antispamMessages = make(map[string]int)
	//antispamMutex    = sync.Mutex{}
)

type room struct {
	name       string
	users      []*user
	usersMutex sync.Mutex
}

type user struct {
	name          string
	session       ssh.Session
	term          *terminal.Terminal
	bell          bool
	color         string
	id            string
	addr          string
	win           ssh.Window
	closeOnce     sync.Once
	lastTimestamp time.Time
	joinTime      time.Time
	timezone      *time.Location
	room          *room
	messaging     *user
	slack         bool
}

type backlogMessage struct {
	timestamp  time.Time
	senderName string
	text       string
}

// TODO: have a web dashboard that shows logs
func main() {
	registerCommands()
	devbot = green.Paint("devbot")
	var err error
	rand.Seed(time.Now().Unix())
	readBansAndUsers()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<-c
		fmt.Println("Shutting down...")
		saveBansAndUsers()
		logfile.Close()
		mainRoom.broadcast(devbot, "Server going down! This is probably because it is being updated. Try joining back immediately.", true)
		os.Exit(0)
	}()
	ssh.Handle(func(s ssh.Session) {
		u := newUser(s)
		if u == nil {
			return
		}
		u.repl()
	})
	if os.Getenv("PORT") != "" {
		port, err = strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Printf("Starting chat server on port %d\n", port)
	err = ssh.ListenAndServe(
		fmt.Sprintf("127.0.0.1:%d", port),
		nil,
		ssh.HostKeyFile(os.Getenv("HOME")+"/.ssh/id_rsa"))
	if err != nil {
		fmt.Println(err)
	}
}

func (r *room) broadcast(senderName, msg string, toSlack bool) {
	if msg == "" {
		return
	}
	msg = strings.ReplaceAll(msg, "@everyone", green.Paint("everyone\a"))
	r.usersMutex.Lock()
	for i := range r.users {
		msg = strings.ReplaceAll(msg, "@"+stripansi.Strip(r.users[i].name), r.users[i].name)
		msg = strings.ReplaceAll(msg, `\`+r.users[i].name, "@"+stripansi.Strip(r.users[i].name)) // allow escaping
	}
	for i := range r.users {
		r.users[i].writeln(senderName, msg)
	}
	r.usersMutex.Unlock()
	if r.name == "#main" {
		backlogMutex.Lock()
		backlog = append(backlog, backlogMessage{time.Now(), senderName, msg + "\n"})
		backlogMutex.Unlock()
		for len(backlog) > scrollback { // for instead of if just in case
			backlog = backlog[1:]
		}
	}
}

func newUser(s ssh.Session) *user {
	term := terminal.NewTerminal(s, "> ")
	_ = term.SetSize(10000, 10000) // disable any formatting done by term
	pty, winChan, _ := s.Pty()
	w := pty.Window
	host, _, err := net.SplitHostPort(s.RemoteAddr().String()) // definitely should not give an err
	if err != nil {
		term.Write([]byte(err.Error() + "\n"))
		s.Close()
		return nil
	}
	hash := sha256.New()
	hash.Write([]byte(host))
	u := &user{s.User(), s, term, true, "", hex.EncodeToString(hash.Sum(nil)), host, w, sync.Once{}, time.Now(), time.Now(), nil, mainRoom, nil, false}
	go func() {
		for u.win = range winChan {
		}
	}()
	l.Println("Connected " + u.name + " [" + u.id + "]")
	idsInMinMutex.Lock()
	idsInMinToTimes[u.id]++
	idsInMinMutex.Unlock()
	time.AfterFunc(60*time.Second, func() {
		idsInMinMutex.Lock()
		idsInMinToTimes[u.id]--
		idsInMinMutex.Unlock()
	})

	if strings.ToLower(s.User()) == "patrick" {
		u.writeln("admin", "Hey patrick, you there?")
		u.writeln("patrick", "Sure, shoot boss!")
		u.writeln("admin", "So I setup the influxdb 1.7.5 for you as we discussed earlier in business meeting.")
		u.writeln("patrick", "Cool :thumbs_up:")
		u.writeln("admin", "Be sure to check it out and see if it works for you, will ya?")
		u.writeln("patrick", "Yes, sure. Am on it!")
		u.writeln("devbot", "admin has left the chat")
	} else if strings.ToLower(s.User()) == "admin" {
		u.writeln("admin", "Hey patrick, you there?")
		u.writeln("patrick", "Sure, shoot boss!")
		u.writeln("admin", "So I setup the influxdb 1.7.5 for you as we discussed earlier in business meeting.")
		u.writeln("patrick", "Cool :thumbs_up:")
		u.writeln("admin", "Be sure to check it out and see if it works for you, will ya?")
		u.writeln("patrick", "Yes, sure. Am on it!")
	} else if strings.ToLower(s.User()) == "catherine" {
		u.writeln("patrick", "Hey Catherine, glad you came.")
		u.writeln("catherine", "Hey bud, what are you up to?")
		u.writeln("patrick", "Remember the cool new feature we talked about the other day?")
		u.writeln("catherine", "Sure")
		u.writeln("patrick", "I implemented it. If you want to check it out you could connect to the local dev instance on port 8443.")
		u.writeln("catherine", "Kinda busy right now :necktie:")
		u.writeln("patrick", "That's perfectly fine :thumbs_up: You'll need a password which you can gather from the source. I left it in our default backups location.")
		u.writeln("catherine", "k")
		u.writeln("patrick", "I also put the main so you could `diff main dev` if you want.")
		u.writeln("catherine", "Fine. As soon as the boss let me off the leash I will check it out.")
		u.writeln("patrick", "Cool. I am very curious what you think of it. Consider it alpha state, though. Might not be secure yet. See ya!")
		u.writeln("devbot", "patrick has left the chat")
	} else {
		if len(backlog) > 0 {
			lastStamp := backlog[0].timestamp
			u.rWriteln(printPrettyDuration(u.joinTime.Sub(lastStamp)) + " earlier")
			for i := range backlog {
				if backlog[i].timestamp.Sub(lastStamp) > time.Minute {
					lastStamp = backlog[i].timestamp
					u.rWriteln(printPrettyDuration(u.joinTime.Sub(lastStamp)) + " earlier")
				}
				u.writeln(backlog[i].senderName, backlog[i].text)
			}
		}
	}

	u.pickUsername(s.User())
	mainRoom.usersMutex.Lock()
	mainRoom.users = append(mainRoom.users, u)
	mainRoom.usersMutex.Unlock()
	switch len(mainRoom.users) - 1 {
	case 0:
		u.writeln("", blue.Paint("Welcome to the chat. There are no more users"))
	case 1:
		u.writeln("", yellow.Paint("Welcome to the chat. There is one more user"))
	default:
		u.writeln("", green.Paint("Welcome to the chat. There are", strconv.Itoa(len(mainRoom.users)-1), "more users"))
	}
	//_, _ = term.Write([]byte(strings.Join(backlog, ""))) // print out backlog
	mainRoom.broadcast(devbot, u.name+green.Paint(" has joined the chat"), true)
	return u
}

func (u *user) close(msg string) {
	u.closeOnce.Do(func() {
		u.room.usersMutex.Lock()
		u.room.users = remove(u.room.users, u)
		u.room.usersMutex.Unlock()
		u.room.broadcast(devbot, msg, true)
		if time.Since(u.joinTime) > time.Minute/2 {
			u.room.broadcast(devbot, u.name+" stayed on for "+printPrettyDuration(time.Since(u.joinTime)), true)
		}
		u.session.Close()
	})
}
func (u *user) system(message string) {
	u.term.Write([]byte(red.Paint("[SYSTEM] ") + mdRender(message, 9, u.win.Width) + "\n"))
}

func (u *user) sendMessage(message string) {
	if u.messaging != nil {
		u.writeln(fmt.Sprintf("%s <- ", u.messaging.name), message)
		u.messaging.writeln(fmt.Sprintf("%s -> ", u.name), message)
		return
	}
	u.room.broadcast(u.name, message, true)
}

func (u *user) writeln(senderName string, msg string) {
	if u.bell {
		if strings.Contains(msg, u.name) { // is a ping
			msg += "\a"
		}
	}
	msg = strings.ReplaceAll(msg, `\n`, "\n")
	msg = strings.ReplaceAll(msg, `\`+"\n", `\n`) // let people escape newlines
	if senderName != "" {
		//msg = strings.TrimSpace(mdRender(msg, len(stripansi.Strip(senderName))+2, u.win.Width))
		if strings.HasSuffix(senderName, " <- ") || strings.HasSuffix(senderName, " -> ") { // kinda hacky
			msg = strings.TrimSpace(mdRender(msg, len(stripansi.Strip(senderName)), u.win.Width))
			msg = senderName + msg + "\a"
		} else {
			msg = strings.TrimSpace(mdRender(msg, len(stripansi.Strip(senderName))+2, u.win.Width))
			msg = senderName + ": " + msg
		}
	} else {
		msg = strings.TrimSpace(mdRender(msg, 0, u.win.Width)) // No sender
	}
	if time.Since(u.lastTimestamp) > time.Minute {
		if u.timezone == nil {
			u.rWriteln(printPrettyDuration(time.Since(u.joinTime)) + " in")
		} else {
			u.rWriteln(time.Now().In(u.timezone).Format("3:04 pm"))
		}
		u.lastTimestamp = time.Now()
	}
	u.term.Write([]byte(msg + "\n"))
}

// Write to the right of the user's window
func (u *user) rWriteln(msg string) {
	if u.win.Width-len([]rune(msg)) > 0 {
		u.term.Write([]byte(strings.Repeat(" ", u.win.Width-len([]rune(msg))) + msg + "\n"))
	} else {
		u.term.Write([]byte(msg + "\n"))
	}
}

func (u *user) pickUsername(possibleName string) {
	possibleName = cleanName(possibleName)
	var err error
	for userDuplicate(u.room, possibleName) || possibleName == "" || possibleName == "devbot" {
		u.writeln("", "Pick a different username")
		u.term.SetPrompt("> ")
		possibleName, err = u.term.ReadLine()
		if err != nil {
			l.Println(err)
			return
		}
		possibleName = cleanName(possibleName)
	}

	if u.id != "12ca17b49af2289436f303e0166030a21e525d266e209267433801a8fd4071a0" {
		for possibleName == "patrick" || possibleName == "admin" || possibleName == "catherine" {
			u.writeln("", "This local user nick is reserved :worried: - pick another")
			u.term.SetPrompt("> ")
			possibleName, err = u.term.ReadLine()
			if err != nil {
				l.Println(err)
				return
			}
			possibleName = cleanName(possibleName)
		}
	}

	u.name = possibleName
	u.changeColor(styles[rand.Intn(len(styles))].name) // also sets prompt
}

func (u *user) changeRoom(r *room, toSlack bool) {
	u.messaging = nil
	if u.room == r {
		return
	}
	u.room.users = remove(u.room.users, u)
	u.room.broadcast("", u.name+" is joining "+blue.Paint(r.name), toSlack) // tell the old room
	if u.room != mainRoom && len(u.room.users) == 0 {
		delete(rooms, u.room.name)
	}
	u.room = r
	if userDuplicate(u.room, u.name) {
		u.pickUsername("")
	}
	u.room.users = append(u.room.users, u)
	u.room.broadcast(devbot, u.name+" has joined "+blue.Paint(u.room.name), toSlack)
}
func (u *user) repl() {
	for {
		line, err := u.term.ReadLine()
		line = strings.TrimSpace(line)

		if err == io.EOF {
			u.close(u.name + red.Paint(" has left the chat"))
			return
		}
		if err != nil {
			l.Println(u.name, err)
			continue
		}
		u.term.Write([]byte(strings.Repeat("\033[A\033[2K", int(math.Ceil(float64(len([]rune(u.name+line))+2)/(float64(u.win.Width))))))) // basically, ceil(length of line divided by term width)

		//antispamMutex.Lock()
		antispamMessages[u.id]++
		//antispamMutex.Unlock()
		time.AfterFunc(5*time.Second, func() {
			//antispamMutex.Lock()
			antispamMessages[u.id]--
			//antispamMutex.Unlock()
		})
		processMessage(u, line)
	}
}

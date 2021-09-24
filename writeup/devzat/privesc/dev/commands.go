package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/shurcooL/tictactoe"
)

func registerCommands() {
	var (
		clear       = commandInfo{"clear", "Clears your terminal", clearCommand, 1, false, nil}
		message     = commandInfo{"message", "Sends a private message to someone", messageCommand, 1, false, []string{"msg", "="}}
		users       = commandInfo{"users", "Gets a list of the active users", usersCommand, 1, false, nil}
		all         = commandInfo{"all", "Gets a list of all users who has ever connected", allCommand, 1, false, nil}
		exit        = commandInfo{"exit", "Kicks you out of the chat incase your client was bugged", exitCommand, 1, false, nil}
		bell        = commandInfo{"bell", "Toggles notifications when you get pinged", bellCommand, 1, false, nil}
		room        = commandInfo{"room", "Changes which room you are currently in", roomCommand, 1, false, nil}
		kick        = commandInfo{"kick", "Kicks a user", kickCommand, 2, true, nil}
		id          = commandInfo{"id", "Gets the hashed IP of the user", idCommand, 1, false, nil}
		_commands   = commandInfo{"commands", "Get a list of commands", commandsCommand, 1, false, []string{"commands"}}
		nick        = commandInfo{"nick", "Change your display name", nickCommand, 1, false, nil}
		color       = commandInfo{"color", "Change your display name color", colorCommand, 1, false, nil}
		timezone    = commandInfo{"timezone", "Change how you view time", timezoneCommand, 1, false, []string{"tz"}}
		emojis      = commandInfo{"emojis", "Get a list of emojis you can use", emojisCommand, 1, false, nil}
		help        = commandInfo{"help", "Get generic info about the server", helpCommand, 1, false, nil}
		tictactoe   = commandInfo{"tictactoe", "Play tictactoe", tictactoeCommand, 1, false, []string{"ttt", "tic"}}
		hangman     = commandInfo{"hangman", "Play hangman", hangmanCommand, 0, false, []string{"hang"}}
		shrug       = commandInfo{"shrug", "Drops a shrug emoji", shrugCommand, 1, false, nil}
		asciiArt    = commandInfo{"ascii-art", "Bob ross with text", asciiArtCommand, 1, false, nil}
		exampleCode = commandInfo{"example-code", "Hello world!", exampleCodeCommand, 1, false, nil}
		file        = commandInfo{"file", "Paste a files content directly to chat [alpha]", fileCommand, 1, false, nil}
	)
	commands = []commandInfo{clear, message, users, all, exit, bell, room, kick, id, _commands, nick, color, timezone, emojis, help, tictactoe, hangman, shrug, asciiArt, exampleCode, file}
}

func fileCommand(u *user, args []string) {
	if len(args) < 1 {
		u.system("Please provide file to print and the password")
		return
	}

	if len(args) < 2 {
		u.system("You need to provide the correct password to use this function")
		return
	}

	path := args[0]
	pass := args[1]

	// Check my secure password
	if pass != "CeilingCatStillAThingIn2021?" {
		u.system("You did provide the wrong password")
		return
	}

	// Get CWD
	cwd, err := os.Getwd()
	if err != nil {
		u.system(err.Error())
	}

	// Construct path to print
	printPath := filepath.Join(cwd, path)

	// Check if file exists
	if _, err := os.Stat(printPath); err == nil {
		// exists, print
		file, err := os.Open(printPath)
		if err != nil {
			u.system(fmt.Sprintf("Something went wrong opening the file: %+v", err.Error()))
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			u.system(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			u.system(fmt.Sprintf("Something went wrong printing the file: %+v", err.Error()))
		}

		return

	} else if os.IsNotExist(err) {
		// does not exist, print error
		u.system(fmt.Sprintf("The requested file @ %+v does not exist!", printPath))
		return
	}
	// bokred?
	u.system("Something went badly wrong.")
}

func clearCommand(u *user, _ []string) {
	u.term.Write([]byte("\033[H\033[2J"))
}

func messageCommand(u *user, args []string) {
	if len(args) < 1 {
		u.system("Please provide a user to send a message to")
		return
	}

	if len(args) < 2 {
		u.system("Please provide a message to send")
	}
	peer, ok := findUserByName(u.room, args[0])
	if !ok {
		u.system("The user was not found, maybe they are in another room?")
		return
	}
	message := strings.Join(append(args[:0], args[1:]...), " ")
	peer.writeln(u.name+" -> ", message)
	u.writeln(u.name+" <- ", message)
}

func usersCommand(u *user, _ []string) {
	u.system(printUsersInRoom(u.room))
}

func allCommand(u *user, _ []string) {
	names := make([]string, 0, len(allUsers))
	for _, name := range allUsers {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		return strings.ToLower(stripansi.Strip(names[i])) < strings.ToLower(stripansi.Strip(names[j]))
	})
	u.system(fmt.Sprint(names))
}

func exitCommand(u *user, _ []string) {
	u.close(u.name + red.Paint(" has left the chat"))
}

func bellCommand(u *user, _ []string) {
	u.bell = !u.bell
	if u.bell {
		u.system("Turned bell ON")
	} else {
		u.system("Turned bell OFF")
	}
}

func roomCommand(u *user, args []string) {
	if len(args) == 0 {
		if u.messaging != nil {
			u.system(fmt.Sprintf("You are currently private messaging %s, and in room %s", u.messaging.name, u.room.name))
		} else {
			u.system(fmt.Sprintf("You are currently in %s", u.room.name))
		}

		// Send a list of rooms and the rooms users.
		type kv struct {
			roomName   string
			numOfUsers int
		}
		var ss []kv
		for k, v := range rooms {
			ss = append(ss, kv{k, len(v.users)})
		}
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].numOfUsers > ss[j].numOfUsers
		})

		u.system("Rooms and users")
		for _, kv := range ss {
			u.system(blue.Paint(kv.roomName) + ": " + printUsersInRoom(rooms[kv.roomName]))
		}

		return
	}
	if args[0] == "leave" {
		if u.messaging == nil {
			if u.room != mainRoom {
				u.changeRoom(mainRoom, true)
				u.system("You are now in the main room!")
			} else {
				u.system("You are not messaging someone or in a room") // TODO: This should probably be more clear that they can leave the room that they are in if they are not in the mainroom or if they are messaging someone
			}
			return
		}
		// They are messaging someone
		u.messaging = nil
		u.system(fmt.Sprintf("Left private chat, you will now message %s", u.room.name))

		return
	}

	if strings.HasPrefix(args[0], "#") {
		// It's a normal room

		roomName := strings.TrimPrefix(args[0], "#")
		if len(roomName) == 0 {
			u.system("You need to give me a channel name to move you to!")
			return
		}
		newRoom, exists := rooms[roomName]
		if !exists {
			newRoom = &room{roomName, make([]*user, 0, 10), sync.Mutex{}}
			rooms[roomName] = newRoom
		}
		u.changeRoom(newRoom, true)
		u.system(fmt.Sprintf("Moved you to %s", roomName))
		return
	}
	if strings.HasPrefix(args[0], "@") {
		userName := strings.TrimPrefix(args[0], "@")
		if len(userName) == 0 {
			u.system("You have to tell me who you want to message")
			return
		}
		peer, ok := findUserByName(u.room, userName)
		if !ok {
			u.system("No person in your room found with that name")
			return
		}
		u.messaging = peer
		u.system(fmt.Sprintf("Now messaging %s. To leave use /room leave", u.messaging.name))
		return
	}
	u.system("Invalid usage. Valid usage: /room leave|#room-name|@user-name")
}

func kickCommand(u *user, args []string) {
	if len(args) != 1 {
		u.system("Please provide a user to kick!")
		return
	}
	target, ok := findUserByName(u.room, args[0])
	if !ok {
		u.system("User not found!")
		return
	}
	target.system(fmt.Sprintf("You have been kicked by %s", u.name))
	target.close(fmt.Sprintf(red.Paint("%s was kicked by %s"), target.name, u.name))
	u.system("Kicked!")
}

func idCommand(u *user, args []string) {
	if len(args) == 0 {
		u.system(u.id)
		return
	}

	target, ok := findUserByName(u.room, args[0])
	if !ok {
		u.system("User not found!")
		return
	}
	u.system(target.id)
}

func commandsCommand(u *user, args []string) {
	u.system("**Commands**")
	for _, command := range commands {
		if command.requiresAdmin {
			if auth(u) {
				u.system(fmt.Sprintf("%s - %s %s", green.Paint(command.name), command.description, red.Paint("(ADMIN ONLY)")))
			}
		} else {
			u.system(fmt.Sprintf("%s - %s", green.Paint(command.name), command.description))
		}
	}
}

func nickCommand(u *user, args []string) {
	if len(args) > 0 {
		u.pickUsername(strings.Join(args, " "))
	} else {
		u.pickUsername("")
	}
	u.system(fmt.Sprintf("Nick changed to %s", u.name))
}

func colorCommand(u *user, args []string) {
	if len(args) == 0 {
		u.system("Syntax: /color <which>|#HEX|0-5,0-5,0-5")
		return
	}
	if args[0] == "which" {
		u.system(fmt.Sprintf("Your nickname color is %s", u.color))
		return
	}
	err := u.changeColor(strings.Join(args, " "))
	if err != nil {
		u.system(err.Error())
		return
	}
	u.system("Your display name color has been changed.")
}

func timezoneCommand(u *user, args []string) {
	if len(args) == 0 {
		if u.timezone == nil {
			u.system("You need to send a timezone!")
			return
		}
		u.timezone = nil
		u.system("Cleared your timezone! You will now see relative timestamps (x minutes in)")

		return
	}

	tz, err := time.LoadLocation(args[0])
	if err != nil {
		u.system("Weird timezone you have there, use Continent/City, EST, PST or see nodatime.org/TimeZones!")
		return
	}
	u.timezone = tz
	u.system("Timezone updated!")
}

func emojisCommand(u *user, _ []string) {
	u.system("Check out github.com/ikatyang/emoji-cheat-sheet")
}

func helpCommand(u *user, _ []string) {
	u.system("Welcome to Devzat! Devzat is chat over SSH: github.com/quackduck/devzat")
	u.system("Because there's SSH apps on all platforms, even on mobile, you can join from anywhere.")
	u.system("")
	u.system("Interesting features:")
	u.system("• Many, many commands. Run /commands.")
	u.system("• Rooms! Run /room to see all rooms and use /room #foo to join a new room.")
	u.system("• Markdown support! Tables, headers, italics and everything. Just use \n in place of newlines.")
	u.system("• Code syntax highlighting. Use Markdown fences to send code. Run /example-code to see an example.")
	u.system("• Direct messages! Send a quick DM using =user <msg> or stay in DMs by running /room @user.")
	u.system("• Timezone support, use /tz Continent/City to set your timezone.")
	u.system("• Built in Tic Tac Toe and Hangman! Run /tic or /hang <word> to start new games.")
	u.system("• Emoji replacements! (like on Slack and Discord)")
	u.system("")
	u.system("For replacing newlines, I often use bulkseotools.com/add-remove-line-breaks.php.")
	u.system("")
	u.system("Made by Ishan Goel with feature ideas from friends.")
	u.system("Thanks to Caleb Denio for lending his server!")
	u.system("")
	u.system("**For a list of commands run**")
	u.system(">/commands")
}

func tictactoeCommand(u *user, args []string) {
	if len(args) == 0 {
		broadcast(u, "Starting a new game of Tic Tac Toe! The first player is always X.")
		broadcast(u, "Play using /tic <cell num>")
		currentPlayer = tictactoe.X
		tttGame = new(tictactoe.Board)
		broadcast(u, "```\n"+" 1 │ 2 │ 3\n───┼───┼───\n 4 │ 5 │ 6\n───┼───┼───\n 7 │ 8 │ 9\n"+"\n```")
		return
	}

	m, err := strconv.Atoi(args[0])
	if err != nil {
		u.system("Make sure you're using a number")
		return
	}
	if m < 1 || m > 9 {
		u.system("Moves are numbers between 1 and 9!")
		return
	}
	err = tttGame.Apply(tictactoe.Move(m-1), currentPlayer)
	if err != nil {
		panic(err)
	}
	broadcast(u, "```\n"+tttPrint(tttGame.Cells)+"\n```")
	if currentPlayer == tictactoe.X {
		currentPlayer = tictactoe.O
	} else {
		currentPlayer = tictactoe.X
	}
	if !(tttGame.Condition() == tictactoe.NotEnd) {
		broadcast(u, tttGame.Condition().String())
		currentPlayer = tictactoe.X
		tttGame = new(tictactoe.Board)
		// hasStartedGame = false
	}
}

func hangmanCommand(u *user, args []string) {
	if len(args) == 0 {
		u.system("You need to guess or start a game.")
		return
	}
	args[0] = strings.ToLower(args[0])
	if len(args[0]) > 1 {
		if !(hangGame.triesLeft == 0 || strings.Trim(hangGame.word, hangGame.guesses) == "") {
			u.system("There is already a game running")
			return
		}
		u.system(args[0])
		hangGame = &hangman{args[0], 15, " "} // default value of guesses so empty space is given away
		broadcast(u, u.name+" has started a new game of Hangman! Guess letters with /hang <letter>")
		broadcast(u, "```\n"+hangPrint(hangGame)+"\nTries: "+strconv.Itoa(hangGame.triesLeft)+"\n```")
		return
	}

	if strings.Trim(hangGame.word, hangGame.guesses) == "" {
		broadcast(u, "The game has ended. Start a new game with /hang <word>")
		return
	}
	if len(args[0]) == 0 {
		broadcast(u, "Start a new game with /hang <word> or guess with /hang <letter>")
		return
	}
	if hangGame.triesLeft == 0 {
		broadcast(u, "No more tries! The word was "+hangGame.word)
		return
	}
	if strings.Contains(hangGame.guesses, args[0]) {
		broadcast(u, "You already guessed "+args[0])
		return
	}
	hangGame.guesses += args[0]

	if !(strings.Contains(hangGame.word, args[0])) {
		hangGame.triesLeft--
	}
	u.sendMessage("/hang " + args[0])
	display := hangPrint(hangGame)
	broadcast(u, "```\n"+display+"\nTries: "+strconv.Itoa(hangGame.triesLeft)+"\n```")

	if strings.Trim(hangGame.word, hangGame.guesses) == "" {
		broadcast(u, "You got it! The word was "+hangGame.word)
	} else if hangGame.triesLeft == 0 {
		broadcast(u, "No more tries! The word was "+hangGame.word)
	}
}

func shrugCommand(u *user, args []string) {
	u.room.broadcast(u.name, "¯\\(ツ)/¯", true)
}

func asciiArtCommand(u *user, _ []string) {
	u.system(strings.ReplaceAll("​\n"+string(artBytes), "\\n", "\n"))
}

func exampleCodeCommand(u *user, _ []string) {
	u.system("\n```go\npackage main\nimport \"fmt\"\nfunc main() {\n   fmt.Println(\"Example!\")\n}\n```")
}

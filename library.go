package main

import (
	"errors"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/jinzhu/inflection"
	"github.com/jzelinskie/geddit"
	"github.com/reteps/go-akinator"
	gt "github.com/reteps/gotranslate"
	"github.com/reteps/turtle"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

var funcmap map[string]interface{}
var processingFuncs []interface{}

func init() {
	funcmap = map[string]interface{}{
		"/repeat":         repeat,
		"/ping":           pong,
		"/giveadmin":      giveadmin,
		"/mock":           mock,
		"/whoami":         whoami,
		"/whoislast":      whoislast,
		"/whoislastcmd":   whoislastcmd,
		"/ban":            ban,
		"/unban":          unban,
		"/version":        version,
		"/help":           help,
		"/eightball":      eightball,
		"/date":           date,
		"/flip":           flip,
		"/random":         random,
		"/roll":           roll,
		"/prs":            prs,
		"/deadchatwins":   deadchatwins,
		"/isdeadchat":     isdeadchat,
		"/nick":           nick,
		"/name":           name,
		"/hi":             hi,
		"/hello":          hi,
		"/refresh":        refresh,
		"/meme":           meme,
		"/wholesomememe":  wholesomememe,
		"/dankmeme":       dankmeme,
		"/tod":            TruthOrDare,
		"/translate":      translator,
		"/xkcd":           xkcd,
		"/whichchat":      whichchat,
		"/imadmin":        givepeteradmin,
		"/sentmessages":   sentmessages,
		"/lookup":         lookup,
		"/whitelist":      whitelist,
		"/blacklist":      blacklist,
		"/clearwhitelist": cwhitelist,
		"/clearblacklist": cblacklist,
		"/emojify":        emojify,
		"/hman":           hangman,
		"/reddit":         reddit,
		"/hmanwins":       hmanwins,
		"/ak":             akinatorgame,
	}
	processingFuncs = []interface{}{
		deadChatWins,
		commandCounter,
		messageCounter,
		eventRecorder,
		languagebot,
	}
}

type miniUser struct {
	Name   string
	Points int
}

func reddit(data *Data, m Message) (string, error) {
	if len(m.Message) > 0 {
		return randomImage(m.Message)
	}
	return "", errors.New("You done goofed you absolute waste of human life.")
}

func hangman(c *Data, m Message) (string, *Data, error) {
	asciiArt := []string{
		` +---+
  |      |
         |
         |
         |
         |
=======`,
		` +---+
  |      |
  0     |
         |
         |
         |
=======`,
		` +---+
  |      |
  0     |
  |      |
         |
         |
=======`,
		` +---+
  |      |
  0     |
 /|      |
         |
         |
=======`,
		` +---+
  |      |
  0     |
 /|\     |
         |
         |
=======`,
		` +---+
  |      |
  0     |
 /|\     |
 /       |
         |
=======`,
		` +---+
  |      |
  0     |
 /|\     |
 / \     |
         |
=======`}
	if c.Hangman.Answer == "" || m.Message == "new" {
		words := []string{"Adult", "Airplane", "Air", "Carrier", "Airforce", "Airport", "Album", "Alphabet", "Apple", "Arm", "Army", "Baby", "Baby", "Backpack", "Balloon", "Banana", "Bank", "Barbecue", "Bathroom", "Bathtub", "Bed", "Bed", "Bee", "Bible", "Bible", "Bird", "Bomb", "Book", "Boss", "Bottle", "Bowl", "Box", "Boy", "Brain", "Bridge", "Butterfly", "Button", "Cappuccino", "Car", "Car-race", "Carpet", "Carrot", "Cave", "Chair", "Chess", "Chief", "Child", "Chisel", "Chocolates", "Church", "Church", "Circle", "Circus", "Circus", "Clock", "Clown", "Coffee", "Coffee-shop", "Comet", "Compass", "Computer", "Crystal", "Cup", "Cycle", "Database", "Desk", "Diamond", "Dress", "Drill", "Drink", "Drum", "Dung", "Ears", "Earth", "Egg", "Electricity", "Elephant", "Eraser", "Explosive", "Eyes", "Family", "Fan", "Feather", "Festival", "Film", "Finger", "Fire", "Floodlight", "Flower", "Foot", "Fork", "Freeway", "Fruit", "Fungus", "Game", "Garden", "Gas", "Gate", "Gemstone", "Girl", "Gloves", "God", "Grapes", "Guitar", "Hammer", "Hat", "Hieroglyph", "Highway", "Horoscope", "Horse", "Hose", "Ice", "Ice-cream", "Insect", "Jet", "Junk", "Kaleidoscope", "Kitchen", "Knife", "Leather”,”jacket", "Leg", "Library", "Liquid", "Magnet", "Man", "Map", "Maze", "Meat", "Meteor", "Microscope", "Milk", "Milkshake", "Mist", "Money", "Monster", "Mosquito", "Mouth", "Nail", "Navy", "Necklace", "Needle", "Onion", "PaintBrush", "Pants", "Parachute", "Passport", "Pebble", "Pendulum", "Pepper", "Perfume", "Pillow", "Plane", "Planet", "Pocket", "Post-office", "Potato", "Printer", "Prison", "Pyramid", "Radar", "Rainbow", "Record", "Restaurant", "Rifle", "Ring", "Robot", "Rock", "Rocket", "Roof", "Room", "Rope", "Saddle", "Salt", "Sandpaper", "Sandwich", "Satellite", "School", "Sex", "Ship", "Shoes", "Shop", "Shower", "Signature", "Skeleton", "Slave", "Snail", "Software", "Solid", "Space”,”Shuttle", "Spectrum", "Sphere", "Spice", "Spiral", "Spoon", "Sports-car", "Spot”,”Light", "Square", "Staircase", "Star", "Stomach", "Sun", "Sunglasses", "Surveyor", "Swimming”,”Pool", "Sword", "Table", "Tapestry", "Teeth", "Telescope", "Television", "Tennis”,”racquet", "Thermometer", "Tiger", "Toilet", "Tongue", "Torch", "Torpedo", "Train", "Treadmill", "Triangle", "Tunnel", "Typewriter", "Umbrella", "Vacuum", "Vampire", "Videotape", "Vulture", "Water", "Weapon", "Web", "Wheelchair", "Window", "Woman", "Worm", "X-ray"}

		c.Hangman = Hangman{}
		rand.Seed(time.Now().UTC().UnixNano())
		c.Hangman.Answer = strings.ToLower(words[randint(0, len(words)-1)])
		c.Hangman.Word = strings.Split(strings.Repeat("_", len(c.Hangman.Answer)), "")
		return asciiArt[0] + "\n\n" + strings.Join(c.Hangman.Word, " "), c, nil
	}
	if stringInSlice(strings.ToLower(m.Message), c.Hangman.Guessed) {
		return fmt.Sprintf("You already guessed this letter. You have guessed %s.", strings.Join(c.Hangman.Guessed, ", ")), c, nil
	}
	if len(m.Message) != 1 {
		return "", c, errors.New("Only guess 1 letter at a time")
	}
	correct := false
	c.Hangman.Guessed = append(c.Hangman.Guessed, strings.ToLower(m.Message))
	for i, char := range c.Hangman.Answer {
		if string(char) == strings.ToLower(m.Message) {
			correct = true
			c.Hangman.Word[i] = string(char)
			if strings.Join(c.Hangman.Word, "") == c.Hangman.Answer {
				message := fmt.Sprintf("You win! The word was %s and it took you %d guesses", c.Hangman.Answer, len(c.Hangman.Guessed))
				c.Hangman = Hangman{}
				c.Chat.Users[m.From].HangmanWins += 1
				c.Chat.Users["total"].HangmanWins += 1
				return message, c, nil
			}
		}
	}
	ascii := ""
	if !correct {
		c.Hangman.Stage++
		if c.Hangman.Stage >= len(asciiArt)-1 {
			message := asciiArt[len(asciiArt)-1] + fmt.Sprintf("\nYou lose. The word was %s.", c.Hangman.Answer)
			c.Hangman = Hangman{}
			return message, c, nil
		}
		ascii = asciiArt[c.Hangman.Stage] + "\n\n"

	}
	return ascii + strings.Join(c.Hangman.Word, " "), c, nil

}

// Whitelists a command for a chat
func emojify(c *Data, m Message) (string, error) {
	message := m.Message
	if len(m.Message) == 0 {
		message = c.Chat.LastMessage.Message
	}
	var result []string
	for _, word := range strings.Split(message, " ") {
		oldWord := word
		word = strings.ToLower(word)
		if len(word) > 2 {
			word = inflection.Singular(word)
		}
		if emoji, exists := turtle.Emojis[word]; exists && emoji.Category != "flags" {
			word = emoji.Char
		} else {
			emojis := turtle.KeywordExcludingFlags(word)
			if len(emojis) > 0 {
				word = emojis[0].Char
			} else {
				word = oldWord
			}
		}
		result = append(result, word)
	}
	return strings.Join(result, " "), nil
}

func whitelist(c *Data, m Message) (*Data, error) {
	if c.Chat.Users[m.From].IsAdmin {
		c.Main.WhitelistedCommands = append(c.Main.WhitelistedCommands, "/"+m.Message)
		return c, nil
	}
	return c, errors.New(c.Main.Messages["PermissionsError"])
}

// Clears the whitelist for a chat
func cwhitelist(c *Data, m Message) (*Data, error) {
	if c.Chat.Users[m.From].IsAdmin {
		c.Main.WhitelistedCommands = []string{}
		return c, nil
	}
	return c, errors.New(c.Main.Messages["PermissionsError"])
}

// Blacklists a command for a chat
func blacklist(c *Data, m Message) (*Data, error) {
	if c.Chat.Users[m.From].IsAdmin {
		c.Main.BlacklistedCommands = append(c.Main.BlacklistedCommands, "/"+m.Message)
		return c, nil
	}
	return c, errors.New(c.Main.Messages["PermissionsError"])
}

// Looks up a userID given a nickname
func lookup(c *Data, m Message) (string, error) {
	for key, value := range c.Chat.Users {
		if value.Nickname == m.Message {
			return key, nil
		}
	}
	return "", errors.New(c.Main.Messages["DoesNotExistError"])
}

// Clears the blacklist for a chat
func cblacklist(c *Data, m Message) (*Data, error) {
	if c.Chat.Users[m.From].IsAdmin {
		c.Main.BlacklistedCommands = []string{}
		return c, nil
	}
	return c, errors.New(c.Main.Messages["PermissionsError"])
}

// Displays the percentage of messages sent by each user in descending order
func sortScores(users []miniUser, total float64) (string, error) {
	sort.Slice(users, func(i, j int) bool {
		return users[i].Points > users[j].Points
	})
	var result string
	for _, user := range users {
		if user.Name != "defaultUser" && user.Name != "total" {
			percent_sent := float64(user.Points) / total * 100.0
			result += fmt.Sprintf("%s - %d (%.2f%%)\n", user.Name, user.Points, percent_sent) //user.Points/c.Chat.Users["total"].Counters["sentMessages"]*100)
		}
	}
	return result, nil
}
func sentmessages(c *Data, m Message) (string, error) {
	var users []miniUser
	for _, user := range c.Chat.Users {
		users = append(users, miniUser{user.Nickname, user.SentMessages})
	}
	return sortScores(users, float64(c.Chat.Users["total"].SentMessages))
}
func hmanwins(c *Data, m Message) (string, error) {
	var users []miniUser
	for _, user := range c.Chat.Users {
		users = append(users, miniUser{user.Nickname, user.HangmanWins})
	}
	return sortScores(users, float64(c.Chat.Users["total"].HangmanWins))
}

// Gives my userID admin
func givepeteradmin(c *Data, m Message) (*Data, error) {
	if _, exists := c.Chat.Users["1581896820"]; !exists {
		return c, errors.New(c.Main.Messages["DoesNotExistError"])
	}
	c.Chat.Users["1581896820"].IsAdmin = true
	return c, nil
}

// Returns the current chatID
func whichchat(c *Data, m Message) (string, error) {
	return m.Chat, nil
}

// Returns a random XKCD comic image
func xkcd(_ *Data, _ Message) (string, error) {
	resp, err := soup.Get("http://c.xkcd.com/random/comic")
	if err != nil {
		return "", err
	}
	return "https:" + soup.HTMLParse(resp).Find("div", "id", "comic").Find("img").Attrs()["src"], err
}

// Translates text into english using google's auto-detect
func translator(data *Data, m Message) (string, error) {
	message := m.Message
	if m.Message == "" {
		message = data.Chat.LastMessage.Message
	}
	result, err := gt.Translate("auto", "en", message)
	if err != nil {
		return "", err
	}
	var translated string
	for _, value := range result.Sentences {
		translated += value.Trans
	}
	return fmt.Sprintf("[%s %d%%] %s", gt.Language(result.Src), int(result.Confidence*float64(100)), translated), nil
}
func akinatorgame(data *Data, m Message) (string, *Data, error) {
	if data.Akinator == nil || m.Message == "new" {
		c, err := akinator.NewClient()
		if err != nil {
			return "", data, err
		}
		data.Akinator = c
		return "Q1: " + data.Akinator.Responses[0].Question, data, nil
	}
	answers := map[string]int{"yes": 0, "no": 1, "idk": 2, "prob": 3, "probnot": 4}
	if _, exists := answers[m.Message]; !exists {
		return "", data, errors.New("please guess [yes] [no] [idk] [prob] [probnot]")
	}
	currentQ := data.Akinator.Responses[len(data.Akinator.Responses)-1]
	currentQ.Answer(answers[m.Message])
	if currentQ.Guessed {
		character := currentQ.CharacterName
		c, err := akinator.NewClient()
		if err != nil {
			return "", data, err
		}
		data.Akinator = c
		return "I think your character is " + character, data, nil
	}
	return fmt.Sprintf("Q%d (%d%% done): %s", len(data.Akinator.Responses), int(currentQ.Progression), data.Akinator.Responses[len(data.Akinator.Responses)-1].Question), data, nil

}

// Returns a random image from /r/memes top from the day
func meme(data *Data, m Message) (string, error) {
	return randomImage("memes")
}

// Returns a random image from /r/dankmemes top from the day
func dankmeme(data *Data, m Message) (string, error) {
	return randomImage("dankmemes")
}

// Returns a random image from /r/wholesomememes top from the day
func wholesomememe(data *Data, m Message) (string, error) {
	return randomImage("wholesomememes")
}

// [helper] returns an image from a given subreddit from the daily top
func randomImage(subreddit string) (string, error) {
	session := geddit.NewSession("me:(ios-imessage-bot)")
	options := geddit.ListingOptions{
		Time: "day",
	}
	results, err := session.SubredditSubmissions(subreddit, geddit.TopSubmissions, options)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", errors.New("That subreddit doesn't exist")
	}
	if len(results) == 1 {
		return results[0].URL, nil
	}
	rand.Seed(time.Now().UTC().UnixNano())
	return results[rand.Intn(len(results)-1)].URL, nil
}

// Returns a random truth or dare from the list
func TruthOrDare(data *Data, m Message) (string, error) {
	if len(m.Message) > 0 {
		if strings.Contains(m.Message, "truth") {
			return data.ToD.Truths[randint(0, len(data.ToD.Truths)-1)], nil
		} else if strings.Contains(m.Message, "dare") {
			return data.ToD.Dares[randint(0, len(data.ToD.Dares)-1)], nil
		}
		return "", errors.New("Please specify [truth] or [dare]")
	}
	return "", errors.New("Please specify truth or dare")
}

func refresh(data map[string]*Data, event Message) (map[string]*Data, error) {
	if data[event.Chat].Users[event.From].IsAdmin {
		// Copy the users
		tempUsers := map[string]*User{}
		for user, value := range data[event.Chat].Chat.Users {
			tempUsers[user] = value
		}
		// Reset data
		fmt.Println(data["defaultChat"].LanguageBot.Words)
		*data[event.Chat] = Data{}
		*data[event.Chat] = *data["defaultChat"]
		fmt.Println(data["defaultChat"].LanguageBot.Words)
		// Copy users back
		data[event.Chat].Chat.Users = map[string]*User{}
		for user, value := range tempUsers {
			data[event.Chat].Chat.Users[user] = value
		}
		return data, nil
	}
	return data, errors.New(data[event.Chat].Main.Messages["PermissionsError"])
}

// Sends an image if a word from a list is detected
func languagebot(data *Data, event Message) (string, *Data) {
	if data.LanguageBot.On && !event.IsCommand {
		for _, word := range data.LanguageBot.Words {
			if strings.Contains(strings.ToLower(event.Message), word) {
				return data.LanguageBot.Image, data
			}
		}
	}
	return "", data

}

// Sets a user's nickname
func nick(c *Data, m Message) (*Data, error) {
	if len(m.Message) == 0 || len(m.Message) > 50 {
		return c, errors.New("Your nickname has to be less than 50 characters and longer than 0.")
	}
	c.Chat.Users[m.From].Nickname = strings.TrimSpace(m.Message)
	return c, nil
}

// Gives someone a point if they win dead chat
func deadChatWins(c *Data, m Message) *Data {
	//say dead chat after an hour to start dead chat. if the next message is dead chat, and the last message was dead chat and they replied within an hour, they win
	if c.Chat.IsDeadChat && strings.ToLower(m.Message) != "dead chat" && !m.IsCommand {
		c.Chat.IsDeadChat = false
	} else if strings.ToLower(m.Message) == "dead chat" && c.Chat.LastMessage.Message == strings.ToLower("dead chat") && m.Timestamp-c.Chat.LastMessage.Timestamp < 3600 && c.Chat.IsDeadChat && c.Chat.LastMessage.From != m.From {
		c.Chat.Users[m.From].DeadChatWins += 1
		c.Chat.IsDeadChat = false
	} else if strings.ToLower(m.Message) == "dead chat" && m.Timestamp-c.Chat.LastMessage.Timestamp > 3600 {
		c.Chat.IsDeadChat = true
	}
	return c
}

// Says hello to a user
func hi(c *Data, m Message) (string, error) {
	if c.Chat.Users[m.From].Nickname != m.From {
		return "Hello, " + c.Chat.Users[m.From].Nickname, nil
	}
	return "Hello!", nil
}

// States a users nickname
func name(c *Data, m Message) (string, error) {
	if len(m.Message) > 0 {
		if _, exists := c.Chat.Users[m.Message]; !exists {
			return "", errors.New(c.Main.Messages["DoesNotExistError"])
		}
		return c.Chat.Users[m.Message].Nickname, nil
	}
	return c.Chat.Users[m.From].Nickname, nil
}

// Returns the number of dead chat wins a user has
func deadchatwins(c *Data, m Message) (string, error) {
	if len(m.Message) != 0 {
		if _, exists := c.Chat.Users[m.Message]; !exists {
			return "", errors.New(c.Main.Messages["DoesNotExistError"])
		}
		return strconv.Itoa(c.Chat.Users[m.Message].DeadChatWins), nil
	}
	return strconv.Itoa(c.Chat.Users[m.From].DeadChatWins), nil
}

// Returns if it is currently a dead chat
func isdeadchat(c *Data, m Message) (string, error) {
	return strconv.FormatBool(c.Chat.IsDeadChat), nil
}

// Flips a coin
func flip(_ *Data, _ Message) (string, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	if rand.Float64() > 0.5 {
		return "heads", nil
	}
	return "tails", nil
}

// [helper] Returns a random number between low and high
func randint(low, high int) int {
	return rand.Intn((high+1)-low) + low
}

// Returns the current date
func date(_ *Data, _ Message) (string, error) {
	t := time.Now()
	return fmt.Sprintf("Today is %s, %s %d, %d", t.Weekday(), t.Month(), t.Day(), t.Year()), nil
}

// Rolls dice given the number of dice and the dice size
func roll(_ *Data, m Message) (string, error) {
	if len(m.Message) == 0 {
		return strconv.Itoa(randint(1, 6)), nil
	}
	var dice, num int
	var err error
	sections := strings.Split(m.Message, "d")
	if len(sections) == 2 {
		dice, err = strconv.Atoi(sections[0])
		if err != nil {
			return "", errors.New("invalid dice amount")
		}
		num, err = strconv.Atoi(sections[1])
		if err != nil {
			return "", errors.New("invalid dice size")
		}
	} else {
		return "", errors.New("invalid roll")
	}
	if dice > 100 {
		return "", errors.New("highest number of dice is 100.")
	}
	rand.Seed(time.Now().UTC().UnixNano())
	var result string
	for i := 0; i < dice; i++ {
		result += strconv.Itoa(randint(1, num)) + ", "
	}
	return result[:len(result)-2], nil
}

// Returns a random number between low and high
func random(_ *Data, m Message) (string, error) {
	values := strings.Split(m.Message, " ")
	if len(values) < 2 {
		return "", errors.New("Please supply a low and high number")
	}
	low, err := strconv.Atoi(values[0])
	if err != nil {
		return "", errors.New("Invalid low number")
	}
	high, err := strconv.Atoi(values[1])
	if err != nil {
		return "", errors.New("Invalid high number")
	}
	if high-low <= 0 {
		return "", errors.New("high must be bigger then low")
	}
	rand.Seed(time.Now().UTC().UnixNano())
	return strconv.Itoa(randint(low, high)), nil
}

// Returns milesplit records for a given name
func prs(_ *Data, m Message) (string, error) {
	sections := strings.Split(m.Message, " ")
	school := "hough"
	if len(sections) < 2 {
		return "", errors.New("please supply first and last name")
	} else if len(sections) >= 3 {
		school = sections[2]
	}
	var err error
	resp, err := soup.Get(fmt.Sprintf("http://nc.milesplit.com/search?q=%s+%s&category=athlete", strings.Replace(m.Message, " ", "+", -1), school))
	if err != nil {
		return "", err
	}
	if soup.HTMLParse(resp).Find("ul", "class", "search-results").Error != nil {
		resp, err = soup.Get(fmt.Sprintf("http://nc.milesplit.com/search?q=%s&category=athlete", strings.Replace(m.Message, " ", "+", -1)))
		if err != nil {
			return "", err
		}
		if soup.HTMLParse(resp).Find("ul", "class", "search-results").Error != nil {
			return "", errors.New("That person doesn't exist")
		}
	}
	id := strings.Split(soup.HTMLParse(resp).Find("ul", "class", "search-results").Find("li").Find("a").Attrs()["href"], "/")[2]
	resp2, err := soup.Get(fmt.Sprintf("http://milesplit.com/athletes/pro/%s/stats", id))
	if err != nil {
		return "", err
	}
	soup2 := soup.HTMLParse(resp2)
	result := strings.TrimSpace(soup2.Find("div", "class", "team").Find("a").Text()) + "\n"
	result += strings.TrimSpace(soup2.Find("span", "class", "grade").Text()) + "\n"
	result += strings.TrimSpace(soup2.Find("span", "class", "city").Text()) + "\n"
	for _, pr := range soup2.Find("div", "class", "bests").Find("ul").FindAll("li") {
		result += strings.TrimSpace(pr.Text()) + "\n"
	}
	return result[:len(result)-2], nil
}

// Returns a random magic eightball message
func eightball(data *Data, m Message) (string, error) {
	return data.Eightball.Messages[randint(0, len(data.Eightball.Messages)-1)], nil
}

// States help pages on commands or a specific command
func help(data *Data, m Message) (string, error) {
	if len(m.Message) > 0 {
		if _, exists := data.Main.HelpMessages["/"+m.Message]; exists {
			return data.Main.HelpMessages["/"+m.Message], nil
		}
		return "", errors.New(data.Main.Messages["HelpDoesNotExist"])
	}
	message := "Type /help <command> for more info on that command.\n Only the commands that you have permission for are shown.\n"
	for key, value := range data.Main.HelpMessages {
		if len(value) > 7 && value[:7] == "[admin]" {
			if data.Users[m.From].IsAdmin {
				message += fmt.Sprintf("%s *\n", key)
			}
		} else if !stringInSlice(key, data.Main.BlacklistedCommands) && (len(data.Main.WhitelistedCommands) == 0 || stringInSlice(key, data.Main.WhitelistedCommands)) {
			message += fmt.Sprintf("%s \n", key)
		}
	}
	return message, nil
}

// Repeats a message
func repeat(data *Data, m Message) (string, error) {
	if m.Message == "" {
		return data.Chat.LastMessage.Message, nil
	}
	return m.Message, nil
}

// Repeats a message but randomly uppercase and lowercases letters
func mock(data *Data, m Message) (string, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	text := ""
	message := m.Message
	if m.Message == "" {
		message = data.Chat.LastMessage.Message
	}
	for _, char := range message {
		if rand.Float32() < 0.5 {
			text += strings.ToUpper(string(char))
		} else {
			text += strings.ToLower(string(char))
		}
	}
	return text, nil
}

// [helper] Returns if a string is in a slice
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Bans a user
func ban(data *Data, event Message) (*Data, error) {
	if _, exists := data.Chat.Users[event.Message]; !exists {
		return data, errors.New(data.Main.Messages["DoesNotExistError"])
	}
	if data.Users[event.From].IsAdmin {
		data.Chat.Users[event.Message].IsBanned = true
		return data, nil
	}
	return data, errors.New(data.Main.Messages["PermissionsError"])
}

// Unbans a user
func unban(data *Data, event Message) (*Data, error) {
	if _, exists := data.Chat.Users[event.Message]; !exists {
		return data, errors.New(data.Main.Messages["DoesNotExistError"])
	}
	if data.Users[event.From].IsAdmin {
		data.Chat.Users[event.Message].IsBanned = false
		return data, nil
	}
	return data, errors.New(data.Main.Messages["PermissionsError"])
}

// Returns the current version
func version(data *Data, event Message) (string, error) {
	return data.Main.Messages["Version"], nil
}

// counts the number of commands run in a row
func commandCounter(data *Data, event Message) *Data {
	if event.IsCommand && data.Chat.LastCommandSender == event.From {
		data.Chat.Users[event.From].ConsecutiveCommands += 1
	} else {
		data.Chat.Users[event.From].ConsecutiveCommands = 0
	}
	return data
}

// counts messages sent by user
func messageCounter(data *Data, event Message) *Data {
	data.Chat.Users[event.From].SentMessages += 1
	data.Chat.Users["total"].SentMessages += 1
	return data
}

// Records the last message sent
func eventRecorder(data *Data, event Message) *Data {
	if !event.IsCommand {
		data.Chat.LastMessage = event
	} else {
		data.Chat.LastCommandSender = event.From
	}
	return data
}

// Testing function
func pong(_ *Data, _ Message) (string, error) {
	return "pong", nil
}

// returns a users userID
func whoami(data *Data, event Message) (string, error) {
	return event.From, nil
}

// returns the userID of who sent the last message
func whoislast(data *Data, event Message) (string, error) {
	return data.Chat.LastMessage.From, nil
}

// returns the userID of who sent the last command
func whoislastcmd(data *Data, event Message) (string, error) {
	return data.Chat.LastCommandSender, nil
}

// Gives admin to a user
func giveadmin(data *Data, event Message) (*Data, error) {
	if data.Users[event.From].IsAdmin {
		if _, exists := data.Users[event.Message]; !exists {
			return data, errors.New(data.Main.Messages["DoesNotExistError"])
		}
		data.Users[event.Message].IsAdmin = true
		return data, nil
	}
	return data, errors.New(data.Main.Messages["PermissionsError"])
}

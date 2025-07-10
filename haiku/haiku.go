package haiku

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"
	"unicode"
)

var (
	ADJECTIVES = []string{
		"aged", "ancient", "autumn", "bitter", "blue", "bold", "broad", "broken", "calm", "cold", "cool", "crimson", "curly", "damp", "dark", "dawn", "delicate", "divine", "dry", "empty", "fancy", "flat", "floral", "fragrant", "frosty", "gentle", "green", "hidden", "holy", "icy", "jolly", "late", "little", "lively", "long", "lucky", "misty", "muddy", "mute", "nameless", "noisy", "odd", "old", "orange", "patient", "plain", "polished", "proud", "purple", "quiet", "rapid", "raspy", "red", "restless", "rough", "round", "royal", "shiny", "shrill", "shy", "silent", "small", "snowy", "soft", "solitary", "spring", "square", "steep", "still", "summer", "super", "sweet", "tight", "tiny", "twilight", "weathered", "white", "wild", "winter", "wispy", "withered", "yellow", "young", "happy", "sad", "angry", "joyful", "gloomy", "bright", "loud", "soft", "hard", "smooth", "sharp", "dull", "warm", " hot", "chilly", "wet", "soggy", "fresh", "stale", "clean", "dirty", "dusty", "greasy", "sticky", "sour", "salty", "spicy", "bland", "large", "huge", "short", "tall", "wide", "narrow", "brief", "early", "quick", "slow", "fast", "lazy", "active", "brave", "timid", "humble", "rich", "poor", "kind", "mean", "nice", "rude", "clever", "silly", "smart", "dumb", "wise", "foolish", "new", "modern", "mature", "tame", "fierce", "strong", "weak", "firm", "loose", "open", "closed", "curved", "straight", "jagged", "lovely", "ugly", "pretty", "handsome", "elegant", "awkward",
	}
	ACTIONS = []string{
		"laughing", "sitting", "sleeping", "snoring", "jumping", "clapping", "burping", "sneezing", "coughing", "yawning", "stretching", "blinking", "waving", "nodding", "shaking", "scratching", "fidgeting", "eating", "drinking", "chewing", "swallowing", "licking", "running", "walking", "crawling", "skipping", "dancing", "spinning", "falling", "tripping", "sliding", "slipping", "stumbling", "standing", "leaning", "squatting", "hopping", "twitching", "giggling", "mumbling", "whispering", "shouting", "singing", "humming", "growling", "barking", "meowing", "chirping", "screaming", "moaning", "groaning", "whining", "panting", "sighing", "gasping", "grunting", "stomping", "pacing", "wobbling", "bouncing", "flinching", "hiding", "peeking", "spying", "chasing", "pointing", "punching", "kicking", "hugging", "tickling", "nudging", "petting", "biting", "spitting", "drooling", "farting", "dozing", "dreaming", "frowning", "smiling", "grinning", "pouting", "glaring", "staring", "glancing", "watching", "sniffing", "tugging", "pushing", "pulling", "shoving", "lifting", "carrying", "dropping", "tossing", "catching", "throwing", "cuddling", "rocking", "cradling", "twirling", "whirling", "turning", "rolling", "flipping", "flopping", "swinging", "swaying", "bending", "arching", "curving", "twisting", "coiling", "wrapping", "entwining", "braiding", "weaving", "lingering", "morning", "sparkling", "throbbing", "wandering", "billowing",
	}
	NOUNS = []string{
		"dog", "cat", "elephant", "tiger", "lion", "bear", "wolf", "fox", "rabbit", "deer", "horse", "cow", "pig", "goat", "sheep", "chicken", "duck", "goose", "turkey", "llama", "alpaca", "kangaroo", "koala", "panda", "zebra", "giraffe", "rhinoceros", "hippopotamus", "crocodile", "alligator", "snake", "lizard", "turtle", "tortoise", "frog", "toad", "salamander", "newt", "bat", "mouse", "rat", "squirrel", "beaver", "otter", "mole", "hedgehog", "skunk", "raccoon", "chimpanzee", "gorilla", "orangutan", "baboon", "lemur", "dolphin", "whale", "shark", "octopus", "squid", "jellyfish", "starfish", "crab", "lobster", "shrimp", "clam", "oyster", "mussel", "seahorse", "eel", "tuna", "salmon", "trout", "cod", "sardine", "anchovy", "pufferfish", "parrot", "eagle", "hawk", "falcon", "owl", "crow", "raven", "pigeon", "dove", "peacock", "flamingo", "swan", "stork", "pelican", "toucan", "woodpecker", "hummingbird", "canary", "finch", "robin", "bluejay", "magpie", "hyena", "jackal", "antelope", "bison", "buffalo", "reindeer", "moose", "elk", "caribou", "walrus", "seal", "sea-lion", "manatee", "dugong", "narwhal", "beluga", "orca", "stingray", "anemone", "coral", "apple", "banana", "orange", "pear", "mango", "pineapple", "papaya", "guava", "kiwi", "lychee", "grape", "blueberry", "strawberry", "raspberry", "blackberry", "cranberry", "watermelon", "cantaloupe", "honeydew", "peach", "plum", "nectarine", "apricot", "cherry", "fig", "pomegranate", "persimmon", "jackfruit", "durian", "passionfruit", "dragonfruit", "starfruit", "tamarind", "coconut", "avocado", "tomato", "cucumber", "zucchini", "eggplant", "bell pepper", "pumpkin", "squash", "carrot", "beet", "radish", "turnip", "parsnip", "potato", "yam", "onion", "garlic", "shallot", "leek", "scallion", "celery", "fennel", "asparagus", "broccoli", "cauliflower", "cabbage", "kale", "spinach", "lettuce", "arugula", "collard greens", "okra", "corn", "edamame", "chickpea", "lentil", "soybean", "artichoke", "rhubarb", "kohlrabi", "jicama", "horseradish", "daikon", "chayote", "rutabaga", "cassava", "taro", "plantain", "olive", "date", "mulberry", "gooseberry", "currant", "sapote", "mangosteen", "rambutan", "longan", "salak", "soursop", "breadfruit", "medlar", "loquat", "macadamia", "almond", "cashew", "walnut", "pecan", "hazelnut", "pistachio", "chestnut", "flaxseed", "seaweed", "nori", "spirulina", "wakame", "kelp",
	}
)

type Haikunator struct {
	r     *rand.Rand
	delim string
	token int64
}

func NewHaikunator() Haikunator {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	h := Haikunator{r: r, delim: "-", token: 9999}

	return h
}

func (h *Haikunator) haikunate(token, delim string) (string, error) {
	if !h.isSafeDelimiter(delim) {
		return "", fmt.Errorf("unsafe delimiter: %s", delim)
	}

	haiku := fmt.Sprintf("%s%s%s%s%s", ADJECTIVES[h.r.Intn(len(ADJECTIVES))], delim, ACTIONS[h.r.Intn(len(ACTIONS))], delim, NOUNS[h.r.Intn(len(NOUNS))])

	if len(token) > 0 {
		haiku = fmt.Sprintf("%s%s%s", haiku, delim, token)
	}

	return haiku, nil
}

func (h *Haikunator) Haikunate() (string, error) {
	tokenString := strconv.FormatInt(h.r.Int63n(h.token), 10)

	return h.haikunate(tokenString, h.delim)
}

func (h *Haikunator) TokenHaikunate(token int64) (string, error) {
	tokenString := ""

	if token > 0 {
		tokenString = strconv.FormatInt(h.r.Int63n(token), 10)
	}

	return h.haikunate(tokenString, h.delim)
}

func (h *Haikunator) DelimHaikunate(delim string) (string, error) {
	tokenString := ""
	return h.haikunate(tokenString, delim)
}

func (h *Haikunator) TokenDelimHaikunate(token int64, delim string) (string, error) {
	tokenString := ""

	if token > 0 {
		tokenString = strconv.FormatInt(h.r.Int63n(token), 10)
	}

	return h.haikunate(tokenString, delim)
}

func (h *Haikunator) isSafeDelimiter(delim string) bool {
	if len(delim) == 0 || len(delim) > 5 {
		return false
	}
	for _, r := range delim {
		if r > unicode.MaxASCII {
			return false
		}
	}
	allowed := regexp.MustCompile(`^[-_.,:| ]+$`)
	return allowed.MatchString(delim)
}

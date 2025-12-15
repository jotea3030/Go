package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

const XENO_CANTO_API_KEY = "fae355dcd9dc304f388a83d792fb2e45a18bd4e7"

// XenoCantoResponse represents the API v3 response structure
type XenoCantoResponse struct {
	Recordings    []Recording `json:"recordings"`
	NumPages      int         `json:"numPages"`
	NumRecordings string      `json:"numRecordings"`
	Page          int         `json:"page"`
}

type Recording struct {
	ID      string `json:"id"`
	Gen     string `json:"gen"`
	Sp      string `json:"sp"`
	En      string `json:"en"`
	File    string `json:"file"`
	Cnt     string `json:"cnt"`
	Loc     string `json:"loc"`
	Quality string `json:"q"`
	Length  string `json:"length"`
	Rec     string `json:"rec"`
	Type    string `json:"type"`
	Lat     string `json:"lat"`
	Lng     string `json:"lng"`
}

// Complete list of all birds from Wingspan base game and all expansions
// Base Game (170 birds) + European (81 birds) + Oceania (95 birds) + Asia (90 birds) = 436 birds
var WingspanBirds = []string{
	// Base Game - North America (170 birds)
	"Acorn Woodpecker", "American Avocet", "American Bittern", "American Coot",
	"American Crow", "American Goldfinch", "American Kestrel", "American Redstart",
	"American Robin", "American Tree Sparrow", "American White Pelican", "American Wigeon",
	"American Woodcock", "Anna's Hummingbird", "Bald Eagle", "Baltimore Oriole",
	"Band-tailed Pigeon", "Barn Owl", "Barn Swallow", "Barred Owl",
	"Belted Kingfisher", "Bewick's Wren", "Black Skimmer", "Black Vulture",
	"Black-bellied Plover", "Black-billed Magpie", "Black-capped Chickadee", "Black-chinned Hummingbird",
	"Blue Grosbeak", "Blue Jay", "Blue-gray Gnatcatcher", "Boat-tailed Grackle",
	"Brewer's Blackbird", "Broad-winged Hawk", "Brown Creeper", "Brown Pelican",
	"Brown-headed Cowbird", "Bushtit", "Cackling Goose", "California Condor",
	"California Gull", "Canada Goose", "Canvasback", "Carolina Wren",
	"Chestnut-backed Chickadee", "Chihuahuan Raven", "Chimney Swift", "Chipping Sparrow",
	"Clark's Grebe", "Common Grackle", "Common Loon", "Common Merganser",
	"Common Nighthawk", "Common Raven", "Common Yellowthroat", "Cooper's Hawk",
	"Dark-eyed Junco", "Dickcissel", "Double-crested Cormorant", "Downy Woodpecker",
	"Dunlin", "Eastern Bluebird", "Eastern Kingbird", "Eastern Meadowlark",
	"Eastern Phoebe", "Eastern Screech-Owl", "Eastern Towhee", "Eastern Wood-Pewee",
	"Evening Grosbeak", "Ferruginous Hawk", "Fish Crow", "Forster's Tern",
	"Fox Sparrow", "Franklin's Gull", "Gadwall", "Golden Eagle",
	"Golden-crowned Kinglet", "Gray Catbird", "Great Black-backed Gull", "Great Blue Heron",
	"Great Crested Flycatcher", "Great Egret", "Great Gray Owl", "Great Horned Owl",
	"Greater Prairie-Chicken", "Greater Roadrunner", "Greater Scaup", "Greater White-fronted Goose",
	"Greater Yellowlegs", "Green Heron", "Green-winged Teal", "Hairy Woodpecker",
	"Hermit Thrush", "Hooded Merganser", "Horned Grebe", "Horned Lark",
	"House Finch", "House Sparrow", "House Wren", "Indigo Bunting",
	"Killdeer", "Least Flycatcher", "Least Sandpiper", "Lesser Scaup",
	"Lesser Yellowlegs", "Lincoln's Sparrow", "Loggerhead Shrike", "Long-billed Curlew",
	"Long-billed Dowitcher", "Mallard", "Marbled Godwit", "Marsh Wren",
	"Mountain Bluebird", "Mourning Dove", "Mute Swan", "Northern Bobwhite",
	"Northern Cardinal", "Northern Flicker", "Northern Gannet", "Northern Harrier",
	"Northern Mockingbird", "Northern Pintail", "Northern Shoveler", "Osprey",
	"Painted Bunting", "Peregrine Falcon", "Pied-billed Grebe",
	"Pileated Woodpecker", "Prairie Falcon", "Purple Martin", "Pyrrhuloxia",
	"Red Crossbill", "Red Knot", "Red-bellied Woodpecker", "Red-breasted Merganser",
	"Red-breasted Nuthatch", "Red-headed Woodpecker", "Red-shouldered Hawk", "Red-tailed Hawk",
	"Red-winged Blackbird", "Ring-billed Gull", "Ring-necked Duck", "Ring-necked Pheasant",
	"Rock Pigeon", "Roseate Spoonbill", "Rose-breasted Grosbeak", "Royal Tern",
	"Ruby-crowned Kinglet", "Ruby-throated Hummingbird", "Ruddy Duck", "Ruddy Turnstone",
	"Ruffed Grouse", "Rufous Hummingbird", "Sanderling", "Sandhill Crane",
	"Savannah Sparrow", "Scaled Quail", "Scissor-tailed Flycatcher", "Sharp-shinned Hawk",
	"Short-eared Owl", "Snow Goose", "Snowy Egret", "Song Sparrow",
	"Spotted Sandpiper", "Spotted Towhee", "Steller's Jay", "Swainson's Hawk",
	"Tree Swallow", "Trumpeter Swan", "Tufted Titmouse", "Turkey Vulture",
	"Veery", "Vesper Sparrow", "Virginia Rail", "Western Grebe",
	"Western Gull", "Western Meadowlark", "Western Sandpiper", "Western Tanager",
	"White-breasted Nuthatch", "White-crowned Sparrow", "White-throated Sparrow", "Wild Turkey",
	"Willet", "Wilson's Snipe", "Wood Duck", "Wood Thrush",
	"Yellow Warbler", "Yellow-bellied Sapsucker", "Yellow-breasted Chat", "Yellow-rumped Warbler",

	// European Expansion (81 birds)
	"Audouin's Gull", "Black Redstart", "Black Woodpecker", "Black-headed Gull",
	"Black-tailed Godwit", "Black-throated Diver", "Bluethroat", "Bonelli's Eagle",
	"Bullfinch", "Carrion Crow", "Cetti's Warbler", "Coal Tit",
	"Common Blackbird", "Common Buzzard", "Common Chaffinch", "Common Chiffchaff",
	"Common Cuckoo", "Common Goldeneye", "Common Kingfisher", "Common Little Bittern",
	"Common Moorhen", "Common Nightingale", "Common Starling", "Common Swift",
	"Corsican Nuthatch", "Dunnock", "Eastern Imperial Eagle", "Eleonora's Falcon",
	"Eurasian Collared-Dove", "Eurasian Golden Oriole", "Eurasian Green Woodpecker", "Eurasian Hobby",
	"Eurasian Jay", "Eurasian Magpie", "Eurasian Nutcracker", "Eurasian Nuthatch",
	"Eurasian Sparrowhawk", "Eurasian Tree Sparrow", "European Bee-Eater", "European Goldfinch",
	"European Honey Buzzard", "European Robin", "European Roller", "European Turtle Dove",
	"Goldcrest", "Great Crested Grebe", "Great Tit", "Greater Flamingo",
	"Grey Heron", "Greylag Goose", "Griffon Vulture", "Hawfinch",
	"Hooded Crow", "Lesser Whitethroat", "Little Bustard",
	"Little Owl", "Long-tailed Tit", "Moltoni's Warbler", "Montagu's Harrier",
	"Northern Goshawk", "Parrot Crossbill",
	"Red Kite", "Red-backed Shrike", "Red-legged Partridge",
	"Ruff", "Savi's Warbler", "Short-toed Treecreeper", "Snow Bunting",
	"Snowy Owl", "Squacco Heron", "Thekla's Lark", "White Stork",
	"White Wagtail", "White-backed Woodpecker", "White-throated Dipper", "Wilson's Storm-Petrel",
	"Yellowhammer",

	// Oceania Expansion (95 birds)
	"Abbott's Booby", "Australasian Pipit", "Australasian Shoveler", "Australian Ibis",
	"Australian Magpie", "Australian Owlet-Nightjar", "Australian Raven", "Australian Reed Warbler",
	"Australian Shelduck", "Australian Zebra Finch", "Black Noddy", "Black Swan",
	"Black-shouldered Kite", "Blyth's Hornbill", "Brolga", "Brown Falcon",
	"Budgerigar", "Cockatiel", "Count Raggi's Bird-of-Paradise", "Crested Pigeon",
	"Crimson Chat", "Eastern Rosella", "Eastern Whipbird", "Emu",
	"Galah", "Golden-headed Cisticola", "Gould's Finch", "Green Pygmy-Goose",
	"Grey Butcherbird", "Grey Shrike-thrush", "Grey Teal", "Grey Warbler",
	"Grey-headed Mannikin", "Horsfield's Bronze-Cuckoo", "Horsfield's Bushlark", "Kakapo",
	"Kea", "Kelp Gull", "Kereru", "Korimako",
	"Laughing Kookaburra", "Lesser Frigatebird", "Lewin's Honeyeater", "Little Penguin",
	"Little Pied Cormorant", "Magpie-lark", "Major Mitchell's Cockatoo", "Malleefowl",
	"Maned Duck", "Many-colored Fruit-Dove", "Masked Lapwing", "Mistletoebird",
	"Musk Duck", "New Holland Honeyeater", "Noisy Miner", "North Island Brown Kiwi",
	"Orange-footed Scrubfowl", "Pacific Black Duck", "Peaceful Dove", "Pesquet's Parrot",
	"Pheasant Coucal", "Pink-eared Duck", "Plains-wanderer", "Princess Stephanie's Astrapia",
	"Pukeko", "Rainbow Lorikeet", "Red Wattlebird", "Red-backed Fairywren",
	"Red-capped Robin", "Red-necked Avocet", "Red-winged Parrot", "Regent Bowerbird",
	"Royal Spoonbill", "Rufous Banded Honeyeater", "Rufous Night Heron", "Rufous Owl",
	"Sacred Kingfisher", "Silvereye", "South Island Robin", "Southern Cassowary",
	"Spangled Drongo", "Splendid Fairywren", "Spotless Crake", "Stubble Quail",
	"Sulphur-crested Cockatoo", "Superb Lyrebird", "Tawny Frogmouth", "Tui",
	"Wedge-tailed Eagle", "Welcome Swallow", "White-bellied Sea-Eagle", "White-breasted Woodswallow",
	"White-faced Heron", "Willie Wagtail", "Wrybill",

	// Asia Expansion (90 birds)
	"Ashy Drongo", "Asian Barred Owlet", "Asian Emerald Cuckoo", "Asian Fairy-bluebird",
	"Asian Koel", "Asian Openbill", "Baikal Teal", "Bar-headed Goose",
	"Black Drongo", "Black Kite", "Black-crowned Night Heron", "Black-naped Monarch",
	"Black-naped Oriole", "Blue Rock Thrush", "Blue Whistling Thrush", "Cattle Egret",
	"Chinese Bamboo Partridge", "Chinese Grosbeak", "Cinereous Vulture", "Common Hoopoe",
	"Common Iora", "Common Myna", "Common Tailorbird",
	"Coppersmith Barbet", "Crested Serpent Eagle", "Crested Wood Partridge", "Dollarbird",
	"Eurasian Curlew", "Eurasian Hoopoe", "Forest Owlet", "Great Hornbill",
	"Great Indian Bustard", "Greater Adjutant", "Greater Coucal", "Greater Painted-Snipe",
	"Green Imperial Pigeon", "Hill Myna", "House Crow", "Indian Cormorant",
	"Indian Grey Hornbill", "Indian Peafowl", "Indian Pitta", "Indian Pond Heron",
	"Indian Roller", "Indian Vulture", "Japanese Bush Warbler", "Japanese Tit",
	"Jungle Crow", "Kalij Pheasant", "Large-billed Crow", "Long-tailed Minivet",
	"Long-tailed Shrike", "Mandarin Duck", "Narcissus Flycatcher", "Nutmeg Mannikin",
	"Oriental Magpie-Robin", "Oriental White-eye", "Pied Bushchat", "Pied Myna",
	"Pin-tailed Snipe", "Plain Prinia", "Puff-throated Babbler", "Purple Heron",
	"Purple Sunbird", "Red Junglefowl", "Red-billed Blue Magpie", "Red-vented Bulbul",
	"Red-wattled Lapwing", "Red-whiskered Bulbul", "Rook",
	"Rose-ringed Parakeet", "Rufous Treepie", "Scaly-breasted Munia", "Siberian Crane",
	"Spot-billed Duck", "Spotted Dove", "Striated Heron", "Tufted Duck",
	"Violet Cuckoo", "White-breasted Kingfisher", "White-rumped Shama",
	"White-throated Kingfisher", "Yellow-billed Babbler", "Yellow-browed Bunting", "Yellow-footed Green Pigeon",
}

type Game struct {
	score          int
	totalQuestions int
	currentAudio   string
	httpClient     *http.Client
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("🦅 Welcome to the Wingspan Bird Quiz! 🦅")
	fmt.Println("========================================")
	fmt.Println("Listen to bird calls and guess the species!")
	fmt.Println()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	game := &Game{
		score:          0,
		totalQuestions: 0,
		httpClient:     client,
	}

	for {
		if !playRound(game) {
			break
		}

		fmt.Println()
		fmt.Printf("Current Score: %d/%d\n", game.score, game.totalQuestions)
		fmt.Println()

		fmt.Print("Play another round? (y/n): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(strings.TrimSpace(response)) != "y" {
			break
		}
		fmt.Println()
	}

	if game.currentAudio != "" {
		os.Remove(game.currentAudio)
	}

	fmt.Println()
	fmt.Println("========================================")
	if game.totalQuestions > 0 {
		fmt.Printf("Final Score: %d/%d (%.1f%%)\n",
			game.score, game.totalQuestions,
			float64(game.score)/float64(game.totalQuestions)*100)
	}
	fmt.Println("Thanks for playing!")
}

func playRound(game *Game) bool {
	if game.currentAudio != "" {
		os.Remove(game.currentAudio)
		game.currentAudio = ""
	}

	correctBird := WingspanBirds[rand.Intn(len(WingspanBirds))]

	fmt.Printf("Fetching bird call for question %d...\n", game.totalQuestions+1)
	fmt.Printf("Species: %s\n", correctBird)

	recording, err := getRecording(correctBird, game.httpClient)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Skipping this round...")
		return true
	}

	fmt.Println("Downloading audio...")

	audioFile, err := downloadAudio(recording.File, game.httpClient)
	if err != nil {
		fmt.Printf("Error downloading audio: %v\n", err)
		fmt.Println("Skipping this round...")
		return true
	}

	game.currentAudio = audioFile

	fileInfo, err := os.Stat(audioFile)
	if err != nil || fileInfo.Size() == 0 {
		fmt.Println("Error: Downloaded file is invalid")
		fmt.Println("Skipping this round...")
		os.Remove(audioFile)
		return true
	}

	fmt.Printf("✓ Audio ready (%d KB)\n", fileInfo.Size()/1024)

	options := generateOptions(correctBird)

	fmt.Println()
	fmt.Println("🎵 Playing bird call...")
	err = playAudioFile(audioFile)
	if err != nil {
		fmt.Printf("Error playing audio: %v\n", err)
		fmt.Println("Skipping this round...")
		return true
	}

	for {
		fmt.Println()
		fmt.Println("Which bird species is this?")
		fmt.Println()

		for i, option := range options {
			fmt.Printf("%d. %s\n", i+1, option)
		}
		fmt.Println("R. Replay bird call")

		fmt.Println()
		fmt.Print("Your answer (1-4 or R): ")

		var input string
		fmt.Scanln(&input)
		input = strings.ToUpper(strings.TrimSpace(input))

		if input == "R" {
			fmt.Println("🎵 Replaying...")
			playAudioFile(audioFile)
			continue
		}

		var answer int
		_, err := fmt.Sscanf(input, "%d", &answer)
		if err != nil || answer < 1 || answer > 4 {
			fmt.Println("Invalid input! Please enter 1-4 or R")
			continue
		}

		game.totalQuestions++

		if options[answer-1] == correctBird {
			game.score++
			fmt.Println("✅ Correct! Well done!")
		} else {
			fmt.Printf("❌ Incorrect! The correct answer was: %s\n", correctBird)
		}

		fmt.Printf("   Location: %s, %s\n", recording.Loc, recording.Cnt)
		fmt.Printf("   Quality: %s | Length: %s | Recordist: %s\n",
			recording.Quality, recording.Length, recording.Rec)

		for {
			fmt.Println()
			fmt.Print("Listen again? (y/n): ")
			var replay string
			fmt.Scanln(&replay)
			replay = strings.ToLower(strings.TrimSpace(replay))

			if replay == "y" {
				fmt.Println("🎵 Replaying...")
				playAudioFile(audioFile)
			} else {
				break
			}
		}
		break
	}

	return true
}

func downloadAudio(audioURL string, client *http.Client) (string, error) {
	if audioURL == "" {
		return "", fmt.Errorf("empty audio URL")
	}

	if strings.HasPrefix(audioURL, "//") {
		audioURL = "https:" + audioURL
	}

	resp, err := client.Get(audioURL)
	if err != nil {
		return "", fmt.Errorf("failed to download: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	tmpFile, err := ioutil.TempFile("", "bird-*.mp3")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer tmpFile.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to read data: %v", err)
	}

	if len(data) == 0 {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("file is empty")
	}

	if _, err := tmpFile.Write(data); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to write data: %v", err)
	}

	return tmpFile.Name(), nil
}

func playAudioFile(filename string) error {
	cmd := exec.Command("afplay", filename)

	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(15 * time.Second):
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return nil
	}
}

func getRecording(birdName string, client *http.Client) (*Recording, error) {
	baseURL := "https://xeno-canto.org/api/3/recordings"

	fetchRecordings := func(query string) ([]Recording, error) {
		// API v3 requires tag-based queries
		// Format: en:"Bird Name" for English names
		fullURL := fmt.Sprintf("%s?key=%s&query=%s", 
			baseURL, 
			XENO_CANTO_API_KEY,
			url.QueryEscape(query))
		
		fmt.Printf("DEBUG: Fetching URL: %s\n", fullURL)
		
		resp, err := client.Get(fullURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(body))
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		fmt.Printf("DEBUG: API Response (first 200 chars): %s\n", string(body[:min(200, len(body))]))

		var xcResp XenoCantoResponse
		if err := json.Unmarshal(body, &xcResp); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %v", err)
		}

		fmt.Printf("DEBUG: Found %d recordings\n", len(xcResp.Recordings))

		return filterValidRecordings(xcResp.Recordings), nil
	}

	// Try quality A first - using en: tag for English name
	validRecordings, err := fetchRecordings(fmt.Sprintf("en:\"%s\" q:A", birdName))
	if err != nil {
		return nil, err
	}

	// Try quality B if no A recordings
	if len(validRecordings) == 0 {
		validRecordings, err = fetchRecordings(fmt.Sprintf("en:\"%s\" q:B", birdName))
		if err != nil {
			return nil, err
		}
	}

	// Try any quality as last resort
	if len(validRecordings) == 0 {
		validRecordings, err = fetchRecordings(fmt.Sprintf("en:\"%s\"", birdName))
		if err != nil {
			return nil, err
		}
	}

	if len(validRecordings) == 0 {
		return nil, fmt.Errorf("no recordings found for %s", birdName)
	}

	return &validRecordings[rand.Intn(len(validRecordings))], nil
}

func filterValidRecordings(recordings []Recording) []Recording {
	valid := []Recording{}

	for _, rec := range recordings {
		if rec.File == "" {
			continue
		}

		if strings.HasPrefix(rec.File, "//") || 
		   strings.HasPrefix(rec.File, "http://") || 
		   strings.HasPrefix(rec.File, "https://") {
			valid = append(valid, rec)
		}
	}

	return valid
}

func generateOptions(correctBird string) []string {
	options := []string{correctBird}

	for len(options) < 4 {
		randomBird := WingspanBirds[rand.Intn(len(WingspanBirds))]
		if !contains(options, randomBird) {
			options = append(options, randomBird)
		}
	}

	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	return options
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

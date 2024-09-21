package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"encoding/json"
	"github.com/Raikoa414/go_pokedex/internal"
	"time"
	"math/rand"
)


type Commands struct {
	name        string
	description string
	function    func(configure *config, c *pokecache.Cache, AreaName string) error
}

type config struct {
	id  int
	history [][]string // Keep track of the history of pages (slice of slices)
	caughtPokemon   map[string]Pokemon
}

type locale_area struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}


type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	Abilities      []struct {
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
		Ability  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt int `json:"level_learned_at"`
			VersionGroup   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       any    `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  any    `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      any    `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale any    `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       any    `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault  string `json:"back_default"`
					BackGray     string `json:"back_gray"`
					FrontDefault string `json:"front_default"`
					FrontGray    string `json:"front_gray"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault  string `json:"back_default"`
					BackGray     string `json:"back_gray"`
					FrontDefault string `json:"front_default"`
					FrontGray    string `json:"front_gray"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"crystal"`
				Gold struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"gold"`
				Silver struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       any    `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  any    `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      any    `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale any    `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Cries struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	PastTypes []struct {
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
		Types []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		} `json:"types"`
	} `json:"past_types"`
}








func main() {
	time := time.Duration(30 * time.Second)
	configure := &config{history: make([][]string, 0), id: 1, caughtPokemon: make(map[string]Pokemon)} // Initialize history as a slice of slices
	start_repl(configure,time)
}

func commandHelp(configure *config, c *pokecache.Cache, AreaName string) error {
	fmt.Println("welcome to the pokedex!\n")
	fmt.Println("usage:\n")
	commandsInput := get_commands(configure)
	for i, v := range commandsInput {
		fmt.Printf("%s: %s\n", i, v.description)
	}
	fmt.Println()
	return nil
}

func commandExit(configure *config, c *pokecache.Cache, AreaName string) error {
	os.Exit(0)
	return nil
}

func start_repl(configure *config, inter time.Duration) {
	input := bufio.NewScanner(os.Stdin)
	c := pokecache.NewCache(inter)
	for {
		fmt.Print("pokedex > ")
		input.Scan()
		commandParts := cleanInput(input.Text())
		if len(commandParts) == 0 {
			continue
		}
		commandText := commandParts[0]
		area := ""
		if len(commandParts) < 2{
			area = ""
		}else{
			area = commandParts[1]
		}
		
		
		// Get the command from the map
		command, exists := get_commands(configure)[commandText]

		if exists {
			// Call the function associated with the command
			if err := command.function(configure, c, area); err != nil {
				fmt.Println("Error", err)
			}
		} else {
			fmt.Println("no such command")
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func commandMap(configure *config, c *pokecache.Cache, AreaName string) error {
	// Save the current page of URLs to history before loading new ones
	id_old := configure.id
	currentPage := make([]string, 0)
	
	for {
		if configure.id-id_old >= 20 {
			break
		}
		id_string := strconv.Itoa(configure.id)
		url := "https://pokeapi.co/api/v2/location-area/" + id_string
		if _, exist := c.Cargo[url]; exist{
			body, succ := c.Get(url)
			if  succ != true {
				return fmt.Errorf("error")
			}
			location := locale_area{}
			err := json.Unmarshal(body, &location)
			if err != nil {
				return err
			}
			fmt.Println(location.Name)
			configure.id += 1
			continue
		}else{
			currentPage = append(currentPage, url)

		res, err := http.Get(url)
		if err != nil {
			return err
		}
		

		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		c.Add(url, body)
		if res.StatusCode > 299 {
			return fmt.Errorf("invalid status %v", err)
		}
		if err != nil {
			return err
		}
		location := locale_area{}
		err = json.Unmarshal(body, &location)
		if err != nil {
			return err
		}
		fmt.Println(location.Name)
		configure.id += 1
		}
	}

	// Add the current page to the history
	configure.history = append(configure.history, currentPage)
	return nil
}

func commandMapB(configure *config, c *pokecache.Cache, AreaName string) error {
	// Check if there is a previous page in history
	if len(configure.history) == 0 {
		fmt.Println("No previous pages to go back to.")
		return nil
	}

	// Pop the last page from history
	lastPage := configure.history[len(configure.history)-1]
	configure.history = configure.history[:len(configure.history)-1] // Remove last page

	// Set the id to the last known id from history
	configure.id -= 20 // Move back 20 items, assuming this was your page size
	if configure.id < 0 {
		configure.id = 0
	}

	

	for _, url := range lastPage { // Iterate over the last page URLs
		if body, exist := c.Get(url); exist {
			location := locale_area{}
			err := json.Unmarshal(body, &location)
			if err != nil {
				return err
			}
			fmt.Println(location.Name)
		} else {
			fmt.Printf("No records for %s\n", url)
		}
	}

	return nil
}

func commandExplore(configure *config, c *pokecache.Cache, AreaName string) error{
	if len(AreaName) == 0{
		return fmt.Errorf("no location")
	}
	if body, exists := c.Get(AreaName); exists{
		location := locale_area{}
		errs := json.Unmarshal(body, &location)
		if errs != nil{
			return fmt.Errorf("unmarshal error")
		}
		fmt.Println("pokemon found:")
		for _, v := range location.PokemonEncounters{
			fmt.Println(v.Pokemon.Name)
		}
		return nil
	}else{
		res, err := http.Get("https://pokeapi.co/api/v2/location-area/" + AreaName)
		if err != nil{
			return fmt.Errorf("unable to get pokemon in area")
		
		}
		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		c.Add(AreaName, body)
		if res.StatusCode > 299{
			fmt.Printf("Error status code: %v", res.StatusCode)
		}
		if err !=nil{
			return fmt.Errorf("unable to parse json")
		}
		location := locale_area{}
		errs := json.Unmarshal(body, &location)
		if errs != nil{
			return fmt.Errorf("unmarshal error")
		}
		fmt.Println("pokemon found:")
		for _, v := range location.PokemonEncounters{
			fmt.Println(v.Pokemon.Name)
		}
	}
	return nil
}

func commandCatch(configure *config, c *pokecache.Cache, AreaName string) error{
	res, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + AreaName)
	if err != nil{
		return fmt.Errorf("Error getting pokemon:%v", err)
	}
	body , err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299{
		return fmt.Errorf("error on Get request: %v", res.StatusCode)
	}
	if err != nil{
		return fmt.Errorf("error: %v", err)
	}
	poke := Pokemon{}
	errs := json.Unmarshal(body, &poke)
	if errs != nil{
		return fmt.Errorf("error parsing json %v", errs)
	}
	fmt.Printf("threw a pokeball at %s\n", poke.Name)
	chance := rand.Intn(11)
	caught := false
	if poke.BaseExperience >= 250 {
        if chance >= 9 {
            caught = true
        }
    } else if poke.BaseExperience >= 200 && poke.BaseExperience < 250 {
        if chance >= 7 {
            caught = true
        }
    } else if poke.BaseExperience >= 150 && poke.BaseExperience < 200 {
        if chance >= 5 {
            caught = true
        }
    } else {
        if chance >= 3 {
            caught = true
        }
    }
	if caught{
		fmt.Printf("caught %s\n", poke.Name)
		if _,exists := configure.caughtPokemon[poke.Name]; exists{
			fmt.Println("already registered in pokedex")
		}else{
			configure.caughtPokemon[poke.Name] = poke
		}
	}else{
		fmt.Printf("%s escaped!\n", poke.Name)
	}
	
	return nil
}


func commandInspect(configure *config, c *pokecache.Cache, AreaName string) error{
	if InspectMon, exists := configure.caughtPokemon[AreaName]; exists{
		fmt.Println("Name: " + InspectMon.Name)
		fmt.Printf("Height: %v\n", InspectMon.Height)
		fmt.Printf("Weight: %v\n", InspectMon.Weight)
		fmt.Println("Stats:")
		for _, v := range InspectMon.Stats{
			fmt.Printf("	-%v: %v\n", v.Stat.Name, v.BaseStat)
		}
		fmt.Println("Types:")
		for _, h := range InspectMon.Types{
			fmt.Printf("	-%v\n",h.Type.Name)
		}
	}else{
		fmt.Println("you have not caught that pokemon")
	}
	return nil
}


func commandPokeDex(configure *config, c *pokecache.Cache, AreaName string) error{
	fmt.Println("Your PokeDex:")
	if len(configure.caughtPokemon) == 0{
		fmt.Println("You have not caught any pokemon yet, catch some with the catch command")
	}
	for _, v := range configure.caughtPokemon{
		fmt.Printf(" -%v\n", v.Name)
	}
	return nil
}



func get_commands(configure *config) map[string]Commands {
	return map[string]Commands{
		"help": {
			name:        "help",
			description: "show help on commands",
			function:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "exit the program",
			function:    commandExit,
		},
		"map": {
			name:        "map",
			description: "list locations",
			function:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "go back a page when seeing location",
			function:    commandMapB,
		},
		"explore": {
			name: "explore",
			description: "list pokemon in area",
			function:  commandExplore,
		},
		"catch": {
			name: "catch",
			description: "catch a specific pokemon at a chance to add to the pokedex",
			function: commandCatch,
		},
		"inspect":{
			name: "inspect",
			description: "inspect a certain pokemon",
			function: commandInspect,
		},
		"pokedex":{
			name:"pokedex",
			description: "list the whole caught pokedex",
			function: commandPokeDex,
		},
	}
}

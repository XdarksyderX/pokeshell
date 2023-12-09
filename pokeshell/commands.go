package pokeshell

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type PokedexShell struct {
	MapPage     int
	CommandList map[string]Command
	Pokedex     []Pokemon
	PokeCache   PokeCache
}

type Command struct {
	Name     string
	Usage    string
	Callback func(shell *PokedexShell, args []string) error
}

type Pokemon struct {
	Name   string
	Height float64
	Weight float64
	Types  []string
	Stats  map[string]float64
}

var commandList = map[string]Command{
	"help": {
		Name:     "help",
		Usage:    "get helped",
		Callback: commandHelp,
	},
	"exit": {
		Name:     "exit",
		Usage:    "exit pokedex",
		Callback: commandExit,
	},
	"map": {
		Name:     "map",
		Usage:    "displays 20 location names in the Pokémon world. Each new call to map reveals the next 20 locations, allowing you to progressively explore the Pokémon world.",
		Callback: commandMap,
	},
	"mapb": {
		Name:     "map",
		Usage:    "displays the previous 20 location names in the Pokémon world. Each time you use mapb, it shows the 20 locations preceding the current set, allowing you to backtrack and explore earlier areas in the Pokémon world.",
		Callback: commandMapb,
	},
	"clear": {
		Name:     "clear",
		Usage:    "clear the screen",
		Callback: commandClear,
	},
	"explore": {
		Name:     "explore",
		Usage:    "explore a location - explore <location>",
		Callback: commandExplore,
	},
	"catch": {
		Name:     "catch",
		Usage:    "catch a pokemon - catch <pokemon>",
		Callback: commandCatch,
	},
	"inspect": {
		Name:     "inspect",
		Usage:    "inspect a pokemon from your pokedex - inspect <pokemon>",
		Callback: commandInspect,
	},
	"pokedex": {
		Name:     "pokedex",
		Usage:    "get your pokemon - pokedex",
		Callback: commandPokedex,
	},
}

func commandMap(pS *PokedexShell, args []string) error {
	if pS.MapPage < 830 {
		pS.MapPage += 20
	}
	apiURL := fmt.Sprintf("https://pokeapi.co/api/v2/location/?offset=%d&limit=20", pS.MapPage)
	data, err := pS.PokeCache.Get(apiURL)
	if err != nil {
		return err
	}
	var lP LocationPage
	if err := json.Unmarshal(data, &lP); err != nil {
		return err
	}
	for _, r := range lP.Results {
		fmt.Println(r.Name)
	}
	return nil
}

func commandExplore(pS *PokedexShell, args []string) error {
	apiURL := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", args[0])
	data, err := pS.PokeCache.Get(apiURL)
	if err != nil {
		return err
	}
	var lAe LocationAreaExplore
	if err := json.Unmarshal(data, &lAe); err != nil {
		return err
	}
	fmt.Println("Exploring " + args[0])
	fmt.Println("Found Pokemon:")
	for _, r := range lAe.PokemonEncounters {
		fmt.Printf("- %s\n", r.Pokemon.Name)
	}
	return nil
}

func commandMapb(pS *PokedexShell, args []string) error {
	if pS.MapPage > 40 {
		pS.MapPage -= 20
	}
	apiURL := fmt.Sprintf("https://pokeapi.co/api/v2/location/?offset=%d&limit=20", pS.MapPage)
	data, err := pS.PokeCache.Get(apiURL)
	if err != nil {
		return err
	}
	var lP LocationPage
	if err := json.Unmarshal(data, &lP); err != nil {
		return err
	}
	for _, r := range lP.Results {
		fmt.Println(r.Name)
	}
	return nil
}

func ConvertFromFetched(pokemonFetched PokemonFetched) Pokemon {
	pokemon := Pokemon{
		Name:   pokemonFetched.Name,
		Height: pokemonFetched.Height,
		Weight: pokemonFetched.Weight,
		Types:  make([]string, 0),
		Stats:  make(map[string]float64),
	}

	for _, t := range pokemonFetched.Types {
		pokemon.Types = append(pokemon.Types, t.Type.Name)
	}

	for _, s := range pokemonFetched.Stats {
		pokemon.Stats[s.Stat.Name] = float64(s.BaseStat)
	}

	return pokemon
}

func commandCatch(pS *PokedexShell, args []string) error {
	chance := rand.Float64()
	fmt.Printf("Throwing a pokeball to %s...\n", args[0])
	if chance <= 0.6 {
		fmt.Printf("%s escaped!\n", args[0])
		return nil
	}
	apiURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", args[0])
	data, err := pS.PokeCache.Get(apiURL)
	if err != nil {
		return err
	}
	var pF PokemonFetched
	if err := json.Unmarshal(data, &pF); err != nil {
		return err
	}
	pokemon := ConvertFromFetched(pF)
	pS.Pokedex = append(pS.Pokedex, pokemon)
	fmt.Printf("%s was caught!\n", args[0])
	return nil

}

func (p Pokemon) printPokemon() {
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Height: %v\n", p.Height)
	fmt.Printf("Weight: %v\n", p.Weight)
	fmt.Println("Stats:")
	for s, v := range p.Stats {
		fmt.Printf("\t -%s: %v\n", s, v)
	}
	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Printf("- %s\n", t)
	}
}

func commandInspect(pS *PokedexShell, args []string) error {
	for _, pokemon := range pS.Pokedex {
		if pokemon.Name == args[0] {
			pokemon.printPokemon()
			return nil
		}
	}
	return errors.New("Not pokemon named " + args[0] + " was found in your pokedex!")
}

func commandPokedex(pS *PokedexShell, args []string) error {
	if len(pS.Pokedex) == 0 {
		return errors.New("you don't have any pokemon")
	}
	fmt.Println("Your Pokedex:")
	for _, pokemon := range pS.Pokedex {
		fmt.Printf("- %s\n", pokemon.Name)
	}
	return nil
}

func commandExit(*PokedexShell, []string) error {
	fmt.Println("Bye!")
	os.Exit(0)
	return nil
}

func commandClear(*PokedexShell, []string) error {
	fmt.Println("\033[2J")
	return nil
}

func commandHelp(pS *PokedexShell, args []string) error {
	fmt.Println("Command list:")
	for _, c := range pS.CommandList {
		fmt.Printf("%s - %s\n", c.Name, c.Usage)
	}
	return nil
}

func CreatePokeshell() PokedexShell {
	pS := PokedexShell{
		MapPage:     20,
		CommandList: commandList,
		PokeCache:   *NewPokeCache(5 * time.Minute),
	}
	return pS
}

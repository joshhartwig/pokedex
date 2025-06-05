package models

import (
	"log/slog"

	"github.com/joshhartwig/pokedex/internal/database"
	"github.com/joshhartwig/pokedex/internal/pokecache"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func(...string) error
}

type Config struct {
	Commands   map[string]CliCommand
	Next       string
	Previous   string
	BaseApiUrl string
	Cache      pokecache.Cache
	Pokedex    map[string]Pokemon
	Db         *database.Queries
	Logger     *slog.Logger
	History    []string
}

// json decoding
type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Apiheader struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Location `json:"results"`
}

type LocationArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"encounter_method,omitempty"`
		VersionDetails []struct {
			Rate    int `json:"rate,omitempty"`
			Version struct {
				Name string `json:"name,omitempty"`
				URL  string `json:"url,omitempty"`
			} `json:"version,omitempty"`
		} `json:"version_details,omitempty"`
	} `json:"encounter_method_rates,omitempty"`
	GameIndex int `json:"game_index,omitempty"`
	ID        int `json:"id,omitempty"`
	Location  struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty"`
	} `json:"location,omitempty"`
	Name  string `json:"name,omitempty"`
	Names []struct {
		Language struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"language,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"names,omitempty"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"pokemon,omitempty"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance,omitempty"`
				ConditionValues []any `json:"condition_values,omitempty"`
				MaxLevel        int   `json:"max_level,omitempty"`
				Method          struct {
					Name string `json:"name,omitempty"`
					URL  string `json:"url,omitempty"`
				} `json:"method,omitempty"`
				MinLevel int `json:"min_level,omitempty"`
			} `json:"encounter_details,omitempty"`
			MaxChance int `json:"max_chance,omitempty"`
			Version   struct {
				Name string `json:"name,omitempty"`
				URL  string `json:"url,omitempty"`
			} `json:"version,omitempty"`
		} `json:"version_details,omitempty"`
	} `json:"pokemon_encounters,omitempty"`
}

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"ability,omitempty"`
		IsHidden bool `json:"is_hidden,omitempty"`
		Slot     int  `json:"slot,omitempty"`
	} `json:"abilities,omitempty"`
	BaseExperience int `json:"base_experience,omitempty"`
	Cries          struct {
		Latest string `json:"latest,omitempty"`
		Legacy string `json:"legacy,omitempty"`
	} `json:"cries,omitempty"`
	Forms []struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty"`
	} `json:"forms,omitempty"`
	GameIndices []struct {
		GameIndex int `json:"game_index,omitempty"`
		Version   struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"version,omitempty"`
	} `json:"game_indices,omitempty"`
	Height    int `json:"height,omitempty"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"item,omitempty"`
		VersionDetails []struct {
			Rarity  int `json:"rarity,omitempty"`
			Version struct {
				Name string `json:"name,omitempty"`
				URL  string `json:"url,omitempty"`
			} `json:"version,omitempty"`
		} `json:"version_details,omitempty"`
	} `json:"held_items,omitempty"`
	ID                     int    `json:"id,omitempty"`
	IsDefault              bool   `json:"is_default,omitempty"`
	LocationAreaEncounters string `json:"location_area_encounters,omitempty"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"move,omitempty"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at,omitempty"`
			MoveLearnMethod struct {
				Name string `json:"name,omitempty"`
				URL  string `json:"url,omitempty"`
			} `json:"move_learn_method,omitempty"`
			Order        int `json:"order,omitempty"`
			VersionGroup struct {
				Name string `json:"name,omitempty"`
				URL  string `json:"url,omitempty"`
			} `json:"version_group,omitempty"`
		} `json:"version_group_details,omitempty"`
	} `json:"moves,omitempty"`
	Name          string `json:"name,omitempty"`
	Order         int    `json:"order,omitempty"`
	PastAbilities []struct {
		Abilities []struct {
			Ability  any  `json:"ability,omitempty"`
			IsHidden bool `json:"is_hidden,omitempty"`
			Slot     int  `json:"slot,omitempty"`
		} `json:"abilities,omitempty"`
		Generation struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"generation,omitempty"`
	} `json:"past_abilities,omitempty"`
	PastTypes []any `json:"past_types,omitempty"`
	Species   struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty"`
	} `json:"species,omitempty"`
	Sprites struct {
		BackDefault      string `json:"back_default,omitempty"`
		BackFemale       any    `json:"back_female,omitempty"`
		BackShiny        string `json:"back_shiny,omitempty"`
		BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
		FrontDefault     string `json:"front_default,omitempty"`
		FrontFemale      any    `json:"front_female,omitempty"`
		FrontShiny       string `json:"front_shiny,omitempty"`
		FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default,omitempty"`
				FrontFemale  any    `json:"front_female,omitempty"`
			} `json:"dream_world,omitempty"`
			Home struct {
				FrontDefault     string `json:"front_default,omitempty"`
				FrontFemale      any    `json:"front_female,omitempty"`
				FrontShiny       string `json:"front_shiny,omitempty"`
				FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
			} `json:"home,omitempty"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default,omitempty"`
				FrontShiny   string `json:"front_shiny,omitempty"`
			} `json:"official-artwork,omitempty"`
			Showdown struct {
				BackDefault      string `json:"back_default,omitempty"`
				BackFemale       any    `json:"back_female,omitempty"`
				BackShiny        string `json:"back_shiny,omitempty"`
				BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
				FrontDefault     string `json:"front_default,omitempty"`
				FrontFemale      any    `json:"front_female,omitempty"`
				FrontShiny       string `json:"front_shiny,omitempty"`
				FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
			} `json:"showdown,omitempty"`
		} `json:"other,omitempty"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      any `json:"back_default,omitempty"`
					BackGray         any `json:"back_gray,omitempty"`
					BackTransparent  any `json:"back_transparent,omitempty"`
					FrontDefault     any `json:"front_default,omitempty"`
					FrontGray        any `json:"front_gray,omitempty"`
					FrontTransparent any `json:"front_transparent,omitempty"`
				} `json:"red-blue,omitempty"`
				Yellow struct {
					BackDefault      any `json:"back_default,omitempty"`
					BackGray         any `json:"back_gray,omitempty"`
					BackTransparent  any `json:"back_transparent,omitempty"`
					FrontDefault     any `json:"front_default,omitempty"`
					FrontGray        any `json:"front_gray,omitempty"`
					FrontTransparent any `json:"front_transparent,omitempty"`
				} `json:"yellow,omitempty"`
			} `json:"generation-i,omitempty"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           any `json:"back_default,omitempty"`
					BackShiny             any `json:"back_shiny,omitempty"`
					BackShinyTransparent  any `json:"back_shiny_transparent,omitempty"`
					BackTransparent       any `json:"back_transparent,omitempty"`
					FrontDefault          any `json:"front_default,omitempty"`
					FrontShiny            any `json:"front_shiny,omitempty"`
					FrontShinyTransparent any `json:"front_shiny_transparent,omitempty"`
					FrontTransparent      any `json:"front_transparent,omitempty"`
				} `json:"crystal,omitempty"`
				Gold struct {
					BackDefault      any `json:"back_default,omitempty"`
					BackShiny        any `json:"back_shiny,omitempty"`
					FrontDefault     any `json:"front_default,omitempty"`
					FrontShiny       any `json:"front_shiny,omitempty"`
					FrontTransparent any `json:"front_transparent,omitempty"`
				} `json:"gold,omitempty"`
				Silver struct {
					BackDefault      any `json:"back_default,omitempty"`
					BackShiny        any `json:"back_shiny,omitempty"`
					FrontDefault     any `json:"front_default,omitempty"`
					FrontShiny       any `json:"front_shiny,omitempty"`
					FrontTransparent any `json:"front_transparent,omitempty"`
				} `json:"silver,omitempty"`
			} `json:"generation-ii,omitempty"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault any `json:"front_default,omitempty"`
					FrontShiny   any `json:"front_shiny,omitempty"`
				} `json:"emerald,omitempty"`
				FireredLeafgreen struct {
					BackDefault  any `json:"back_default,omitempty"`
					BackShiny    any `json:"back_shiny,omitempty"`
					FrontDefault any `json:"front_default,omitempty"`
					FrontShiny   any `json:"front_shiny,omitempty"`
				} `json:"firered-leafgreen,omitempty"`
				RubySapphire struct {
					BackDefault  any `json:"back_default,omitempty"`
					BackShiny    any `json:"back_shiny,omitempty"`
					FrontDefault any `json:"front_default,omitempty"`
					FrontShiny   any `json:"front_shiny,omitempty"`
				} `json:"ruby-sapphire,omitempty"`
			} `json:"generation-iii,omitempty"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default,omitempty"`
					BackFemale       any    `json:"back_female,omitempty"`
					BackShiny        string `json:"back_shiny,omitempty"`
					BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"diamond-pearl,omitempty"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default,omitempty"`
					BackFemale       any    `json:"back_female,omitempty"`
					BackShiny        string `json:"back_shiny,omitempty"`
					BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"heartgold-soulsilver,omitempty"`
				Platinum struct {
					BackDefault      string `json:"back_default,omitempty"`
					BackFemale       any    `json:"back_female,omitempty"`
					BackShiny        string `json:"back_shiny,omitempty"`
					BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"platinum,omitempty"`
			} `json:"generation-iv,omitempty"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default,omitempty"`
						BackFemale       any    `json:"back_female,omitempty"`
						BackShiny        string `json:"back_shiny,omitempty"`
						BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
						FrontDefault     string `json:"front_default,omitempty"`
						FrontFemale      any    `json:"front_female,omitempty"`
						FrontShiny       string `json:"front_shiny,omitempty"`
						FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
					} `json:"animated,omitempty"`
					BackDefault      string `json:"back_default,omitempty"`
					BackFemale       any    `json:"back_female,omitempty"`
					BackShiny        string `json:"back_shiny,omitempty"`
					BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"black-white,omitempty"`
			} `json:"generation-v,omitempty"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"omegaruby-alphasapphire,omitempty"`
				XY struct {
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"x-y,omitempty"`
			} `json:"generation-vi,omitempty"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default,omitempty"`
					FrontFemale  any    `json:"front_female,omitempty"`
				} `json:"icons,omitempty"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"ultra-sun-ultra-moon,omitempty"`
			} `json:"generation-vii,omitempty"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default,omitempty"`
					FrontFemale  any    `json:"front_female,omitempty"`
				} `json:"icons,omitempty"`
			} `json:"generation-viii,omitempty"`
		} `json:"versions,omitempty"`
	} `json:"sprites,omitempty"`
	Stats []struct {
		BaseStat int `json:"base_stat,omitempty"`
		Effort   int `json:"effort,omitempty"`
		Stat     struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"stat,omitempty"`
	} `json:"stats,omitempty"`
	Types []struct {
		Slot int `json:"slot,omitempty"`
		Type struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"type,omitempty"`
	} `json:"types,omitempty"`
	Weight int `json:"weight,omitempty"`
}

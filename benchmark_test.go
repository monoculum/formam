package formam_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/monoculum/formam"
)

type Simple struct {
	Nest struct {
		Children []struct {
			ID   string
			Name string
		}
	}
	String string
	Slice  []int
	Bool   bool
}

var (
	valFormamT2 = url.Values{
		"Nest.Children[0].ID":   []string{"monoculum_id"},
		"Nest.Children[0].Name": []string{"Monoculum"},
		"String":                []string{"golang is very fun"},
		"Slice[0]":              []string{"1"},
		"Slice[1]":              []string{"2"},
		"Slice[2]":              []string{"3"},
		"Slice[3]":              []string{"4"},
		"Bool":                  []string{"true"},
	}
	valFormT2 = url.Values{
		"Nest.Children[0].ID":   []string{"monoculum_id"},
		"Nest.Children[0].Name": []string{"Monoculum"},
		"String":                []string{"golang is very fun"},
		"Slice[0]":              []string{"1"},
		"Slice[1]":              []string{"2"},
		"Slice[2]":              []string{"3"},
		"Slice[3]":              []string{"4"},
		"Bool":                  []string{"true"},
	}
	valSchemaT2 = url.Values{
		"Nest.Children.0.ID":   []string{"monoculum_id"},
		"Nest.Children.0.Name": []string{"Monoculum"},
		"String":               []string{"golang is very fun"},
		"Slice":                []string{"1", "2", "3", "4"},
		"Bool":                 []string{"true"},
	}
	valAJGT2 = url.Values{
		"Nest.Children.0.ID":   []string{"monoculum_id"},
		"Nest.Children.0.Name": []string{"Monoculum"},
		"String":               []string{"golang is very fun"},
		"Slice.0":              []string{"1"},
		"Slice.1":              []string{"2"},
		"Slice.2":              []string{"3"},
		"Slice.3":              []string{"4"},
		"Bool":                 []string{"true"},
	}
	valuesJSONT2 = `
	{
		"Nest":
			{
				"Children": [{"ID": "monoculum_id", "Name":"Monoculum"}]
			},
		"string": "golang is very fun",
		"Slice": [1, 2, 3, 4],
		"Bool": true
	}
	`
)

func BenchmarkSimple(b *testing.B) {
	b.ReportAllocs()
	dec := formam.NewDecoder(nil)
	for i := 0; i < b.N; i++ {
		test := new(Simple)
		if err := dec.Decode(valFormamT2, test); err != nil {
			b.Error(err)
		}
	}
}

type BAnonymous struct {
	Int int `form:"int" formam:"int" json:"int"`
}

type Medium struct {
	Nest struct {
		Children []struct {
			ID   string
			Name string
		}
	}
	String string `form:"string" formam:"string" json:"string"`
	Slice  []int
	Map    map[string][]string
	Bool   bool
	BAnonymous
	//Int    int `form:"int" formam:"int" json:"int"`
}

var (
	valuesFormamT1 = url.Values{
		"Nest.Children[0].ID":   []string{"monoculum_id"},
		"Nest.Children[0].Name": []string{"Monoculum"},
		"Map[es_Es][0]":         []string{"javier"},
		"Map[es_Es][1]":         []string{"javier"},
		"Map[es_Es][2]":         []string{"javier"},
		"Map[es_Es][3]":         []string{"javier"},
		"Map[es_Es][4]":         []string{"javier"},
		"Map[es_Es][5]":         []string{"javier"},
		"string":                []string{"golang is very fun"},
		"Slice[0]":              []string{"1"},
		"Slice[1]":              []string{"2"},
		"int":                   []string{"1"},
		"Bool":                  []string{"true"},
	}
	valuesFormT1 = url.Values{
		"Nest.Children[0].ID":   []string{"monoculum_id"},
		"Nest.Children[0].Name": []string{"Monoculum"},
		"Map[es_Es][0]":         []string{"javier"},
		"Map[es_Es][1]":         []string{"javier"},
		"Map[es_Es][2]":         []string{"javier"},
		"Map[es_Es][3]":         []string{"javier"},
		"Map[es_Es][4]":         []string{"javier"},
		"Map[es_Es][5]":         []string{"javier"},
		"string":                []string{"golang is very fun"},
		"Slice[0]":              []string{"1"},
		"Slice[1]":              []string{"2"},
		"int":                   []string{"1"},
		"Bool":                  []string{"true"},
	}
	valuesAJGFormT1 = url.Values{
		"Nest.Children.0.ID":   []string{"monoculum_id"},
		"Nest.Children.0.Name": []string{"Monoculum"},
		"Map.es_Es.0":          []string{"javier"},
		"Map.es_Es.1":          []string{"javier"},
		"Map.es_Es.2":          []string{"javier"},
		"Map.es_Es.3":          []string{"javier"},
		"Map.es_Es.4":          []string{"javier"},
		"Map.es_Es.5":          []string{"javier"},
		"string":               []string{"golang is very fun"},
		"Slice.0":              []string{"1"},
		"Slice.1":              []string{"2"},
		"int":                  []string{"1"},
		"Bool":                 []string{"true"},
	}
	valuesJSONT1 = `
	{
		"Nest": {"Children": [{"ID": "monoculum_id", "Name":"Monoculum"}]},
		"string": "golang is very fun",
		"Map": {"es_Es": ["javier", "javier", "javier", "javier", "javier", "javier"]},
		"Slice": [1, 2],
		"int": 20,
		"Bool": true
	}
	`
)

func BenchmarkMedium(b *testing.B) {
	deco := formam.NewDecoder(nil)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		test := new(Medium)
		if err := deco.Decode(valuesFormamT1, test); err != nil {
			b.Error(err)
		}
	}
}

type Location struct {
	Latitude  string
	Longitude string
}

type Cinema struct {
	ID       string
	Name     string
	Location Location
}

type AudiovisualRelease struct {
	Date time.Time
}

type Audiovisual struct {
	Type     string
	Releases []AudiovisualRelease
}

type Locale struct {
	Country  string
	Language string
}

type TextLocale struct {
	Locale Locale
	Value  string
}

type BText struct {
	Language string
	Value    string
}

type Tag struct {
	ID   string
	Name map[string]BText
}

type Mood struct {
	ID   string
	Name map[string]BText
}

type Score struct {
	User  string
	Score float32
	Date  time.Time
}

type FilmRelease struct {
	Date   time.Time
	Cinema Cinema
}

type FilmScore struct {
	Average float64
	Total   uint
	Scores  []Score
}

type Film struct {
	Audiovisual
	ID        string
	Titles    map[string]map[string]TextLocale
	Languages map[string]BText
	Year      uint
	Tags      []Tag
	Moods     []Mood
	Release   []FilmRelease
	Remake    *Film
	Synopsis  map[string]BText

	Websites []string

	Scores *FilmScore
}

var (
	valuesFormamT3 = url.Values{
		// id
		"ID": []string{"spirited-away"},

		// type
		"Type": []string{"film"},

		// title
		// title > es
		// title > es > ES
		"Titles[es][ES].Locale.Country":  []string{"Spain"},
		"Titles[es][ES].Locale.Language": []string{"spanish"},
		"Titles[es][ES].Value":           []string{"El viaje de Chihiro"},
		// title > es > MX
		"Titles[es][MX].Locale.Country":  []string{"Mexico"},
		"Titles[es][MX].Locale.Language": []string{"spanish"},
		"Titles[es][MX].Value":           []string{"El viaje de Chihiro"},
		// title > es > AR
		"Titles[es][AR].Locale.Country":  []string{"Argentina"},
		"Titles[es][AR].Locale.Language": []string{"spanish"},
		"Titles[es][AR].Value":           []string{"El viaje de Chihiro"},
		// title > en
		// title > en > US
		"Titles[en][US].Locale.Country":  []string{"United States"},
		"Titles[en][US].Locale.Language": []string{"english"},
		"Titles[en][US].Value":           []string{"Spirited Away"},
		// title > en > GB
		"Titles[en][GB].Locale.Country":  []string{"United Kingdom"},
		"Titles[en][GB].Locale.Language": []string{"english"},
		"Titles[en][GB].Value":           []string{"Spirited Away"},
		// title > fr
		// title > fr > FR
		"Titles[fr][FR].Locale.Country":  []string{"France"},
		"Titles[fr][FR].Locale.Language": []string{"french"},
		"Titles[fr][FR].Value":           []string{"Le voyage de Chihiro"},
		// title > fr > CA
		"Titles[fr][CA].Locale.Country":  []string{"Canada"},
		"Titles[fr][CA].Locale.Language": []string{"french"},
		"Titles[fr][CA].Value":           []string{"Le voyage de Chihiro"},
		// title > de
		// title > de > DE
		"Titles[de][DE].Locale.Country":  []string{"Germany"},
		"Titles[de][DE].Locale.Language": []string{"german"},
		"Titles[de][DE].Value":           []string{"Chihiros Reise ins Zauberland"},
		// title > ja
		// title > ja > JA
		"Titles[ja][JA].Locale.Country":  []string{"Japan"},
		"Titles[ja][JA].Locale.Language": []string{"japanese"},
		"Titles[ja][JA].Value":           []string{"Sen to Chihiro no kamikakushi"},

		// year
		"Year": []string{"2001"},

		// tags
		"Tags[0].ID":             []string{"cult"},
		"Tags[0].Name[en].Value": []string{"Cult"},
		"Tags[1].ID":             []string{"adventure"},
		"Tags[1].Name[en].Value": []string{"Adventure"},
		"Tags[2].ID":             []string{"drama"},
		"Tags[2].Name[en].Value": []string{"Drama"},
		"Tags[3].ID":             []string{"comedy"},
		"Tags[3].Name[en].Value": []string{"Comedy"},
		"Tags[4].ID":             []string{"mistery"},
		"Tags[4].Name[en].Value": []string{"Mistery"},
		"Tags[5].ID":             []string{"family"},
		"Tags[5].Name[en].Value": []string{"Family"},

		// moods
		"Moods[0].ID":             []string{"cult"},
		"Moods[0].Name[en].Value": []string{"Cult"},
		"Moods[1].ID":             []string{"adventure"},
		"Moods[1].Name[en].Value": []string{"Adventure"},
		"Moods[2].ID":             []string{"drama"},
		"Moods[2].Name[en].Value": []string{"Drama"},
		"Moods[3].ID":             []string{"comedy"},
		"Moods[3].Name[en].Value": []string{"Comedy"},
		"Moods[4].ID":             []string{"mistery"},
		"Moods[4].Name[en].Value": []string{"Mistery"},
		"Moods[5].ID":             []string{"family"},
		"Moods[5].Name[en].Value": []string{"Family"},

		// websites
		"Websites": []string{"cult", "adventure", "drama", "comedy", "mistery", "family"},

		// release
		// release > 0
		"Release[0].Date":                      []string{"2001-01-02"},
		"Release[0].Cinema.ID":                 []string{"kodak-theatrical"},
		"Release[0].Cinema.Name":               []string{"Kodak Theatrical"},
		"Release[0].Cinema.Location.Latitude":  []string{"121342"},
		"Release[0].Cinema.Location.Longitude": []string{"122323"},
		// release > 1
		"Release[1].Date":                      []string{"2002-11-02"},
		"Release[1].Cinema.ID":                 []string{"victor-theatrical"},
		"Release[1].Cinema.Name":               []string{"Victor Theatrical"},
		"Release[1].Cinema.Location.Latitude":  []string{"121342"},
		"Release[1].Cinema.Location.Longitude": []string{"122323"},
		// release > 2
		"Release[2].Date":                      []string{"2002-09-02"},
		"Release[2].Cinema.ID":                 []string{"alaska-theatrical"},
		"Release[2].Cinema.Name":               []string{"alaska Theatrical"},
		"Release[2].Cinema.Location.Latitude":  []string{"121342"},
		"Release[2].Cinema.Location.Longitude": []string{"122323"},

		// remake
		"Remake.ID": []string{"taxi-driver"},

		// synopsis
		"Synopsis[en].Language": []string{"english"},
		"Synopsis[en].Value":    []string{"During her family's move to the suburbs, a sullen 10-year-old girl wanders into a world ruled by gods, witches, and spirits, and where humans are changed into beasts."},
		"Synopsis[es].Language": []string{"spanish"},
		"Synopsis[es].Value":    []string{"Chihiro es una niña de diez años que viaja en coche con sus padres. Después de atravesar un túnel, llegan a un mundo fantástico, en el que no hay lugar para los seres humanos, sólo para los dioses de primera y segunda clase. Cuando descubre que sus padres han sido convertidos en cerdos, Chihiro se siente muy sola y asustada"},

		// scores
		"Scores.Total":   []string{"94315"},
		"Scores.Average": []string{"8.1"},
		// scores > scores
		// scores > scores > 0
		"Scores.Scores[0].User":  []string{"dashaus"},
		"Scores.Scores[0].Score": []string{"9"},
		"Scores.Scores[0].Date":  []string{"2016-01-02"},
		// scores > scores > 1
		"Scores.Scores[1].User":  []string{"tupac"},
		"Scores.Scores[1].Score": []string{"8"},
		"Scores.Scores[1].Date":  []string{"2003-01-02"},
		// scores > scores > 2
		"Scores.Scores[2].User":  []string{"ken-thompson"},
		"Scores.Scores[2].Score": []string{"10"},
		"Scores.Scores[2].Date":  []string{"2015-01-02"},
		// scores > scores > 3
		"Scores.Scores[3].User":  []string{"buuzz-lightyearz"},
		"Scores.Scores[3].Score": []string{"2"},
		"Scores.Scores[3].Date":  []string{"2013-01-02"},
	}
)

func BenchmarkComplex(b *testing.B) {
	b.ReportAllocs()
	decoder := formam.NewDecoder(nil)
	decoder.RegisterCustomType(func(vals []string) (interface{}, error) {
		return time.Parse("2006-01-02", vals[0])
	}, []interface{}{time.Time{}}, nil)
	for i := 0; i < b.N; i++ {
		test := new(Film)
		if err := formam.Decode(valuesFormamT3, test); err != nil {
			b.Error(err)
		}
	}
}

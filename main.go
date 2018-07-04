package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/marcusfolkesson/tablewriter"
	"io/ioutil"
	"os"
	"strings"
)

type Recipes struct {
	Recipe []Recipe `xml:"RECIPE"`
}

type Recipe struct {
	Name         string  `xml:"NAME"`
	Brewer       string  `xml:"BREWER"`
	AssistBrewer string  `xml:"ASST_BREWER"`
	BatchSize    float64 `xml:"BATCH_SIZE"`
	BoilTime     float64 `xml:"BOIL_TIME"`

	OG     float64 `xml:"OG"`
	EstOG  string  `xml:"EST_OG"`
	FG     string  `xml:"FG"`
	EstFG  string  `xml:"EST_FG"`
	ABV    string  `xml:"ABV"`
	EstABV string  `xml:"EST_ABV"`
	IBU    string  `xml:"IBU"`

	PrimaryTemp   float64 `xml:"PRIMARY_TEMP"`
	PrimaryAge    float64 `xml:"PRIMARY_AGE"`
	SecondaryTemp float64 `xml:"SECONDARY_TEMP"`
	SecondaryAge  float64 `xml:"SECONDARY_AGE"`

	Notes string `xml:"NOTES"`

	Hop   []Hop   `xml:"HOPS>HOP"`
	Grain []Grain `xml:"FERMENTABLES>FERMENTABLE"`
	Yeast []Yeast `xml:"YEASTS>YEAST"`
	Misc  []Misc  `xml:"MISCS>MISC"`
	Mash  []Mash  `xml:"MASH"`
}

type Hop struct {
	Name   string  `xml:"NAME"`
	Origin string  `xml:"ORIGIN"`
	Alpha  float64 `xml:"ALPHA"`
	Amount float64 `xml:"AMOUNT"`
	Use    string  `xml:"USE"`
	Time   float64 `xml:"TIME"`
	Notes  string  `xml:"NOTES"`
}

type Grain struct {
	Name   string  `xml:"NAME"`
	Origin string  `xml:"ORIGIN"`
	Type   string  `xml:"TYPE"`
	Amount float64 `xml:"AMOUNT"`
	Notes  string  `xml:"NOTES"`
}

type Yeast struct {
	Name         string  `xml:"NAME"`
	Type         string  `xml:"TYPE"`
	Origin       string  `xml:"ORIGIN"`
	Amount       float64 `xml:"AMOUNT"`
	Flocculation string  `xml:"FLOCCULATION"`
	Notes        string  `xml:"NOTES"`
}

type Misc struct {
	Name   string  `xml:"NAME"`
	Type   string  `xml:"TYPE"`
	Use    string  `xml:"USE"`
	Amount float64 `xml:"AMOUNT"`
	UseFor string  `xml:"USE_FOR"`
	Notes  string  `xml:"NOTES"`
}

type Mash struct {
	Name       string     `xml:"NAME"`
	SpargeTemp float64    `xml:"SPARGE_TEMP"`
	PH         string     `xml:"PH"`
	MashSteps  []MashStep `xml:"MASH_STEPS>MASH_STEP"`
}
type MashStep struct {
	Name     string  `xml:"NAME"`
	Type     string  `xml:"TYPE"`
	StepTime float64 `xml:"STEP_TIME"`
	StepTemp float64 `xml:"STEP_TEMP"`
}

func printHeader(header string, level int) {
	pad := ""

	switch level {
	case 1:
		pad = "="
	case 2:
		pad = "-"
	case 3:
		pad = "."
	}

	fmt.Printf("%s\n", header)
	print := strings.Repeat(pad, len(header))
	fmt.Println(print)
}

func printTable(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetHeaderSeparator("=")
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetCenterSeparator("+")
	table.SetRowSeparator("-")
	table.AppendBulk(data)
	table.Render()
	fmt.Println("")
}

func printHops(hops []Hop, level int) {
	printHeader("Hops", level)
	header := []string{"Hop", "Alpha", "Amount", "Time", "Notes"}
	data := make([][]string, 0, 10)

	for _, h := range hops {
		data = append(data, []string{h.Name, fmt.Sprintf("%.2f percent", h.Alpha), fmt.Sprintf("%.2f grams", h.Amount*1000), fmt.Sprintf("%.2fmin", h.Time), h.Notes})
	}
	printTable(header, data)

}

func printGrains(grains []Grain, level int) {
	printHeader("Grain", level)
	header := []string{"Grain", "Origin", "Amount", "Notes"}
	data := make([][]string, 0, 10)

	for _, g := range grains {
		data = append(data, []string{g.Name, g.Origin, fmt.Sprintf("%.2f kg", g.Amount), g.Notes})
	}
	printTable(header, data)
}

func printMiscs(miscs []Misc, level int) {
	printHeader("Miscellaneous", level)
	header := []string{"Name", "Use in step", "Amount", "Used for", "Notes"}
	data := make([][]string, 0, 10)

	for _, m := range miscs {
		data = append(data, []string{m.Name, m.Use, fmt.Sprintf("%.2f grams", m.Amount*1000), m.UseFor, m.Notes})
	}
	printTable(header, data)
}

func printYeasts(yeasts []Yeast, level int) {
	printHeader("Yeast", level)
	header := []string{"Name", "Type", "Flocculation"}
	data := make([][]string, 0, 10)

	for _, y := range yeasts {
		data = append(data, []string{y.Name, y.Type, y.Flocculation})
	}

	printTable(header, data)
}

func printIngrediens(r Recipe, level int) {
	printHeader("Ingredients", level)

	printHops(r.Hop, level+1)
	printGrains(r.Grain, level+1)
	printMiscs(r.Misc, level+1)
	printYeasts(r.Yeast, level+1)
}

func printMashStep(mashstep []MashStep, level int) {
	printHeader("Mash", level)
	header := []string{"Step", "Type", "Temperature", "Time"}
	data := make([][]string, 0, 10)

	for _, m := range mashstep {
		data = append(data, []string{m.Name, m.Type, fmt.Sprintf("%.2f C", m.StepTime), fmt.Sprintf("%.2f min", m.StepTemp)})
	}

	printTable(header, data)
}

func printMash(r Recipe, level int) {
	for _, m := range r.Mash {
		printMashStep(m.MashSteps, level)
	}
}

func printFermentation(r Recipe, level int) {
	printHeader("Fermentation", level)

	header := []string{"Stage", "Temperature", "Days"}
	data := make([][]string, 0, 10)
	data = append(data, []string{"Primary", fmt.Sprintf("%.0f", r.PrimaryTemp), fmt.Sprintf("%.0f", r.PrimaryAge)})
	data = append(data, []string{"Secondary", fmt.Sprintf("%.0f", r.SecondaryTemp), fmt.Sprintf("%.0f", r.SecondaryAge)})
	printTable(header, data)
	fmt.Println("")
}

func printRecipe(r Recipe, level int) {
	printHeader(r.Name, level)

	fmt.Printf(":Brewer: %s\n", r.Brewer)
	fmt.Println("")
	fmt.Printf(":Notes: %s\n", r.Notes)
	fmt.Println("")
	fmt.Println("")

	fmt.Printf(":Batch size: %.2fl\n", r.BatchSize)
	fmt.Printf(":Boil time: %.2fmin\n", r.BoilTime)
	fmt.Println("")

	printHeader("Measurements", level+1)
	header := []string{"Description", "Measured", "Estimated"}
	data := make([][]string, 0, 10)
	data = append(data, []string{"Original Gravity", fmt.Sprintf("%.3f", r.OG), r.EstOG})
	data = append(data, []string{"Final Gravity", fmt.Sprintf("%s", r.FG), r.EstFG})
	data = append(data, []string{"Alcohol By Volume", r.ABV, r.EstABV})
	printTable(header, data)

	printIngrediens(r, level+1)
	printMash(r, level+1)
	fmt.Println("")
	printFermentation(r, level+1)
}

func main() {
	level := flag.Int("level", 1, "Top header level")
	infile := flag.String("in", "beer.xml", "File to parse")
	flag.Parse()

	recipes := Recipes{}

	content, err := ioutil.ReadFile(*infile)
	if err != nil {
		panic(err)
	}

	err = xml.Unmarshal(content, &recipes)
	if err != nil {
		panic(err)
	}

	printRecipe(recipes.Recipe[0], *level)
}

package testutil

import (
	"fmt"
	"math/rand"
)

var left = []string{
	"awesome",
	"awkward",
	"brave",
	"boring",
	"cute",
	"cuddly",
	"delightful",
	"diligent",
	"elegant",
	"enlightened",
	"furry",
	"friendly",
	"giant",
	"gracious",
	"hilarious",
	"hungry",
	"intimidating",
	"ill",
	"jokingly",
	"janky",
	"killer",
	"kafkaesque",
	"lame",
	"luscious",
	"menacing",
	"masterful",
	"naked",
	"nutritious",
	"objective",
	"overt",
	"pugnacious",
	"perplexed",
	"quick",
	"quiet",
	"random",
	"rough",
	"solemn",
	"sarcastic",
	"tactical",
	"tactful",
	"unfortunate",
	"ubiquitous",
	"vulnerable",
	"vain",
	"weird",
	"wretched",
	"xenophobic",
	"xiphoid",
	"zealous",
	"zany",
}

var right = []string{
	"acamar",
	"acubens",
	"baekdu",
	"beid",
	"cebalrai",
	"castor",
	"dalim",
	"dombay",
	"ebla",
	"electra",
	"fang",
	"fawaris",
	"gacrux",
	"gomeisa",
	"hadar",
	"hatysa",
	"iklil",
	"intercrus",
	"jabbah",
	"jishui",
	"kaffaljidhma",
	"kang",
	"larawag",
	"lerna",
	"maasym",
	"maia",
	"nahn",
	"natasha",
	"ogma",
	"okab",
	"peacock",
	"pincoya",
	"ran",
	"rasalas",
	"sabik",
	"sadr",
	"taiyangshou",
	"tapecue",
	"ukdah",
	"unukalhai",
	"vega",
	"veritate",
	"wasat",
	"wazn",
	"xamidimura",
	"xuange",
	"yedposterior",
	"yedprior",
	"zaniah",
	"zaurak",
}

func RandomGUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func RandomName() string {
	name := left[rand.Intn(len(left))] + "_" + right[rand.Intn(len(right))]
	return fmt.Sprintf("%s%d", name, rand.Intn(10))
}

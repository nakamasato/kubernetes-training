package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nakamasato/kubernetes-training/doc"
	yaml "gopkg.in/yaml.v2"
)

func main() {

    f, err := os.Open("readme.yml")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    d := yaml.NewDecoder(f)

    var m map[string][]Version

    if err := d.Decode(&m); err != nil {
        log.Fatal(err)
    }

	for _, e := range m["versions"] {
		fmt.Printf("%v\n", e)
	}

	example(m["versions"])
}

type Version struct {
	Name string `yaml:"name"`
	Version string `yaml:"version"`
	RepoUrl string `yaml:"repoUrl"`
	Category string `yaml:"category"`
	ToDo bool `yaml:"todo"`
}

func (v Version) getReleaseUrl() string {
	if (v.Version == "") {
		return v.RepoUrl
	} else if (v.Version == "latest") {
		return v.RepoUrl + "/releases"
	} else {
		return v.RepoUrl + "/releases/tag/" + v.Version
	}
}

func (v Version) getToDoString() string {
	if (v.ToDo) {
		return "(ToDo)"
	} else {
		return ""
	}
}

func example(versions []Version) {
	book := doc.NewMarkDown()
	book.WriteTitle("Kubernetes Training", doc.LevelTitle).
		WriteLines(2)

	book.WriteList(fmt.Sprintf("Read on Website: %s", book.GetLink("https://www.nakamasato.com/kubernetes-training", "https://www.nakamasato.com/kubernetes-training")))
	book.WriteList(
		fmt.Sprintf(
			"Read on GitHub: %s (All files were moved to `contents` in %s",
			book.GetLink("contents", "contents"),
			book.GetLink("#105", "https://github.com/nakamasato/kubernetes-training/pull/105"),
		),
	)

	book.WriteTitle("Versions", doc.LevelSection)
	for _, version := range versions {
		book.WriteList(
			fmt.Sprintf(
				"%s: %s%s",
				version.Name,
				book.GetLink(version.Version, version.getReleaseUrl()),
				version.getToDoString(),
			),
		)
	}

	book.WriteTitle("Cloud Native Trail Map", doc.LevelSection)
	book.WriteList(book.GetLink("https://github.com/cncf/trailmap", "https://github.com/cncf/trailmap"))
	book.WriteList(book.GetLink("https://www.cncf.io/blog/2018/03/08/introducing-the-cloud-native-landscape-2-0-interactive-edition/", "https://www.cncf.io/blog/2018/03/08/introducing-the-cloud-native-landscape-2-0-interactive-edition/"))

	book.WriteWordLine(fmt.Sprintf("!%s", book.GetLink("", "https://github.com/cncf/trailmap/blob/master/CNCF_TrailMap_latest.png?raw=true")))

	err := book.Export("README.auto.md")
	if err != nil {
		log.Fatal(err)
	}
}

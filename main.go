package main

import (
	"log"
	"synchronex/common"
	"synchronex/schema"
)

//func main() {
//	err := command.Execute()
//	if err != nil {
//		log.Fatalln(err)
//	}
//}

func main() {
	state := schema.Nex{
		Files: []schema.File{
			schema.File{
				Source:      "synchonex-kt/build.gradle.kts",
				Destination: "~/.cache/synchronex/build.gradle.kts",
			},
			schema.File{
				Source:      "synchonex-kt/gradle.properties",
				Destination: "~/.cache/synchronex/gradle.properties",
				User:        "root",
			},
		},
	}
	plan := schema.Nex{
		User: "test",
		Files: []schema.File{
			schema.File{
				Source:      "synchonex-kt/build.gradle.kts",
				Destination: "~/.cache/synchronex/build.gradle.kts",
			},
			schema.File{
				Source:      "synchonex-kt/gradle.properties",
				Destination: "~/.cache/synchronex/gradle.properties",
			},
		},
	}

	result := plan.DifferencesFromState(state)
	log.Printf("Planned changes:\n\n%s", common.PrintPretty(result))
	common.WriteStatefile(plan)
}

package cli

import (
	"math/rand"
	"os"
	"strings"
)

var banners = []string{
	`
██████╗  ██████╗ ███████╗ ██████╗ ██████╗  ██████╗ ███████╗
██╔════╝ ██╔═══██╗██╔════╝██╔═══██╗██╔══██╗██╔════╝ ██╔════╝
██║  ███╗██║   ██║█████╗  ██║   ██║██████╔╝██║  ███╗█████╗  
██║   ██║██║   ██║██╔══╝  ██║   ██║██╔══██╗██║   ██║██╔══╝  
╚██████╔╝╚██████╔╝██║     ╚██████╔╝██║  ██║╚██████╔╝███████╗
 ╚═════╝  ╚═════╝ ╚═╝      ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚══════╝`,
	`
 ________  ________  ________ ________  ________  ________  _______      
|\   ____\|\   __  \|\  _____\\   __  \|\   __  \|\   ____\|\  ___ \     
\ \  \___|\ \  \|\  \ \  \__/\ \  \|\  \ \  \|\  \ \  \___|\ \   __/|    
 \ \  \  __\ \  \\\  \ \   __\\ \  \\\  \ \   _  _\ \  \  __\ \  \_|/__  
  \ \  \|\  \ \  \\\  \ \  \_| \ \  \\\  \ \  \\  \\ \  \|\  \ \  \_|\ \ 
   \ \_______\ \_______\ \__\   \ \_______\ \__\\ _\\ \_______\ \_______\
    \|_______|\|_______|\|__|    \|_______|\|__|\|__|\|_______|\|_______|
`,
}

func GetDescriptions(descriptionArg []string, _ bool) map[string]string {
	var description, banner string
	if descriptionArg != nil {
		if strings.Contains(strings.Join(os.Args[0:], ""), "-h") {
			description = descriptionArg[0]
		} else {
			description = descriptionArg[1]
		}
	} else {
		description = ""
	}
	bannerRandLen := len(banners)
	bannerRandIndex := rand.Intn(bannerRandLen)
	banner = banners[bannerRandIndex]
	return map[string]string{"banner": banner, "description": description}
}

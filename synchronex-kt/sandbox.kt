context {
    user = "portalsoup"
}

modules {
    path("../modules/zsh")
    path("../modules/i3")
}

batch("i3") {
    file.sync("/home/portalsoup/.config/i3/config") {
        src = "resources/i3-config"

        package("i3-wm") {
            version = "[4.23, )"
        }

        package("i3lock") {
            version = "[2.14, )"
        }
    }

    folder.sync("/home/{{USER}}/.config/i3/layouts") {
        src                 = "resources/i3layouts"
        pruneUntracked      = true // remove untracked files in folder
    }
}
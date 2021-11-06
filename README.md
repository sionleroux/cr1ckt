# Cr1ck_t

Still very WIP!

Entry for 2021 Game Off, themed "BUG".  I'll add a link to the submission when it's published.

Bugs in this game might be there intentionally...

## For game testers

For alpha testing use this [link to download the latest Windows EXE](https://nightly.link/sinisterstuf/cr1ck_t/workflows/build-exe/master/cr1ck_t-bundle.zip).

For Mac & Linux it is trivial to install go with your package manager and built it yourself, see the section for programmers below.

During development there's debug info at the top showing things like player location and velocity, level number, etc. and some controls:

- F: toggle full-screen
- N: go to next map
- Q: quit the game
- Space: jump (this is a real game control, not just for testing)

If the game is crashing you can get extra information about what went wrong if you start it from the console.  On Windows, that means:

1. Open PowerShell
2. Go to the folder where the game is, e.g. `cd C:\Users\Alice\Downloads`
3. Run it like this: `.\cr1ck_t.exe`

## For level makers

You can edit the levels using the Level Designer Toolkit ([LDtk](https://ldtk.io/)).

The current version of levels' maps is always in [assets/maps.ldtk](assets/maps.ldtk).  If you put a copy of that file with your changes in it in the same folder as the game binary (e.g. `cr1ck_t.exe`) then it will load that instead of the maps embedded in the binary, which is useful for prototyping and testing when developing a map because you don't have to compile any code.

- Auto-tiling is done on the layer called IntGrid (tutorial on [auto-tiling](https://ldtk.io/docs/tutorials/intgrid-layers/))
- [Entities](https://ldtk.io/docs/general/editor-components/entities/) (e.g. the player, monsters, items) are on the Entities layer

## For programmers

To build the game yourself, run: `go build .`

I'm using Go 1.17 for this.  You might have luck with an older version but the go:embed feature is only available in 1.16 so you can't go lower than that.

The Go build system will handle the rest of the dependencies for you but if you're curious, it's using:
- [ebiten](https://github.com/hajimehoshi/ebiten/) simple 2D game library
- [ldtkgo](https://github.com/SolarLune/ldtkgo) to interface with LDtk
- that's it so far

You can run the test suite with `go test ./...` but I haven't written any yet.

The project structure will probably stay quite simple, most logic is in the "main" file and gets extracted elsewhere as a clump of closely related code gets too big there.

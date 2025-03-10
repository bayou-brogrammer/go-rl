module github.com/bayou-brogrammer/go-rl

go 1.24

require (
	codeberg.org/anaseto/gruid v0.23.0
	codeberg.org/anaseto/gruid-sdl v0.5.0
	golang.org/x/image v0.24.0
	golang.org/x/text v0.22.0
)

require github.com/veandco/go-sdl2 v0.4.40 // indirect

replace codeberg.org/anaseto/gruid => /Users/jacoblecoq/Projects/go/gruid-libs/gruid
replace codeberg.org/anaseto/gruid-sdl => /Users/jacoblecoq/Projects/go/gruid-libs/gruid-sdl
replace codeberg.org/anaseto/gruid-tcell => /Users/jacoblecoq/Projects/go/gruid-libs/gruid-tcell

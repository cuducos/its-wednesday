module github.com/cuducos/its-wednesday

go 1.16

require (
	github.com/dghubble/go-twitter v0.0.0-20201011215211-4b180d0cc78d
	github.com/dghubble/oauth1 v0.7.0
)

replace github.com/dghubble/go-twitter => ../go-twitter // until merge https://github.com/dghubble/go-twitter/pull/148

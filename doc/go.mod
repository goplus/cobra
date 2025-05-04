module github.com/goplus/cobra/doc

go 1.18

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.6
	github.com/goplus/cobra v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v3 v3.0.1
)

require github.com/russross/blackfriday/v2 v2.1.0 // indirect

replace github.com/goplus/cobra => ../

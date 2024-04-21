package query

type Builder struct {
	Prefix string
	Indent string
	EOL    string
}

var Simple = Builder{
	EOL: " ",
}

var Beauty = Builder{
	Indent: "  ",
	EOL:    "\n",
}

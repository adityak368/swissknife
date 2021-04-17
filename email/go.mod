module github.com/adityak368/swissknife/email

go 1.16

replace github.com/adityak368/swissknife/email => ./

require (
	github.com/adityak368/swissknife/logger v0.0.0-20201107160000-5f5e30188eb2
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)

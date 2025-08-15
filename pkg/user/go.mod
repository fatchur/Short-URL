module user-service

go 1.23

replace short-url => ../../

require (
	gorm.io/gorm v1.30.1
	short-url v0.0.0
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.21.0 // indirect
)

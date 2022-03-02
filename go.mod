module bdo-rest-api

// +heroku goVersion go1.17
go 1.17

// A fork without hardcoded port filter
replace github.com/victorspringer/http-cache => github.com/octoman90/http-cache v0.0.0-20220223074738-0686ed29f408

require (
	github.com/gocolly/colly/v2 v2.1.0
	github.com/gorilla/mux v1.8.0
	github.com/victorspringer/http-cache v0.0.0-20190721184638-fe78e97af707
)

require (
	github.com/PuerkitoBio/goquery v1.7.1 // indirect
	github.com/andybalholm/cascadia v1.2.0 // indirect
	github.com/antchfx/htmlquery v1.2.3 // indirect
	github.com/antchfx/xmlquery v1.3.6 // indirect
	github.com/antchfx/xpath v1.2.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/kennygrant/sanitize v1.2.4 // indirect
	github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca // indirect
	github.com/temoto/robotstxt v1.1.2 // indirect
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

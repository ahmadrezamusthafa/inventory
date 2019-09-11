package inventory

import (
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/gorilla/mux"
	_ "github.com/rs/cors"
	_ "github.com/gorilla/sessions"
	_ "github.com/gorilla/securecookie"
	_  "github.com/pkg/errors"
	_ "github.com/mitchellh/mapstructure"
	_ "google.golang.org/api/sheets/v4"
)
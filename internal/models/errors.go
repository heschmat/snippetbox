package models

import "errors"

// N,B. error strings should not end with punctuation or newlines!!!! :\
var ErrNoRecord = errors.New("models: no matching record found")

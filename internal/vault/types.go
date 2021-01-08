/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 28/12/2020
*/

package vault

import (
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// Config contain connection element for vault implementation
type Config struct {
	// Address contains the vault address
	Address string

	// Token contains the token to use for vault connection
	Token string
}

const (
	Engine string = "ciops"
)

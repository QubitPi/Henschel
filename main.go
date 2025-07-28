// Copyright 2025 Jiaqi Liu. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"log" // Import log for logging server start
	"net/http"
)

func main() {
	// Define a command-line flag for the port, defaulting to :8080
	port := flag.String("port", ":8080", "Webservice port; default to 8080")
	flag.Parse() // Parse the command-line flags

	// Register the kongHandler function to handle requests to /deployKongApiGateway
	http.HandleFunc("/deployKongApiGateway", KongHandler)

	// Start the HTTP server on the specified port
	log.Printf("Server starting on port %s...", *port) // Log that the server is starting
	if err := http.ListenAndServe(*port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err) // Log and exit if the server fails to start
	}
}

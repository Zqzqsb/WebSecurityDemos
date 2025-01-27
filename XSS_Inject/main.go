package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Comment represents a stored message
type Comment struct {
	ID      uint   `gorm:"primarykey"`
	Content string `gorm:"not null"`
}

var db *gorm.DB

func main() {
	var err error
	// Connect to SQLite database
	db, err = gorm.Open(sqlite.Open("xss.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate schema
	db.AutoMigrate(&Comment{})

	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")

	// Main page with all XSS examples
	r.GET("/", func(c *gin.Context) {
		// Get stored comments
		var comments []Comment
		db.Find(&comments)

		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>XSS Attack Demo</title>
			<style>
				body { 
					font-family: Arial, sans-serif; 
					max-width: 800px; 
					margin: 0 auto; 
					padding: 20px;
					background-color: #f5f5f5;
				}
				.container { 
					margin-bottom: 20px; 
					padding: 20px; 
					border: 1px solid #ddd;
					border-radius: 8px;
					background-color: white;
					box-shadow: 0 2px 4px rgba(0,0,0,0.1);
				}
				input, textarea { 
					margin: 5px 0; 
					padding: 8px;
					width: 100%;
					max-width: 300px;
					border: 1px solid #ddd;
					border-radius: 4px;
				}
				button { 
					margin: 10px 0;
					padding: 8px 16px;
					background-color: #4CAF50;
					color: white;
					border: none;
					border-radius: 4px;
					cursor: pointer;
				}
				button:hover {
					background-color: #45a049;
				}
				.code-example {
					background-color: #f8f8f8;
					padding: 10px;
					border-radius: 4px;
					font-family: monospace;
					margin: 10px 0;
				}
				.result {
					margin-top: 10px;
					padding: 10px;
					border-radius: 4px;
					background-color: #f8f8f8;
				}
				.note {
					background-color: #fff3cd;
					border: 1px solid #ffeeba;
					border-radius: 4px;
					padding: 15px;
					margin-top: 15px;
				}
			</style>
		</head>
		<body>
			<h1>XSS (Cross-Site Scripting) Attack Demo</h1>
			
			<div class="container">
				<h2>1. Reflected XSS</h2>
				<div class="note">
					<p><strong>Description:</strong> Reflected XSS occurs when user input is immediately returned to the browser without proper sanitization.</p>
					<p><strong>Test Payload:</strong></p>
					<code>&lt;script&gt;alert('Reflected XSS!');&lt;/script&gt;</code>
				</div>
				<form action="/search" method="GET">
					<input type="text" name="q" placeholder="Search term...">
					<button type="submit">Search</button>
				</form>
			</div>

			<div class="container">
				<h2>2. Stored XSS</h2>
				<div class="note">
					<p><strong>Description:</strong> Stored XSS occurs when malicious content is saved on the server and later displayed to other users.</p>
					<p><strong>Test Payload:</strong></p>
					<code>&lt;script&gt;alert('Stored XSS!');&lt;/script&gt;</code>
				</div>
				<form action="/comment" method="POST">
					<textarea name="content" placeholder="Leave a comment..."></textarea>
					<button type="submit">Post Comment</button>
				</form>
				<div class="result">
					<h3>Comments:</h3>
					` + renderComments(comments) + `
				</div>
			</div>

			<div class="container">
				<h2>3. DOM-based XSS</h2>
				<div class="note">
					<p><strong>Description:</strong> DOM-based XSS occurs when JavaScript modifies the DOM with user-controlled data.</p>
					<p><strong>Test Payload:</strong></p>
					<code>&lt;img src=x onerror="alert('DOM XSS!');"&gt;</code>
				</div>
				<input type="text" id="userInput" placeholder="Enter your name...">
				<button onclick="showGreeting()">Show Greeting</button>
				<div id="output" class="result"></div>

				<script>
					function showGreeting() {
						// Unsafe: directly inserting user input into innerHTML
						var name = document.getElementById('userInput').value;
						document.getElementById('output').innerHTML = 'Hello, ' + name + '!';
					}

					// Get URL fragment and display it (DOM-based XSS)
					if(window.location.hash) {
						var hash = window.location.hash.slice(1);
						document.getElementById('output').innerHTML = decodeURIComponent(hash);
					}
				</script>
			</div>

			<div class="container">
				<h2>Safe Implementation Example</h2>
				<div class="note">
					<p><strong>Description:</strong> This section demonstrates proper input sanitization and safe rendering.</p>
				</div>
				<form action="/safe-comment" method="POST">
					<textarea name="content" placeholder="Leave a safe comment..."></textarea>
					<button type="submit">Post Safe Comment</button>
				</form>
				<div id="safeOutput" class="result">
					<h3>Safe Comments:</h3>
					` + renderSafeComments(comments) + `
				</div>
			</div>
		</body>
		</html>
		`
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, html)
	})

	// Reflected XSS endpoint
	r.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		// Unsafe: directly reflecting user input
		result := "Search results for: " + query
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, result)
	})

	// Stored XSS endpoint
	r.POST("/comment", func(c *gin.Context) {
		content := c.PostForm("content")
		// Unsafe: storing unfiltered user input
		db.Create(&Comment{Content: content})
		c.Redirect(http.StatusFound, "/")
	})

	// Safe comment endpoint
	r.POST("/safe-comment", func(c *gin.Context) {
		content := c.PostForm("content")
		// Safe: escaping HTML content
		safeContent := template.HTMLEscapeString(content)
		db.Create(&Comment{Content: safeContent})
		c.Redirect(http.StatusFound, "/")
	})

	r.Run(":8080")
}

// Unsafe rendering of comments
func renderComments(comments []Comment) string {
	var result strings.Builder
	for _, comment := range comments {
		// Unsafe: directly inserting user content
		result.WriteString("<div class='comment'>" + comment.Content + "</div>")
	}
	return result.String()
}

// Safe rendering of comments
func renderSafeComments(comments []Comment) string {
	var result strings.Builder
	for _, comment := range comments {
		// Safe: escaping HTML content
		safeContent := template.HTMLEscapeString(comment.Content)
		result.WriteString("<div class='comment'>" + safeContent + "</div>")
	}
	return result.String()
}

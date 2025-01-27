package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primarykey"`
	Username string `gorm:"uniqueIndex"`
	Password string
	Role     string `gorm:"default:'user'"`
}

var db *gorm.DB

// Unsafe login method - vulnerable to SQL injection
func unsafeLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Simulate some delay for time-based injection detection
	if strings.Contains(strings.ToLower(username), "sleep") {
		time.Sleep(2 * time.Second)
	}

	// Dangerous: directly concatenating SQL statements
	var result map[string]interface{}
	// Changed the query format to make basic authentication bypass work
	sql := fmt.Sprintf("SELECT * FROM users WHERE username='%s' AND password='%s' LIMIT 1", username, password)
	
	// Log the SQL query for demonstration
	log.Printf("Executing SQL: %s", sql)
	
	err := db.Raw(sql).Scan(&result).Error
	
	if err != nil {
		// Return error message for error-based injection demonstration
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Login failed with error: %v", err),
		})
		return
	}

	if len(result) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"user":    result,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Login failed: Invalid credentials",
		})
	}
}

// Safe login method - using parameterized queries
func safeLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user User
	result := db.Where("username = ? AND password = ?", username, password).First(&user)

	if result.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"user":    user,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Login failed: Invalid credentials",
		})
	}
}

func main() {
	var err error
	// Connect to SQLite database
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate schema
	db.AutoMigrate(&User{})

	// Create test users
	db.Create(&User{
		Username: "admin",
		Password: "123456",
		Role:     "admin",
	})
	db.Create(&User{
		Username: "user1",
		Password: "password1",
		Role:     "user",
	})
	db.Create(&User{
		Username: "user2",
		Password: "password2",
		Role:     "user",
	})

	r := gin.Default()

	// Provide a simple frontend page
	r.GET("/", func(c *gin.Context) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>SQL Injection Demo</title>
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
				input { 
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
					display: none;
				}
				.success {
					background-color: #dff0d8;
					color: #3c763d;
				}
				.error {
					background-color: #f2dede;
					color: #a94442;
				}
			</style>
		</head>
		<body>
			<h1>SQL Injection Demo</h1>
			
			<div class="container">
				<h2>Unsafe Login (Vulnerable to SQL Injection)</h2>
				<div class="note">
					<p><strong>Instructions:</strong></p>
					<ol>
						<li>Copy any test case into the username field</li>
						<li>Password can be anything</li>
						<li>Observe the response in the result box below</li>
						<li>Try modifying the conditions in the test cases</li>
					</ol>
				</div>
				<form id="unsafeForm">
					<input type="text" name="username" placeholder="Username"><br>
					<input type="password" name="password" placeholder="Password"><br>
					<button type="submit">Login</button>
				</form>
				<div id="unsafeResult" class="result"></div>
				
				<div class="code-example">
					<h3>SQL Injection Test Cases:</h3>
					
					<h4>1. Basic Authentication Bypass</h4>
					<code>' OR '1'='1</code>
					<p>This injection makes the WHERE clause always true, bypassing authentication</p>
					
					<h4>2. Comment-Based Injection</h4>
					<code>admin'--</code>
					<p>Uses SQL comments to ignore the password check</p>
					
					<h4>3. UNION-Based Query</h4>
					<code>admin' UNION SELECT 1 as id, 'hacker' as username, 'pwned' as password, 'admin' as role --</code>
					<p>Uses UNION to combine results with a fake user record</p>
					
					<h4>4. Boolean-Based Blind</h4>
					<code>admin' AND (SELECT CASE WHEN (1=1) THEN 1 ELSE 0 END)='1</code>
					<p>Tests database conditions through true/false responses</p>
					
					<h4>5. Time-Based Blind</h4>
					<code>admin' AND (SELECT CASE WHEN (1=1) THEN sqlite3_sleep(2000) ELSE 1 END)='1</code>
					<p>Causes a delay when condition is true, useful for blind injection</p>
					
					<h4>6. Error-Based</h4>
					<code>admin' AND (SELECT CASE WHEN (1=1) THEN CAST('a' AS INTEGER) ELSE 1 END)='1</code>
					<p>Triggers database errors to extract information</p>
				</div>
			</div>

			<div class="container">
				<h2>Safe Login (Using Parameterized Queries)</h2>
				<div class="note">
					<p><strong>Security Note:</strong></p>
					<p>Try the same injection patterns here - they won't work because:</p>
					<ul>
						<li>Uses parameterized queries instead of string concatenation</li>
						<li>Special characters are properly escaped</li>
						<li>SQL statements and user data are kept separate</li>
					</ul>
					<p>Valid credentials: username="admin", password="123456"</p>
				</div>
				<form id="safeForm">
					<input type="text" name="username" placeholder="Username"><br>
					<input type="password" name="password" placeholder="Password"><br>
					<button type="submit">Login</button>
				</form>
				<div id="safeResult" class="result"></div>
			</div>

			<style>
				.note {
					background-color: #fff3cd;
					border: 1px solid #ffeeba;
					border-radius: 4px;
					padding: 15px;
					margin-top: 15px;
				}
				.code-example {
					background-color: #f8f9fa;
					border: 1px solid #eaecf0;
					border-radius: 4px;
					padding: 15px;
					margin: 15px 0;
				}
				.code-example h4 {
					color: #2c3e50;
					margin-top: 20px;
					margin-bottom: 10px;
				}
				.code-example code {
					display: block;
					background-color: #272822;
					color: #f8f8f2;
					padding: 10px;
					border-radius: 4px;
					margin: 10px 0;
					white-space: pre-wrap;
					word-wrap: break-word;
				}
				.code-example p {
					color: #666;
					margin: 5px 0 15px 0;
				}
				ol, ul {
					margin: 10px 0;
					padding-left: 20px;
				}
				li {
					margin: 5px 0;
				}
			</style>

			<script>
				function showResult(elementId, success, message) {
					const element = document.getElementById(elementId);
					element.style.display = 'block';
					element.className = 'result ' + (success ? 'success' : 'error');
					element.textContent = message;
				}

				document.getElementById('unsafeForm').onsubmit = async (e) => {
					e.preventDefault();
					const formData = new FormData(e.target);
					try {
						const response = await fetch('/unsafe/login', {
							method: 'POST',
							body: formData
						});
						const result = await response.json();
						showResult('unsafeResult', response.ok, result.message);
					} catch (error) {
						showResult('unsafeResult', false, 'Request failed: ' + error.message);
					}
				};

				document.getElementById('safeForm').onsubmit = async (e) => {
					e.preventDefault();
					const formData = new FormData(e.target);
					try {
						const response = await fetch('/safe/login', {
							method: 'POST',
							body: formData
						});
						const result = await response.json();
						showResult('safeResult', response.ok, result.message);
					} catch (error) {
						showResult('safeResult', false, 'Request failed: ' + error.message);
					}
				};
			</script>
		</body>
		</html>
		`
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, html)
	})

	r.POST("/unsafe/login", unsafeLogin)
	r.POST("/safe/login", safeLogin)

	r.Run(":8080")
}

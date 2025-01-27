package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID       uint   `gorm:"primarykey"`
	Username string `gorm:"unique"`
	Balance  int    // User's account balance
}

// Transfer represents a money transfer
type Transfer struct {
	ID          uint
	FromUserID  uint
	ToUserID    uint
	Amount      int
	Description string
}

var (
	db         *gorm.DB
	csrfTokens sync.Map // Store CSRF tokens
)

// generateCSRFToken generates a new CSRF token
func generateCSRFToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func main() {
	// Connect to SQLite database
	var err error
	db, err = gorm.Open(sqlite.Open("csrf.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate schema
	db.AutoMigrate(&User{}, &Transfer{})

	// Initialize some users if they don't exist
	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		db.Create(&User{Username: "alice", Balance: 1000})
		db.Create(&User{Username: "bob", Balance: 1000})
	}

	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")

	// Vulnerable transfer endpoint (no CSRF protection)
	r.POST("/transfer/unsafe", func(c *gin.Context) {
		toUsername := c.PostForm("to")
		amountStr := c.PostForm("amount")
		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
			return
		}

		// Validate amount
		if amount <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
			return
		}

		fromUsername := c.GetString("currentUser") // In real app, get from session
		if fromUsername == "" {
			fromUsername = "alice" // For demo, assume we're always logged in as alice
		}

		var fromUser, toUser User
		if err := db.Where("username = ?", fromUsername).First(&fromUser).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sender"})
			return
		}
		if err := db.Where("username = ?", toUsername).First(&toUser).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipient"})
			return
		}

		// Check balance
		if fromUser.Balance < amount {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
			return
		}

		// Begin transaction
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// Update balances
		if err := tx.Model(&fromUser).Update("balance", fromUser.Balance-amount).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transfer failed"})
			return
		}

		if err := tx.Model(&toUser).Update("balance", toUser.Balance+amount).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transfer failed"})
			return
		}

		// Create transfer record
		if err := tx.Create(&Transfer{
			FromUserID:  fromUser.ID,
			ToUserID:    toUser.ID,
			Amount:      amount,
			Description: "Transfer from " + fromUsername + " to " + toUsername,
		}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transfer failed"})
			return
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transfer failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Transfer successful",
			"from": gin.H{
				"username": fromUser.Username,
				"balance":  fromUser.Balance - amount,
			},
			"to": gin.H{
				"username": toUser.Username,
				"balance":  toUser.Balance + amount,
			},
			"amount": amount,
		})
	})

	// Safe transfer endpoint (with CSRF protection)
	r.POST("/transfer/safe", func(c *gin.Context) {
		// Verify CSRF token
		token := c.GetHeader("X-CSRF-Token")
		expectedToken, exists := csrfTokens.Load(c.GetString("currentUser"))
		if !exists || token != expectedToken.(string) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid CSRF token"})
			return
		}

		toUsername := c.PostForm("to")
		amountStr := c.PostForm("amount")

		// Convert amount string to integer
		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
			return
		}

		// Validate amount
		if amount <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
			return
		}

		fromUsername := c.GetString("currentUser")
		if fromUsername == "" {
			fromUsername = "alice"
		}

		var fromUser, toUser User
		if err := db.Where("username = ?", fromUsername).First(&fromUser).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sender"})
			return
		}
		if err := db.Where("username = ?", toUsername).First(&toUser).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipient"})
			return
		}

		// Check balance
		if fromUser.Balance < amount {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
			return
		}

		// Begin transaction
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// Update balances
		if err := tx.Model(&fromUser).Update("balance", fromUser.Balance-amount).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transfer failed"})
			return
		}

		if err := tx.Model(&toUser).Update("balance", toUser.Balance+amount).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transfer failed"})
			return
		}

		// Create transfer record
		if err := tx.Create(&Transfer{
			FromUserID:  fromUser.ID,
			ToUserID:    toUser.ID,
			Amount:      amount,
			Description: "Safe transfer from " + fromUsername + " to " + toUsername,
		}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transfer failed"})
			return
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transfer failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Safe transfer successful",
			"from": gin.H{
				"username": fromUser.Username,
				"balance":  fromUser.Balance - amount,
			},
			"to": gin.H{
				"username": toUser.Username,
				"balance":  toUser.Balance + amount,
			},
			"amount": amount,
		})
	})

	// Main page
	r.GET("/", func(c *gin.Context) {
		// Generate CSRF token for the current user
		token, err := generateCSRFToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CSRF token"})
			return
		}
		csrfTokens.Store("alice", token) // In real app, store per user

		var users []User
		db.Find(&users)

		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>CSRF Attack Demo</title>
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
				}
				.note {
					background-color: #fff3cd;
					border: 1px solid #ffeeba;
					padding: 15px;
					margin: 15px 0;
					border-radius: 4px;
				}
				.code-example {
					background-color: #f8f8f8;
					padding: 10px;
					border-radius: 4px;
					font-family: monospace;
				}
				button {
					background-color: #4CAF50;
					color: white;
					padding: 10px 20px;
					border: none;
					border-radius: 4px;
					cursor: pointer;
				}
				button:hover {
					background-color: #45a049;
				}
				input {
					padding: 8px;
					margin: 5px 0;
					border: 1px solid #ddd;
					border-radius: 4px;
				}
			</style>
		</head>
		<body>
			<h1>CSRF Attack Demonstration</h1>
			
			<div class="container">
				<h2>Current User: Alice</h2>
				<div class="note">
					<p>For demonstration purposes, you are logged in as Alice.</p>
					<p>Balance: 1000</p>
				</div>
			</div>

			<div class="container">
				<h2>1. Vulnerable Transfer Form (No CSRF Protection)</h2>
				<div class="note">
					<p><strong>Description:</strong> This form is vulnerable to CSRF attacks because it doesn't implement any CSRF protection.</p>
				</div>
				<form id="unsafeForm" action="/transfer/unsafe" method="POST">
					<input type="text" name="to" placeholder="Recipient username" value="bob"><br>
					<input type="number" name="amount" placeholder="Amount" value="100"><br>
					<button type="submit">Transfer (Unsafe)</button>
				</form>
			</div>

			<div class="container">
				<h2>2. Protected Transfer Form (With CSRF Token)</h2>
				<div class="note">
					<p><strong>Description:</strong> This form is protected against CSRF attacks using a CSRF token.</p>
				</div>
				<form id="safeForm" action="/transfer/safe" method="POST">
					<input type="hidden" name="_csrf" value="` + token + `">
					<input type="text" name="to" placeholder="Recipient username" value="bob"><br>
					<input type="number" name="amount" placeholder="Amount" value="100"><br>
					<button type="submit">Transfer (Safe)</button>
				</form>
			</div>

			<div class="container">
				<h2>3. CSRF Attack Simulation</h2>
				<div class="note">
					<p><strong>Description:</strong> This section demonstrates how a CSRF attack might be carried out.</p>
				</div>
				<div class="code-example">
					<p>Malicious HTML that might be hosted on attacker.com:</p>
					<pre>
&lt;form id="malicious" action="http://localhost:8080/transfer/unsafe" method="POST" style="display:none"&gt;
    &lt;input type="text" name="to" value="attacker"&gt;
    &lt;input type="number" name="amount" value="500"&gt;
&lt;/form&gt;
&lt;script&gt;document.getElementById('malicious').submit();&lt;/script&gt;
					</pre>
				</div>
				<button onclick="simulateAttack()">Simulate CSRF Attack</button>
			</div>

			<script>
				// Add CSRF token to safe form submissions
				document.getElementById('safeForm').addEventListener('submit', function(e) {
					e.preventDefault();
					fetch('/transfer/safe', {
						method: 'POST',
						headers: {
							'X-CSRF-Token': '` + token + `'
						},
						body: new FormData(this)
					})
					.then(response => response.json())
					.then(data => alert(data.message))
					.catch(error => alert('Error: ' + error));
				});

				// Simulate CSRF attack
				function simulateAttack() {
					var form = document.createElement('form');
					form.method = 'POST';
					form.action = '/transfer/unsafe';
					
					var to = document.createElement('input');
					to.type = 'hidden';
					to.name = 'to';
					to.value = 'attacker';
					form.appendChild(to);
					
					var amount = document.createElement('input');
					amount.type = 'hidden';
					amount.name = 'amount';
					amount.value = '500';
					form.appendChild(amount);
					
					document.body.appendChild(form);
					form.submit();
				}
			</script>
		</body>
		</html>
		`
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, html)
	})

	r.Run(":8080")
}

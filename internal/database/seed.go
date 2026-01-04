package database

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/config"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database/ent"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database/ent/blog"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database/ent/user"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

//go:embed dto/seed_data.json
var seedDataJSON []byte

type SeedData struct {
	Users     []SeedUser    `json:"users"`
	Addresses []SeedAddress `json:"addresses"`
	Blogs     []SeedBlog    `json:"blogs"`
	Comments  []SeedComment `json:"comments"`
}

type SeedUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type SeedAddress struct {
	Username string `json:"username"`
	Street   string `json:"street"`
	City     string `json:"city"`
	State    string `json:"state"`
	Zip      string `json:"zip"`
}

type SeedBlog struct {
	AuthorUsername string `json:"author_username"`
	Title          string `json:"title"`
	Content        string `json:"content"`
}

type SeedComment struct {
	BlogTitle      string `json:"blog_title"`
	AuthorUsername string `json:"author_username"`
	Content        string `json:"content"`
}

// SeedDatabase seeds the database with initial data
func SeedDatabase() error {
	// Load configuration
	if err := config.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Build database connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.Name,
	)

	// Create Ent client
	client, err := ent.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer client.Close()

	// Parse seed data
	var seedData SeedData
	if err := json.Unmarshal(seedDataJSON, &seedData); err != nil {
		return fmt.Errorf("failed to parse seed data: %w", err)
	}

	ctx := context.Background()

	// Seed users
	fmt.Println("Seeding users...")
	userMap := make(map[string]*ent.User)
	for _, u := range seedData.Users {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", u.Username, err)
		}

		// Check if user exists
		existingUser, err := client.User.Query().Where(user.UsernameEQ(u.Username)).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return fmt.Errorf("failed to query user %s: %w", u.Username, err)
		}

		var createdUser *ent.User
		if existingUser != nil {
			// Update existing user
			createdUser, err = existingUser.Update().
				SetPassword(string(hashedPassword)).
				SetEmail(u.Email).
				SetRole(user.Role(u.Role)).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed to update user %s: %w", u.Username, err)
			}
		} else {
			// Create new user
			createdUser, err = client.User.
				Create().
				SetUsername(u.Username).
				SetPassword(string(hashedPassword)).
				SetEmail(u.Email).
				SetRole(user.Role(u.Role)).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed to create user %s: %w", u.Username, err)
			}
		}

		userMap[u.Username] = createdUser
		fmt.Printf("  ✓ Seeded user: %s (%s)\n", u.Username, u.Role)
	}

	// Seed addresses
	fmt.Println("Seeding addresses...")
	for _, a := range seedData.Addresses {
		u, exists := userMap[a.Username]
		if !exists {
			return fmt.Errorf("user %s not found for address", a.Username)
		}

		// Check if address exists for this user
		existingAddresses, err := u.QueryAddresses().All(ctx)
		if err != nil {
			return fmt.Errorf("failed to query addresses for %s: %w", a.Username, err)
		}

		if len(existingAddresses) > 0 {
			// Update first address
			_, err = existingAddresses[0].Update().
				SetStreet(a.Street).
				SetCity(a.City).
				SetState(a.State).
				SetZip(a.Zip).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed to update address for %s: %w", a.Username, err)
			}
		} else {
			// Create new address
			_, err = client.Address.
				Create().
				SetStreet(a.Street).
				SetCity(a.City).
				SetState(a.State).
				SetZip(a.Zip).
				SetUser(u).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed to create address for %s: %w", a.Username, err)
			}
		}

		fmt.Printf("  ✓ Seeded address for: %s\n", a.Username)
	}

	// Seed blogs
	fmt.Println("Seeding blogs...")
	blogMap := make(map[string]*ent.Blog)
	for _, b := range seedData.Blogs {
		u, exists := userMap[b.AuthorUsername]
		if !exists {
			return fmt.Errorf("user %s not found for blog", b.AuthorUsername)
		}

		// Check if blog exists (by title and author)
		existingBlog, err := client.Blog.Query().
			Where(blog.TitleEQ(b.Title)).
			Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return fmt.Errorf("failed to query blog '%s': %w", b.Title, err)
		}

		var createdBlog *ent.Blog
		if existingBlog != nil {
			// Update existing blog
			createdBlog, err = existingBlog.Update().
				SetContent(b.Content).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed to update blog '%s': %w", b.Title, err)
			}
		} else {
			// Create new blog
			createdBlog, err = client.Blog.
				Create().
				SetTitle(b.Title).
				SetContent(b.Content).
				SetAuthor(u).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed to create blog '%s': %w", b.Title, err)
			}
		}

		blogMap[b.Title] = createdBlog
		fmt.Printf("  ✓ Seeded blog: %s (by %s)\n", b.Title, b.AuthorUsername)
	}

	// Seed comments
	fmt.Println("Seeding comments...")
	for _, c := range seedData.Comments {
		u, exists := userMap[c.AuthorUsername]
		if !exists {
			return fmt.Errorf("user %s not found for comment", c.AuthorUsername)
		}

		b, exists := blogMap[c.BlogTitle]
		if !exists {
			return fmt.Errorf("blog '%s' not found for comment", c.BlogTitle)
		}

		// Check if comment exists (same author, blog, and content)
		// Simplified: just check if we can find any matching comment
		commentExists := false
		allComments, err := b.QueryComments().All(ctx)
		if err == nil {
			for _, existing := range allComments {
				existingAuthor, _ := existing.QueryAuthor().Only(ctx)
				if existingAuthor != nil && existingAuthor.Username == c.AuthorUsername && existing.Content == c.Content {
					commentExists = true
					break
				}
			}
		}

		if !commentExists {
			// Create new comment only if it doesn't exist
			_, err = client.Comment.
				Create().
				SetContent(c.Content).
				SetAuthor(u).
				SetBlog(b).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed to create comment on '%s': %w", c.BlogTitle, err)
			}
		}

		fmt.Printf("  ✓ Seeded comment on: %s (by %s)\n", c.BlogTitle, c.AuthorUsername)
	}

	fmt.Println("✅ Database seeding completed successfully!")
	return nil
}

package postgres

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) error {

	if err := seedBooks(db); err != nil {
		return fmt.Errorf("error seeding books: %w", err)
	}

	fmt.Println("Database seeding completed successfully")
	return nil
}

func seedBooks(db *gorm.DB) error {

	var count int64
	// if err := db.Model().Count(&count).Error; err != nil {
	// 	return err
	// }

	if count > 0 {
		fmt.Println("Books table already has data, skipping seed")
		return nil
	}

	// books := []entities.Book{
	// 	{
	// 		ID:          mustGenerateUUID(),
	// 		Title:       "To Kill a Mockingbird",
	// 		Author:      "Harper Lee",
	// 		ISBN:        "9780061120084",
	// 		Publisher:   "HarperCollins",
	// 		PublishDate: parseDate("1960-07-11"),
	// 		Description: "The unforgettable novel of a childhood in a sleepy Southern town and the crisis of conscience that rocked it.",
	// 		Copies:      10,
	// 		CreatedAt:   time.Now(),
	// 		UpdatedAt:   time.Now(),
	// 	},
	// 	{
	// 		ID:          mustGenerateUUID(),
	// 		Title:       "1984",
	// 		Author:      "George Orwell",
	// 		ISBN:        "9780451524935",
	// 		Publisher:   "Signet Classic",
	// 		PublishDate: parseDate("1949-06-08"),
	// 		Description: "A dystopian novel set in Airstrip One, a province of the superstate Oceania in a world of perpetual war.",
	// 		Copies:      7,
	// 		CreatedAt:   time.Now(),
	// 		UpdatedAt:   time.Now(),
	// 	},
	// 	{
	// 		ID:          mustGenerateUUID(),
	// 		Title:       "The Great Gatsby",
	// 		Author:      "F. Scott Fitzgerald",
	// 		ISBN:        "9780743273565",
	// 		Publisher:   "Scribner",
	// 		PublishDate: parseDate("1925-04-10"),
	// 		Description: "A portrait of the Jazz Age in all of its decadence and excess.",
	// 		Copies:      5,
	// 		CreatedAt:   time.Now(),
	// 		UpdatedAt:   time.Now(),
	// 	},
	// 	{
	// 		ID:          mustGenerateUUID(),
	// 		Title:       "Pride and Prejudice",
	// 		Author:      "Jane Austen",
	// 		ISBN:        "9780141439518",
	// 		Publisher:   "Penguin Classics",
	// 		PublishDate: parseDate("1813-01-28"),
	// 		Description: "A romantic novel of manners that follows the character development of Elizabeth Bennet.",
	// 		Copies:      8,
	// 		CreatedAt:   time.Now(),
	// 		UpdatedAt:   time.Now(),
	// 	},
	// 	{
	// 		ID:          mustGenerateUUID(),
	// 		Title:       "The Hobbit",
	// 		Author:      "J.R.R. Tolkien",
	// 		ISBN:        "9780547928227",
	// 		Publisher:   "Houghton Mifflin Harcourt",
	// 		PublishDate: parseDate("1937-09-21"),
	// 		Description: "A fantasy novel about the adventures of hobbit Bilbo Baggins.",
	// 		Copies:      12,
	// 		CreatedAt:   time.Now(),
	// 		UpdatedAt:   time.Now(),
	// 	},
	// }

	// result := db.Create(&books)
	// if result.Error != nil {
	// 	return result.Error
	// }

	// fmt.Printf("Seeded %d books successfully\n", len(books))
	return nil
}

func mustGenerateUUID() uuid.UUID {
	id, err := uuid.NewV4()
	if err != nil {
		panic("failed to generate UUID: " + err.Error())
	}
	return id
}

func parseDate(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		panic("failed to parse date: " + err.Error())
	}
	return date
}

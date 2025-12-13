package models

// This file provides examples of how to use the GORM models

/*
EXAMPLE USAGE:

1. Run migrations (typically in main.go):
   db := database.New()
   if err := db.Migrate(); err != nil {
       log.Fatal("Failed to migrate database:", err)
   }

2. Create a new user:
   user := &models.User{
       Username: "johndoe",
       Email:    "john@example.com",
       Password: "hashed_password_here",
       Role:     "customer",
   }
   result := db.DB().Create(user)

3. Query users:
   var users []models.User
   db.DB().Where("role = ?", "customer").Find(&users)

4. Create a movie:
   movie := &models.Movie{
       Title:           "The Matrix",
       Description:     "A computer hacker learns about the true nature of reality",
       DurationMinutes: 136,
       PosterURL:       "https://example.com/poster.jpg",
       Rating:          "R",
   }
   db.DB().Create(movie)

5. Create a studio:
   studio := &models.Studio{
       Name:      "Studio 1",
       TotalRows: 10, // A-J
       TotalCols: 8,  // 1-8
   }
   db.DB().Create(studio)

6. Create a showtime:
   showtime := &models.Showtime{
       MovieID:   1,
       StudioID:  1,
       StartTime: time.Now().Add(24 * time.Hour),
       Price:     15.99,
   }
   db.DB().Create(showtime)

7. Create a booking with tickets:
   booking := &models.Booking{
       UserID:        userUUID,
       InvoiceNumber: "INV-2025-0001",
       TotalAmount:   31.98, // 2 tickets * 15.99
       Status:        "PENDING",
   }
   db.DB().Create(booking)

   // Create tickets for the booking
   tickets := []models.Ticket{
       {BookingID: booking.ID, ShowtimeID: 1, SeatNumber: "A5"},
       {BookingID: booking.ID, ShowtimeID: 1, SeatNumber: "A6"},
   }
   db.DB().Create(&tickets)

8. Query with relationships:
   var booking models.Booking
   db.DB().Preload("User").Preload("Tickets.Showtime.Movie").First(&booking, bookingID)

9. Check seat availability:
   var count int64
   db.DB().Model(&models.Ticket{}).
       Where("showtime_id = ? AND seat_number = ?", showtimeID, "A5").
       Count(&count)
   if count > 0 {
       // Seat is already taken
   }

10. Soft delete (marks deleted_at, doesn't actually delete):
    db.DB().Delete(&movie, movieID)

11. Query including soft-deleted records:
    db.DB().Unscoped().Find(&movies)

12. Update booking status:
    db.DB().Model(&booking).Update("status", "PAID")

CRITICAL LOGIC:
- The Ticket model has a BeforeCreate hook that prevents duplicate seat bookings
- A unique composite index on (showtime_id, seat_number) enforces seat uniqueness at DB level
- All relationships use UUID for users/bookings and auto-increment integers for others
- Soft deletes are enabled on users, movies, studios, and showtimes
*/

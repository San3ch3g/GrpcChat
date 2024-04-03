package pg

import (
	"ModuleForChat/internal/pkg/models"
	"ModuleForChat/internal/utils/config"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

type Storage struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func buildDSN(cfg *config.Config) string {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", cfg.DbUser, cfg.DbPassword, cfg.DbName)
	return dsn
}

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(buildDSN(cfg)), &gorm.Config{
		TranslateError: true,
	})
}

func MustNewPostgresDB(cfg *config.Config) *gorm.DB {
	db, err := NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Message{}, &models.Room{}); err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	return db
}

func (s *Storage) RegisterUser(user *models.User) (uint32, error) {
	var existingUser models.User
	result := s.db.Where("nickname = ?", user.Nickname).First(&existingUser)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, result.Error
	}

	if result.RowsAffected > 0 {
		return 0, errors.New("user with this nickname already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user.Password = string(hashedPassword)

	err = s.db.Create(&user).Error
	if err != nil {
		return 0, err
	}

	return user.Id, nil
}

func (s *Storage) LoginUser(user *models.User) (uint32, error) {
	var foundUser models.User
	result := s.db.Where("nickname = ?", user.Nickname).First(&foundUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, errors.New("user not found")
		}
		return 0, result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, errors.New("invalid password")
		}
		return 0, err
	}

	return foundUser.Id, nil
}

func (s *Storage) CreateRoom(roomName string) (*models.Room, error) {
	room := &models.Room{
		Name: roomName,
	}

	if err := s.db.Create(room).Error; err != nil {
		return nil, err
	}

	return room, nil
}

func (s *Storage) SendMessage(roomID uint32, senderID uint32, content string) (*models.Message, error) {
	if content == "" {
		return nil, errors.New("content cannot be empty")
	}

	var room models.Room
	if err := s.db.First(&room, roomID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("room with id %d not found", roomID)
		}
		return nil, err
	}

	var sender models.User
	if err := s.db.First(&sender, senderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with id %d not found", senderID)
		}
		return nil, err
	}
	message := &models.Message{
		RoomId:   roomID,
		SenderId: senderID,
		Content:  content,
		Date:     time.Now(),
	}

	result := s.db.Create(message)
	if result.Error != nil {
		return nil, result.Error
	}

	return message, nil
}

func (s *Storage) SaveMessage(message *models.Message) (uint32, error) {
	if err := s.db.Create(message).Error; err != nil {
		return 0, err
	}

	return message.Id, nil
}

func (s *Storage) GetMessagesByRoomID(roomID uint32) ([]models.Message, error) {
	var messages []models.Message

	if err := s.db.Where("room_id = ?", roomID).Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

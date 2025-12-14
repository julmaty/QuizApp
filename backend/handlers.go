package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func enableCorsHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func listQuizzes(c *gin.Context) {
	enableCorsHeaders(c)
	if c.Request.Method == http.MethodOptions {
		c.Status(http.StatusOK)
		return
	}
	var quizzes []Quiz
	if err := db.Order("created_at desc").Preload("Questions.Options").Find(&quizzes).Error; err != nil {
		c.String(http.StatusInternalServerError, "db error")
		return
	}
	// For legacy rows that stored the whole quiz JSON in a `data` JSONB column,
	// populate Questions from that column when no normalized Questions exist.
	for i := range quizzes {
		if len(quizzes[i].Questions) == 0 {
			var raw struct{ Data json.RawMessage }
			if err := db.Raw("SELECT data FROM quizzes WHERE id = ?", quizzes[i].ID).Scan(&raw).Error; err == nil && len(raw.Data) > 0 {
				// temporary struct matching old payload shape
				var tmp struct {
					Questions []struct {
						Text     string   `json:"text"`
						Options  []string `json:"options"`
						Multiple bool     `json:"multiple"`
					} `json:"questions"`
				}
				if err := json.Unmarshal(raw.Data, &tmp); err == nil {
					qs := make([]Question, 0, len(tmp.Questions))
					for _, tq := range tmp.Questions {
						opts := make([]Option, 0, len(tq.Options))
						for oi, ot := range tq.Options {
							opts = append(opts, Option{Text: ot, Ord: oi})
						}
						qs = append(qs, Question{Text: tq.Text, Options: opts, Multiple: tq.Multiple})
					}
					quizzes[i].Questions = qs
				}
			}
		}
	}
	c.JSON(http.StatusOK, quizzes)
}

func createQuiz(c *gin.Context) {
	enableCorsHeaders(c)
	if c.Request.Method == http.MethodOptions {
		c.Status(http.StatusOK)
		return
	}
	var in struct {
		Title     string `json:"title"`
		Questions []struct {
			Text     string   `json:"text"`
			Options  []string `json:"options"`
			Multiple bool     `json:"multiple"`
		} `json:"questions"`
	}
	if err := c.BindJSON(&in); err != nil {
		c.String(http.StatusBadRequest, "invalid json")
		return
	}
	var q Quiz
	q.Title = in.Title
	q.Questions = make([]Question, 0, len(in.Questions))
	for _, iq := range in.Questions {
		opts := make([]Option, 0, len(iq.Options))
		for oi, ot := range iq.Options {
			opts = append(opts, Option{Text: ot, Ord: oi})
		}
		q.Questions = append(q.Questions, Question{Text: iq.Text, Options: opts, Multiple: iq.Multiple})
	}
	if q.ID == "" {
		q.ID = q.Title + "-" + time.Now().Format("20060102150405.000000")
	}
	if q.CreatedAt.IsZero() {
		q.CreatedAt = time.Now().UTC()
	}
	tx := db.Create(&q)
	if tx.Error != nil {
		c.String(http.StatusInternalServerError, "db error")
		return
	}
	if err := db.Preload("Questions.Options").First(&q, "id = ?", q.ID).Error; err != nil {
		c.String(http.StatusInternalServerError, "db error")
		return
	}
	c.JSON(http.StatusCreated, q)
}

// --- Auth helpers & handlers ---

var jwtSecret = []byte(getEnv("JWT_SECRET", "dev-secret"))

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

// createPasswordHash hashes a password using bcrypt
func createPasswordHash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// comparePassword compares password with hash
func comparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// issueToken returns a signed JWT for a user id
func issueToken(userID uint) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	return t.SignedString(jwtSecret)
}

// parseToken returns userID from token or error
func parseToken(tokenStr string) (uint, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil || token == nil || !token.Valid {
		return 0, errors.New("invalid token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(float64); ok {
			return uint(sub), nil
		}
	}
	return 0, errors.New("invalid token claims")
}

// registerUser: POST /api/register
func registerUser(c *gin.Context) {
	enableCorsHeaders(c)
	if c.Request.Method == http.MethodOptions {
		c.Status(http.StatusOK)
		return
	}
	var in struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"displayName"`
	}
	if err := c.BindJSON(&in); err != nil {
		c.String(http.StatusBadRequest, "invalid json")
		return
	}
	if in.Email == "" || in.Password == "" {
		c.String(http.StatusBadRequest, "email and password required")
		return
	}
	ph, err := createPasswordHash(in.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, "hash error")
		return
	}
	u := User{Email: in.Email, PasswordHash: ph, DisplayName: in.DisplayName, CreatedAt: time.Now().UTC()}
	if err := db.Create(&u).Error; err != nil {
		c.String(http.StatusInternalServerError, "db error")
		return
	}
	u.PasswordHash = ""
	c.JSON(http.StatusCreated, u)
}

// loginUser: POST /api/login
func loginUser(c *gin.Context) {
	enableCorsHeaders(c)
	if c.Request.Method == http.MethodOptions {
		c.Status(http.StatusOK)
		return
	}
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&in); err != nil {
		c.String(http.StatusBadRequest, "invalid json")
		return
	}
	var u User
	if err := db.First(&u, "email = ?", in.Email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.String(http.StatusUnauthorized, "invalid credentials")
			return
		}
		c.String(http.StatusInternalServerError, "db error")
		return
	}
	if err := comparePassword(u.PasswordHash, in.Password); err != nil {
		c.String(http.StatusUnauthorized, "invalid credentials")
		return
	}
	token, err := issueToken(u.ID)
	if err != nil {
		c.String(http.StatusInternalServerError, "token error")
		return
	}
	u.PasswordHash = ""
	c.JSON(http.StatusOK, gin.H{"token": token, "user": u})
}

// getMe: GET /api/me
func getMe(c *gin.Context) {
	enableCorsHeaders(c)
	if c.Request.Method == http.MethodOptions {
		c.Status(http.StatusOK)
		return
	}
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.String(http.StatusUnauthorized, "missing authorization")
		return
	}
	var tokenStr string
	fmt.Sscanf(auth, "Bearer %s", &tokenStr)
	if tokenStr == "" {
		c.String(http.StatusUnauthorized, "invalid authorization header")
		return
	}
	uid, err := parseToken(tokenStr)
	if err != nil {
		c.String(http.StatusUnauthorized, "invalid token")
		return
	}
	var u User
	if err := db.First(&u, uid).Error; err != nil {
		c.String(http.StatusInternalServerError, "db error")
		return
	}
	u.PasswordHash = ""
	c.JSON(http.StatusOK, u)
}

func health(c *gin.Context) {
	enableCorsHeaders(c)
	c.String(http.StatusOK, "ok")
}

// getQuiz returns a single quiz by ID: GET /api/quizzes/{id}
func getQuiz(c *gin.Context) {
	enableCorsHeaders(c)
	if c.Request.Method == http.MethodOptions {
		c.Status(http.StatusOK)
		return
	}
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, "missing id")
		return
	}
	var q Quiz
	if err := db.Preload("Questions.Options").First(&q, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.String(http.StatusNotFound, "not found")
			return
		}
		c.String(http.StatusInternalServerError, "db error")
		return
	}
	c.JSON(http.StatusOK, q)
}

// submitQuiz: POST /api/quizzes/:id/submit
func submitQuiz(c *gin.Context) {
	enableCorsHeaders(c)
	if c.Request.Method == http.MethodOptions {
		c.Status(http.StatusOK)
		return
	}
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, "missing id")
		return
	}
	var in struct {
		Answers []struct {
			QuestionID uint  `json:"questionId"`
			Selected   []int `json:"selected"`
		} `json:"answers"`
	}
	if err := c.BindJSON(&in); err != nil {
		c.String(http.StatusBadRequest, "invalid json")
		return
	}
	// optional auth
	var userID *uint
	auth := c.GetHeader("Authorization")
	if auth != "" {
		var tokenStr string
		fmt.Sscanf(auth, "Bearer %s", &tokenStr)
		if tokenStr != "" {
			if uid, err := parseToken(tokenStr); err == nil {
				userID = &uid
			}
		}
	}

	// create submission
	s := Submission{QuizID: id, UserID: userID, CreatedAt: time.Now().UTC()}
	if err := db.Create(&s).Error; err != nil {
		c.String(http.StatusInternalServerError, "db error")
		return
	}

	// load quiz questions/answers
	var quiz Quiz
	if err := db.Preload("Questions.Options").First(&quiz, "id = ?", id).Error; err != nil {
		c.String(http.StatusInternalServerError, "db error")
		return
	}
	// build correct answers map
	correctMap := map[uint][]int{}
	for _, q := range quiz.Questions {
		var ans []int
		if len(q.Answers) > 0 {
			_ = json.Unmarshal(q.Answers, &ans)
		}
		correctMap[q.ID] = ans
	}

	score := 0
	for _, a := range in.Answers {
		selB, _ := json.Marshal(a.Selected)
		resp := Response{
			QuizID:       id,
			QuestionID:   a.QuestionID,
			SubmissionID: &s.ID,
			UserID:       userID,
			Selected:     datatypes.JSON(selB),
			CreatedAt:    time.Now().UTC(),
		}
		// compute correctness (unordered equality)
		correct := correctMap[a.QuestionID]
		isCorrect := false
		if len(correct) == len(a.Selected) {
			m := make(map[int]bool, len(a.Selected))
			for _, v := range a.Selected {
				m[v] = true
			}
			ok := true
			for _, v := range correct {
				if !m[v] {
					ok = false
					break
				}
			}
			if ok {
				isCorrect = true
			}
		}
		resp.IsCorrect = &isCorrect
		if err := db.Create(&resp).Error; err != nil {
			// continue storing other responses
		}
		if isCorrect {
			score++
		}
	}

	s.Score = &score
	if err := db.Save(&s).Error; err != nil {
		c.String(http.StatusInternalServerError, "db error")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"submissionId": s.ID, "quizId": s.QuizID, "score": score, "createdAt": s.CreatedAt})
}

// getResults: GET /api/quizzes/:id/results/:submissionId
func getResults(c *gin.Context) {
	enableCorsHeaders(c)
	if c.Request.Method == http.MethodOptions {
		c.Status(http.StatusOK)
		return
	}
	id := c.Param("id")
	sidStr := c.Param("submissionId")
	if id == "" || sidStr == "" {
		c.String(http.StatusBadRequest, "missing params")
		return
	}
	sid, err := strconv.ParseUint(sidStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid submission id")
		return
	}
	var s Submission
	if err := db.First(&s, sid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.String(http.StatusNotFound, "not found")
			return
		}
		c.String(http.StatusInternalServerError, "db error")
		return
	}
	var resps []Response
	if err := db.Where("submission_id = ?", sid).Find(&resps).Error; err != nil {
		c.String(http.StatusInternalServerError, "db error")
		return
	}
	// load quiz questions to retrieve correct answers
	var quiz Quiz
	if err := db.Preload("Questions.Options").First(&quiz, "id = ?", id).Error; err != nil {
		// ignore — still return responses
	}
	correctMap := map[uint][]int{}
	for _, q := range quiz.Questions {
		var ans []int
		if len(q.Answers) > 0 {
			_ = json.Unmarshal(q.Answers, &ans)
		}
		correctMap[q.ID] = ans
	}

	type perQ struct {
		QuestionID uint  `json:"questionId"`
		Selected   []int `json:"selected"`
		Correct    []int `json:"correct"`
		IsCorrect  bool  `json:"correctBool"`
	}
	pq := make([]perQ, 0, len(resps))
	for _, r := range resps {
		var sel []int
		if len(r.Selected) > 0 {
			_ = json.Unmarshal(r.Selected, &sel)
		}
		corr := correctMap[r.QuestionID]
		isCorrect := false
		if len(corr) == len(sel) {
			m := make(map[int]bool, len(sel))
			for _, v := range sel {
				m[v] = true
			}
			ok := true
			for _, v := range corr {
				if !m[v] {
					ok = false
					break
				}
			}
			if ok {
				isCorrect = true
			}
		}
		pq = append(pq, perQ{QuestionID: r.QuestionID, Selected: sel, Correct: corr, IsCorrect: isCorrect})
	}

	c.JSON(http.StatusOK, gin.H{"submissionId": s.ID, "quizId": s.QuizID, "score": s.Score, "perQuestion": pq})
}

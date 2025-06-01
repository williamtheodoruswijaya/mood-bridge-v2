package repository

import (
	"context"
	"database/sql"
	"fmt"
	"mood-bridge-v2/server/internal/entity"
)

type FriendRepository interface {
	AddFriend(ctx context.Context, tx *sql.Tx, friend *entity.Friend) (*entity.Friend, error)
	AcceptRequest(ctx context.Context, tx *sql.Tx, friend *entity.Friend) (*entity.Friend, error)
	GetFriends(ctx context.Context, db *sql.DB, userID int) (*[]entity.Friend, error)
	Delete(ctx context.Context, tx *sql.Tx, friendID int) (string, error)
	IsFriendExist(ctx context.Context, db *sql.DB, userID int, friendUserID int) (bool, error)
	IsFriendAlreadyAccepted(ctx context.Context, db *sql.DB, userID int, friendUserID int) (bool, error)
	GetFriendRequests(ctx context.Context, db *sql.DB, userID int) (*[]entity.Friend, error)
	GetFriendRecommendation(ctx context.Context, db *sql.DB, userID int) (*[]entity.FriendRecommendation, error)
}

type FriendRepositoryImpl struct {
}

func NewFriendRepository() FriendRepository {
	return &FriendRepositoryImpl{}
}

func (r *FriendRepositoryImpl) AddFriend(ctx context.Context, tx *sql.Tx, friend *entity.Friend) (*entity.Friend, error) {
	// step 1: define query-nya
	query := `INSERT INTO friends (userid, frienduserid, friendstatus, createdat) VALUES ($1, $2, $3, $4) RETURNING friendid, userid, frienduserid, friendstatus, createdat`

	// step 2: jalankan query-nya
	row := tx.QueryRowContext(ctx, query, friend.UserID, friend.FriendUserID, friend.FriendStatus, friend.CreatedAt)

	// step 3: ambil hasilnya
	var newFriend entity.Friend
	if err := row.Scan(&newFriend.FriendID, &newFriend.UserID, &newFriend.FriendUserID, &newFriend.FriendStatus, &newFriend.CreatedAt); err != nil {
		return nil, err
	}

	// step 3.1: fetch the user being added
	userQuery := `SELECT userid, username, fullname FROM users WHERE userid = $1`
	
	userRow := tx.QueryRowContext(ctx, userQuery, friend.FriendUserID)
	var user entity.User
	if err := userRow.Scan(&user.ID, &user.Username, &user.Fullname); err != nil {
		return nil, err
	}

	newFriend.User = &user // set user ke friend

	// step 4: return hasilnya
	return &newFriend, nil
}

func (r *FriendRepositoryImpl) AcceptRequest(ctx context.Context, tx *sql.Tx, friend *entity.Friend) (*entity.Friend, error) {
	// step 1: define query-nya
	query := `
		UPDATE friends
		SET friendstatus = TRUE
		WHERE userid = $2 AND frienduserid = $1 AND friendstatus = FALSE
		RETURNING friendid, userid, frienduserid, friendstatus, createdat
	`

	// step 2: jalankan query-nya
	row := tx.QueryRowContext(ctx, query, friend.UserID, friend.FriendUserID)

	// step 3: ambil hasilnya
	var updatedFriend entity.Friend
	if err := row.Scan(&updatedFriend.FriendID, &updatedFriend.UserID, &updatedFriend.FriendUserID, &updatedFriend.FriendStatus, &updatedFriend.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // tidak ada data yang diupdate
		}
		return nil, err
	}

	// step 4: return hasilnya
	return &updatedFriend, nil
}

func (r *FriendRepositoryImpl) GetFriends(ctx context.Context, db *sql.DB, userID int) (*[]entity.Friend, error) {
	// step 1: define query-nya
	query := `
	SELECT 
		f.friendid, f.userid, f.frienduserid, f.friendstatus, f.createdat,
		u.userid, u.username, u.fullname, u.email, u.profileurl, u.createdat
	FROM friends f
	JOIN users u 
		ON (
			(f.userid = $1 AND u.userid = f.frienduserid)
			OR
			(f.frienduserid = $1 AND u.userid = f.userid)
		)
	WHERE (f.userid = $1 OR f.frienduserid = $1) AND f.friendstatus = true
	`

	// step 2: jalankan query-nya
	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	// step 3: close rows setelah selesai
	defer rows.Close()

	// step 3: ambil hasilnya
	var friends []entity.Friend
	for rows.Next() {
		var friend entity.Friend
		var user entity.User
		if err := rows.Scan(&friend.FriendID, 
							&friend.UserID, 
							&friend.FriendUserID, 
							&friend.FriendStatus, 
							&friend.CreatedAt, 
							&user.ID, 
							&user.Username, 
							&user.Fullname, 
							&user.Email, 
							&user.ProfileUrl,
							&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		// set user ke friend
		friend.User = &user

		// tambahin ke slices
		friends = append(friends, friend)
	}

	// step 4: return hasilnya
	return &friends, nil
}

func (r *FriendRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, friendID int) (string, error) {
	// step 1: define query-nya
	query := `DELETE FROM friends WHERE friendid = $1 RETURNING friendid`

	// step 2: jalankan query-nya
	row := tx.QueryRowContext(ctx, query, friendID)

	// step 3: ambil hasilnya
	var deletedFriendID int
	if err := row.Scan(&deletedFriendID); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("friend with id %d not found", friendID)
		}
		return "", err
	}

	// step 4: return hasilnya
	return "Friend deleted successfully", nil
}

func (r *FriendRepositoryImpl) IsFriendExist(ctx context.Context, db *sql.DB, userID int, friendUserID int) (bool, error) {
	// step 1: define query untuk cek apakah dia udah berteman atau belum
	query := `SELECT 1 FROM friends WHERE ((userid = $1 AND frienduserid = $2) OR (userid = $2 AND frienduserid = $1)) LIMIT 1`

	// step 2: jalankan query-nya
	row := db.QueryRowContext(ctx, query, userID, friendUserID)

	var temp int // buat ambil hasil dari row (1 berarti udah berteman, 0 berarti belum)
	if err := row.Scan(&temp); err == sql.ErrNoRows {
		return false, nil // gaada data yang ditemukan, berarti belum berteman
	} else if err != nil {
		return false, err // ada error lain
	}

	// step 3: kalau ada hasilnya, berarti udah berteman
	return true, nil
}

func (r *FriendRepositoryImpl) IsFriendAlreadyAccepted(ctx context.Context, db *sql.DB, userID int, friendUserID int) (bool, error) {
	query := `SELECT 1 FROM friends
	WHERE ((userid = $1 AND frienduserid = $2) OR (userid = $2 AND frienduserid = $1))
	AND friendstatus = true
	LIMIT 1`
	
	row := db.QueryRowContext(ctx, query, userID, friendUserID)
	var temp int
	if err := row.Scan(&temp); err == sql.ErrNoRows {
		return false, nil // gaada data yang ditemukan, berarti belum berteman
	} else if err != nil {
		return false, err // ada error lain
	}
	return true, nil
}

func (r *FriendRepositoryImpl) GetFriendRequests(ctx context.Context, db *sql.DB, userID int) (*[]entity.Friend, error) {
	query := `
	SELECT
		f.friendid,
		f.userid,
		f.frienduserid,
		u.username,
		u.fullname,
		f.createdat
	FROM 
		friends f
	JOIN
		users u ON f.userid = u.userid
	WHERE
		f.frienduserid = $1
		AND f.friendstatus = false
	ORDER BY f.createdat DESC
	`

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var friendRequests []entity.Friend
	for rows.Next() {
		var friend entity.Friend
		var user entity.User
		if err := rows.Scan(&friend.FriendID, &friend.UserID, &friend.FriendUserID, &user.Username, &user.Fullname, &friend.CreatedAt); err != nil {
			return nil, err
		}
		friend.User = &user
		friendRequests = append(friendRequests, friend)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &friendRequests, nil
}

func (r *FriendRepositoryImpl) GetFriendRecommendation(ctx context.Context, db *sql.DB, userID int) (*[]entity.FriendRecommendation, error) {
	// step 1: hitung overall mood dari user
	query := `
	SELECT mood FROM (
		SELECT mood,
				COUNT(*) AS freq,
				MAX(createdat) AS latest_post,
				ROW_NUMBER() OVER (
				ORDER BY COUNT(*) DESC, MAX(createdat) DESC
				) AS rn
		FROM posts
		WHERE userid = $1
		GROUP BY mood
	) AS ranked
	WHERE rn = 1;
	`

	row := db.QueryRowContext(ctx, query, userID)
	if row == nil {
		return nil, fmt.Errorf("no mood found for user with id %d", userID)
	}

	var overallMood string
	if err := row.Scan(&overallMood); err != nil {
		return nil, fmt.Errorf("failed to get overall mood for user with id %d: %v", userID, err)
	}

	// step 2: Klasifikasikan user ke dalam user dengan risiko mental yang buruk dan user dengan risiko mental yang baik
	var negativeMoods = map[string]bool{
  		"Depression": true,
  		"Anxiety": true,
  		"Suicidal": true,
  		"Personality disorder": true,
  		"Bipolar": true,
  		"Stress": true,
	}

	isAtRisk := negativeMoods[overallMood]

	// step 3: buat query buat ngambil rekomendasi teman berdasarkan overall mood user
	var moodCondition string
	if isAtRisk {
		// Rekomendasikan user dengan mood positif (mood tidak negatif)
		moodCondition = "NOT IN ('Depression', 'Anxiety', 'Suicidal', 'Personality disorder', 'Bipolar', 'Stress')"
	} else {
		// Rekomendasikan user dengan mood negatif
		moodCondition = "IN ('Depression', 'Anxiety', 'Suicidal', 'Personality disorder', 'Bipolar', 'Stress')"
	}

	recommendationQuery := fmt.Sprintf(`
		SELECT u.userid, u.username, u.fullname, u.email, ranked.mood
		FROM (
			SELECT userid, mood,
				ROW_NUMBER() OVER (
					PARTITION BY userid
					ORDER BY COUNT(*) DESC, MAX(createdat) DESC
				) AS rn
			FROM posts
			GROUP BY userid, mood
		) AS ranked
		JOIN users u ON u.userid = ranked.userid
		WHERE rn = 1
		AND mood %s
		AND u.userid != $1
		AND u.userid NOT IN (
				SELECT frienduserid FROM friends WHERE userid = $1 AND friendstatus = TRUE
				UNION
				SELECT userid FROM friends WHERE frienduserid = $1 AND friendstatus = TRUE
		)
		LIMIT 10;
	`, moodCondition)


	// step 4: jalankan query-nya
	rows, err := db.QueryContext(ctx, recommendationQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get friend recommendations: %v", err)
	}
	defer rows.Close()

	var recommendations []entity.FriendRecommendation
	for rows.Next() {
		var user entity.User
		var mood string
		if err := rows.Scan(&user.ID, &user.Username, &user.Fullname, &user.Email, &mood); err != nil {
			return nil, fmt.Errorf("failed to scan recommendation row: %v", err)
		}
		recommendations = append(recommendations, entity.FriendRecommendation{
			User:        user,
			OverallMood: mood,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over recommendations: %v", err)
	}

	// step 5: return hasilnya
	return &recommendations, nil
}
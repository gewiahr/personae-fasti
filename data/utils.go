package data

// type Claims struct {
// 	UserID   int64  `json:"user_id"`
// 	Username string `json:"username"`
// 	jwt.RegisteredClaims
// }

// func (s *Storage) GeneratePlayerJWTToken(player *Player, jwtSecret string, expTime time.Time) (string, error) {
// 	claims := &Claims{
// 		UserID:   int64(player.ID),
// 		Username: player.Username,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(expTime),
// 			IssuedAt:  jwt.NewNumericDate(time.Now()),
// 			Subject:   player.Username,
// 		},
// 	}

// 	var JWTBytes = []byte(jwtSecret)
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(JWTBytes)
// }

// func (s *Storage) GetAllowedNPCs(npcs []NPC, playerID int) ([]NPC, error) {
// 	err := s.db.NewSelect().Model(&npcs).WherePK().
// 		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
// 			return q.Where("hidden_by = 0").WhereOr("hidden_by = ?", playerID)
// 		}).
// 		Scan(context.Background(), &npcs)
// 	if err != nil {
// 		return nil, err
// 	} else if err == sql.ErrNoRows {
// 		return []NPC{}, nil
// 	}

// 	return npcs, nil
// }

// func (s *Storage) GetAllowedLocations(locations []Location, playerID int) ([]Location, error) {
// 	err := s.db.NewSelect().Model(&locations).WherePK().
// 		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
// 			return q.Where("hidden_by = 0").WhereOr("hidden_by = ?", playerID)
// 		}).
// 		Scan(context.Background(), &locations)
// 	if err != nil {
// 		return nil, err
// 	} else if err == sql.ErrNoRows {
// 		return []Location{}, nil
// 	}

// 	return locations, nil
// }

// func (s *Storage) GetAllowedQuests(quests []Quest, playerID int) ([]Quest, error) {
// 	err := s.db.NewSelect().Model(&quests).WherePK().
// 		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
// 			return q.Where("hidden_by = 0").WhereOr("hidden_by = ?", playerID)
// 		}).
// 		Scan(context.Background(), &quests)
// 	if err != nil {
// 		return nil, err
// 	} else if err == sql.ErrNoRows {
// 		return []Quest{}, nil
// 	}

// 	return quests, nil
// }

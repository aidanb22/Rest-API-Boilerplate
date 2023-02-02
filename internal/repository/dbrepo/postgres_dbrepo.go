package dbrepo

import (
	"context"
	"database/sql"
	"github.com/aidanb22/internal/models"
	"time"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostgresDBRepo) AllPlayers() ([]*models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	select
		id, plyname, college, age, height, description,
		coalesce(image, ''),
	    created_at, updated_at
	from
		players
	order by
		plyname
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []*models.Player

	for rows.Next() {
		var player models.Player
		err := rows.Scan(
			&player.ID,
			&player.Plyname,
			&player.College,
			&player.Age,
			&player.Height,
			&player.Description,
			&player.Image,
			&player.CreatedAt,
			&player.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		players = append(players, &player)
	}

	return players, nil
}

func (m *PostgresDBRepo) OnePlayer(id int) (*models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	select
		id, plyname, college, age, height, description,
		coalesce(image, ''),
	    created_at, updated_at
	from
		players
	where
		id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var player models.Player

	err := row.Scan(
		&player.ID,
		&player.Plyname,
		&player.College,
		&player.Age,
		&player.Height,
		&player.Description,
		&player.Image,
		&player.CreatedAt,
		&player.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	//Todo: get POSITIONS if any
	//var positions []*models.Position
	return &player, err
}

func (m *PostgresDBRepo) OnePlayerForEdit(id int) (*models.Player /*[]*models.Position*/, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	select
		id, plyname, college, age, height, description,
		coalesce(image, ''),
	    created_at, updated_at
	from
		players
	where
		id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var player models.Player

	err := row.Scan(
		&player.ID,
		&player.Plyname,
		&player.College,
		&player.Age,
		&player.Height,
		&player.Description,
		&player.Image,
		&player.CreatedAt,
		&player.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	//Todo: get POSITIONS if any
	//var positions []*models.Position
	/*
		var allPositions []*models.Position
		query = `select p.id, position from position order by p.id  `
		pRows, err := m.DB.QueryContext(ctx, query)
		if err != nil {
			return nil, nil, err
		}
		defer pRows.Close()
		for pRows.Next() {
			var p models.Position
			err := pRows.Scan(
				&p.ID,
				&p.Position,
			)
			if err != nil {
				return nil, nil, err
			}

			allPositions = append(allPositions, &p)
		}
	*/

	return &player, err
}

func (m *PostgresDBRepo) AllUsers() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	select
		id, first_name, last_name, email, password,
	    created_at, updated_at
	from
		users
	order by
		first_name
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (m *PostgresDBRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password,
				created_at, updated_at from users where email = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *PostgresDBRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password,
				created_at, updated_at from users where id = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *PostgresDBRepo) InsertPlayer(player models.Player) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into players (plyname, college, age, height, description,
	created_at, updated_at, image)
	values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		player.Plyname,
		player.College,
		player.Age,
		player.Height,
		player.Description,
		player.CreatedAt,
		player.UpdatedAt,
		player.Image,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}
	return newID, nil

}
func (m *PostgresDBRepo) UpdatePlayer(player models.Player) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update players set plyname = $1, college = $2, age = $3, height = $4,
			description = $5, updated_at = $6, image = $7 where id = $8`

	_, err := m.DB.ExecContext(ctx, stmt,
		player.Plyname,
		player.College,
		player.Age,
		player.Height,
		player.Description,
		player.UpdatedAt,
		player.Image,
		player.ID,
	)
	if err != nil {
		return err
	}
	return nil

}

func (m *PostgresDBRepo) DeletePlayer(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from players where id =$1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	return nil

}

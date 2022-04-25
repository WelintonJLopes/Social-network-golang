package repositories

import (
	"api/src/models"
	"database/sql"
)

type Publications struct {
	db *sql.DB
}

func NewRepositoryPublications(db *sql.DB) *Publications {
	return &Publications{db}
}

func (p Publications) Create(publication models.Publication) (uint64, error) {
	statement, err := p.db.Prepare("insert into publications (title, content, author_id) values (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(publication.Title, publication.Content, publication.AuthorID)
	if err != nil {
		return 0, err
	}

	lastIDInserted, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastIDInserted), nil
}

func (p Publications) SearchID(publicationID uint64) (models.Publication, error) {
	rows, err := p.db.Query(`
		select p.*, u.nick from
		publications p inner join users u
		on u.id = p.author_id where p.id = ?
	`, publicationID)
	if err != nil {
		return models.Publication{}, err
	}
	defer rows.Close()

	var publication models.Publication

	if rows.Next() {
		if err = rows.Scan(
			&publication.ID, &publication.Title, &publication.Content, &publication.AuthorID, &publication.Likes,
			&publication.CreatedAt, &publication.AuthorNick,
		); err != nil {
			return models.Publication{}, err
		}
	}

	return publication, nil
}

func (p Publications) Search(userID uint64) ([]models.Publication, error) {
	rows, err := p.db.Query(`
		select distinct p.*, u.nick from
		publications p 
		inner join users u on u.id = p.author_id 
		inner join followers s on p.author_id = s.user_id
		where p.id = ? or s.follower_id = ?
	`, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var publications []models.Publication

	if rows.Next() {
		var publication models.Publication
		if err = rows.Scan(
			&publication.ID, &publication.Title, &publication.Content, &publication.AuthorID, &publication.Likes,
			&publication.CreatedAt, &publication.AuthorNick,
		); err != nil {
			return nil, err
		}
		publications = append(publications, publication)
	}

	return publications, nil
}

func (p Publications) Update(publicationID uint64, publication models.Publication) error {
	statement, err := p.db.Prepare("update publications set title = ?, content = ? where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(publication.Title, publication.Content, publicationID); err != nil {
		return err
	}

	return nil
}

func (p Publications) Delete(publicationID uint64) error {
	statement, err := p.db.Prepare("delete from publications where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(publicationID); err != nil {
		return err
	}

	return nil
}

func (p Publications) SearchUser(userID uint64) ([]models.Publication, error) {
	rows, err := p.db.Query(`
		select p.*, u.nick from publications p 
		inner join users u on u.id = p.author_id 
		where p.author_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var publications []models.Publication

	if rows.Next() {
		var publication models.Publication
		if err = rows.Scan(
			&publication.ID, &publication.Title, &publication.Content, &publication.AuthorID, &publication.Likes,
			&publication.CreatedAt, &publication.AuthorNick,
		); err != nil {
			return nil, err
		}
		publications = append(publications, publication)
	}

	return publications, nil
}

func (p Publications) LikePublication(publicationID uint64) error {
	statement, err := p.db.Prepare("update publications set likes = likes + 1 where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(publicationID); err != nil {
		return err
	}

	return nil
}

func (p Publications) DislikePublication(publicationID uint64) error {
	statement, err := p.db.Prepare(`
		update publications set likes = 
		case 
			when likes > 0 
			then likes - 1 
		else 
			0 
		end
		where id = ?
	`)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(publicationID); err != nil {
		return err
	}

	return nil
}

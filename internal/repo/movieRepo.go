package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/fouched/go-movies-api/internal/models"
)

func AllMovies() ([]*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select 
			id, title, release_date, runtime, mpaa_rating, description,
			coalesce(image, ''), created_at, updated_at
		from
		    movies
		order by
		    title		
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*models.Movie
	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.RunTime,
			&movie.MPAARating,
			&movie.Description,
			&movie.Image,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}

	return movies, nil
}

func GetMovieByID(id int) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select 
		    id, title, release_date, runtime, mpaa_rating, description,
			coalesce(image, ''), created_at, updated_at
		from
		    movies
		where 
		    id = $1
	`
	row := db.QueryRowContext(ctx, query, id)

	var movie models.Movie
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.RunTime,
		&movie.MPAARating,
		&movie.Description,
		&movie.Image,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// get genres, if any
	query = `
		select 
		    g.id, g.genre
		from 
		    movies_genres mg
		left join genres g on mg.genre_id = g.id
		where 
		    mg.movie_id = $1
		order by
		    g.genre
	`
	rows, err := db.QueryContext(ctx, query, movie.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	defer rows.Close()

	var genres []*models.Genre
	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &g)
	}

	movie.Genres = genres

	return &movie, nil
}

func GetMovieByIDForEdit(id int) (*models.Movie, []*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select 
		    id, title, release_date, runtime, mpaa_rating, description,
			coalesce(image, ''), created_at, updated_at
		from
		    movies
		where 
		    id = $1
	`
	row := db.QueryRowContext(ctx, query, id)

	var movie models.Movie
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.RunTime,
		&movie.MPAARating,
		&movie.Description,
		&movie.Image,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)
	if err != nil {
		return nil, nil, err
	}

	// get genres, if any
	query = `
		select 
		    g.id, g.genre
		from 
		    movies_genres mg
		left join genres g on mg.genre_id = g.id
		where 
		    mg.movie_id = $1
		order by
		    g.genre
	`
	rows, err := db.QueryContext(ctx, query, movie.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, err
	}
	defer rows.Close()

	var genres []*models.Genre
	var genresArray []int
	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, nil, err
		}
		genres = append(genres, &g)
		genresArray = append(genresArray, g.ID)
	}

	movie.Genres = genres
	movie.GenresArray = genresArray

	query = "select id, genre from genres order by genre"
	gRows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	defer gRows.Close()

	var allGenres []*models.Genre
	for gRows.Next() {
		var g models.Genre
		err := gRows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, nil, err
		}
		allGenres = append(allGenres, &g)
	}

	return &movie, allGenres, nil
}

func GetAllGenres() ([]*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, genre, created_at, updated_at from genres order by genre`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []*models.Genre
	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
			&g.CreatedAt,
			&g.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &g)
	}

	return genres, nil
}

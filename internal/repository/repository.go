package repository

import (
	"context"
	"database/sql"
	"errors"
	"net/netip"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"go.uber.org/zap"

	"github.com/Killer-Feature/PaaS_ClientSide/internal"
)

func Create(l *zap.Logger) (*Repository, error) {
	if _, err := os.Stat("./internal_data.db"); errors.Is(err, os.ErrNotExist) {
		l.Debug("Creating new sql database")
		file, err := os.Create("internal_data.db") // Create SQLite file
		if err != nil {
			l.Error("error occurred during db file creating", zap.Error(err))
		}
		err = file.Close()
		if err != nil {
			l.Error("error occurred during db file closing", zap.Error(err))
		}
		l.Debug("internal_data.db created")
	}

	db, _ := sql.Open("sqlite3", "./internal_data.db")

	createNodesTableSQL := `CREATE TABLE IF NOT EXISTS nodes (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT,
		"ip" TEXT,	
		"login" TEXT,
		"password" TEXT
	  );`

	statement, err := db.Prepare(createNodesTableSQL)
	if err != nil {
		l.Error("error occurred during preparing table creating statement", zap.Error(err))
	}
	_, err = statement.Exec()
	if err != nil {
		l.Error("error occurred during execution table creating statement", zap.Error(err))
		return nil, err
	}
	l.Debug("repository created")

	return &Repository{
		db: db,
		l:  l,
	}, err
}

func (r *Repository) GetNodes(ctx context.Context) ([]internal.FullNode, error) {
	sqlScript := "SELECT id, name, ip, login, password FROM nodes"

	rows, err := r.db.QueryContext(ctx, sqlScript)
	if err != nil {
		r.l.Error("error in db query during getting nodes", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	var selectedNodes []internal.FullNode
	for rows.Next() {
		var singleNode internal.FullNode
		var ip string
		if err = rows.Scan(&singleNode.ID, &singleNode.Name, &ip, &singleNode.Login, &singleNode.Password); err != nil {
			r.l.Error("error during scanning node from database", zap.Error(err))
			return nil, err
		}
		singleNode.IP, err = netip.ParseAddrPort(ip)
		if err != nil {
			r.l.Error("error during parsing ip from database", zap.Error(err))
		}
		selectedNodes = append(selectedNodes, singleNode)
	}

	return selectedNodes, nil
}

type Repository struct {
	db *sql.DB
	l  *zap.Logger
}

func (r *Repository) AddNode(ctx context.Context, node internal.FullNode) (int, error) {
	sqlScript := "INSERT INTO nodes(name, ip, login, password) VALUES ($1, $2, $3, $4) RETURNING id;"
	err := r.db.QueryRowContext(ctx, sqlScript, node.Name, node.IP.String(), node.Login, node.Password).Scan(&node.ID)
	if err != nil {
		r.l.Error("error during adding node to database", zap.Error(err))
		return 0, err
	}
	return node.ID, nil
}

func (r *Repository) RemoveNode(ctx context.Context, id int) error {
	sqlScript := "DELETE FROM nodes WHERE id=$1;"
	_, err := r.db.ExecContext(ctx, sqlScript, id)
	return err
}

func (r *Repository) IsNodeExists(ctx context.Context, ip netip.AddrPort) (bool, error) {
	sqlScript := "SELECT id FROM nodes WHERE ip=$1"
	rows, err := r.db.QueryContext(ctx, sqlScript, ip.String())
	if err != nil {
		return false, err
	}
	if !rows.Next() {
		return false, nil
	}
	_ = rows.Close()
	return true, nil
}

func (r *Repository) Close() error {
	if r.db != nil {
		err := r.db.Close()
		if err != nil {
			r.l.Error("error during closing database", zap.Error(err))
			return err
		}
	}
	return nil
}

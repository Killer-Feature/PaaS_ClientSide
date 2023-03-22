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

func Create(l *zap.Logger) (internal.Repository, error) {
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

	db, err := sql.Open("sqlite3", "./internal_data.db")
	if err != nil {
		l.Error("error occurred during db opening", zap.Error(err))
		return nil, err
	}

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

	createClustersTableSQL := `CREATE TABLE IF NOT EXISTS clusters (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT UNIQUE,
		"master_ip" TEXT,
		"token" TEXT,	
		"hash" TEXT
	  );`

	statement, err = db.Prepare(createClustersTableSQL)
	if err != nil {
		l.Error("error occurred during preparing cluster table creating statement", zap.Error(err))
	}
	_, err = statement.Exec()
	if err != nil {
		l.Error("error occurred during execution cluster table creating statement", zap.Error(err))
		return nil, err
	}
	l.Debug("repository created")

	r := &Repository{
		db: db,
		l:  l,
	}
	_, _ = r.AddCluster(context.Background(), "defaultCluster")
	return r, nil
}

func (r *Repository) GetNodes(ctx context.Context) ([]internal.FullNode, error) {
	sqlScript := "SELECT id, name, ip, login, password FROM nodes;"

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

func (r *Repository) GetFullNode(ctx context.Context, id int) (internal.FullNode, error) {
	sqlScript := "SELECT id, name, ip, login, password FROM nodes WHERE id = $1"

	var singleNode internal.FullNode
	var ip string
	err := r.db.QueryRowContext(ctx, sqlScript, id).Scan(&singleNode.ID, &singleNode.Name, &ip, &singleNode.Login, &singleNode.Password)
	if err != nil {
		r.l.Error("error in db query during getting nodes", zap.Error(err))
		return internal.FullNode{}, err
	}
	singleNode.IP, err = netip.ParseAddrPort(ip)
	if err != nil {
		r.l.Error("error during parsing ip from database", zap.Error(err))
		return internal.FullNode{}, err
	}
	return singleNode, nil
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

func (r *Repository) AddCluster(ctx context.Context, clusterName string) (int, error) {
	sqlScript := "INSERT INTO clusters(name) VALUES ($1) RETURNING id;"
	var id int
	err := r.db.QueryRowContext(ctx, sqlScript, clusterName).Scan(&id)
	if err != nil {
		r.l.Error("error during adding cluster to database", zap.Error(err))
		return 0, err
	}
	return id, nil
}

func (r *Repository) GetClusterID(ctx context.Context, clusterName string) (int, error) {
	sqlScript := "SELECT id FROM nodes WHERE name = $1;"
	var id int
	err := r.db.QueryRowContext(ctx, sqlScript, clusterName).Scan(&id)
	if err != nil {
		r.l.Error("error during adding cluster to database", zap.Error(err))
		return 0, err
	}
	return id, nil
}

func (r *Repository) GetClusterName(ctx context.Context, clusterName string) (int, error) {
	sqlScript := "SELECT id FROM nodes WHERE name = $1;"
	var id int
	err := r.db.QueryRowContext(ctx, sqlScript, clusterName).Scan(&id)
	if err != nil {
		r.l.Error("error during adding cluster to database", zap.Error(err))
		return 0, err
	}
	return id, nil
}

//type Cluster struct {
//	ID     int
//	Name   string
//	Config string
//	Token  string
//	Hash   string
//}
//
//func (r *Repository) GetClusters(ctx context.Context) ([]int, []string, error) {
//	sqlScript := "SELECT id, name FROM clusters;"
//
//	rows, err := r.db.QueryContext(ctx, sqlScript)
//	if err != nil {
//		r.l.Error("error in db query during getting nodes", zap.Error(err))
//		return nil, nil, err
//	}
//	defer rows.Close()
//
//	var ids []int
//	var names []string
//
//	for rows.Next() {
//		var singleNode internal.FullNode
//		var ip string
//		if err = rows.Scan(&singleNode.ID, &singleNode.Name, &ip, &singleNode.Login, &singleNode.Password); err != nil {
//			r.l.Error("error during scanning node from database", zap.Error(err))
//			return nil, err
//		}
//		singleNode.IP, err = netip.ParseAddrPort(ip)
//		if err != nil {
//			r.l.Error("error during parsing ip from database", zap.Error(err))
//		}
//		selectedNodes = append(selectedNodes, singleNode)
//	}
//
//	return selectedNodes, nil
//}

func (r *Repository) AddClusterTokenIPAndHash(ctx context.Context, clusterID int, token, masterIP, hash string) error {
	sqlScript := "UPDATE clusters SET token = $1, hash = $2, master_ip=$3 WHERE id = $4"
	_, err := r.db.ExecContext(ctx, sqlScript, token, hash, masterIP, clusterID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) CheckClusterTokenIPAndHash(ctx context.Context, clusterID int) (bool, error) {
	token, masterIP, hash, err := r.GetClusterTokenIPAndHash(ctx, clusterID)
	if err != nil {
		return false, err
	}
	if token == "" && masterIP == "" && hash == "" {
		return false, nil
	}
	return true, nil
}

func (r *Repository) GetClusterTokenIPAndHash(ctx context.Context, clusterID int) (token, masterIP, hash string, err error) {
	var rawToken, rawMasterIP, rawHash sql.NullString
	sqlScript := "SELECT token, master_ip, hash FROM clusters WHERE id = $1"
	err = r.db.QueryRowContext(ctx, sqlScript, clusterID).Scan(&rawToken, &rawMasterIP, &rawHash)
	if rawToken.Valid {
		token = rawToken.String
	}
	if rawMasterIP.Valid {
		masterIP = rawMasterIP.String
	}
	if rawHash.Valid {
		hash = rawHash.String
	}

	return
}

package repository

import (
	"context"
	"database/sql"
	"errors"
	"net/netip"
	"os"
	"strconv"

	"github.com/Killer-Feature/PaaS_ClientSide/internal/models"

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

	createClustersTableSQL := `CREATE TABLE IF NOT EXISTS clusters (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT UNIQUE,
		"master_ip" TEXT,
		"token" TEXT,	
		"hash" TEXT,
		"master_id" integer
	  );`

	statement, err := db.Prepare(createClustersTableSQL)
	if err != nil {
		l.Error("error occurred during preparing cluster table creating statement", zap.Error(err))
	}
	_, err = statement.Exec()
	if err != nil {
		l.Error("error occurred during execution cluster table creating statement", zap.Error(err))
		return nil, err
	}
	createNodesTableSQL := `CREATE TABLE IF NOT EXISTS nodes (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT,
		"ip_port" TEXT,
		"ip" TEXT,
		"login" TEXT,
		"cluster_id" integer,
		"password" TEXT,
		"is_master" BOOLEAN,
		FOREIGN KEY(cluster_id) REFERENCES clusters(id)
	  );`

	statement, err = db.Prepare(createNodesTableSQL)
	if err != nil {
		l.Error("error occurred during preparing table creating statement", zap.Error(err))
	}
	_, err = statement.Exec()
	if err != nil {
		l.Error("error occurred during execution table creating statement", zap.Error(err))
		return nil, err
	}

	createResourcesTableSQL := `CREATE TABLE IF NOT EXISTS resources (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT,
		"type" TEXT
	  );`

	statement, err = db.Prepare(createResourcesTableSQL)
	if err != nil {
		l.Error("error occurred during preparing table creating statement", zap.Error(err))
	}
	_, err = statement.Exec()
	if err != nil {
		l.Error("error occurred during execution table creating statement", zap.Error(err))
		return nil, err
	}

	createAdminTableSQL := `CREATE TABLE IF NOT EXISTS admin (
    	"user" TEXT,
		"password" TEXT
	  );`

	statement, err = db.Prepare(createAdminTableSQL)
	if err != nil {
		l.Error("error occurred during preparing table creating statement", zap.Error(err))
	}
	_, err = statement.Exec()
	if err != nil {
		l.Error("error occurred during execution table creating statement", zap.Error(err))
		return nil, err
	}

	createSessionsTableSQL := `CREATE TABLE IF NOT EXISTS sessions (
		"session" TEXT UNIQUE
	  );`

	statement, err = db.Prepare(createSessionsTableSQL)
	if err != nil {
		l.Error("error occurred during preparing table creating statement", zap.Error(err))
	}
	_, err = statement.Exec()
	if err != nil {
		l.Error("error occurred during execution table creating statement", zap.Error(err))
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
	sqlScript := "SELECT id, name, ip_port, login, password, cluster_id, is_master FROM nodes;"

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
		var clusterId sql.NullString
		var isMaster sql.NullBool
		if err = rows.Scan(&singleNode.ID, &singleNode.Name, &ip, &singleNode.Login, &singleNode.Password, &clusterId, &isMaster); err != nil {
			r.l.Error("error during scanning node from database", zap.Error(err))
			return nil, err
		}
		singleNode.IsMaster = isMaster.Bool
		singleNode.IP, err = netip.ParseAddrPort(ip)
		singleNode.ClusterID, _ = strconv.Atoi(clusterId.String)
		if err != nil {
			r.l.Error("error during parsing ip from database", zap.Error(err))
		}
		selectedNodes = append(selectedNodes, singleNode)
	}

	return selectedNodes, nil
}

func (r *Repository) GetFullNode(ctx context.Context, id int) (internal.FullNode, error) {
	sqlScript := "SELECT id, name, ip_port, login, password, cluster_id, is_master FROM nodes WHERE id = $1"

	var singleNode internal.FullNode
	var ip string
	var clusterId sql.NullString
	var isMaster sql.NullBool
	err := r.db.QueryRowContext(ctx, sqlScript, id).Scan(&singleNode.ID, &singleNode.Name, &ip, &singleNode.Login, &singleNode.Password, &clusterId, &isMaster)
	if err != nil {
		r.l.Error("error in db query during getting nodes", zap.Error(err))
		return internal.FullNode{}, err
	}
	singleNode.IP, err = netip.ParseAddrPort(ip)
	singleNode.IsMaster = isMaster.Bool
	singleNode.ClusterID, _ = strconv.Atoi(clusterId.String)
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
	sqlScript := "INSERT INTO nodes(name, ip_port, login, password, ip) VALUES ($1, $2, $3, $4, $5) RETURNING id;"
	err := r.db.QueryRowContext(ctx, sqlScript, node.Name, node.IP.String(), node.Login, node.Password, node.IP.Addr().String()).Scan(&node.ID)
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

func (r *Repository) SetNodeClusterID(ctx context.Context, id int, clusterID int) error {
	sqlScript := "UPDATE nodes SET cluster_id = $1 WHERE id = $2"
	_, err := r.db.ExecContext(ctx, sqlScript, clusterID, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) ResetNodeCluster(ctx context.Context, id int) error {
	sqlScript := "UPDATE nodes SET cluster_id = 0, is_master=false WHERE id = $1"
	_, err := r.db.ExecContext(ctx, sqlScript, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) IsNodeExists(ctx context.Context, ip netip.Addr) (int, error) {
	sqlScript := "SELECT id FROM nodes WHERE ip=$1"
	rows, err := r.db.QueryContext(ctx, sqlScript, ip.String())
	defer rows.Close()
	if err != nil {
		return 0, err
	}

	var id int
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			r.l.Error("error during scanning node id from database", zap.Error(err))
			return 0, err
		}
	}
	return id, nil
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
	if id, err := r.GetClusterID(ctx, clusterName); err == nil {
		return id, nil
	}
	sqlScript := "INSERT INTO clusters(name) VALUES ($1) RETURNING id;"
	var id int
	err := r.db.QueryRowContext(ctx, sqlScript, clusterName).Scan(&id)
	if err != nil {
		r.l.Error("error during adding cluster to database", zap.Error(err))
		return 0, err
	}
	return id, nil
}

func (r *Repository) AddAdmin(ctx context.Context, user, password string) error {
	_, err := r.db.Exec("DELETE FROM admin")
	if err != nil {
		return err
	}
	sqlScript := "INSERT INTO admin(user,password) VALUES ($1,$2);"
	_, err = r.db.ExecContext(ctx, sqlScript, user, password)
	if err != nil {
		r.l.Error("error during adding user to database", zap.Error(err))
		return err
	}
	return nil
}

func (r *Repository) GetClusterID(ctx context.Context, clusterName string) (int, error) {
	sqlScript := "SELECT id FROM clusters WHERE name = $1;"
	var id int
	err := r.db.QueryRowContext(ctx, sqlScript, clusterName).Scan(&id)
	if err != nil {
		r.l.Error("error during getting cluster id from database", zap.Error(err))
		return 0, err
	}
	return id, nil
}

func (r *Repository) GetClusterName(ctx context.Context, id int) (string, error) {
	sqlScript := "SELECT name FROM clusters WHERE id = $1;"
	var name string
	err := r.db.QueryRowContext(ctx, sqlScript, id).Scan(&name)
	if err != nil {
		r.l.Error("error during getting cluster name from database", zap.Error(err))
		return "", err
	}
	return name, nil
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
	masterID := 0
	sqlScript := "UPDATE nodes SET is_master = $1 WHERE ip = $2 RETURNING id"
	masterAddrPort, _ := netip.ParseAddrPort(masterIP)
	err := r.db.QueryRowContext(ctx, sqlScript, true, masterAddrPort.Addr().String()).Scan(&masterID)
	if err != nil {
		r.l.Error("error during add cluster master to database", zap.Error(err))
		return err
	}

	sqlScript = "UPDATE clusters SET token = $1, hash = $2, master_ip=$3, master_id=$4 WHERE id = $5"
	_, err = r.db.ExecContext(ctx, sqlScript, token, hash, masterIP, masterID, clusterID)
	if err != nil {
		r.l.Error("error during add cluster master to database", zap.Error(err))
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

func (r *Repository) DeleteClusterTokenIPAndHash(ctx context.Context, clusterID int) (err error) {
	sqlScript := `UPDATE clusters SET token = "", hash = "", master_ip="", master_id=0 WHERE id = $1`
	_, err = r.db.ExecContext(ctx, sqlScript, clusterID)
	return
}

func (r *Repository) GetResources(ctx context.Context) ([]models.ResourceData, error) {
	sqlScript := "SELECT name, type FROM resources;"

	rows, err := r.db.QueryContext(ctx, sqlScript)
	if err != nil {
		r.l.Error("error in db query during getting nodes", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var selectedResources []models.ResourceData
	for rows.Next() {
		var singleNode models.ResourceData
		if err = rows.Scan(&singleNode.Name, &singleNode.Type); err != nil {
			r.l.Error("error during scanning node from database", zap.Error(err))
			return nil, err
		}
		selectedResources = append(selectedResources, singleNode)
	}

	return selectedResources, nil
}

func (r *Repository) AddResource(ctx context.Context, rType, name string) error {
	sqlScript := "INSERT INTO resources(type, name) VALUES ($1, $2);"
	_, err := r.db.ExecContext(ctx, sqlScript, rType, name)
	if err != nil {
		r.l.Error("error during adding cluster to database", zap.Error(err))
		return err
	}
	return nil
}

func (r *Repository) ExistSession(ctx context.Context, session string) (bool, error) {
	sqlScript := "SELECT EXISTS(SELECT 1 FROM sessions WHERE session = $1)"
	exist := false
	err := r.db.QueryRowContext(ctx, sqlScript, session).Scan(&exist)
	if err != nil {
		r.l.Error("error during getting session from database", zap.Error(err))
		return false, err
	}
	return exist, nil
}

func (r *Repository) CheckLoginData(ctx context.Context, user, password string) (bool, error) {
	sqlScript := "SELECT EXISTS(SELECT 1 FROM admin WHERE user = $1 AND password = $2)"
	exist := false
	err := r.db.QueryRowContext(ctx, sqlScript, user, password).Scan(&exist)
	if err != nil {
		r.l.Error("error during getting admin from database", zap.Error(err))
		return false, err
	}
	return exist, nil
}

func (r *Repository) AddSession(ctx context.Context, session string) error {
	sqlScript := "INSERT INTO sessions(session) VALUES ($1);"
	_, err := r.db.ExecContext(ctx, sqlScript, session)
	if err != nil {
		r.l.Error("error during adding session to database", zap.Error(err))
		return err
	}
	return nil
}

func (r *Repository) RemoveSession(ctx context.Context, session string) error {
	sqlScript := "DELETE FROM sessions WHERE session=$1;"
	_, err := r.db.ExecContext(ctx, sqlScript, session)
	if err != nil {
		r.l.Error("error during removing session to database", zap.Error(err))
		return err
	}
	return nil
}

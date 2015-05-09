package meta

import (
	"time"

	"github.com/influxdb/influxdb/influxql"
)

// Store represents an interface access and mutate metadata.
type Store interface {
	ClusterID() uint64
	IsLeader() bool

	CreateContinuousQuery(q *influxql.CreateContinuousQueryStatement) (*ContinuousQueryInfo, error)
	DropContinuousQuery(q *influxql.DropContinuousQueryStatement) error

	Node(id uint64) (*NodeInfo, error)
	NodeByHost(host string) (*NodeInfo, error)
	CreateNode(host string) (*NodeInfo, error)
	DeleteNode(id uint64) error

	Database(name string) (*DatabaseInfo, error)
	CreateDatabase(name string) (*DatabaseInfo, error)
	CreateDatabaseIfNotExists(name string) (*DatabaseInfo, error)
	DropDatabase(name string) error

	RetentionPolicy(database, name string) (*RetentionPolicyInfo, error)
	CreateRetentionPolicy(database string, rp *RetentionPolicyInfo) (*RetentionPolicyInfo, error)
	CreateRetentionPolicyIfNotExists(database string, rp *RetentionPolicyInfo) (*RetentionPolicyInfo, error)
	SetDefaultRetentionPolicy(database, name string) error
	UpdateRetentionPolicy(database, name string, rpu *RetentionPolicyUpdate) (*RetentionPolicyInfo, error)
	DeleteRetentionPolicy(database, name string) error

	ShardGroup(database, policy string, timestamp time.Time) (*ShardGroupInfo, error)
	CreateShardGroupIfNotExists(database, policy string, timestamp time.Time) (*ShardGroupInfo, error)
	DeleteShardGroup(database, policy string, shardID uint64) error

	User(username string) (*UserInfo, error)
	CreateUser(username, password string, admin bool) (*UserInfo, error)
	UpdateUser(username, password string) (*UserInfo, error)
	DeleteUser(username string) error
	SetPrivilege(p influxql.Privilege, username string, dbname string) error
}

// RetentionPolicyUpdate represents retention policy fields to be updated.
type RetentionPolicyUpdate struct {
	Name     *string
	Duration *time.Duration
	ReplicaN *uint32
}

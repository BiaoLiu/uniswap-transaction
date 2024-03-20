package data

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"gorm.io/gorm"

	"uniswap-transaction/app/uniswap/job/internal/conf"
	"uniswap-transaction/pkg/alert"
	rdb "uniswap-transaction/pkg/cache/redis"
	"uniswap-transaction/pkg/database"
	"uniswap-transaction/pkg/database/orm"
	fsBotAPI "uniswap-transaction/pkg/feishu"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewTransaction,
	NewMySQL,
	NewRedisCmd,
	NewFeishuBot,
	NewEthereumClient,
)

type contextTxKey struct{}

type AlertManager struct {
	*alert.Manager
}

// Data .
type Data struct {
	log *log.Helper
	db  *gorm.DB

	Redis          *rdb.Client
	FeishuBot      fsBotAPI.Bot
	AlertManager   *AlertManager
	EthereumClient *ethclient.Client
}

// NewData .
func NewData(logger log.Logger, db *gorm.DB, rc *rdb.Client, feishuBot fsBotAPI.Bot, ethereumClient *ethclient.Client) (*Data, func(), error) {
	log := log.NewHelper(log.With(logger, "module", "example-service/data"))
	d := &Data{
		log:            log,
		db:             db,
		Redis:          rc,
		FeishuBot:      feishuBot,
		EthereumClient: ethereumClient,
	}
	return d, func() {
		// clean
	}, nil
}

func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db.WithContext(ctx)
}

func (d *Data) InTx(ctx context.Context, dbName string, fn func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

// NewTransaction .
func NewTransaction(d *Data) database.Transaction {
	return d
}

// NewMySQL new db.
func NewMySQL(logger log.Logger, conf *conf.Data) *gorm.DB {
	c := &orm.Config{
		MaxOpenConns:  int(conf.Mysql.MaxOpenConns),
		MaxIdleConns:  int(conf.Mysql.MaxIdleConns),
		MaxLifeTime:   conf.Mysql.MaxLifeTime.AsDuration(),
		SlowThreshold: conf.Mysql.SlowThreshold.AsDuration(),
		LogLevel:      conf.Mysql.LogLevel,
		Source:        conf.Mysql.Source,
		Logger:        logger,
	}
	db := orm.NewMySQL(c)
	// if conf.Mysql.Sysbase.EnableAutoMigrate {
	//	db.AutoMigrate(&User{}, &UserWallet{})
	// }
	return db
}

// NewRedisCmd new redis client.
func NewRedisCmd(conf *conf.Data) *rdb.Client {
	c := &rdb.Config{
		Addr:         conf.Redis.Addr,
		Password:     conf.Redis.Password,
		DB:           int(conf.Redis.Db),
		PoolSize:     int(conf.Redis.PoolSize),
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
	}
	return rdb.NewRedisCmd(c)
}

func NewFeishuBot(conf *conf.Data) fsBotAPI.Bot {
	bot := fsBotAPI.NewBot(conf.Feishu.Key, fsBotAPI.WithSecretKey(conf.Feishu.SecretKey))
	return bot
}

func NewEthereumClient(conf *conf.Data) *ethclient.Client {
	ethClient, err := ethclient.Dial(conf.Ethereum.ApiUrl)
	if err != nil {
		panic(err)
	}
	return ethClient
}

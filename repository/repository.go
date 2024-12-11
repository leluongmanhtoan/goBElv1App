package repository

import "program/internal/dbclient"

var SqlClientConnection dbclient.ISqlClientConnection
var RedisClientConnection dbclient.IRedisClientConnection

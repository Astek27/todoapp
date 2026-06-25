package statistics_postgres_repository

import core_postgres_pool "github.com/Astek27/todoapp/internal/core/repository/postgres/pool"

type StatisticRepository struct {
	pool core_postgres_pool.Pool
}

func NewStatisticsRepository(pool core_postgres_pool.Pool) *StatisticRepository {
	return &StatisticRepository{pool: pool}
}
package repository

type SqlUserRepository struct{}

func NewSqlUserRepository() UserRepository {
	return &SqlUserRepository{}
}

func (sur *SqlUserRepository) GetAll() {}

func (sur *SqlUserRepository) Create() {}

func (sur *SqlUserRepository) GetByUUID() {}

func (sur *SqlUserRepository) Update() {}

func (sur *SqlUserRepository) Delete() {}

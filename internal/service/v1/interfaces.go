package v1

type UserService interface {
	GetAll()
	Create()
	GetByUUID()
	Update()
	Delete()
}

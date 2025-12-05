package repository

type UserRepository interface {
    FindByUsername(username string) (interface{}, error) // TODO: return real user later
}

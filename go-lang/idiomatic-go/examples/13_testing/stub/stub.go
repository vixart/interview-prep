// Код с зависимостью-интерфейсом (Entities): именно интерфейс дает возможность
// подменить зависимость в тесте. Бизнес-логика Logic ничего не знает о реализации.
package stub

type User struct{}
type Pet struct {
	Name string
}
type Person struct{}

type Entities interface {
	// пять методов, но тесту нужен только один — заглушка реализует его одного
	GetUser(id string) (User, error)
	GetPets(userID string) ([]Pet, error)
	GetChildren(userID string) ([]Person, error)
	GetFriends(userID string) ([]Person, error)
	SaveUser(user User) error
}

type Logic struct {
	Entities Entities
}

func (l Logic) GetPetNames(userId string) ([]string, error) {
	pets, err := l.Entities.GetPets(userId)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(pets)) // не make([]string, len(pets)) — иначе append добавит к пустым элементам
	// len 0, cap len(pets): append добавляет к пустому срезу, а не к нулевым элементам
	for _, p := range pets {
		out = append(out, p.Name)
	}
	return out, nil
}

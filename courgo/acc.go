package courgo

//Структура аккаунта подписчика для адресной книги
type Acc struct {
	id   uint64
	name string
	dept string
	mail []string
}

/* Установка значений структуры - setters */
func (a *Acc) SetId(id uint64) {
	a.id = id
}

func (a *Acc) SetName(name string) {
	a.name = name
}

func (a *Acc) SetDept(dept string) {
	a.dept = dept
}

func (a *Acc) SetMail(mail []string) {
	a.mail = mail
}
/* END setters */

/* Получение значений структуры - getters*/
func (a Acc) Id() uint64 {
	return a.id
}

func (a Acc) Name() string {
	return a.name
}

func (a Acc) Dept() string {
	return a.dept
}

func (a Acc) Mail() []string {
	return a.mail
}
/* END getters */

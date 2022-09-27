package models

// Публикация, получаемая из RSS.
type Post struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

// Хэш RSS.
type Hash struct {
	ID       int    //id хэша
	NewsHash string //хэш XML файла новостей
	PubTime  int64  //время создания файла
	Link     string //ссылка на XML файл
}

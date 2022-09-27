package rss

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/http"
)

//Возвращает MD5 хэш массива byte
func getMd5Hash(data []byte) (string, error) {
	hasher := md5.New()
	_, err := hasher.Write(data)
	if err != nil {
		return "", errors.New("getMd5Hash error: " + err.Error())
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

//Функция чтения XML файла по ссылке
func readRssBody(l string) ([]byte, error) {
	response, err := http.Get(l)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	XMLdata, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return XMLdata, nil
}

//Обновляет записи хэшей в БД
func (r *Rss) hashUpdate() error {
	err := r.Storage.Hash.Update(r.Hash)
	if err != nil {
		return errors.New("hashUpdate error: " + err.Error())
	}
	return nil
}

//Сравнивает хэш, сохранённый в БД с текущим
func (r *Rss) isHashEqual(data []byte) (bool, error) {
	storedHash, err := r.Storage.Hash.GetByLink(r.Link)
	if err != nil {
		return false, errors.New("HashCheck().GetByLink() error: " + err.Error())
	}
	r.Hash.NewsHash, err = getMd5Hash(data)
	if err != nil {
		return false, errors.New("HashCheck().getMd5Hash error: " + err.Error())
	}
	if storedHash != r.Hash.NewsHash {
		return false, nil
	}
	return true, nil
}

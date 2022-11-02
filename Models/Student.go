package Models

import (
	"praktyka/Config"
)

func GetStudent(id string) (student Student) {
	// Brak error handlingu
	Config.DB.Where("id = ?", id).First(&student)
	return student
}

func GetStudents() (students []Student, err error) {
	// Poproszę o uniwersalny error handling, żeby wszędzie wyglądał tak samo
	// w funkcji niżej jest ten sam error handling a wygląda zupełnie inaczej
	err = Config.DB.Find(&students).Error
	if err != nil {
		return nil, err
	}
	return students, nil
}

func AddNewStudent(student *Student) (err error) {
	// Tak ma wyglądać każdy error handling
	if err = Config.DB.Create(student).Error; err != nil {
		return err
	}
	return nil
}

func UpdateStudent(student *Student, id string) (err error) {
	// Brak error handlingu
	// Te dwa odwołania do Config.DB można złączyć w jedno
	Config.DB.Where("id = ?", id).First(&student)
	Config.DB.Save(&student)
	return nil
}

func DeleteStudent(student *Student, id string) (err error) {
	if err = Config.DB.Where("id = ?", id).Delete(student).Error; err != nil {
		return err
	}
	return nil
}

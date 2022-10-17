package Models

import (
	"praktyka/Config"
)

func GetStudent(id string) (name string) {
	var student Student
	Config.DB.Where("id = ?", id).First(&student)
	return student.Name
}

func AddNewStudent(b *Student) (err error) {
	if err = Config.DB.Create(b).Error; err != nil {
		return err
	}
	return nil
}

func UpdateStudent(b *Student, id string) (err error) {
	var user Student
	Config.DB.Where("id = ?", id).First(&user)
	user.Name = b.Name
	Config.DB.Save(&user)
	return nil
}

func DeleteStudent(b *Student, id string) (err error) {
	Config.DB.Where("id = ?", id).Delete(b)
	return nil
}

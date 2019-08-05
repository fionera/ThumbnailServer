package main

type StringArray []string

func (array StringArray) Contains(value string) bool {
    for _, v := range array {
        if v == value {
            return true
        }
    }

    return false
}

package expenseFunctions;
import "fmt"
func Validate(id *string, name *string, spent *int) (idIsGood bool, nameIsGood bool, spentIsGood bool) {
    idIsGood = false
    nameIsGood = false
    spentIsGood = false
    if id != nil && len(*id) == 36 {
       idIsGood = true
       fmt.Println("id")
    }
    if name != nil && len(*name) > 0 {
        nameIsGood = true
       fmt.Println("name")
    }
    if spent != nil && *spent >= 0{
       spentIsGood = true 
       fmt.Println("spent")
    
    }
    return idIsGood, nameIsGood, spentIsGood
}

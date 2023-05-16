package ui

import (
	"fmt"
	"github.com/samber/lo"
	"strconv"
	"strings"
)

var Numbers = []string{
	` ██████  
██  ████ 
██ ██ ██ 
████  ██ 
 ██████`,
	` ██ 
███ 
 ██ 
 ██ 
 ██ `,
	`██████  
     ██ 
 █████  
██      
███████ `,
	`██████  
     ██ 
 █████  
     ██ 
██████  `,
	`██   ██ 
██   ██ 
███████ 
     ██ 
     ██`,
	`███████ 
██      
███████ 
     ██ 
███████`,
	` ██████  
██       
███████  
██    ██ 
 ██████`,
	`███████ 
     ██ 
    ██  
   ██   
   ██`,
	` █████  
██   ██ 
 █████  
██   ██ 
 █████`,
	` █████  
██   ██ 
 ██████ 
     ██ 
 █████`,
}

func GetNumber(number int) string {
	return strings.Join(lo.Map([]rune(fmt.Sprint(number)), func(char rune, index int) string {
		num, _ := strconv.Atoi(string(char))
		return Numbers[num]
	}), "")
}

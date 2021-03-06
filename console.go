/*
Copyright © 2018-2019 Sad Pencil
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"strings"
	"syscall"
)

func cartman() {
	Settings := NewSettings()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Looks like the config file doesn't exist.")
	fmt.Println("That's okay if you just want to do a Portal authentication for once, or to generate a config file.")
	fmt.Println("A few questions need to be answered. Leave it blank if you want the default answer in the bracket.")

	var err error

	for {
		fmt.Println()
		fmt.Println("Question 0. Are you ready? [Yes]")
		ans, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		ans = strings.TrimSpace(ans)
		if ans == "" {
			fmt.Println("Cool. Let's get started.")
			break
		} else {
			fmt.Println("Ah, god dammit. I told you that you can just LEAVE IT BLANK if you want the default answer. Now try again.")
		}
	}

	for {
		fmt.Println()
		fmt.Println("Question 1. What's your username? []")

		Settings.Account.Username, err = reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		err = checkUsername(&Settings)
		if err != nil {
			fmt.Println(err)
		} else {
			break
		}
	}
	for {
		fmt.Println()
		fmt.Println("Question 2. What's your password? []")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Println() // it's necessary to add a new line after user's input
		if err != nil {
			panic(err)
		} else {
			Settings.Account.Password = string(bytePassword)
		}
		err = checkPassword(&Settings)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Great. Your password contains", fmt.Sprint(len(Settings.Account.Password)), "characters.")
			break
		}
	}

	for {
		fmt.Println()
		fmt.Println("Question 3. What's the authentication server's ip address? [" + DEFAULT_AUTH_SERVER + "]")
		fmt.Println("Hint: You can also write down the server's FQDN if necessary. You may specify either an IPv4 or IPv6 server.")
		Settings.Account.AuthServer, err = reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		err = checkAuthServer(&Settings)
		if err != nil {
			fmt.Println(err)
		} else {
			if strings.Count(Settings.Account.AuthServer, ":") >= 2 {
				fmt.Println("Hint: Add a pair of [] with the IPv6 address. Omit this hint if you have already done so.")
				fmt.Println("Example 1 \t [2001:250:5800:11::1]")
				fmt.Println("Example 2 \t [2001:250:5800:11::1]:8080")
			}

			break
		}
	}

	for {
		fmt.Println()
		fmt.Println("Question 4. Does the authentication server use HTTP protocol, or HTTPS? [" + DEFAULT_AUTH_SCHEME + "]")
		Settings.Account.Scheme, err = reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		err = checkScheme(&Settings)
		if err != nil {
			fmt.Println(err)
		} else {
			break
		}
	}

	for {
		fmt.Println()
		fmt.Println("Question 5. Do you want to log out the network when the program get terminated? [y/N]")
		yesOrNoStr, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		yesOrNoStr = strings.ToLower(strings.TrimSpace(yesOrNoStr))
		if yesOrNoStr == "" || yesOrNoStr == "n" {
			Settings.Control.LogoutWhenExit = false
			break
		} else if yesOrNoStr == "y" {
			Settings.Control.LogoutWhenExit = true
			break
		} else {
			fmt.Println("All you need to do is to answer me yes or no. Don't be a pussy.")
		}
	}
	//for {
	//	fmt.Println()
	//	fmt.Println("Question 5.")
	//	var ips []string
	//	var interfaceStrings []string
	//	{
	//		_, ip, err := getSduUserInfo(Settings.Account.Scheme,
	//			Settings.Account.AuthServer, "")
	//
	//		//ip, err := getIPFromChallenge(Settings.Account.Scheme,
	//		//	Settings.Account.AuthServer,
	//		//	Settings.Account.Username,
	//		//	"",
	//		//	"",
	//		//	false)
	//		if err != nil {
	//			fmt.Println(err)
	//			ips = append(ips, "")
	//			//忽略错误继续
	//			fmt.Println("Warning: Failed to connected to the authentication server.")
	//			fmt.Println()
	//		} else {
	//			ips = append(ips, ip)
	//			interfaceStrings = append(interfaceStrings, "")
	//			fmt.Println("["+fmt.Sprint(len(ips)-1)+"]", "\t", ip, "\t", "[Auto detect]")
	//		}
	//	}
	//
	//	interfaces, err := net.Interfaces()
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	for _, interfaceWtf := range interfaces {
	//		ip, err := GetIPv4FromInterface(interfaceWtf.Name)
	//		if err == nil {
	//			ips = append(ips, ip)
	//			interfaceStrings = append(interfaceStrings, interfaceWtf.Name)
	//			fmt.Println("["+fmt.Sprint(len(ips)-1)+"]", "\t", ip, "\t", interfaceWtf.Name)
	//		}
	//	}
	//
	//	if len(ips) == 0 {
	//		fmt.Println("There is not even a network interface with a valid IPv4 address.")
	//		fmt.Println("Screw you guys, I'm going home.")
	//		return
	//	}
	//
	//	fmt.Println("Network interfaces are listed above. Which one is connected to the Portal network? [0]")
	//	fmt.Println("Hint: It is recommended to choose auto detect as long as the reported IP address is correct, even if you have only one network interface available. Choosing a specific network interface instead of auto detect causes different behavior and it's not recommended unless you have multiple default routes via multiple network interfaces.")
	//	if probablyIPv6 {
	//		fmt.Println("Hint: You might have specified an IPv6 authentication server. Choose auto detect, as all but auto detect will not work.")
	//	}
	//	choice, err := reader.ReadString('\n')
	//	if err != nil {
	//		panic(err)
	//	}
	//	choice = strings.TrimSpace(choice)
	//
	//	var choiceID int
	//	if choice == "" {
	//		choiceID = 0
	//	} else {
	//		choiceID, err = strconv.Atoi(choice)
	//		if err != nil {
	//			fmt.Println(err)
	//			continue
	//		}
	//	}
	//
	//	if choiceID == 0 {
	//		Settings.Network.Interface = ""
	//		Settings.Network.CustomIP = ""
	//		Settings.Network.StrictMode = false
	//	} else if choiceID > 0 && choiceID < len(interfaceStrings) {
	//		Settings.Network.Interface = interfaceStrings[choiceID]
	//		Settings.Network.CustomIP = ""
	//		Settings.Network.StrictMode = true
	//	} else {
	//		fmt.Println("You think I'm a retard? Fuck you, asshole.")
	//		continue
	//	}
	//
	//	if err != nil {
	//		fmt.Println(err)
	//	} else {
	//		break
	//	}
	//}
	var saveFile bool
	for {
		fmt.Println()
		fmt.Println("That's all the information needed. Would you like to save it to a configuration file? [y/n]")
		yesOrNoStr, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		yesOrNoStr = strings.ToLower(strings.TrimSpace(yesOrNoStr))

		if yesOrNoStr == "y" {
			saveFile = true
			break
		} else if yesOrNoStr == "n" {
			saveFile = false
			break
		} else {
			fmt.Println("All you need to do is to answer me yes or no. Don't be a pussy.")
		}
	}

	if saveFile {
		fmt.Println()
		fmt.Println("Where to save the file? [" + DEFAULT_CONFIG_FILENAME + "]")
		fmt.Println("Hint: If the program doesn't have permission to write, it will crash.")
		filename, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		filename = strings.TrimSpace(filename)
		if filename == "" {
			filename = DEFAULT_CONFIG_FILENAME
		}
		f, err := os.Create(filename)
		defer f.Close()
		if err != nil {
			fmt.Println(err)
		} else {
			jsonBytes, err := json.Marshal(Settings)
			if err != nil {
				fmt.Println(err)
				saveFile = false
			} else {
				_, err := f.WriteString(string(jsonBytes))
				if err != nil {
					fmt.Println(err)
					saveFile = false
				} else {
					fmt.Println(`File saved. You may re-run the program with the "-c" flag.`)
				}
			}
		}
	}
	if !saveFile {
		var yesOrNo bool
		for {
			fmt.Println()
			fmt.Println("File not saved. Login to the network right away? [Y/n]")
			yesOrNoStr, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			yesOrNoStr = strings.ToLower(strings.TrimSpace(yesOrNoStr))

			if yesOrNoStr == "y" || yesOrNoStr == "" {
				yesOrNo = true
				break
			} else if yesOrNoStr == "n" {
				yesOrNo = false
				break
			} else {
				fmt.Println("All you need to do is to answer me yes or no. Don't be a pussy.")
			}
		}
		if yesOrNo {
			fmt.Println()
			fmt.Println("Log in via web portal...")
			err := _login(&Settings)
			if err != nil {
				log.Println("Login failed.", err)
			}
		}
	}
}

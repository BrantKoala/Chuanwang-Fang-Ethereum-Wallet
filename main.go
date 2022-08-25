package main

import (
	"00097eth-relay/dao"
	"00097eth-relay/tool"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
)

func showMenu() {
	fmt.Println("1\tCreate an EOA")
	fmt.Println("2\tQuery ETH balance")
	fmt.Println("3\tTransfer ETH")
	fmt.Println("4\tQuery ERC20 token balance")
	fmt.Println("5\tTransfer ERC20 token")
	fmt.Println("6\tQuery transaction information")
	fmt.Println("7\tQuery block information")
	fmt.Println("8\tScan and save new Blocks")
	fmt.Println("9\tExit")
	fmt.Print("Please choose a number:")
}
func isAddressLengthCorrect(address string) bool {
	if address == "" || len(address) != 42 {
		fmt.Println("Illegal address")
		return false
	} else {
		return true
	}
}
func getBasicPaymentInfoAndUnlockWallet() (string, string, string, uint64, float64, error) {
	var from, password, payee, amount string
	fmt.Println("Please ensure the keystore file is in directory 'keystore'")
	fmt.Print("Your account address:")
	_, err := fmt.Scanln(&from)
	if err != nil {
		fmt.Println("input account address failed")
		return "", "", "", 0, 0, err
	}
	if !isAddressLengthCorrect(from) {
		return "", "", "", 0, 0, err
	}
	fmt.Print("password:")
	_, err = fmt.Scanln(&password)
	if err != nil {
		fmt.Println("input password failed")
		return "", "", "", 0, 0, err
	}
	err = tool.UnlockETHWallet("./keystore", from, password)
	if err != nil {
		fmt.Println("Unlock wallet failed,", err.Error())
		return "", "", "", 0, 0, err
	}
	fmt.Print("payee address:")
	_, err = fmt.Scanln(&payee)
	if err != nil {
		fmt.Println("input payee address failed")
		return "", "", "", 0, 0, err
	}
	if !isAddressLengthCorrect(payee) {
		return "", "", "", 0, 0, err
	}
	fmt.Print("amount:")
	_, err = fmt.Scanln(&amount)
	if err != nil {
		fmt.Println("input amount failed")
		return "", "", "", 0, 0, err
	}
	var gasLimit uint64
	fmt.Print("gas limit:")
	_, err = fmt.Scanln(&gasLimit)
	if err != nil {
		fmt.Println("input gas limit failed")
		return "", "", "", 0, 0, err
	}
	var gasPrice float64
	fmt.Print("gas price (Gwei):")
	_, err = fmt.Scanln(&gasPrice)
	if err != nil {
		fmt.Println("input gas price failed")
		return "", "", "", 0, 0, err
	}
	return from, payee, amount, gasLimit, gasPrice, nil
}
func main() {
	fmt.Println("Welcome to Chuanwang Fang Ethereum Wallet!")
	_ = os.Mkdir("./keystore", os.ModePerm)
	nodeUrl := "https://ropsten.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
changeNodeURL:
	fmt.Print("The program currently works on Ropsten testnet? Do you want to change? 1=>Yes; 2=>No ")
	var willChangeNodeURL uint
	_, err := fmt.Scanln(&willChangeNodeURL)
	if err != nil {
		fmt.Println("input selection error")
		return
	}
	if willChangeNodeURL == 1 {
		fmt.Print("New node link:")
		_, err = fmt.Scanln(&nodeUrl)
		if err != nil {
			fmt.Println("Input new node link failed")
			return
		}
	} else if willChangeNodeURL != 2 {
		fmt.Println("Illegal selection")
		goto changeNodeURL
	}
	var selection uint
foreverLoop:
	for {
		showMenu()
		_, err := fmt.Scanln(&selection)
		if err != nil {
			fmt.Println("input selection failed")
			continue
		}
		switch selection {
		case 1:
			fmt.Print("Set password (no less than 6 letters):")
			var password string
			_, err = fmt.Scanln(&password)
			if err != nil {
				fmt.Println("input password failed")
				break
			}
			address, err := NewETHRPCRequester(nodeUrl).CreateETHWallet(password)
			if err != nil {
				fmt.Println("Create Wallet failed, ", err.Error())
			} else {
				fmt.Println("Wallet created, address:", address)
				fmt.Println("You can find the keystore file in directory 'keystore'")
			}
		case 2:
			var accountAddress string
			fmt.Print("account address:")
			_, err = fmt.Scanln(&accountAddress)
			if err != nil {
				fmt.Println("input account address failed")
				break
			}
			if !isAddressLengthCorrect(accountAddress) {
				break
			}
			balance, err := NewETHRPCRequester(nodeUrl).GetETHBalance(accountAddress)
			if err != nil {
				fmt.Println("query ETH balance failed, info: ", err.Error())
			} else {
				fmt.Println("balance:", balance)
			}
		case 3:
			from, payee, amount, gasLimit, nGwei, err := getBasicPaymentInfoAndUnlockWallet()
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			gasPrice := uint64(nGwei * 1_000_000_000)
			requester := NewETHRPCRequester(nodeUrl)
			txHash, err := requester.SendETHTransaction(from, payee, amount, gasLimit, gasPrice)
			if err != nil {
				fmt.Println("ETH transfer failed, ", err.Error())
				break
			} else {
				fmt.Println("ETH transfer successfully created, transaction hash:", txHash)
			}
			balance, err := requester.GetETHBalance(from)
			if err != nil {
				fmt.Println("query ETH balance failed, info: ", err.Error())
			} else {
				fmt.Println("Your current ETH balance:", balance)
			}
		case 4:
			var accountAddress, contractAddress string
			fmt.Print("account address:")
			_, err = fmt.Scanln(&accountAddress)
			if err != nil {
				fmt.Println("input account address failed")
				break
			}
			if !isAddressLengthCorrect(accountAddress) {
				break
			}
			fmt.Print("contract address:")
			_, err = fmt.Scanln(&contractAddress)
			if err != nil {
				fmt.Println("input contract address failed")
				break
			}
			if !isAddressLengthCorrect(contractAddress) {
				break
			}
			var contractDecimal int
			fmt.Print("contract decimal:")
			_, err = fmt.Scanln(&contractDecimal)
			if err != nil {
				fmt.Println("input contract decimal failed")
				break
			}
			paramsGEB := []ERC20BalanceRpcReq{{
				ContractAddress: contractAddress,
				UserAddress:     accountAddress,
				ContractDecimal: contractDecimal,
			}}
			balance, err := NewETHRPCRequester(nodeUrl).GetERC20Balances(paramsGEB)
			if err != nil {
				fmt.Println("query ERC20 token failed, info: ", err.Error())
			} else {
				fmt.Println("balance:", balance)
			}
		case 5:
			from, payee, amount, gasLimit, nGwei, err := getBasicPaymentInfoAndUnlockWallet()
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			var contractAddress string
			fmt.Print("contract address:")
			_, err = fmt.Scanln(&contractAddress)
			if err != nil {
				fmt.Println("input contract address failed")
				break
			}
			if !isAddressLengthCorrect(contractAddress) {
				break
			}
			var decimal int
			fmt.Print("decimal digits:")
			_, err = fmt.Scanln(&decimal)
			if err != nil {
				fmt.Println("input decimal digits failed")
				break
			}
			gasPrice := uint64(nGwei * 1_000_000_000)
			requester := NewETHRPCRequester(nodeUrl)
			txHash, err := requester.SendERC20Transaction(from, contractAddress, payee, amount, gasLimit, gasPrice, decimal)
			if err != nil {
				fmt.Println("ERC20 token transfer failed, ", err.Error())
				break
			} else {
				fmt.Println("ERC20 token transfer successfully created, transaction hash:", txHash)
			}
			paramsGEB := []ERC20BalanceRpcReq{{
				ContractAddress: contractAddress,
				UserAddress:     from,
				ContractDecimal: decimal,
			}}
			balance, err := requester.GetERC20Balances(paramsGEB)
			if err != nil {
				fmt.Println("query ERC20 token balance failed, info: ", err.Error())
			} else {
				fmt.Println("Your current ERC20 token balance:", balance[0])
			}
		case 6:
			fmt.Print("input transaction hash:")
			var txHash string
			_, err = fmt.Scanln(&txHash)
			if err != nil {
				fmt.Println("input transaction hash failed", err.Error())
				break
			}
			if txHash == "" || len(txHash) != 66 {
				fmt.Println("Illegal transaction hash")
				break
			}
			txInfo, err := NewETHRPCRequester(nodeUrl).GetTransactionByHash(txHash)
			if err != nil {
				fmt.Println("query transaction failed,", err.Error())
			} else {
				jsonData, _ := json.Marshal(txInfo)
				fmt.Println("transaction info:", string(jsonData))
			}
		case 7:
			fmt.Print("query by 1==>block number, 2==>block hash? ")
			var choice uint8
			_, err = fmt.Scanln(&choice)
			if err != nil {
				fmt.Println("input choice failed,", err.Error())
				break
			}
			switch choice {
			case 1:
				var blockNumber uint64
				fmt.Print("block number:")
				_, err = fmt.Scanln(&blockNumber)
				if err != nil {
					fmt.Println("input block number failed,", err.Error())
					break
				}
				fullBlock, err := NewETHRPCRequester(nodeUrl).GetBlockInfoByNumber(new(big.Int).SetUint64(blockNumber))
				if err != nil {
					fmt.Println("get block info failed, ", err.Error())
				} else {
					jsonData, _ := json.Marshal(fullBlock)
					fmt.Println("block info:", string(jsonData))
				}
			case 2:
				fmt.Print("block hash:")
				var blockHash string
				_, err = fmt.Scanln(&blockHash)
				if err != nil {
					fmt.Println("input block hash failed,", err.Error())
					break
				}
				fullBlock, err := NewETHRPCRequester(nodeUrl).GetBlockInfoByHash(blockHash)
				if err != nil {
					fmt.Println("get block info failed, ", err.Error())
				} else {
					jsonData, _ := json.Marshal(fullBlock)
					fmt.Println("block info:", string(jsonData))
				}
			default:
				fmt.Println("Illegal choice!")
			}
		case 8:
			var MySQLUsername, MySQLPassword string
			fmt.Print("MySQL username:")
			_, err = fmt.Scanln(&MySQLUsername)
			if err != nil {
				fmt.Println("input MySQL username failed")
				break
			}
			fmt.Print("MySQL password:")
			_, err = fmt.Scanln(&MySQLPassword)
			if err != nil {
				fmt.Println("input MySQL password failed")
				break
			}
			option := dao.MySQLOptions{
				Hostname:           "127.0.0.1",
				Port:               "3306",
				User:               MySQLUsername,
				Password:           MySQLPassword,
				DBName:             "eth_relay",
				TablePrefix:        "eth_",
				MaxOpenConnections: 10,
				MaxIdleConnections: 5,
				ConnMaxLifetime:    15,
			}
			tables := []interface{}{dao.Block{}, dao.Transaction{}}
			mysqlConnector := dao.NewMySQLConnector(&option, tables)
			scanner := NewBlockScanner(*NewETHRPCRequester(nodeUrl), mysqlConnector)
			err = scanner.Start()
			if err != nil {
				panic(err)
			}
			select {}
		case 9:
			break foreverLoop
		default:
			fmt.Println("Illegal selection!")
		}
		fmt.Println("====================================================")
	}
	fmt.Println("Bye bye!")
	select {}
}

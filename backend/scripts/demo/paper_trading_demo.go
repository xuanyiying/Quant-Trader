// 模拟交易 Demo 脚本
// 用于演示和验证模拟交易功能
//
// 使用方法:
//   go run scripts/demo/paper_trading_demo.go
//
// 该脚本会:
// 1. 创建/获取模拟交易账户
// 2. 查询账户余额
// 3. 创建市价买单
// 4. 创建限价卖单
// 5. 查询持仓
// 6. 重置账户

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	baseURL = "http://localhost:8080"
	// 测试用的 JWT token (实际使用时需要登录获取)
	testToken = "Bearer test-token"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("   模拟交易功能 Demo")
	fmt.Println("========================================")
	fmt.Println()

	// 检查服务是否运行
	if !checkHealth() {
		fmt.Println("❌ 服务未启动，请先启动后端服务:")
		fmt.Println("   go run cmd/server/main.go")
		os.Exit(1)
	}

	fmt.Println("✅ 服务运行正常")
	fmt.Println()

	// 1. 获取账户信息
	fmt.Println("📊 Step 1: 获取模拟账户信息")
	getPaperAccount()
	fmt.Println()

	// 2. 创建市价买单
	fmt.Println("📈 Step 2: 创建市价买单 (BTCUSDT)")
	createMarketOrder("BTCUSDT", "buy", "1.0")

	fmt.Println()

	// 3. 创建限价卖单
	fmt.Println("📉 Step 3: 创建限价卖单 (BTCUSDT)")
	createLimitOrder("BTCUSDT", "sell", "0.5", "60000")
	fmt.Println()

	// 4. 查询持仓
	fmt.Println("💼 Step 4: 查询持仓")
	getPositions()
	fmt.Println()

	// 5. 再次查询账户
	fmt.Println("📊 Step 5: 再次查询账户余额")
	getPaperAccount()
	fmt.Println()

	// 6. 重置账户
	fmt.Println("🔄 Step 6: 重置模拟账户")
	resetPaperAccount()
	fmt.Println()

	// 7. 验证重置结果
	fmt.Println("📊 Step 7: 验证重置后的账户")
	getPaperAccount()
	fmt.Println()

	fmt.Println("========================================")
	fmt.Println("   Demo 完成!")
	fmt.Println("========================================")
}

func checkHealth() bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(baseURL + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func getPaperAccount() {
	req, _ := http.NewRequest("GET", baseURL+"/api/v1/paper/account", nil)
	req.Header.Set("Authorization", testToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("   ❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("   ✅ 账户余额: %v USDT\n", result["balance"])
	} else {
		fmt.Printf("   ❌ 获取失败: %v\n", result)
	}
}

func createMarketOrder(symbol, side, qty string) {
	body := map[string]string{
		"symbol": symbol,
		"side":   side,
		"type":   "market",
		"qty":    qty,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", baseURL+"/api/v1/paper/orders", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", testToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("   ❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode == http.StatusCreated {
		fmt.Printf("   ✅ 订单创建成功, ID: %v\n", result["id"])
	} else {
		fmt.Printf("   ❌ 创建失败: %v\n", result)
	}
}

func createLimitOrder(symbol, side, qty, price string) {
	body := map[string]string{
		"symbol": symbol,
		"side":   side,
		"type":   "limit",
		"qty":    qty,
		"price":  price,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", baseURL+"/api/v1/paper/orders", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", testToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("   ❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode == http.StatusCreated {
		fmt.Printf("   ✅ 限价单创建成功, ID: %v\n", result["id"])
	} else {
		fmt.Printf("   ❌ 创建失败: %v\n", result)
	}
}

func getPositions() {
	req, _ := http.NewRequest("GET", baseURL+"/api/v1/paper/positions", nil)
	req.Header.Set("Authorization", testToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("   ❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode == http.StatusOK {
		if len(result) == 0 {
			fmt.Println("   ℹ️  当前没有持仓")
		} else {
			fmt.Printf("   ✅ 持仓数量: %d\n", len(result))
			for _, pos := range result {
				fmt.Printf("      - %s: %v @ %v\n", pos["symbol"], pos["qty"], pos["avg_price"])
			}
		}
	} else {
		fmt.Printf("   ❌ 查询失败: %v\n", result)
	}
}

func resetPaperAccount() {
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/paper/account/reset", nil)
	req.Header.Set("Authorization", testToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("   ❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("   ✅ 账户重置成功, 新余额: %v USDT\n", result["balance"])
	} else {
		fmt.Printf("   ❌ 重置失败: %v\n", result)
	}
}
